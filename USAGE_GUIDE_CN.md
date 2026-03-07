# CliGool 使用指南

一个基于Go和WebSocket的跨平台远程终端解决方案，支持30种操作系统和架构。

## 🏗️ 系统架构

```
CLI客户端            中继服务器              Web浏览器
┌──────────────┐              ┌──────────┐              ┌─────────────┐
│ 真实PTY环境    │──WebSocket──▶│ 消息转发器│◀──WebSocket───│  xterm.js    │
│              │              │          │              │  终端界面     │
│ 18个平台     │              │          │              │              │
└──────────────┘              └──────────┘              └─────────────┘
```

**核心特点**：
- ✅ **真实PTY**：CLI客户端提供完整的终端环境
- ✅ **消息转发**：中继服务器负责WebSocket消息路由
- ✅ **独立界面**：Web界面基于xterm.js，支持任意浏览器
- ✅ **无状态设计**：内存中维护会话，无需数据库依赖

## 🚀 快速开始

### 方法一：使用Docker Compose（推荐）

#### 1. 启动中继服务器

```bash
# 克隆仓库
git clone https://github.com/topcheer/cligool.git
cd cligool

# 生产环境（使用预构建镜像）
docker-compose up -d

# 或：开发环境（本地构建）
docker-compose -f docker-compose.dev.yml up -d --build

# 检查服务状态
docker-compose ps
```

**预期输出**：
```
NAME                STATUS
cligool-relay       Up (healthy)
```

#### 2. 启动CLI客户端

**在需要远程控制的机器上**：

```bash
# 下载对应平台的客户端
# 从下载页面获取：http://localhost:8081/

# 启动客户端（连接到本地服务器）
./cligool -server http://localhost:8081

# 或连接到远程服务器
./cligool -server https://your-server.com
```

**输出示例**：
```
╔═══════════════════════════════════════════════════════════╗
║                    🚀 CliGool 远程终端                      ║
╠═══════════════════════════════════════════════════════════╣
║ 📋 会话ID: abc123-def456-7890-abcd-ef1234567890         ║
║ 🌐 Web访问: http://localhost:8081/session/abc123-...     ║
║ 🔗 连接状态: 🟢 已连接                                    ║
╚═══════════════════════════════════════════════════════════╝
```

#### 3. 打开Web界面

在浏览器中访问：
```
http://localhost:8081/session/[会话ID]
```

或直接访问下载页面：
```
http://localhost:8081/
```

### 方法二：手动构建

#### 1. 构建中继服务器

```bash
# 构建服务器
go build -o bin/relay-server ./cmd/relay

# 运行服务器
./bin/relay-server
```

#### 2. 构建客户端

```bash
# 构建当前平台的客户端
go build -o cligool ./cmd/client

# 交叉编译其他平台
GOOS=linux GOARCH=amd64 go build -o cligool-linux-amd64 ./cmd/client
GOOS=windows GOARCH=amd64 go build -o cligool-windows-amd64.exe ./cmd/client
GOOS=darwin GOARCH=arm64 go build -o cligool-darwin-arm64 ./cmd/client
```

## 🌍 支持的平台

### Windows (2个)
- Windows amd64 (Intel/AMD 64位)
- Windows arm64 (Surface Pro X等ARM设备)

### Linux (8个)
- Linux amd64 (Ubuntu、Debian、CentOS等64位系统)
- Linux arm64 (树莓派4/5、ARM服务器)
- Linux 386 (32位x86系统)
- Linux arm (树莓派等32位ARM设备)
- Linux ppc64le (PowerPC系统)
- Linux riscv64 (RISC-V架构)
- Linux s390x (IBM System z大型机)
- Linux mips64le (MIPS架构)

### *BSD系统 (6个)
- FreeBSD amd64/arm64
- OpenBSD amd64/arm64
- NetBSD amd64
- DragonFlyBSD amd64

### macOS (2个)
- macOS Intel (Intel处理器)
- macOS ARM (Apple M1/M2/M3)

