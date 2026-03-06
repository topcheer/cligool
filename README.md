# CliGool - 远程终端解决方案

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个基于Go和WebSocket的远程终端解决方案，支持多用户协作和真实PTY体验。

## 🎉 架构修复完成！

✅ **重大更新**：已修复JSON消息格式不匹配问题，采用Base64编码确保数据正确传递！

## 🏗️ 正确的系统架构

```
用户A的电脑           中继服务器              用户B的浏览器
┌──────────────┐              ┌──────────┐              ┌─────────────┐
│ CLI客户端     │──WebSocket──▶│ 消息转发器│◀──WebSocket───│  独立HTML     │
│ (真实PTY)     │              │          │              │  (本地文件)   │
└──────────────┘              └──────────┘              └─────────────┘
```

**核心特点**：
- ✅ **独立组件**：CLI客户端、中继服务器、Web界面完全分离
- ✅ **Base64编码**：解决了消息格式兼容性问题
- ✅ **真实PTY**：只有CLI客户端提供真实的终端环境
- ✅ **本地Web界面**：HTML文件可以在任何地方打开使用

## 🚀 快速开始

### 方法一：使用演示脚本（推荐）

```bash
# 一键启动完整演示环境
./demo.sh

# 停止演示环境
./stop-demo.sh
```

### 方法二：手动启动

#### 1. 启动中继服务器

```bash
# 构建并启动所有服务
docker-compose up -d

# 检查服务状态
docker-compose ps
```

#### 2. 启动CLI客户端

```bash
# 连接到本地服务器
./bin/cligool-simple -server http://localhost:8081 -connect-only

# 连接到远程服务器
./bin/cligool-simple -server https://cligool.zty8.cn -connect-only
```

#### 3. 打开Web界面

```bash
# 在浏览器中打开
open web-client.html
```

## 📋 项目结构

## 🚀 功能特性

- ✅ 实时终端访问和控制
- ✅ 完整的终端特性支持（颜色、光标控制等）
- ✅ 多用户协作功能
- ✅ 会话管理和权限控制
- ✅ 端到端加密通信
- ✅ Cloudflare Tunnel支持 (零配置HTTPS)
- ✅ 容器化部署，一键启动

## 🛠️ 技术栈

- **中继服务**: Go + WebSocket + PostgreSQL + Redis
- **CLI客户端**: Go
- **Web界面**: 原生JavaScript + xterm.js
- **部署**: Docker + Cloudflare Tunnel

## 📋 核心组件

### 中继服务器 (`cmd/relay/`)

基于Gin框架的WebSocket中继服务器，负责：
- 管理终端会话
- 转发客户端消息
- 处理多用户连接
- 支持健康检查

**启动命令**：
```bash
go run ./cmd/relay
# 或使用Docker
docker-compose up relay-server
```

### CLI客户端 (`cmd/client/`)

提供真实PTY环境的命令行客户端：

**简化版** (`simple.go`)：
```bash
./bin/cligool-simple -server http://localhost:8081 -connect-only
```

**完整版** (`main.go`)：
```bash
./bin/cligool-client -server http://localhost:8081 -shell /bin/bash
```

### Web界面 (`web-client.html`)

独立的HTML文件，包含：
- xterm.js终端仿真
- WebSocket客户端
- 完整的UI界面

**特点**：
- 单一文件，无需额外依赖
- 可本地打开或托管到任何地方
- 支持多终端同步

## 🔧 技术栈

- **后端**：Go 1.21+, Gin框架
- **前端**：xterm.js, 原生JavaScript
- **数据库**：PostgreSQL 15
- **缓存**：Redis 7
- **容器化**：Docker, docker-compose
- **反向代理**：Cloudflare Tunnel（可选）

## 📖 详细文档

- [完整使用指南](USAGE_GUIDE_CN.md) - 详细的中文使用说明
- [PTY故障排除](docs/PTY_TROUBLESHOOTING.md) - PTY相关问题解决

## 🌟 使用场景

### 1. 远程访问

```bash
# 在家里的Mac上启动客户端
./bin/cligool-simple -connect-only

# 在办公室的浏览器中连接
# 使用生成的会话ID
```

### 2. 技术支持

```bash
# 朋友的电脑出现问题
# 让朋友运行客户端
# 你在浏览器中远程协助
```

### 3. 团队协作

```bash
# 多人同时连接同一会话
# 实时查看和操作
```

## 🛠️ 开发指南

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- (可选) Cloudflare Tunnel

### 构建项目

```bash
# 构建中继服务器
go build -o bin/relay-server ./cmd/relay

# 构建CLI客户端
go build -o bin/cligool-client ./cmd/client
go build -o bin/cligool-simple ./cmd/client/simple.go

# 构建Docker镜像
docker build -t cligool-relay-server .
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

**3. PTY权限问题**
- 在真实终端中运行（非IDE内置终端）
- 检查 `/dev/ptmx` 权限
- 参考 [PTY故障排除指南](docs/PTY_TROUBLESHOOTING.md)

## 🚀 生产环境部署

### 使用Docker

```bash
# 构建生产镜像
docker build -t cligool-relay:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/dbname \
  -e REDIS_URL=redis://host:6379 \
  --name cligool-relay \
  cligool-relay:latest
```

### 使用Cloudflare Tunnel

1. 复制配置文件：
```bash
cp cloudflare-tunnel.yml.example ~/.cloudflared/config.yml
```

2. 修改域名和服务地址

3. 启动隧道：
```bash
cloudflared tunnel run
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

---

**注意**：本项目仅用于授权的远程访问和管理使用。请遵守相关法律法规，不得用于非法用途。