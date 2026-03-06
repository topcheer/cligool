package relay

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cligool/cligool/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查Origin
	},
}

// TerminalMessage 终端消息类型
type TerminalMessage struct {
	Type    string `json:"type"`    // "input", "output", "resize", "close"
	Data    string `json:"data"`    // 终端数据（Base64编码的字符串）
	Rows    int    `json:"rows"`    // 终端行数
	Cols    int    `json:"cols"`    // 终端列数
	Session string `json:"session"` // 会话ID
	UserID  string `json:"user_id"` // 用户ID
}

// Session 终端会话
type Session struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Owner     string
	Clients   map[string]*websocket.Conn // Web客户端连接
	ClientCon *websocket.Conn            // CLI客户端连接
	Mutex     sync.RWMutex
	Active    bool
}

// Service 中继服务
type Service struct {
	Config     Config
	DB         *database.DB
	Sessions   map[string]*Session
	SessionsMu sync.RWMutex
}

// Config 服务配置
type Config struct {
	DB   *database.DB
	Host string
	Port string
}

// NewService 创建新的中继服务
func NewService(config Config) *Service {
	return &Service{
		Config:   config,
		DB:       config.DB,
		Sessions: make(map[string]*Session),
	}
}

// HandleTerminalConnection 处理WebSocket终端连接
func (s *Service) HandleTerminalConnection(c *gin.Context) {
	// 获取会话ID（优先使用路径参数，如果没有则使用查询参数）
	sessionID := c.Param("session_id")
	if sessionID == "" {
		sessionID = c.Query("session_id")
	}
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session_id"})
		return
	}

	// 获取用户ID
	userID := c.Query("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	// 升级到WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// 获取或创建会话
	session := s.getOrCreateSession(sessionID, userID)

	// 判断连接类型
	connType := c.Query("type") // "web" 或 "client"

	if connType == "client" {
		// CLI客户端连接
		session.Mutex.Lock()
		session.ClientCon = conn
		session.Active = true
		session.Mutex.Unlock()

		log.Printf("CLI client connected to session: %s", sessionID)

		// 等待客户端准备就绪
		s.handleClientConnection(session)
	} else {
		// Web客户端连接
		session.Mutex.Lock()
		session.Clients[userID] = conn
		session.Mutex.Unlock()

		log.Printf("Web client connected to session: %s, user: %s", sessionID, userID)

		// 处理Web客户端
		s.handleWebClient(session, userID)
	}
}

// handleClientConnection 处理CLI客户端连接
func (s *Service) handleClientConnection(session *Session) {
	conn := session.ClientCon

	for {
		var msg TerminalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Client connection error: %v", err)
			break
		}

		// 广播消息给所有Web客户端
		s.broadcastToWebClients(session, msg)
	}

	// 清理连接
	session.Mutex.Lock()
	session.ClientCon = nil
	session.Active = false
	session.Mutex.Unlock()
}

// handleWebClient 处理Web客户端连接
func (s *Service) handleWebClient(session *Session, userID string) {
	conn := session.Clients[userID]

	for {
		var msg TerminalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Web client connection error: %v", err)
			break
		}

		// 转发消息到CLI客户端
		session.Mutex.RLock()
		if session.ClientCon != nil {
			err := session.ClientCon.WriteJSON(msg)
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}
		session.Mutex.RUnlock()
	}

	// 清理连接
	session.Mutex.Lock()
	delete(session.Clients, userID)
	session.Mutex.Unlock()
}

// broadcastToWebClients 广播消息给所有Web客户端
func (s *Service) broadcastToWebClients(session *Session, msg TerminalMessage) {
	session.Mutex.RLock()
	defer session.Mutex.RUnlock()

	for userID, conn := range session.Clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send message to web client %s: %v", userID, err)
		}
	}
}

// getOrCreateSession 获取或创建会话
func (s *Service) getOrCreateSession(sessionID, owner string) *Session {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()

	if session, exists := s.Sessions[sessionID]; exists {
		return session
	}

	// 创建新会话
	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Owner:     owner,
		Clients:   make(map[string]*websocket.Conn),
		Active:    false,
	}

	s.Sessions[sessionID] = session
	log.Printf("Created new session: %s", sessionID)

	return session
}

// CreateSession 创建新会话
func (s *Service) CreateSession(c *gin.Context) {
	var req struct {
		Owner string `json:"owner"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成会话ID
	sessionID := uuid.New().String()

	// 保存到数据库
	session := &database.Session{
		ID:        sessionID,
		Owner:     req.Owner,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    false,
	}

	if err := s.DB.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// 创建内存会话
	s.getOrCreateSession(sessionID, req.Owner)

	c.JSON(http.StatusCreated, session)
}

// GetSession 获取会话信息
func (s *Service) GetSession(c *gin.Context) {
	sessionID := c.Param("id")

	session, err := s.DB.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteSession 删除会话
func (s *Service) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")

	// 从内存中删除
	s.SessionsMu.Lock()
	if session, exists := s.Sessions[sessionID]; exists {
		// 关闭所有连接
		session.Mutex.Lock()
		for _, conn := range session.Clients {
			conn.Close()
		}
		if session.ClientCon != nil {
			session.ClientCon.Close()
		}
		session.Mutex.Unlock()
		delete(s.Sessions, sessionID)
	}
	s.SessionsMu.Unlock()

	// 从数据库中删除
	if err := s.DB.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
}

// ListSessions 列出所有会话
func (s *Service) ListSessions(c *gin.Context) {
	sessions, err := s.DB.ListSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}