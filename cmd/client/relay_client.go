package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var errRelayClientClosed = errors.New("relay client closed")

type relayConnectionConfig struct {
	Dial           func() (*websocket.Conn, error)
	InitMessage    TerminalMessage
	InboundHandler func(TerminalMessage)
	OnConnected    func(reconnected bool)
	OnDisconnected func(hadConnectedBefore bool, err error)
	BaseRetryDelay time.Duration
	MaxRetryDelay  time.Duration
	WriteTimeout   time.Duration
}

type relayConnectionEventType int

const (
	relayEventInbound relayConnectionEventType = iota
	relayEventReadFailed
	relayEventWriteFailed
)

type relayConnectionEvent struct {
	generation uint64
	eventType  relayConnectionEventType
	message    TerminalMessage
	err        error
}

type relayConnectionManager struct {
	config     relayConnectionConfig
	outbound   chan TerminalMessage
	events     chan relayConnectionEvent
	shutdownCh chan *TerminalMessage
	done       chan struct{}
	closed     atomic.Bool
	closeOnce  sync.Once
	initMu     sync.RWMutex
	initMsg    TerminalMessage
}

func newRelayConnectionManager(config relayConnectionConfig) *relayConnectionManager {
	if config.BaseRetryDelay <= 0 {
		config.BaseRetryDelay = time.Second
	}
	if config.MaxRetryDelay <= 0 {
		config.MaxRetryDelay = 30 * time.Second
	}
	if config.WriteTimeout <= 0 {
		config.WriteTimeout = 5 * time.Second
	}

	manager := &relayConnectionManager{
		config:     config,
		outbound:   make(chan TerminalMessage, 256),
		events:     make(chan relayConnectionEvent, 32),
		shutdownCh: make(chan *TerminalMessage, 1),
		done:       make(chan struct{}),
		initMsg:    config.InitMessage,
	}

	return manager
}

func (m *relayConnectionManager) Start() {
	go m.loop()
}

func (m *relayConnectionManager) Send(message TerminalMessage) error {
	if m.closed.Load() {
		return errRelayClientClosed
	}

	if message.Type == "resize" {
		m.updateInitSize(message.Rows, message.Cols)
	}

	select {
	case m.outbound <- message:
		return nil
	case <-m.done:
		return errRelayClientClosed
	}
}

func (m *relayConnectionManager) Shutdown(closeMessage *TerminalMessage) {
	m.closeOnce.Do(func() {
		m.closed.Store(true)
		select {
		case m.shutdownCh <- closeMessage:
		case <-m.done:
		}
	})

	<-m.done
}

func (m *relayConnectionManager) updateInitSize(rows, cols int) {
	if rows <= 0 || cols <= 0 {
		return
	}

	m.initMu.Lock()
	defer m.initMu.Unlock()

	m.initMsg.Rows = rows
	m.initMsg.Cols = cols
}

func (m *relayConnectionManager) currentInitMessage() TerminalMessage {
	m.initMu.RLock()
	defer m.initMu.RUnlock()
	return m.initMsg
}