## 💡 使用场景

### 1. 远程访问家里的电脑

```bash
# 在家里的Mac上启动客户端
./cligool-darwin-arm64 -server https://your-server.com

# 在办公室的浏览器中连接
# 使用生成的会话ID访问
```

### 2. 服务器管理

```bash
# 在Linux服务器上运行
./cligool-linux-amd64 -server https://your-server.com

# 在手机浏览器中管理服务器
```

### 3. 技术支持

```bash
# 朋友的电脑出现问题
# 让朋友下载并运行对应平台的客户端
# 你在浏览器中远程协助
```

### 4. 团队协作

```bash
# 多人同时连接同一会话
# 实时查看和操作终端
```

### 5. AI CLI工具远程使用

```bash
# 运行Claude CLI
./cligool-darwin-arm64 -cmd claude -server https://your-server.com

# 运行带参数的命令
./cligool-darwin-arm64 -cmd git -args "commit -m '修复bug'" -server https://your-server.com

# Windows运行Gemini
cligool-windows-amd64.exe -cmd gemini -args "chat --model gemini-pro" -server https://your-server.com
```

## 🔧 高级功能

### 心跳保活

系统实现了双向WebSocket心跳机制：
- **服务端 → 客户端**：每30秒发送ping
- **客户端 → 服务端**：自动回复pong
- **超时检测**：90秒无响应自动断开

### 自动重连

Web界面支持自动重连：
- 连接断开时自动尝试重连
- 指数退避策略避免频繁重连
- 可手动点击重连按钮

### 会话管理

会话在内存中自动管理：
- 客户端连接时自动创建会话
- 90秒无活动后自动清理会话
- 无需手动管理会话

## 🔍 故障排除

### 客户端连接失败

**问题**：`WebSocket连接失败`

**解决方案**：
```bash
# 1. 检查服务器是否运行
curl http://localhost:8081/api/health

# 2. 检查服务器日志
docker logs cligool-relay

# 3. 检查防火墙设置
# 确保端口8081可访问
```

### 终端无响应

**问题**：输入命令后无输出

**解决方案**：
1. 检查WebSocket连接状态（浏览器控制台）
2. 确认客户端进程仍在运行
3. 尝试刷新页面重新连接

### Windows客户端乱码

**问题**：中文显示为乱码

**解决方案**：
- ✅ 已自动修复：GBK编码自动转换为UTF-8
- 如仍有问题，检查终端编码设置

### PTY权限问题（Linux/macOS）

**问题**：`Failed to allocate PTY`

**解决方案**：
```bash
# 检查/dev/ptmx权限
ls -l /dev/ptmx

# 确保在真实终端中运行（非IDE内置终端）
# 某些IDE的终端可能不支持PTY
```

## 🚀 生产环境部署

### 使用Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f relay-server

# 停止服务
docker-compose down
```

### 使用Cloudflare Tunnel（HTTPS）

1. 安装cloudflared
2. 创建隧道：`cloudflared tunnel create cligool`
3. 配置`~/.cloudflared/config.yml`：
```yaml
tunnel: <your-tunnel-id>
credentials-file: /path/to/credentials.json

ingress:
  - hostname: cligool.yourdomain.com
    service: http://localhost:8081
  - service: http_status:404
```

4. 启动隧道：`cloudflared tunnel run`

### 环境变量

```bash
# 服务器配置
RELAY_HOST=0.0.0.0
RELAY_PORT=8080
```

## 📚 更多文档

- [README.md](../README.md) - 项目总览
- [DEPLOYMENT.md](DEPLOYMENT.md) - 部署指南
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南
- [PTY_TROUBLESHOOTING.md](PTY_TROUBLESHOOTING.md) - PTY问题排查

## 🤝 贡献

欢迎提交Issue和Pull Request！

1. Fork本仓库
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 开启Pull Request

## 📝 许可证

MIT License - 详见 [LICENSE](../LICENSE) 文件
