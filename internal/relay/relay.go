package relay

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	Type       string `json:"type"`       // "input", "output", "resize", "close", "init"
	Data       string `json:"data"`       // 终端数据
	Rows       int    `json:"rows"`       // 终端行数
	Cols       int    `json:"cols"`       // 终端列数
	Session    string `json:"session"`    // 会话ID
	UserID     string `json:"user_id"`    // 用户ID
	WorkingDir string `json:"working_dir,omitempty"` // 工作目录
	OSInfo     string `json:"os_info,omitempty"`     // 操作系统信息
}

// CachedMessage 缓存的消息（只保存必要字段以节省内存）
type CachedMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Rows int    `json:"rows,omitempty"`
	Cols int    `json:"cols,omitempty"`
}

// Session 终端会话
type Session struct {
	ID               string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Owner            string
	Clients          map[string]*websocket.Conn // Web客户端连接
	ClientCon        *websocket.Conn            // CLI客户端连接
	Mutex            sync.RWMutex
	Active           bool
	LastPing         time.Time // 最后一次收到ping的时间
	WorkingDirectory string   // 客户端当前工作目录
	OSInfo           string   // 客户端操作系统信息
	MessageCache     []CachedMessage // CLI消息缓存（当无Web客户端时）
	CacheSizeLimit   int            // 缓存大小限制（条数）
	TotalCacheSize   int            // 缓存总大小（字节）
}

// Service 中继服务
type Service struct {
	Config     Config
	Sessions   map[string]*Session
	SessionsMu sync.RWMutex
}

// Config 服务配置
type Config struct {
	Host string
	Port string
}

// NewService 创建新的中继服务
func NewService(config Config) *Service {
	return &Service{
		Config:   config,
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
		session.LastPing = time.Now()
		session.Mutex.Unlock()

		log.Printf("CLI client connected to session: %s", sessionID)

		// 设置ping handler
		s.setupPingHandler(conn, session, "client")

		// 等待客户端准备就绪
		s.handleClientConnection(session)
	} else {
		// Web客户端连接
		session.Mutex.Lock()
		session.Clients[userID] = conn
		session.Mutex.Unlock()

		log.Printf("Web client connected to session: %s, user: %s", sessionID, userID)

		// 设置ping handler
		s.setupPingHandler(conn, session, userID)

		// 检查是否有CLI客户端连接
		session.Mutex.RLock()
		cliConnected := session.ClientCon != nil
		session.Mutex.RUnlock()

		if !cliConnected {
			// 没有CLI客户端连接，发送提示消息
			log.Printf("⚠️  Web客户端连接但无CLI客户端，发送提示消息")
			s.sendNoCliClientMessage(session, conn)
			// 注意：不要调用 handleWebClient，因为已经发送了提示消息
			return
		}

		// 先发送缓存的CLI消息（如果有）
		s.sendCachedMessages(session, conn)

		// 如果会话已有初始化信息，发送给新连接的Web客户端
		session.Mutex.RLock()
		hasInit := session.WorkingDirectory != "" || session.OSInfo != ""
		session.Mutex.RUnlock()

		if hasInit {
			session.Mutex.RLock()
			initMsg := TerminalMessage{
				Type:       "init",
				WorkingDir: session.WorkingDirectory,
				OSInfo:     session.OSInfo,
			}
			session.Mutex.RUnlock()

			jsonData, _ := json.Marshal(initMsg)
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("Failed to send init message to web client: %v", err)
			} else {
				log.Printf("📤 发送初始化消息给Web客户端: OS=%s", session.OSInfo)
			}
		}

		// 处理Web客户端
		s.handleWebClient(session, userID)
	}
}

// addToCache 添加消息到缓存（使用简单的FIFO策略）
func (s *Service) addToCache(session *Session, msg TerminalMessage) {
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	// 检查缓存大小限制
	if len(session.MessageCache) >= session.CacheSizeLimit {
		// 移除最旧的消息（第一个）
		removed := session.MessageCache[0]
		session.TotalCacheSize -= len(removed.Data) + len(removed.Type)
		session.MessageCache = session.MessageCache[1:]
	}

	// 添加新消息到缓存
	cachedMsg := CachedMessage{
		Type: msg.Type,
		Data: msg.Data,
		Rows: msg.Rows,
		Cols: msg.Cols,
	}
	session.MessageCache = append(session.MessageCache, cachedMsg)
	session.TotalCacheSize += len(msg.Data) + len(msg.Type)

	log.Printf("📦 消息已缓存: cache_size=%d/%d, total_bytes=%d, type=%s",
		len(session.MessageCache), session.CacheSizeLimit, session.TotalCacheSize, msg.Type)
}

