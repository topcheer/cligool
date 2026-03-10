# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目架构

CliGool是一个三层WebSocket远程终端系统：
- **CLI客户端** → **中继服务器** → **Web浏览器**

### 核心组件

1. **CLI客户端** (`cmd/client/`)
   - Windows: `main_windows.go` - 使用ConPTY（Windows Console Pseudo Terminal）
   - Unix/Linux/macOS: `main_unix.go` - 使用PTY（伪终端）
   - 支持30个操作系统/架构组合（Windows 2个、Linux 13个、*BSD 12个、macOS 2个）
   - 支持 `-cmd` 参数直接执行AI CLI工具
   - 支持 `-args` 参数传递命令行参数

2. **中继服务器** (`cmd/relay/`, `internal/relay/`)
   - 维护WebSocket会话和消息转发
   - 每个session可以有多个Web客户端连接
   - 最多一个CLI客户端连接（控制真实PTY）
   - 内存中维护会话状态，无数据库依赖

3. **Web界面** (`web/`)
   - `landing.html` - 下载页面
   - `terminal.html` - xterm.js终端界面
   - `index.html` - 备用入口页面

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
# Docker构建（包含全部客户端下载产物）
docker-compose build relay-server

# 本地构建所有客户端平台 + relay（含 Windows amd64/arm64）
./build-all.sh

# 仅构建 Windows 客户端（amd64 + arm64）
./build-windows.sh

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

**重要提示**：
- 本地 macOS/Linux 可以直接使用 `./build-all.sh` 交叉编译所有受支持客户端平台
- 如只需 Windows 产物，可使用 `./build-windows.sh`
- Docker构建会将所有客户端打包到`web/downloads/`目录
- Windows客户端会自动压缩为.zip文件以减少下载大小

### 服务端口

- 8081: 中继服务器（主机）
- 8080: 中继服务器（容器内部）

## 平台特定实现细节

### Windows客户端 (`main_windows.go`)

- 使用 **ConPTY**（Windows Console Pseudo Terminal），不是管道
- **自动检测控制台编码**：使用`GetConsoleOutputCP()`检测code page
  - 936: GBK（简体中文）
  - 932: Shift-JIS（日文）
  - 949: EUC-KR（韩文）
  - 950: Big5（繁体中文）
  - 1252: Windows-1252（西欧）
  - 437: CP437（英文）
  - 其他：默认使用Windows-1252
- **自动转换到UTF-8**：所有输出都先转换为UTF-8再发送到WebSocket和本地终端
- **UTF-8有效性检查**：使用`utf8.Valid()`避免双重编码
- **本地终端和Web终端都显示UTF-8编码**（先转换后输出，避免本地乱码）
- 换行符处理：`\r` → `\r\n`（Windows标准换行符）
- **支持完整的终端特性**：颜色、光标控制、屏幕清除等
- **动态窗口大小调整**：后台监控每500ms检查一次

### Unix客户端 (`main_unix.go`)

- 使用PTY (`github.com/creack/pty`)
- 支持完整的终端特性（颜色、光标控制等）
- 数据已是UTF-8，无需转换
- 支持窗口大小动态调整（`SIGWINCH`信号处理）
- 使用更小的缓冲区（1024字节）以减少延迟

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

### WebSocket消息类型

终端消息使用以下类型（`TerminalMessage`结构）：
- `"init"` - 初始化消息，包含工作目录、系统信息、终端大小
- `"input"` - 输入数据，来自用户键盘输入
- `"output"` - 输出数据，来自终端输出
- `"resize"` - 终端窗口大小调整
- `"close"` - 关闭会话

**重要**：心跳消息使用WebSocket控制帧（`PingMessage`/`PongMessage`），不通过上述消息类型。

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

### Web端输入回显逻辑

**关键**：Web端**不应该**本地回显用户输入，应该始终等待真实 PTY/ConPTY 的输出：

- **Windows（ConPTY）**：
  - ConPTY会输出真实的终端回显
  - Web端如果立即本地回显，会出现重复输入

- **Unix（PTY模式）**：
  - PTY会自动回显所有输入
  - Web端同样不应该本地回显

实现位置：`web/terminal.html`中的`terminal.onData`处理器：
```javascript
terminal.onData(data => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: 'input',
            data: data,
            session: sessionId,
            source: 'web'
        }));
    }
});
```

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

### CLI Relay 重连与缓冲

- CLI 客户端在本地 PTY/ConPTY 启动后，会在后台持续尝试连接 relay
- 如果 relay（包括代理链路）暂时不可用，CLI 不应直接退出，而是进入自动重试
- WebSocket 未发送成功的客户端消息需要在本地按顺序缓冲，直到连接恢复后再回放
- 重连成功后需要重新发送 `init`，再继续顺序发送缓冲的 `output` / `resize` 等消息

## 命令行参数支持

### `-cmd` 参数

允许直接执行指定的命令而不是默认shell：

```bash
# Unix/macOS/Linux
./cligool-darwin-arm64 -cmd claude

# Windows
cligool-windows-amd64.exe -cmd claude
```

### `-args` 参数

传递命令行参数给指定的命令：

```bash
# Unix/macOS/Linux
./cligool-darwin-arm64 -cmd git -args "commit -m 'fix bug'"

# Windows
cligool-windows-amd64.exe -cmd git -args "status"
```

### 参数解析规则

- 使用空格分隔多个参数
- 参数会被正确传递给命令
- Windows使用命令行字符串方式
- Unix使用可变参数列表方式

### 自动终端大小检测

