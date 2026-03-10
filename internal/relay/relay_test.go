package relay

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestSendCachedMessagesClearsCacheAfterSuccessfulDelivery(t *testing.T) {
	service := NewService(Config{})
	session := service.getOrCreateSession("session-under-test", "owner")

	service.addToCache(session, TerminalMessage{Type: "output", Data: "hello"})
	service.addToCache(session, TerminalMessage{Type: "resize", Rows: 24, Cols: 80})

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade failed: %v", err)
		}
		defer conn.Close()

		service.sendCachedMessages(session, conn)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	var first TerminalMessage
	if err := conn.ReadJSON(&first); err != nil {
		t.Fatalf("read first message failed: %v", err)
	}
	if first.Type != "output" || first.Data != "hello" || first.Session != session.ID {
		t.Fatalf("unexpected first message: %+v", first)
	}

	var second TerminalMessage
	if err := conn.ReadJSON(&second); err != nil {
		t.Fatalf("read second message failed: %v", err)
	}
	if second.Type != "resize" || second.Rows != 24 || second.Cols != 80 || second.Session != session.ID {
		t.Fatalf("unexpected second message: %+v", second)
	}

	session.Mutex.RLock()
	defer session.Mutex.RUnlock()

	if len(session.MessageCache) != 0 {
		t.Fatalf("expected cache to be cleared, got %d entries", len(session.MessageCache))
	}
	if session.TotalCacheSize != 0 {
		t.Fatalf("expected cache size to reset, got %d", session.TotalCacheSize)
	}
}

func TestWebClientWithoutCliStaysConnectedAndReceivesInit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := NewService(Config{})
	router := gin.New()
	router.GET("/api/terminal/:session_id", service.HandleTerminalConnection)

	server := httptest.NewServer(router)
	defer server.Close()

	wsBaseURL := "ws" + strings.TrimPrefix(server.URL, "http")
	sessionID := "waiting-web-session"

	webConn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/api/terminal/"+sessionID+"?type=web&user_id=web-1", nil)
	if err != nil {
		t.Fatalf("dial web failed: %v", err)
	}
	defer webConn.Close()

	webConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var noCli TerminalMessage
	if err := webConn.ReadJSON(&noCli); err != nil {
		t.Fatalf("read no_cli failed: %v", err)
	}
	if noCli.Type != "no_cli" {
		t.Fatalf("expected no_cli message, got %+v", noCli)
	}

	cliConn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/api/terminal/"+sessionID+"?type=client&user_id=client", nil)
	if err != nil {
		t.Fatalf("dial cli failed: %v", err)
	}
	defer cliConn.Close()

	initMsg := TerminalMessage{
		Type:       "init",
		Session:    sessionID,
		UserID:     "client",
		WorkingDir: "/tmp",
		OSInfo:     "unix",
	}
	if err := cliConn.WriteJSON(initMsg); err != nil {
		t.Fatalf("write init failed: %v", err)
	}

	webConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var receivedInit TerminalMessage
	if err := webConn.ReadJSON(&receivedInit); err != nil {
		t.Fatalf("read init after cli connect failed: %v", err)
	}
	if receivedInit.Type != "init" || receivedInit.OSInfo != "unix" || receivedInit.WorkingDir != "/tmp" {
		t.Fatalf("unexpected init after cli connect: %+v", receivedInit)
	}
}

func TestCliDisconnectDoesNotCloseWebClientConnection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := NewService(Config{})
	router := gin.New()
	router.GET("/api/terminal/:session_id", service.HandleTerminalConnection)

	server := httptest.NewServer(router)
	defer server.Close()

	wsBaseURL := "ws" + strings.TrimPrefix(server.URL, "http")
	sessionID := "disconnect-web-session"

	cliConn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/api/terminal/"+sessionID+"?type=client&user_id=client", nil)
	if err != nil {
		t.Fatalf("dial cli failed: %v", err)
	}

	webConn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/api/terminal/"+sessionID+"?type=web&user_id=web-1", nil)
	if err != nil {
		t.Fatalf("dial web failed: %v", err)
	}
	defer webConn.Close()

	firstInit := TerminalMessage{
		Type:       "init",
		Session:    sessionID,
		UserID:     "client",
		WorkingDir: "/tmp/first",
		OSInfo:     "unix",
	}
	if err := cliConn.WriteJSON(firstInit); err != nil {
		t.Fatalf("write first init failed: %v", err)
	}

	webConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var initialWebInit TerminalMessage
	if err := webConn.ReadJSON(&initialWebInit); err != nil {
		t.Fatalf("read first init failed: %v", err)
	}
	if initialWebInit.Type != "init" || initialWebInit.WorkingDir != "/tmp/first" {
		t.Fatalf("unexpected first init: %+v", initialWebInit)
	}

	if err := cliConn.Close(); err != nil {
		t.Fatalf("close cli failed: %v", err)
	}

	webConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var closeMsg TerminalMessage
	if err := webConn.ReadJSON(&closeMsg); err != nil {
		t.Fatalf("read cli disconnect notification failed: %v", err)
	}
	if closeMsg.Type != "close" {
		t.Fatalf("expected close notification, got %+v", closeMsg)
	}

	cliReconnect, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/api/terminal/"+sessionID+"?type=client&user_id=client", nil)
	if err != nil {
		t.Fatalf("dial reconnect cli failed: %v", err)
	}
	defer cliReconnect.Close()

	secondInit := TerminalMessage{
		Type:       "init",
		Session:    sessionID,
		UserID:     "client",
		WorkingDir: "/tmp/second",
		OSInfo:     "windows",
	}
	if err := cliReconnect.WriteJSON(secondInit); err != nil {
		t.Fatalf("write second init failed: %v", err)
	}

	webConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var reconnectedInit TerminalMessage
	if err := webConn.ReadJSON(&reconnectedInit); err != nil {
		t.Fatalf("read init after cli reconnect failed: %v", err)
	}
	if reconnectedInit.Type != "init" || reconnectedInit.WorkingDir != "/tmp/second" || reconnectedInit.OSInfo != "windows" {
		t.Fatalf("unexpected init after cli reconnect: %+v", reconnectedInit)
	}
}
