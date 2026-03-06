# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目架构

CliGool是一个三层WebSocket远程终端系统：
- **CLI客户端** → **中继服务器** → **Web浏览器**

### 核心组件

1. **CLI客户端** (`cmd/client/`)
   - Windows: `main_windows.go` - 使用cmd.exe管道，需要GBK→UTF-8编码转换
   - Unix/Linux/macOS: `main_unix.go` - 使用PTY（伪终端）
   - 支持18个操作系统/架构组合

2. **中继服务器** (`cmd/relay/`, `internal/relay/`)
   - 维护WebSocket会话和消息转发
   - 每个session可以有多个Web客户端连接
   - 最多一个CLI客户端连接（控制真实PTY）
   - PostgreSQL存储会话信息，Redis用于缓存

3. **Web界面** (`web/`)
   - `landing.html` - 下载页面
   - `terminal.html` - xterm.js终端界面

### 消息流

```
CLI客户端 ──WebSocket──> 中继服务器 <──WebSocket─── Web浏览器
      │                        │                        │
      └────────本地输入输出────┘                        xterm.js
```

**关键点**：CLI客户端控制真实的PTY/Shell，所有终端IO都通过WebSocket转发到Web端

## 构建和部署

### 重要约束

**必须使用docker-compose构建**，不是直接docker build。Docker缓存足够聪明，不需要`--no-cache`。

### 构建命令

```bash
# Docker构建（包含所有18个平台客户端）
docker-compose build relay-server

# 本地构建所有*nix平台客户端（16个，不含Windows）
./build-all.sh

# 本地构建特定平台
GOOS=linux GOARCH=amd64 go build -o bin/cligool-linux-amd64 ./cmd/client
GOOS=windows GOARCH=amd64 go build -o bin/cligool-windows-amd64.exe ./cmd/client

# 启动所有服务
docker-compose up -d

# 重启relay-server（用于代码更新）
docker-compose up -d relay-server

# 查看服务状态
docker-compose ps
```

### 服务端口

- 8081: 中继服务器（主机）
- 8080: 中继服务器（容器内部）

## 平台特定实现细节

### Windows客户端 (`main_windows.go`)

- 使用`cmd.exe`管道，不是PTY
- **必须进行GBK→UTF-8编码转换**：cmd.exe输出GBK编码，需要转换为UTF-8
- **本地终端和Web终端都显示UTF-8编码**（先转换后输出，避免本地乱码）
- 换行符处理：`\r` → `\r\n`

### Unix客户端 (`main_unix.go`)

- 使用PTY (`github.com/creack/pty`)
- 支持完整的终端特性（颜色、光标控制等）
- 数据已是UTF-8，无需转换

## 关键技术约束

### WebSocket并发写入

**gorilla/websocket不支持并发写入**，必须使用channel串行化：

```go
// 创建WebSocket写入channel
wsWriteChan := make(chan []byte, 100)

// 启动WebSocket写入goroutine
go func() {
    for data := range wsWriteChan {
        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            log.Printf("❌ WebSocket写入失败: %v", err)
            return
        }
    }
}()

// 所有WebSocket写入都通过channel
wsWriteChan <- jsonData
```

**不要使用sync.Mutex** - 这无法解决WebSocket并发写入问题。

### WebSocket URL构建

- **CLI客户端**: `ws://host/api/terminal/{session_id}?type=client&user_id=client`
- **Web客户端**: `ws://host/api/terminal/{session_id}?type=web&user_id=web-{timestamp}`

Web端user_id使用时间戳确保每次连接唯一。

### 缓存问题

**Cloudflare会缓存所有内容**，包括下载链接。解决方案：

1. **下载链接**：在`landing.html`中使用JavaScript动态生成时间戳参数
   ```javascript
   document.querySelectorAll('a[href^="/downloads/"]').forEach(link => {
       const cacheBuster = Date.now();
       link.href = link.href.replace(/\?v=\d+/, '') + `?v=${cacheBuster}`;
   });
   ```

2. **WebSocket连接**：user_id参数使用`Date.now()`确保唯一性

## 会话管理

### 会话状态

- `Session.Clients`: map[string]*websocket.Conn - Web客户端连接（多个）
- `Session.ClientCon`: *websocket.Conn - CLI客户端连接（最多一个）
- `Session.Active`: bool - 会话是否活跃
- `Session.LastPing`: time.Time - 最后心跳时间

### 心跳机制

- 服务器每30秒发送ping
- 客户端自动回复pong
- 超过90秒无活动自动断开

## 开发工作流

1. **修改代码后**：
   ```bash
   docker-compose build relay-server
   docker-compose up -d relay-server
   ```

2. **本地测试客户端**：
   ```bash
   # 连接到本地服务器
   ./bin/cligool-darwin-arm64 -server http://localhost:8081

   # 或连接到远程服务器
   ./bin/cligool-darwin-arm64 -server https://your-server.com
   ```

3. **验证部署**：
   ```bash
   # 检查API健康
   curl http://localhost:8081/api/health

   # 查看日志
   docker-compose logs -f relay-server
   ```

## 常见问题

### 问题：Windows客户端CLI无输出或乱码

**无输出原因**：可能在修改时删除了`os.Stdout.Write(data)`

**乱码原因**：直接写入GBK编码数据到本地终端，需要先转换为UTF-8

**解决**：在stdout/stderr读取循环中先转换GBK到UTF-8，再输出到本地终端：
```go
// main_windows.go, stdout读取循环
data := buf[:n]

// 转换为UTF-8
converted, err := convertGBKToUTF8(data)
if err != nil {
    converted = string(data)
}

// 显示到本地终端（UTF-8数据）
os.Stdout.Write([]byte(converted))

// 同时发送UTF-8到WebSocket
msg := TerminalMessage{
    Type: "output",
    Data: converted,
    Session: sessionID,
    UserID: "client",
}
jsonData, _ := json.Marshal(msg)
wsWriteChan <- jsonData
```

### 问题：WebSocket并发写入panic

**错误**：`panic: concurrent write to websocket connection`

**原因**：多个goroutine直接调用`conn.WriteMessage()`

**解决**：使用channel串行化所有WebSocket写入（见上方示例）

### 问题：Cloudflare缓存旧版本

**原因**：下载链接URL是静态的

**解决**：在`landing.html`中添加动态缓存参数（见上方示例）