// sendCachedMessages 发送缓存的消息给新连接的Web客户端
func (s *Service) sendCachedMessages(session *Session, conn *websocket.Conn) {
	session.Mutex.RLock()
	cache := make([]CachedMessage, len(session.MessageCache))
	copy(cache, session.MessageCache)
	session.Mutex.RUnlock()

	if len(cache) == 0 {
		log.Printf("📭 缓存为空，无需发送历史消息")
		return
	}

	log.Printf("📤 开始发送缓存消息: %d 条", len(cache))

	for i, cachedMsg := range cache {
		msg := TerminalMessage{
			Type:   cachedMsg.Type,
			Data:   cachedMsg.Data,
			Rows:   cachedMsg.Rows,
			Cols:   cachedMsg.Cols,
			Session: session.ID,
		}

		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("❌ 发送缓存消息失败 [%d/%d]: %v", i+1, len(cache), err)
			return
		}
	}

	log.Printf("✅ 缓存消息发送完成: %d 条", len(cache))
}

// sendNoCliClientMessage 发送无CLI客户端的提示消息
func (s *Service) sendNoCliClientMessage(session *Session, conn *websocket.Conn) {
	// 构建提示消息，包含命令示例
	hintMsg := TerminalMessage{
		Type:   "no_cli",
		Data:   buildNoCliHintMessage(session.ID),
		Session: session.ID,
	}

	jsonData, _ := json.Marshal(hintMsg)
	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("❌ 发送无CLI提示消息失败: %v", err)
	} else {
		log.Printf("✅ 已发送无CLI提示消息给Web客户端")
	}
}

// buildNoCliHintMessage 构建无CLI客户端的提示消息
func buildNoCliHintMessage(sessionID string) string {
	// 检测操作系统
	ostype := runtime.GOOS

	// 根据操作系统生成不同的命令示例
	var cmdExample string
	switch ostype {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			cmdExample = `./cligool-darwin-arm64 -server http://localhost:8081 -session %s`
		} else {
			cmdExample = `./cligool-darwin-amd64 -server http://localhost:8081 -session %s`
		}
	case "linux":
		if runtime.GOARCH == "amd64" {
			cmdExample = `./cligool-linux-amd64 -server http://localhost:8081 -session %s`
		} else if runtime.GOARCH == "arm64" {
			cmdExample = `./cligool-linux-arm64 -server http://localhost:8081 -session %s`
		} else {
			cmdExample = `./cligool-linux-$(uname -m) -server http://localhost:8081 -session %s`
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			cmdExample = `cligool-windows-amd64.exe -server http://localhost:8081 -session %s`
		} else {
			cmdExample = `cligool-windows-arm64.exe -server http://localhost:8081 -session %s`
		}
	default:
		cmdExample = `./cligool -server http://localhost:8081 -session %s`
	}

	// 格式化提示消息
	hint := fmt.Sprintf(`⚠️  CLI客户端未连接

请先启动CLI客户端，然后再刷新此页面。

启动命令示例：
%s

或者使用配置文件：
1. 编辑 ~/.cligool.json 设置服务器地址
2. 运行: ./cligool -session %s

💡 提示：
- CLI客户端必须在Web客户端之前启动
- 确保使用相同的session ID
- 检查防火墙设置

📥 下载客户端：
- https://cligool.zty8.cn/`, fmt.Sprintf(cmdExample, sessionID), sessionID)

	return hint
}

// notifyWebClientsClientDisconnected 通知所有Web客户端CLI已断开
func (s *Service) notifyWebClientsClientDisconnected(session *Session) {
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	webClientCount := len(session.Clients)
	if webClientCount == 0 {
		log.Printf("📭 没有Web客户端需要通知")
		return
	}

	log.Printf("📡 通知 %d 个Web客户端: CLI已断开", webClientCount)

	// 创建关闭消息
	closeMsg := TerminalMessage{
		Type:   "close",
		Data:   "CLI客户端已断开连接",
		Session: session.ID,
	}
	jsonData, _ := json.Marshal(closeMsg)

	// 向所有Web客户端发送关闭消息
	for userID, conn := range session.Clients {
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Printf("❌ 发送关闭消息失败给Web客户端 %s: %v", userID, err)
		} else {
			log.Printf("✅ 已通知Web客户端: %s", userID)
		}

		// 关闭Web客户端连接
		conn.Close()
	}

	// 清空Web客户端列表
	session.Clients = make(map[string]*websocket.Conn)
	log.Printf("🧹 已清理所有Web客户端连接")
}

// handleClientConnection 处理CLI客户端连接
func (s *Service) handleClientConnection(session *Session) {
	conn := session.ClientCon
	log.Printf("🖥️  开始处理CLI客户端连接")

	for {
		var msg TerminalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ CLI客户端连接断开: %v", err)
			break
		}

		log.Printf("📨 收到CLI消息: type=%s, data_len=%d, session=%s",
			msg.Type, len(msg.Data), msg.Session)

		// 处理初始化消息
		if msg.Type == "init" {
			session.Mutex.Lock()
			session.WorkingDirectory = msg.WorkingDir
			session.OSInfo = msg.OSInfo
			session.Mutex.Unlock()
			log.Printf("📁 客户端初始化: 工作目录=%s, 系统=%s", msg.WorkingDir, msg.OSInfo)

			// 广播初始化消息给所有Web客户端
			s.broadcastToWebClients(session, msg)
			continue
		}

		// 广播消息给所有Web客户端
		s.broadcastToWebClients(session, msg)
	}

	// CLI客户端断开，通知所有Web客户端
	log.Printf("🔔 通知所有Web客户端: CLI已断开")
	s.notifyWebClientsClientDisconnected(session)

	// 清理连接
	session.Mutex.Lock()
	session.ClientCon = nil
	session.Active = false
	session.Mutex.Unlock()
}

