# CliGool 配置指南

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+ (仅开发环境)
- Docker & Docker Compose
- 任意支持Docker的服务器

### 2. 快速部署

```bash
# 克隆项目
git clone https://github.com/cligool/cligool.git
cd cligool

# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps
```

### 3. 访问服务

```bash
# 健康检查
curl http://localhost:8081/api/health

# 访问Web界面
open http://localhost:8081
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `RELAY_HOST` | 中继服务监听地址 | 0.0.0.0 | ❌ |
| `RELAY_PORT` | 中继服务监听端口 | 8080 | ❌ |

**注意**：CliGool Relay Server 现在是无状态的，不需要数据库或 Redis。

### Docker Compose 配置

```yaml
services:
  relay-server:
    image: ghcr.io/topcheer/cligool:latest
    ports:
      - "8081:8080"
    environment:
      - RELAY_HOST=0.0.0.0
      - RELAY_PORT=8080
    restart: unless-stopped
```

### 本地开发配置

```bash
# 运行 relay server
go run cmd/relay/main.go

# 指定端口
RELAY_PORT=9090 go run cmd/relay/main.go

# 指定监听地址
RELAY_HOST=127.0.0.1 go run cmd/relay/main.go
```

## 🔒 安全配置

### 反向代理配置

#### Nginx 示例

```nginx
server {
    listen 80;
    server_name terminal.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### Caddy 示例

```
terminal.yourdomain.com {
    reverse_proxy localhost:8080
}
```

### Cloudflare Tunnel

使用 Cloudflare Tunnel 可以获得：
- ✅ 自动 HTTPS
- ✅ DDoS 保护
- ✅ 全球 CDN
- ✅ 隐藏源站 IP

**安装 cloudflared**：
```bash
# macOS
brew install cloudflare/cloudflare/cloudflared

# Linux
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

**创建 Tunnel**：
```bash
# 登录
cloudflared tunnel login

# 创建 tunnel
cloudflared tunnel create cligool

# 配置 DNS
cloudflared tunnel route dns cligool terminal.yourdomain.com
```

**配置文件** (`cloudflare-tunnel.yml`)：
```yaml
tunnel: <你的-tunnel-id>
credentials-file: /path/to/credentials.json

ingress:
  - hostname: terminal.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404
```

**启动 Tunnel**：
```bash
cloudflared tunnel --config cloudflare-tunnel.yml run
```

## 🌐 网络配置

### 端口说明

- **8080**：容器内部端口
- **8081**：主机端口（可通过 docker-compose.yml 修改）

### 防火墙配置

如果使用 Cloudflare Tunnel，不需要开放 8080 端口到公网。

如果直接暴露服务：
```bash
# 开放 HTTP 端口
ufw allow 80/tcp
ufw allow 443/tcp

# 允许 SSH 管理（可选）
ufw allow 22/tcp
```

## 📊 监控和日志

### 健康检查

```bash
curl http://localhost:8081/api/health
```

响应：
```json
{
  "status": "ok",
  "time": 1699123456
}
```

### 查看日志

```bash
# Docker 日志
docker-compose logs -f relay-server

# 最近 100 行
docker-compose logs --tail=100 relay-server
```

### 性能监控

CliGool Relay Server 是无状态的，可以通过以下方式监控：

1. **健康检查端点**：定期检查 `/api/health`
2. **WebSocket 连接数**：查看日志中的连接信息
3. **内存使用**：`docker stats cligool-relay`

## 🔧 性能优化

### 资源限制

在 `docker-compose.yml` 中设置资源限制：

```yaml
services:
  relay-server:
    image: ghcr.io/topcheer/cligool:latest
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          memory: 128M
```

### 连接优化

Relay Server 会自动管理 WebSocket 连接：
- 心跳检测：每 30 秒发送 ping
- 超时断开：90 秒无响应自动断开
- 自动清理：断开后立即释放资源

## 🚀 部署场景

### 场景 1：个人服务器

```bash
git clone https://github.com/cligool/cligool.git
cd cligool
docker-compose up -d
```

配置 Cloudflare Tunnel 即可公网访问。

### 场景 2：云平台部署

参考 **云平台部署指南** (`docs/CLOUD_DEPLOYMENT_GUIDE.md`)：

- **Render**：完全免费，推荐新手
- **Koyeb**：无冷启动，性能最好
- **Railway**：开发体验最好

### 场景 3：Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cligool-relay
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cligool-relay
  template:
    metadata:
      labels:
        app: cligool-relay
    spec:
      containers:
      - name: relay
        image: ghcr.io/topcheer/cligool:latest
        ports:
        - containerPort: 8080
        env:
        - name: RELAY_HOST
          value: "0.0.0.0"
        - name: RELAY_PORT
          value: "8080"
---
apiVersion: v1
kind: Service
metadata:
  name: cligool-relay
spec:
  selector:
    app: cligool-relay
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 🐛 故障排除

### 常见问题

**1. WebSocket 连接失败**
```bash
# 检查服务状态
curl http://localhost:8081/api/health

# 查看日志
docker-compose logs -f relay-server

# 检查端口是否被占用
lsof -i :8081
```

**2. 客户端无法连接**
- 确认服务正在运行
- 检查防火墙设置
- 验证 WebSocket URL 格式：
  ```
  wss://你的域名/api/terminal/session-id?type=web&user_id=web-123
  ```

**3. 页面无法加载**
- 检查静态文件是否正确挂载
- 查看浏览器控制台错误
- 验证 Docker 镜像是否包含 `web/` 目录

### 调试模式

启用详细日志：
```bash
# 修改 cmd/relay/main.go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### 网络测试

```bash
# 测试本地服务
curl http://localhost:8081/api/health

# 测试 WebSocket 连接
wscat -c ws://localhost:8081/api/terminal/test-session

# 测试远程服务
curl https://你的域名/api/health
```

## 🔐 安全最佳实践

1. **使用 HTTPS**
   - 配置 SSL 证书
   - 使用 Let's Encrypt（免费）
   - 或使用 Cloudflare Tunnel

2. **访问控制**
   - 配置 Cloudflare Access
   - 添加身份验证
   - 设置 IP 白名单

3. **定期更新**
   ```bash
   # 更新 Docker 镜像
   docker-compose pull
   docker-compose up -d
   ```

4. **监控和日志**
   - 启用访问日志
   - 设置告警通知
   - 定期审计日志

5. **备份策略**
   - 备份配置文件
   - 备份 Cloudflare 配置
   - 记录自定义设置

## 📚 相关文档

- **云平台部署指南**: `docs/CLOUD_DEPLOYMENT_GUIDE.md`
- **使用指南**: `USAGE_GUIDE.md`
- **开发指南**: `docs/DEVELOPMENT.md`
- **Windows 支持**: `docs/WINDOWS_SUPPORT.md`

---

**祝配置顺利！🎉**
