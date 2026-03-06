# CliGool - 跨平台远程终端解决方案

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platforms](https://img.shields.io/badge/platforms-18-blue)](#-支持的平台)

一个基于Go和WebSocket的跨平台远程终端解决方案，支持18种操作系统和架构。

## ✨ 核心特性

- 🌍 **跨平台支持**：18个操作系统和架构（Windows、Linux、macOS、*BSD等）
- ⚡ **低延迟**：WebSocket实时通信，毫秒级响应
- 🔒 **安全连接**：端到端加密通信，支持HTTPS/WSS
- 💎 **真实PTY**：完整的终端特性支持（颜色、光标控制等）
- 👥 **多用户协作**：多人可同时连接同一终端会话
- 🚀 **开箱即用**：Docker一键部署，无需复杂配置
- 🎨 **现代Web界面**：基于xterm.js的专业终端UI
- 💓 **心跳保活**：双向心跳机制，自动重连和死连接清理

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
- ✅ **会话管理**：支持会话创建、删除、列表查询
- ✅ **数据库持久化**：PostgreSQL存储会话信息

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

## 🚀 快速开始

### 方法一：Docker Compose部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 启动所有服务
docker-compose up -d

# 3. 检查服务状态
docker-compose ps
```

服务启动后：
- 中继服务器：http://localhost:8081
- 下载页面：http://localhost:8081/
- API健康检查：http://localhost:8081/api/health

### 方法二：手动构建和运行

#### 1. 构建客户端

```bash
# 构建当前平台的客户端
go build -o cligool ./cmd/client

# 或交叉编译其他平台（示例：Linux amd64）
GOOS=linux GOARCH=amd64 go build -o cligool-linux-amd64 ./cmd/client
```

#### 2. 启动中继服务器

```bash
# 使用Docker
docker-compose up -d relay-server

# 或直接运行
go run ./cmd/relay
```

#### 3. 启动CLI客户端

```bash
# 连接到本地服务器
./cligool -server http://localhost:8081

# 连接到远程服务器
./cligool -server https://your-domain.com
```

#### 4. 访问Web界面

客户端启动后会显示：
```
╔═══════════════════════════════════════════════════════════╗
║                    🚀 CliGool 远程终端                      ║
╠═══════════════════════════════════════════════════════════╣
║ 📋 会话ID: [会话ID]                                        ║
║ 🌐 Web访问: http://localhost:8081/session/[会话ID]         ║
║ 🔗 连接状态: 🟢 已连接                                     ║
╚═══════════════════════════════════════════════════════════╝
```

在浏览器中打开显示的地址即可使用远程终端。

## 📂 项目结构

```
cligool/
├── cmd/
│   ├── relay/          # 中继服务器
│   └── client/         # CLI客户端
│       ├── main.go          # Unix/Linux/macOS客户端（PTY）
│       └── main_windows.go  # Windows客户端（cmd.exe）
├── internal/
│   ├── relay/          # 中继服务逻辑
│   └── database/       # 数据库层
├── web/
│   ├── landing.html    # 下载和介绍页面
│   ├── terminal.html   # 终端Web界面
│   ├── lib/           # xterm.js库文件
│   └── downloads/     # 各平台客户端二进制文件
├── Dockerfile         # 多平台构建
├── docker-compose.yml # 服务编排
└── README.md
```

## 🔧 技术栈

- **中继服务**：Go 1.21 + Gin + WebSocket + PostgreSQL + Redis
- **CLI客户端**：Go + PTY (Unix) / 管道 (Windows)
- **Web界面**：xterm.js + 原生JavaScript
- **部署**：Docker + docker-compose
- **反向代理**：Cloudflare Tunnel（可选，零配置HTTPS）

## 📖 使用场景

### 1. 远程访问

```bash
# 在家里的Mac上启动客户端
./cligool-darwin-arm64 -server https://your-server.com

# 在办公室的浏览器中连接
# 使用生成的会话ID访问
```

### 2. 技术支持

```bash
# 朋友的电脑出现问题
# 让朋友下载并运行对应平台的客户端
# 你在浏览器中远程协助
```

### 3. 服务器管理

```bash
# 在Linux服务器上运行
./cligool-linux-amd64 -server https://your-server.com

# 在手机浏览器中管理服务器
```

### 4. 团队协作

```bash
# 多人同时连接同一会话
# 实时查看和操作终端
```

## 🛠️ 开发指南

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- (可选) Cloudflare Tunnel账号

### 构建项目

```bash
# 构建中继服务器
go build -o bin/relay-server ./cmd/relay

# 构建当前平台的客户端
go build -o bin/cligool ./cmd/client

# 构建所有平台（使用Docker）
docker build -t cligool-relay-server .
```

### 本地开发

```bash
# 启动依赖服务
docker-compose up -d postgres redis

# 运行中继服务器
go run ./cmd/relay

# 运行客户端
go run ./cmd/client
```

## 🔍 故障排除

### 常见问题

**1. 端口被占用**
```bash
# 修改docker-compose.yml中的端口映射
ports:
  - "8082:8080"  # 改为其他端口
```

**2. 客户端连接失败**
```bash
# 检查服务器状态
curl http://localhost:8081/api/health

# 查看服务器日志
docker logs cligool-relay
```

**3. Windows客户端输出乱码**
- 已修复：GBK编码自动转换为UTF-8

**4. 终端无响应**
- 检查心跳机制是否正常（30秒间隔）
- 查看浏览器控制台WebSocket连接状态

## 🚀 生产环境部署

### 使用Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

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
# 数据库连接
DATABASE_URL=postgres://user:pass@host:5432/cligool?sslmode=disable

# Redis连接
REDIS_URL=redis://host:6379

# 服务器配置
RELAY_HOST=0.0.0.0
RELAY_PORT=8080
```

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📝 许可证

本项目采用MIT许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [xterm.js](https://xtermjs.org/) - 终端仿真器
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket库
- [creack/pty](https://github.com/creack/pty) - PTY库

---

**注意**：本项目仅用于授权的远程访问和管理使用。请遵守相关法律法规，不得用于非法用途。
