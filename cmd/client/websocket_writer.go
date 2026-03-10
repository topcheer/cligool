package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var errWebSocketWriterStopped = errors.New("websocket writer stopped")

type websocketWriteRequest struct {
	data []byte
	ack  chan error
}

type websocketWriter struct {
	conn          *websocket.Conn
	requests      chan websocketWriteRequest
	done          chan struct{}
	errors        chan error
	writerStop    chan struct{}
	heartbeatStop chan struct{}
	closing       atomic.Bool
	stopOnce      sync.Once
	heartbeatOnce sync.Once
	errorOnce     sync.Once
}

func newWebsocketWriter(conn *websocket.Conn) *websocketWriter {
	w := &websocketWriter{
		conn:          conn,
		requests:      make(chan websocketWriteRequest, 100),
		done:          make(chan struct{}),
		errors:        make(chan error, 1),
		writerStop:    make(chan struct{}),
		heartbeatStop: make(chan struct{}),
	}

	go func() {
		defer close(w.done)
		defer w.closeErrors()

		for {
			select {
			case req := <-w.requests:
				err := w.conn.WriteMessage(websocket.TextMessage, req.data)
				if req.ack != nil {
					req.ack <- err
					close(req.ack)
				}
				if err != nil {
					w.signalError(err)
					if shouldLogWebSocketError(err) {
						log.Printf("❌ WebSocket写入失败: %v", err)
					}
					return
				}
			case <-w.writerStop:
				return
			}
		}
	}()

	return w
}

func (w *websocketWriter) Write(data []byte) error {
	return w.enqueue(websocketWriteRequest{data: data})
}

func (w *websocketWriter) WriteWithTimeout(data []byte, timeout time.Duration) error {
	ack := make(chan error, 1)
	if err := w.enqueue(websocketWriteRequest{data: data, ack: ack}); err != nil {
		return err
	}

	if timeout <= 0 {
		return <-ack
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case err := <-ack:
		return err
	case <-w.done:
		return errWebSocketWriterStopped
	case <-timer.C:
		return fmt.Errorf("websocket write timed out after %s", timeout)
	}
}

func (w *websocketWriter) enqueue(req websocketWriteRequest) error {
	if w.closing.Load() {
		return errWebSocketWriterStopped
	}

	select {
	case w.requests <- req:
		return nil
	case <-w.done:
		return errWebSocketWriterStopped
	}
}

func (w *websocketWriter) StartHeartbeat() {
	w.conn.SetPingHandler(func(appData string) error {
		err := w.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(5*time.Second))
		if err != nil {
			w.signalError(err)
		}
		if shouldLogWebSocketError(err) {
			log.Printf("❌ 发送pong失败: %v", err)
		}
		return err
	})

	w.conn.SetPongHandler(func(string) error {
		return nil
	})

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := w.conn.WriteControl(websocket.PingMessage, []byte("heartbeat"), time.Now().Add(5*time.Second))
				if err != nil {
					w.signalError(err)
				}
				if shouldLogWebSocketError(err) {
					log.Printf("❌ 发送ping失败: %v", err)
				}
				if err != nil {
					return
				}
			case <-w.heartbeatStop:
				return
			}
		}
	}()
}

func (w *websocketWriter) StopHeartbeat() {
	w.heartbeatOnce.Do(func() {
		close(w.heartbeatStop)
	})
}

func (w *websocketWriter) Errors() <-chan error {
	return w.errors
}

func (w *websocketWriter) Shutdown() {
	w.stopOnce.Do(func() {
		w.closing.Store(true)
		w.StopHeartbeat()
		close(w.writerStop)

		select {
		case <-w.done:
		case <-time.After(2 * time.Second):
		}

		if err := w.conn.Close(); shouldLogWebSocketError(err) {
			log.Printf("❌ 关闭WebSocket连接失败: %v", err)
		}
	})
}

func (w *websocketWriter) signalError(err error) {
	if err == nil {
		return
	}

	w.errorOnce.Do(func() {
		w.errors <- err
		close(w.errors)
	})
}

func (w *websocketWriter) closeErrors() {
	w.errorOnce.Do(func() {
		close(w.errors)
	})
}

func shouldLogWebSocketError(err error) bool {
	return err != nil && !errors.Is(err, errWebSocketWriterStopped) && !isExpectedWebSocketShutdown(err)
}

func isExpectedWebSocketShutdown(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, net.ErrClosed) {
		return true
	}

	if websocket.IsCloseError(err,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseNoStatusReceived,
	) {
		return true
	}

	msg := err.Error()
	return strings.Contains(msg, "use of closed network connection") ||
		strings.Contains(msg, "websocket: close sent")
}