func (m *relayConnectionManager) loop() {
	defer close(m.done)

	var (
		conn               *websocket.Conn
		writer             *websocketWriter
		pending            []TerminalMessage
		retryTimer         *time.Timer
		retryCh            <-chan time.Time
		generation         uint64
		retrying           bool
		hadConnectedBefore bool
		retryAttempt       int
	)

	stopRetryTimer := func() {
		if retryTimer == nil {
			retryCh = nil
			return
		}
		if !retryTimer.Stop() {
			select {
			case <-retryTimer.C:
			default:
			}
		}
		retryTimer = nil
		retryCh = nil
	}

	scheduleRetry := func(delay time.Duration) {
		stopRetryTimer()
		retryTimer = time.NewTimer(delay)
		retryCh = retryTimer.C
	}

	scheduleReconnect := func(err error) {
		if !retrying && m.config.OnDisconnected != nil {
			m.config.OnDisconnected(hadConnectedBefore, err)
		}

		retrying = true
		retryAttempt++
		delay := m.config.BaseRetryDelay << (retryAttempt - 1)
		if delay > m.config.MaxRetryDelay {
			delay = m.config.MaxRetryDelay
		}
		scheduleRetry(delay)
	}

	handleConnectionLoss := func(err error) {
		if writer == nil && conn == nil {
			return
		}

		if writer != nil {
			writer.Shutdown()
			writer = nil
		} else if conn != nil {
			_ = conn.Close()
		}
		conn = nil

		scheduleReconnect(err)
	}

	flushPending := func() {
		for writer != nil && len(pending) > 0 {
			data, err := json.Marshal(pending[0])
			if err != nil {
				log.Printf("❌ 序列化待发送消息失败: %v", err)
				pending = pending[1:]
				continue
			}

			if err := writer.WriteWithTimeout(data, m.config.WriteTimeout); err != nil {
				handleConnectionLoss(err)
				return
			}

			pending = pending[1:]
		}
	}

	connect := func() {
		connCandidate, err := m.config.Dial()
		if err != nil {
			scheduleReconnect(err)
			return
		}

		writerCandidate := newWebsocketWriter(connCandidate)
		writerCandidate.StartHeartbeat()

		generation++
		m.startReaderLoop(connCandidate, generation)
		m.startWriterMonitor(writerCandidate, generation)

		initPayload, err := json.Marshal(m.currentInitMessage())
		if err != nil {
			writerCandidate.Shutdown()
			scheduleReconnect(err)
			return
		}

		if err := writerCandidate.WriteWithTimeout(initPayload, m.config.WriteTimeout); err != nil {
			writerCandidate.Shutdown()
			scheduleReconnect(err)
			return
		}

		conn = connCandidate
		writer = writerCandidate
		stopRetryTimer()

		reconnected := hadConnectedBefore
		hadConnectedBefore = true
		retrying = false
		retryAttempt = 0

		if m.config.OnConnected != nil {
			m.config.OnConnected(reconnected)
		}

		flushPending()
	}

	scheduleRetry(0)

	for {
		select {
		case message := <-m.outbound:
			pending = append(pending, message)
			flushPending()
		case event := <-m.events:
			if event.generation != generation {
				continue
			}
			switch event.eventType {
			case relayEventInbound:
				if m.config.InboundHandler != nil {
					m.config.InboundHandler(event.message)
				}
			case relayEventReadFailed, relayEventWriteFailed:
				handleConnectionLoss(event.err)
			}
		case <-retryCh:
			connect()
		case closeMessage := <-m.shutdownCh:
			stopRetryTimer()
			if writer != nil && closeMessage != nil {
				closePayload, err := json.Marshal(closeMessage)
				if err == nil {
					if err := writer.WriteWithTimeout(closePayload, 2*time.Second); shouldLogWebSocketError(err) {
						log.Printf("❌ 发送关闭消息失败: %v", err)
					}
				}
			}
			if writer != nil {
				writer.Shutdown()
			} else if conn != nil {
				_ = conn.Close()
			}
			return
		}
	}
}

func (m *relayConnectionManager) startReaderLoop(conn *websocket.Conn, generation uint64) {
	go func() {
		for {
			var message TerminalMessage
			if err := conn.ReadJSON(&message); err != nil {
				m.emitEvent(relayConnectionEvent{
					generation: generation,
					eventType:  relayEventReadFailed,
					err:        err,
				})
				return
			}

			m.emitEvent(relayConnectionEvent{
				generation: generation,
				eventType:  relayEventInbound,
				message:    message,
			})
		}
	}()
}

func (m *relayConnectionManager) startWriterMonitor(writer *websocketWriter, generation uint64) {
	go func() {
		for err := range writer.Errors() {
			if err == nil {
				continue
			}
			m.emitEvent(relayConnectionEvent{
				generation: generation,
				eventType:  relayEventWriteFailed,
				err:        err,
			})
			return
		}
	}()
}

func (m *relayConnectionManager) emitEvent(event relayConnectionEvent) {
	select {
	case m.events <- event:
	case <-m.done:
	}
}

func shouldLogRelaySendError(err error) bool {
	return err != nil && !errors.Is(err, errRelayClientClosed)
}