客户端启动时自动检测终端大小：
- **Unix**: 使用 `pty.GetsizeFull()`
- **Windows**: 使用 `GetConsoleScreenBufferInfo()` 或 `golang.org/x/term.GetSize()`
- 回退值: 120x80（如果检测失败）

### 动态窗口大小调整

- **Unix**: `SIGWINCH` 信号处理器
- **Windows**: 后台goroutine每500ms检查一次
- **重要**: 调整时不输出日志（避免破坏终端布局）

## 开发工作流

### Make命令

项目使用Makefile简化常用操作：

```bash
# 构建所有组件
make build

# 构建特定组件
make build-relay        # 中继服务器
make build-client       # CLI客户端

# 构建跨平台版本
make build-all-platforms

# 运行服务
make run-relay          # 启动中继服务器
make run-client         # 启动CLI客户端

# Docker操作
make docker-build       # 构建Docker镜像
make docker-up          # 启动Docker服务
make docker-down        # 停止Docker服务
make docker-logs        # 查看Docker日志
make docker-restart     # 重启Docker服务

# 其他
make test               # 运行测试
make clean              # 清理构建文件
make dev                # 启动开发环境
make help               # 显示帮助信息
```

### 代码修改后的工作流

1. **修改代码后**：
   ```bash
   # 方式1：使用Make
   make docker-build
   make docker-up

   # 方式2：直接使用docker-compose
   docker-compose build relay-server
   docker-compose up -d relay-server
   ```

2. **本地测试客户端**：
   ```bash
   # 连接到本地服务器
   ./bin/cligool-darwin-arm64 -server http://localhost:8081

   # 或连接到远程服务器
   ./bin/cligool-darwin-arm64 -server https://your-server.com

   # 指定终端大小
   ./bin/cligool-darwin-arm64 -server http://localhost:8081 -cols 120 -rows 36
   ```

3. **验证部署**：
   ```bash
   # 检查API健康
   curl http://localhost:8081/api/health

   # 查看日志
   docker-compose logs -f relay-server

   # 查看服务状态
   docker-compose ps
   ```

## 常见问题

### 问题：Windows客户端CLI无输出或乱码

**无输出原因**：可能在修改时删除了`os.Stdout.Write(data)`

**乱码原因**：直接写入控制台编码数据到本地终端，需要先转换为UTF-8

**解决**：代码已实现自动编码检测和转换。如果仍有问题：
```go
// main_windows.go, stdout读取循环
data := buf[:n]

// 转换为UTF-8（自动检测控制台编码）
converted, err := convertToUTF8(data)
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

## 测试和调试

### 本地测试脚本

项目提供了多个测试脚本用于调试：

```bash
# PTY测试
./test-pty-simple.sh        # 简单PTY测试
./test-pty-websocket.sh     # PTY WebSocket测试

# 本地终端测试
./test-local-xterm.sh       # xterm.js本地测试
./test-local.sh             # 本地客户端测试

# 延迟测试
./test-latency.sh           # 网络延迟测试

# 完整架构测试
./test-final-architecture.sh # 完整系统测试
```

### 调试技巧

1. **启用详细日志**：
   ```bash
   # 客户端启用调试模式
   ./bin/cligool-darwin-arm64 -server http://localhost:8081 -debug
   ```

2. **浏览器控制台**：
   - 打开开发者工具（F12）
   - 查看Console标签的WebSocket消息
   - 检查Network标签的WebSocket连接状态

3. **服务器日志**：
   ```bash
   docker-compose logs -f relay-server
   ```

4. **测试WebSocket连接**：
   ```javascript
   // 在浏览器控制台中
   const ws = new WebSocket('ws://localhost:8081/api/terminal/test-session?type=web&user_id=test');
   ws.onmessage = (event) => console.log('Received:', event.data);
   ws.onopen = () => console.log('Connected');
   ws.onerror = (error) => console.log('Error:', error);
   ```

## 架构设计决策

### 为什么使用channel而不是Mutex处理WebSocket并发写入？

`gorilla/websocket`的`WriteMessage`方法不是并发安全的。使用`sync.Mutex`**无法**解决问题，因为：
- 问题不在于代码的并发控制，而在于WebSocket连接本身的限制
- 即使使用Mutex，多个goroutine仍然可能并发调用`WriteMessage`
- 正确的做法是确保所有WebSocket写入都通过单个goroutine串行执行

### 为什么Windows使用cmd.exe而不是PTY？

Windows没有原生的PTY支持（直到Windows 10/11的Windows Terminal PTY）。使用cmd.exe管道的权衡：
- ✅ 优点：不需要额外依赖，兼容性好
- ❌ 缺点：不支持完整终端特性（颜色、光标控制等）

未来可以考虑使用Windows 10+的ConPTY。

### 为什么使用小缓冲区（128字节）？

Unix客户端使用128字节缓冲区而不是常见的1024字节或更大：
- ✅ 优点：减少延迟，数据更快发送到WebSocket
- ✅ 优点：对于交互式终端，大多数输出都很短
- ❌ 缺点：稍微增加系统调用次数

对于远程终端场景，低延迟比高吞吐量更重要。

## 相关文档

- **快速开始**: `QUICKSTART.md`
- **使用指南**: `USAGE_GUIDE.md`（英文）、`USAGE_GUIDE_CN.md`（中文）
- **开发指南**: `docs/DEVELOPMENT.md`
- **部署指南**: `docs/DEPLOYMENT.md`
- **配置说明**: `docs/CONFIG.md`
- **Windows支持**: `docs/WINDOWS_SUPPORT.md`
- **PTY故障排除**: `docs/PTY_TROUBLESHOOTING.md`
- **延迟优化**: `LATENCY_OPTIMIZATION.md`
- **平台列表**: `PLATFORMS.md`
