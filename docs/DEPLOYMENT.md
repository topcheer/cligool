# CliGool 部署指南

## 🚀 快速部署

### Docker Compose 一键部署（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 启动所有服务
docker-compose up -d

# 3. 检查服务状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f relay-server
```

就这样！你的服务现在运行在 `http://localhost:8081`

## 📋 详细部署步骤

### 第一步: 准备服务器

#### 最低要求
- **操作系统**: Linux (任何发行版)、macOS、Windows
- **内存**: 512MB+
- **磁盘**: 10GB+
- **网络**: 能够访问Docker Hub

#### 安装Docker

**Ubuntu/Debian**:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**CentOS/RHEL**:
```bash
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
```

**macOS**:
```bash
brew install --cask docker
# 启动Docker Desktop应用
```

**Windows**:
- 下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop)

#### 安装Docker Compose

**Linux**:
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

**macOS/Windows**:
- 已包含在Docker Desktop中

### 第二步: 部署应用

```bash
# 1. 克隆仓库
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 配置环境变量（可选）
cp .env.example .env
# 编辑.env文件设置数据库密码等

# 3. 启动所有服务
docker-compose up -d

# 4. 等待服务启动（约30秒）
sleep 30

# 5. 检查服务状态
docker-compose ps
```

**预期输出**:
```
NAME                STATUS
cligool-postgres    Up (healthy)
cligool-redis       Up (healthy)
cligool-relay       Up (healthy)
```

### 第三步: 验证部署

```bash
# 检查API健康状态
curl http://localhost:8081/api/health

# 预期输出: {"status":"ok","time":1234567890}

# 访问下载页面
open http://localhost:8081/
```

## 🌐 配置HTTPS（可选）

### 方法一：使用Cloudflare Tunnel（推荐）

Cloudflare Tunnel提供免费的HTTPS和DDoS保护。

#### 1. 安装cloudflared

**Linux (amd64)**:
```bash
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

**Linux (arm64)**:
```bash
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64.deb
sudo dpkg -i cloudflared-linux-arm64.deb
```

**macOS**:
```bash
brew install cloudflared
```

**Windows**:
- 下载 [cloudflared-windows-amd64.exe](https://github.com/cloudflare/cloudflared/releases/latest)
- 重命名为`cloudflared.exe`并添加到PATH

#### 2. 登录Cloudflare

```bash
cloudflared tunnel login
```

这会打开浏览器进行授权。

#### 3. 创建隧道

```bash
# 创建隧道（记录返回的隧道ID）
cloudflared tunnel create cligool
```

#### 4. 配置隧道

创建配置文件 `~/.cloudflared/config.yml`:

```yaml
tunnel: <your-tunnel-id>
credentials-file: /root/.cloudflared/<your-tunnel-id>.json

ingress:
  - hostname: cligool.yourdomain.com
    service: http://localhost:8081
  - service: http_status:404
```

#### 5. 启动隧道

```bash
# 测试运行
cloudflared tunnel run

# 或作为系统服务运行
sudo cloudflared service install
sudo cloudflared service start
```

### 方法二：使用Nginx反向代理

#### 1. 安装Nginx

```bash
# Ubuntu/Debian
sudo apt install nginx certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install nginx certbot python3-certbot-nginx
```

#### 2. 配置Nginx

创建 `/etc/nginx/sites-available/cligool`:

```nginx
server {
    listen 80;
    server_name cligool.yourdomain.com;

    location / {
        proxy_pass http://localhost:8081;
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

#### 3. 启用配置

```bash
sudo ln -s /etc/nginx/sites-available/cligool /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

#### 4. 获取SSL证书

```bash
sudo certbot --nginx -d cligool.yourdomain.com
```

## 🔧 环境变量配置

创建 `.env` 文件：

```bash
# PostgreSQL配置
POSTGRES_DB=cligool
POSTGRES_USER=cligool
POSTGRES_PASSWORD=your_secure_password

# 数据库连接URL
DATABASE_URL=postgres://cligool:your_secure_password@postgres:5432/cligool?sslmode=disable

# Redis配置
REDIS_URL=redis://redis:6379

# 服务器配置
RELAY_HOST=0.0.0.0
RELAY_PORT=8080

# 生产模式
GIN_MODE=release
```

## 📊 监控和日志

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs -f relay-server

# 查看最近100行日志
docker-compose logs --tail=100 relay-server
```

### 性能监控

```bash
# 查看资源使用情况
docker stats cligool-relay cligool-postgres cligool-redis

# 查看容器详情
docker inspect cligool-relay
```

### 健康检查

```bash
# API健康检查
curl http://localhost:8081/api/health

# 数据库连接测试
docker exec cligool-postgres pg_isready -U cligool

# Redis连接测试
docker exec cligool-redis redis-cli ping
```

## 🛠️ 维护操作

### 更新服务

```bash
# 1. 拉取最新代码
git pull origin main

# 2. 重新构建镜像
docker-compose build

# 3. 重启服务
docker-compose up -d

# 4. 清理旧镜像
docker image prune -f
```

### 备份数据

```bash
# 备份PostgreSQL数据库
docker exec cligool-postgres pg_dump -U cligool cligool > backup_$(date +%Y%m%d).sql

# 备份Redis数据
docker exec cligool-redis redis-cli SAVE
docker cp cligool-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

### 恢复数据

```bash
# 恢复PostgreSQL数据库
docker exec -i cligool-postgres psql -U cligool cligool < backup_20250101.sql

# 恢复Redis数据
docker cp ./redis_backup_20250101.rdb cligool-redis:/data/dump.rdb
docker-compose restart redis
```

### 清理资源

```bash
# 停止所有服务
docker-compose down

# 删除所有数据（谨慎使用）
docker-compose down -v

# 清理未使用的镜像
docker image prune -a

# 清理未使用的卷
docker volume prune
```

## 🔒 安全建议

1. **更改默认密码**: 修改docker-compose.yml中的数据库密码
2. **启用HTTPS**: 使用Cloudflare Tunnel或Let's Encrypt
3. **防火墙配置**: 只开放必要的端口
4. **定期更新**: 保持Docker和系统更新
5. **监控日志**: 定期检查异常访问
6. **备份数据**: 定期备份重要数据

## 🐛 故障排除

### 端口被占用

```bash
# 检查端口占用
sudo lsof -i :8081

# 修改docker-compose.yml中的端口映射
ports:
  - "8082:8080"  # 改为其他端口
```

### 内存不足

```bash
# 增加Docker内存限制（Docker Desktop）
# Settings > Resources > Memory

# 或减少PostgreSQL内存使用
# 在docker-compose.yml中添加:
command: postgres -c shared_buffers=128MB -c max_connections=50
```

### 容器无法启动

```bash
# 查看详细日志
docker-compose logs relay-server

# 检查容器状态
docker ps -a

# 重新创建容器
docker-compose up -d --force-recreate
```

## 📚 更多文档

- [README.md](../README.md) - 项目总览
- [USAGE_GUIDE_CN.md](../USAGE_GUIDE_CN.md) - 使用指南
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南

## 🤝 获取帮助

- 提交Issue: https://github.com/topcheer/cligool/issues
- 查看文档: https://github.com/topcheer/cligool/wiki
