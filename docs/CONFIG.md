# CliGool 配置指南

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+ (仅开发环境)
- Docker & Docker Compose
- Cloudflare账号 (用于Tunnel)
- 任意支持Docker的服务器

### 2. 快速部署

```bash
# 克隆项目
git clone https://github.com/cligool/cligool.git
cd cligool

# 复制环境变量文件
cp .env.example .env

# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps
```

### 3. Cloudflare Tunnel配置

#### 安装cloudflared
```bash
# macOS
brew install cloudflare/cloudflare/cloudflared

# Linux
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb

# Docker
docker pull cloudflare/cloudflared
```

#### 创建Tunnel
```bash
# 登录Cloudflare
cloudflared tunnel login

# 创建tunnel
cloudflared tunnel create cligool

# 会显示tunnel ID，记录下来
# 例如: Tunnel ID: abc123xyz-def456-ghi789
```

#### 配置Tunnel
创建 `cloudflare-tunnel.yml`:
```yaml
tunnel: abc123xyz-def456-ghi789  # 你的tunnel ID
credentials-file: /path/to/abc123xyz-def456-ghi789.json

ingress:
  - hostname: cligool.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404
```

#### 启动Tunnel
```bash
# 方式1: 直接运行
cloudflared tunnel --config cloudflare-tunnel.yml run

# 方式2: Docker运行
docker run -d --name cloudflared \
  -v /path/to/cloudflare-tunnel.yml:/home/cloudflared/.cloudflared/config.yml \
  -v /path/to/credentials.json:/home/cloudflared/.cloudflared/abc123xyz-def456-ghi789.json \
  cloudflare/cloudflared:latest
```

#### 配置DNS
```bash
# 添加DNS记录
cloudflared tunnel route dns cligool cligool.yourdomain.com
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `DATABASE_URL` | PostgreSQL连接字符串 | - | ✅ |
| `REDIS_URL` | Redis连接字符串 | - | ✅ |
| `RELAY_HOST` | 中继服务监听地址 | 0.0.0.0 | ❌ |
| `RELAY_PORT` | 中继服务监听端口 | 8080 | ❌ |
| `JWT_SECRET` | JWT密钥 | - | ✅ |
| `LOG_LEVEL` | 日志级别 | info | ❌ |
| `SESSION_TIMEOUT` | 会话超时时间(秒) | 3600 | ❌ |

### 数据库配置

**PostgreSQL:**
```bash
DATABASE_URL=postgres://username:password@hostname:5432/database_name?sslmode=disable
```

**Redis:**
```bash
REDIS_URL=redis://hostname:6379
```

## 🔒 安全配置

### JWT认证

生成安全的JWT密钥：
```bash
# 生成随机密钥
openssl rand -base64 32
```

更新 `.env` 文件：
```bash
JWT_SECRET=your-generated-secret-key-here
```

### Cloudflare Tunnel安全优势

使用Cloudflare Tunnel可以获得以下安全特性：

1. **自动HTTPS** - 无需管理证书
2. **DDoS保护** - Cloudflare的防护网络
3. **访问控制** - 可以配置Access策略
4. **隐藏源站** - 不暴露服务器真实IP
5. **全球CDN** - 加速全球访问

### 访问控制

在Cloudflare Zero Trust中配置：

1. **基础认证**
   - 电子邮件验证
   - Google/Microsoft登录
   - 硬件密钥（YubiKey等）

2. **企业认证**
   - SAML 2.0
   - OIDC
   - LDAP集成

3. **设备策略**
   - 设备健康检查
   - 地理位置限制
   - IP白名单

## 🌐 网络配置

### 端口配置

中继服务只需要暴露HTTP端口：
```yaml
ports:
  - "8080:8080"  # 只需要HTTP端口
```

### 防火墙配置

由于使用Cloudflare Tunnel，只需要：

```bash
# 允许SSH管理（可选）
ufw allow 22/tcp

