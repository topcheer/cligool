package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRelayConnectionManagerFlushesBufferedMessagesOnConnect(t *testing.T) {
	var relayAvailable atomic.Bool
	messageCh := make(chan TerminalMessage, 8)

	server, wsURL := startRelayTestServer(t, messageCh, nil)
	defer server.Close()

	dialer := *websocket.DefaultDialer
	manager := newRelayConnectionManager(relayConnectionConfig{
		Dial: func() (*websocket.Conn, error) {
			if !relayAvailable.Load() {
				return nil, errors.New("relay unavailable")
			}
			conn, _, err := dialer.Dial(wsURL, nil)
			return conn, err
		},
		InitMessage: TerminalMessage{
			Type:       "init",
			Session:    "test-session",
			UserID:     "client",
			WorkingDir: "/tmp",
			OSInfo:     "unix",
			Rows:       24,
			Cols:       80,
		},
		BaseRetryDelay: 10 * time.Millisecond,
		MaxRetryDelay:  20 * time.Millisecond,
		WriteTimeout:   500 * time.Millisecond,
	})

	manager.Start()
	defer manager.Shutdown(nil)

	if err := manager.Send(TerminalMessage{Type: "output", Data: "first", Session: "test-session", UserID: "client"}); err != nil {
		t.Fatalf("queue first message: %v", err)
	}
	if err := manager.Send(TerminalMessage{Type: "output", Data: "second", Session: "test-session", UserID: "client"}); err != nil {
		t.Fatalf("queue second message: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	relayAvailable.Store(true)

	initMsg := waitForRelayMessage(t, messageCh)
	firstMsg := waitForRelayMessage(t, messageCh)
	secondMsg := waitForRelayMessage(t, messageCh)

	if initMsg.Type != "init" {
		t.Fatalf("expected init message first, got %#v", initMsg)
	}
	if firstMsg.Type != "output" || firstMsg.Data != "first" {
		t.Fatalf("expected first buffered output second, got %#v", firstMsg)
	}
	if secondMsg.Type != "output" || secondMsg.Data != "second" {
		t.Fatalf("expected second buffered output third, got %#v", secondMsg)
	}
}

func TestRelayConnectionManagerReconnectsAndFlushesBufferedMessages(t *testing.T) {
	var relayAvailable atomic.Bool
	relayAvailable.Store(true)

	messageCh := make(chan TerminalMessage, 8)
	connCh := make(chan *websocket.Conn, 4)

	server, wsURL := startRelayTestServer(t, messageCh, connCh)
	defer server.Close()

	dialer := *websocket.DefaultDialer
	manager := newRelayConnectionManager(relayConnectionConfig{
		Dial: func() (*websocket.Conn, error) {
			if !relayAvailable.Load() {
				return nil, errors.New("relay unavailable")
			}
			conn, _, err := dialer.Dial(wsURL, nil)
			return conn, err
		},
		InitMessage: TerminalMessage{
			Type:       "init",
			Session:    "test-session",
			UserID:     "client",
			WorkingDir: "/tmp",
			OSInfo:     "windows",
			Rows:       30,
			Cols:       100,
		},
		BaseRetryDelay: 10 * time.Millisecond,
		MaxRetryDelay:  20 * time.Millisecond,
		WriteTimeout:   500 * time.Millisecond,
	})

	manager.Start()
	defer manager.Shutdown(nil)

	firstConn := waitForRelayConnection(t, connCh)
	firstInit := waitForRelayMessage(t, messageCh)
	if firstInit.Type != "init" {
		t.Fatalf("expected first init message, got %#v", firstInit)
	}

	relayAvailable.Store(false)
	if err := firstConn.Close(); err != nil {
		t.Fatalf("close first connection: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := manager.Send(TerminalMessage{Type: "output", Data: "after-reconnect", Session: "test-session", UserID: "client"}); err != nil {
		t.Fatalf("queue output after disconnect: %v", err)
	}

	relayAvailable.Store(true)

	_ = waitForRelayConnection(t, connCh)
	reconnectInit := waitForRelayMessage(t, messageCh)
	replayedOutput := waitForRelayMessage(t, messageCh)

	if reconnectInit.Type != "init" {
		t.Fatalf("expected init after reconnect, got %#v", reconnectInit)
	}
	if replayedOutput.Type != "output" || replayedOutput.Data != "after-reconnect" {
		t.Fatalf("expected buffered output after reconnect, got %#v", replayedOutput)
	}
}

func startRelayTestServer(t *testing.T, messageCh chan<- TerminalMessage, connCh chan<- *websocket.Conn) (*httptest.Server, string) {
	t.Helper()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		if connCh != nil {
			connCh <- conn
		}
		defer conn.Close()

		for {
			var message TerminalMessage
			if err := conn.ReadJSON(&message); err != nil {
				return
			}
			messageCh <- message
		}
	}))

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	return server, wsURL
}

func waitForRelayMessage(t *testing.T, messageCh <-chan TerminalMessage) TerminalMessage {
	t.Helper()

	select {
	case message := <-messageCh:
		return message
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for relay message")
		return TerminalMessage{}
	}
}

func waitForRelayConnection(t *testing.T, connCh <-chan *websocket.Conn) *websocket.Conn {
	t.Helper()

	select {
	case conn := <-connCh:
		return conn
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for relay connection")
		return nil
	}
}
