package relay

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