# 允许本地Docker网络通信
ufw allow from 172.16.0.0/12
```

**注意**: 不需要开放8080端口到公网，Cloudflare Tunnel会建立出站连接。

## 📊 监控配置

### 健康检查端点

```bash
curl http://localhost:8080/api/health
```

响应：
```json
{
  "status": "ok",
  "time": 1699123456
}
```

### Cloudflare监控

在Cloudflare Zero Trust控制台可以监控：
- 流量统计
- 连接数
- 响应时间
- 错误率

### 日志配置

```bash
# 日志级别
LOG_LEVEL=info  # debug, info, warn, error

# 日志格式
LOG_FORMAT=json  # json, text
```

查看日志：
```bash
# Docker日志
docker-compose logs -f relay-server

# 特定服务
docker-compose logs -f relay-server
docker-compose logs -f postgres
docker-compose logs -f redis
```

## 🔧 性能优化

### 连接池配置

```go
// PostgreSQL连接池（已优化）
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Redis配置

```bash
# Redis持久化（已配置）
redis-server --appendonly yes
```

### 终端性能

```bash
# 消息大小限制
MAX_MESSAGE_SIZE=1048576  # 1MB

# 缓冲区大小
READ_BUFFER_SIZE=4096
WRITE_BUFFER_SIZE=4096
```

## 🚀 部署场景

### 场景1: 个人服务器

```bash
# 在你的服务器上
git clone https://github.com/cligool/cligool.git
cd cligool
docker-compose up -d

# 配置Cloudflare Tunnel
cloudflared tunnel create cligool-personal
cloudflared tunnel route dns cligool-personal terminal.yourdomain.com
```

### 场景2: 团队协作

```bash
# 部署多个实例
docker-compose -f docker-compose.yml up -d --scale relay-server=3

# 配置负载均衡
# 在Cloudflare中配置负载均衡策略
```

### 场景3: 企业部署

```bash
# 使用私有域名
# 配置企业SSO
# 设置访问策略
# 启用审计日志
```

## 🐛 故障排除

### 常见问题

1. **Cloudflare Tunnel连接失败**
   ```bash
   # 检查tunnel状态
   cloudflared tunnel info cligool

   # 查看tunnel日志
   cloudflared --config cloudflare-tunnel.yml run --loglevel debug
   ```

2. **数据库连接错误**
   ```bash
   # 检查数据库状态
   docker-compose ps postgres

   # 查看数据库日志
   docker-compose logs postgres
   ```

3. **WebSocket连接断开**
   - 检查Cloudflare Tunnel设置
   - 验证防火墙规则
   - 确认服务运行状态

### 日志查看

```bash
# 服务日志
docker-compose logs -f relay-server

# 所有服务日志
docker-compose logs -f

# 最近100行
docker-compose logs --tail=100 relay-server
```

### 网络测试

```bash
# 测试本地服务
curl http://localhost:8080/api/health

# 测试DNS解析
nslookup terminal.yourdomain.com

# 测试Cloudflare连接
curl https://terminal.yourdomain.com/api/health
```

## 📦 备份和恢复

### 数据库备份

```bash
# 备份PostgreSQL
docker exec cligool-postgres pg_dump -U cligool cligool > backup.sql

# 恢复PostgreSQL
docker exec -i cligool-postgres psql -U cligool cligool < backup.sql
```

### Redis备份

```bash
# 备份Redis
docker exec cligool-redis redis-cli BGSAVE

# 备份文件位置
docker cp cligool-redis:/data/dump.rdb ./redis-backup.rdb
```

### 配置备份

```bash
# 备份环境配置
cp .env .env.backup

# 备份Cloudflare配置
cp cloudflare-tunnel.yml cloudflare-tunnel.yml.backup
```

## 🔐 安全最佳实践

1. **使用Cloudflare Access**
   - 配置身份验证
   - 设置访问策略
   - 启用设备验证

2. **定期更新**
   ```bash
   # 更新Docker镜像
   docker-compose pull
   docker-compose up -d
   ```

3. **监控和日志**
   - 启用Cloudflare Analytics
   - 设置告警通知
   - 定期审计访问日志

4. **备份策略**
   - 定期备份数据库
   - 备份配置文件
   - 测试恢复流程