// handleWebClient 处理Web客户端连接
func (s *Service) handleWebClient(session *Session, userID string) {
	conn := session.Clients[userID]
	log.Printf("🌐 开始处理Web客户端: %s", userID)

	for {
		var msg TerminalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Web client connection error: %v", err)
			break
		}

		log.Printf("📨 收到Web消息: type=%s, data_len=%d, session=%s",
			msg.Type, len(msg.Data), msg.Session)

		// 转发消息到CLI客户端
		session.Mutex.RLock()
		if session.ClientCon != nil {
			err := session.ClientCon.WriteJSON(msg)
			session.Mutex.RUnlock()
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			} else {
				log.Printf("✅ 消息已转发到CLI客户端")
			}
		} else {
			session.Mutex.RUnlock()
			log.Printf("⚠️  CLI客户端未连接，消息丢弃")
		}
	}

	// 清理连接
	session.Mutex.Lock()
	delete(session.Clients, userID)
	session.Mutex.Unlock()
}

// broadcastToWebClients 广播消息给所有Web客户端，如果没有Web客户端则缓存
func (s *Service) broadcastToWebClients(session *Session, msg TerminalMessage) {
	session.Mutex.RLock()
	webClientCount := len(session.Clients)
	session.Mutex.RUnlock()

	// 如果没有Web客户端连接，缓存消息
	if webClientCount == 0 {
		s.addToCache(session, msg)
		return
	}

	// 有Web客户端，实时广播
	session.Mutex.RLock()
	defer session.Mutex.RUnlock()

	log.Printf("📡 广播消息到 %d 个Web客户端: type=%s, data_len=%d",
		webClientCount, msg.Type, len(msg.Data))

	for userID, conn := range session.Clients {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Failed to send message to web client %s: %v", userID, err)
		} else {
			log.Printf("✅ 消息已发送到Web客户端: %s", userID)
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
		ID:             sessionID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Owner:          owner,
		Clients:        make(map[string]*websocket.Conn),
		Active:         false,
		MessageCache:   make([]CachedMessage, 0, 1000), // 预分配1000条容量
		CacheSizeLimit: 1000,                            // 最多缓存1000条消息
	}

	s.Sessions[sessionID] = session
	log.Printf("Created new session: %s", sessionID)

	return session
}

// setupPingHandler 设置ping处理器和心跳检测
func (s *Service) setupPingHandler(conn *websocket.Conn, session *Session, peerID string) {
	// 设置pong handler来更新最后活跃时间
	conn.SetPongHandler(func(appData string) error {
		session.Mutex.Lock()
		session.LastPing = time.Now()
		session.Mutex.Unlock()
		log.Printf("💓 收到 %s 的pong", peerID)
		return nil
	})

	// 启动定期ping
	go s.startHeartbeat(conn, session, peerID)

	// 启动超时检测
	go s.monitorHeartbeat(session, peerID)
}

// startHeartbeat 定期发送ping
func (s *Service) startHeartbeat(conn *websocket.Conn, session *Session, peerID string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			session.Mutex.RLock()
			isActive := session.Active && session.ClientCon != nil
			session.Mutex.RUnlock()

			if !isActive {
				log.Printf("🔌 连接已关闭，停止心跳: %s", peerID)
				return
			}

			if err := conn.WriteMessage(websocket.PingMessage, []byte("heartbeat")); err != nil {
				log.Printf("❌ 发送ping失败到 %s: %v", peerID, err)
				return
			}
			log.Printf("💓 发送ping到 %s", peerID)
		}
	}
}

// monitorHeartbeat 监控心跳超时
func (s *Service) monitorHeartbeat(session *Session, peerID string) {
	ticker := time.NewTicker(15 * time.Second) // 每15秒检查一次
	defer ticker.Stop()

	for range ticker.C {
		session.Mutex.RLock()
		lastPing := session.LastPing
		session.Mutex.RUnlock()

		// 如果超过90秒没有收到pong，认为连接已死
		if time.Since(lastPing) > 90*time.Second {
			log.Printf("⚠️  心跳超时，关闭连接: %s (上次活跃: %v 前)",
				peerID, time.Since(lastPing))

			session.Mutex.Lock()
			if peerID == "client" && session.ClientCon != nil {
				session.ClientCon.Close()
				session.ClientCon = nil
				session.Active = false
			} else if conn, ok := session.Clients[peerID]; ok {
				conn.Close()
				delete(session.Clients, peerID)
			}
			session.Mutex.Unlock()
			return
		}
	}
}
