# CliGool - 跨平台远程终端解决方案

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platforms](https://img.shields.io/badge/platforms-30-blue)](#-支持的平台)
[![Demo](https://img.shields.io/badge/demo-online-success.svg)](https://cligool.zty8.cn/)

一个基于Go和WebSocket的跨平台远程终端解决方案，支持30种操作系统和架构。

**[🚀 在线体验](https://cligool.zty8.cn/)** | **[📥 下载客户端](https://cligool.zty8.cn/)**

## ✨ 核心特性

- 🌍 **跨平台支持**：30个操作系统和架构（Windows、Linux、macOS、*BSD等）
- ⚡ **低延迟**：WebSocket实时通信，毫秒级响应
- 🔒 **安全连接**：端到端加密通信，支持HTTPS/WSS
- 💎 **真实PTY**：完整的终端特性支持（颜色、光标控制等）
- 🤖 **AI CLI工具**：完美支持Claude、Gemini、Aider等AI CLI工具
- 👥 **多用户协作**：多人可同时连接同一终端会话
- 📦 **消息缓存**：无Web客户端时自动缓存CLI消息，连接时自动恢复历史
- 🚀 **开箱即用**：Docker一键部署，无需复杂配置
- 🎨 **现代Web界面**：基于xterm.js的专业终端UI
- 💓 **心跳保活**：双向心跳机制，自动重连和死连接清理
- 🖥️ **Windows ConPTY**：Windows版本使用ConPTY，功能与Unix完全对等

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
- ✅ **开箱即用**：单容器部署，无需复杂配置

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

## 🎮 在线体验

不想本地部署？可以先体验在线Demo！

### 访问在线Demo

👉 **[https://cligool.zty8.cn/](https://cligool.zty8.cn/)**

### Demo使用步骤

1. **下载客户端**：在下载页面选择你操作系统的客户端
2. **运行客户端**：启动下载的客户端程序
3. **访问Web终端**：使用客户端显示的会话ID访问Web界面
4. **开始体验**：在浏览器中体验完整的远程终端功能

**注意**：在线Demo仅供体验使用，可能随时关闭。建议自行部署以获得稳定服务。

## 🚀 快速开始

### 方法一：Docker Compose部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 生产环境（使用预构建镜像）
docker-compose up -d

# 或：开发环境（本地构建）
docker-compose -f docker-compose.dev.yml up -d --build

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
# 使用默认shell连接到本地服务器
./cligool -server http://localhost:8081

# 连接到远程服务器
./cligool -server https://your-domain.com

# 使用配置文件（推荐）
# 首次运行会自动创建 ~/.cligool.json 配置文件
./cligool

# 运行AI CLI工具（如Claude）
./cligool -cmd claude

# 运行带参数的命令
./cligool -cmd git -args "status"

# 禁止自动打开浏览器
./cligool -server http://localhost:8081 -no-browser
```

**常用参数**：
- `-server` : 中继服务器地址
- `-session` : 指定会话ID
- `-cols` / `-rows` : 设置终端大小（0=自动检测）
- `-cmd` : 直接执行指定命令
- `-args` : 传递给命令的参数
- `-proxy` : 使用代理服务器
- `-no-browser` : 禁止自动打开浏览器

**配置文件支持**：
CliGool 支持配置文件来设置常用参数：
- 配置文件位置：`./cligool.json` 或 `~/.cligool.json`
- 自动创建：首次运行时自动创建 `~/.cligool.json`
- 可配置项：`server`、`cols`、`rows`
- 命令行参数优先级高于配置文件

详见：[配置文件使用指南](CONFIG_GUIDE.md)

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

### 方法三：云平台一键部署（免费）☁️

不想自己维护服务器？可以使用免费云平台一键部署！

#### 平台对比

| 平台 | 免费额度 | 部署难度 | 冷启动 | 推荐指数 |
|------|---------|---------|--------|---------|
| **Render** | ✅ 完全免费 | ⭐ 最简单 | 15-30秒 | ⭐⭐⭐⭐⭐ |
| **Koyeb** | $5.5/月 | ⭐⭐ 中等 | 无 | ⭐⭐⭐⭐ |
| **Railway** | $5/月 | ⭐⭐ 简单 | 无 | ⭐⭐⭐⭐⭐ |
| **Fly.io** | 部分免费 | ⭐⭐ 中等 | 无 | ⭐⭐⭐⭐ |

#### Render 部署（推荐，完全免费）

```bash
# 1. 访问 Render Blueprint
open https://dashboard.render.com/blueprints/new

# 2. 连接你的 GitHub 账号
# 3. 选择 topcheer/cligool 仓库
# 4. 点击 "Apply Blueprint"
# 5. 等待 3-5 分钟，获得部署 URL
```

**详细指南**: [docs/CLOUD_DEPLOYMENT_GUIDE.md](docs/CLOUD_DEPLOYMENT_GUIDE.md)

#### Railway 部署

```bash
# 1. 访问 Railway
open https://railway.app

# 2. 点击 "Deploy from GitHub"
# 3. 选择 topcheer/cligool 仓库
# 4. Railway 自动检测 railway.toml 配置
# 5. 部署完成，获得 URL
```

#### Koyeb 部署

```bash
# 1. 安装 Koyeb CLI
curl -s https://get.koyeb.com | sh

# 2. 登录并部署
koyeb login
koyeb init

# Koyeb 自动读取 koyeb.yaml 配置并部署
```

#### Fly.io 部署

```bash
# 1. 安装 Fly CLI
curl -L https://fly.io/install.sh | sh

# 2. 登录并部署
flyctl auth login
flyctl launch

# Fly 自动读取 fly.toml 配置并部署
```

**所有云平台配置文件已包含在仓库中，开箱即用！**

## 📂 项目结构

```
cligool/
├── cmd/
│   ├── relay/          # 中继服务器
│   └── client/         # CLI客户端
│       ├── main_unix.go     # Unix/Linux/macOS客户端（PTY）
│       └── main_windows.go  # Windows客户端（ConPTY）
├── internal/
│   └── relay/          # 中继服务逻辑
├── web/
│   ├── landing.html    # 下载和介绍页面
│   ├── terminal.html   # 终端Web界面
│   ├── lib/           # xterm.js库文件
│   └── downloads/     # 各平台客户端二进制文件
├── Dockerfile.multiarch  # 多平台构建
├── docker-compose.yml     # 生产环境配置
├── docker-compose.dev.yml # 开发环境配置
└── README.md
```

## 🔧 技术栈

- **中继服务**：Go 1.21 + Gin + WebSocket
- **CLI客户端**：Go + PTY (Unix/macOS/Linux) / ConPTY (Windows)
- **Web界面**：xterm.js + 原生JavaScript
- **部署**：Docker + docker-compose
- **反向代理**：Cloudflare Tunnel（可选，零配置HTTPS）

**关键技术特性**：
- Windows ConPTY：完整的伪终端支持，功能与Unix版本完全对等
- 自动编码检测：Windows自动检测并转换为UTF-8
- 动态窗口大小：Unix (SIGWINCH) / Windows (监控)
- AI CLI工具支持：完美支持Claude、Gemini、Aider等工具

## 📖 使用场景

### 1. AI CLI工具远程访问

```bash
# 在家里的Mac上启动Claude CLI
./cligool-darwin-arm64 -cmd claude -server https://your-server.com

# 在办公室的浏览器中继续使用Claude
# 使用生成的会话ID访问

# 运行Gemini CLI
./cligool-linux-amd64 -cmd gemini -server https://your-server.com

# 运行带参数的命令
./cligool-darwin-arm64 -cmd git -args "commit -m 'fix bug'" -server https://your-server.com
```

### 2. 远程访问

```bash
# 在家里的Mac上启动客户端
./cligool-darwin-arm64 -server https://your-server.com

# 在办公室的浏览器中连接
# 使用生成的会话ID访问
```

### 3. 技术支持

```bash
# 朋友的电脑出现问题
# 让朋友下载并运行对应平台的客户端
# 你在浏览器中远程协助
```

### 4. 服务器管理

```bash
# 在Linux服务器上运行
./cligool-linux-amd64 -server https://your-server.com

# 在手机浏览器中管理服务器
```

### 5. 团队协作

```bash
# 多人同时连接同一会话
# 实时查看和操作终端
```

### 6. AI CLI工具使用

```bash
# 运行Claude CLI
./cligool-darwin-arm64 -cmd claude -server https://your-domain.com

# 运行带参数的命令
./cligool-darwin-arm64 -cmd git -args "commit -m 'Add new feature'" -server https://your-domain.com

# Windows运行Gemini
cligool-windows-amd64.exe -cmd gemini -args "chat --model gemini-pro" -server https://your-domain.com
```

**详细文档**：
- [命令行参数使用](CMD_ARGS_USAGE.md)
- [AI CLI工具指南](docs/AI_CLI_GUIDE.md)
- [Windows支持说明](docs/WINDOWS_SUPPORT.md)

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
# 使用开发环境配置启动服务
docker-compose -f docker-compose.dev.yml up -d --build

# 运行客户端（另开终端）
go run ./cmd/client

# 查看日志
docker-compose -f docker-compose.dev.yml logs -f relay-server
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
- 已修复：自动检测控制台编码并转换为UTF-8（ConPTY + UTF-8检查）

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
