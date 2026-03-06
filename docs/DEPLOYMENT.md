# CliGool 部署指南

## 🚀 一键部署

### 在任何Docker服务器上部署

```bash
# 1. 克隆项目
git clone https://github.com/cligool/cligool.git
cd cligool

# 2. 运行部署脚本
./scripts/deploy.sh

# 3. 配置Cloudflare Tunnel
./scripts/cloudflare-tunnel.sh
```

就这样！你的服务现在可以通过Cloudflare Tunnel访问了。

## 📋 详细步骤

### 第一步: 准备服务器

#### 最低要求
- **操作系统**: Linux (任何发行版)、macOS、Windows
- **内存**: 512MB+
- **磁盘**: 10GB+
- **网络**: 能够访问Cloudflare

#### 安装Docker
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# CentOS/RHEL
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker

# macOS
brew install docker
```

#### 安装Docker Compose
```bash
# Linux
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# macOS (包含在Docker Desktop中)
# 无需额外安装
```

### 第二步: 部署应用

```bash
# 1. 下载项目
git clone https://github.com/cligool/cligool.git
cd cligool

# 2. 运行部署脚本
./scripts/deploy.sh

# 3. 检查服务状态
docker-compose ps
```

### 第三步: 配置Cloudflare Tunnel

#### 安装cloudflared
```bash
# Linux (amd64)
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb

# Linux (arm64)
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64.deb
sudo dpkg -i cloudflared-linux-arm64.deb

# macOS
brew install cloudflare/cloudflare/cloudflared

# Docker方式
docker pull cloudflare/cloudflared
```

#### 配置Tunnel
```bash
# 方式1: 使用自动化脚本（推荐）
./scripts/cloudflare-tunnel.sh

# 方式2: 手动配置
cloudflared tunnel login
cloudflared tunnel create cligool
cloudflared tunnel route dns <tunnel-id> your-domain.com
cloudflared tunnel run <tunnel-id>
```

#### 后台运行Tunnel
```bash
# 使用systemd（推荐）
sudo cloudflared service install

# 使用screen
screen -dmS cloudflared cloudflared tunnel run <tunnel-id>

# 使用nohup
nohup cloudflared tunnel run <tunnel-id> > /dev/null 2>&1 &
```

## 🌐 不同部署场景

### 场景1: 家庭服务器

```bash
# 在家里的电脑/树莓派上
cd cligool
./scripts/deploy.sh
./scripts/cloudflare-tunnel.sh

# 现在可以从任何地方访问！
```

### 场景2: 云服务器 (VPS)

```bash
# 连接到你的VPS
ssh user@your-vps-ip

# 部署应用
git clone https://github.com/cligool/cligool.git
cd cligool
./scripts/deploy.sh

# 配置Cloudflare Tunnel
./scripts/cloudflare-tunnel.sh
```

### 场景3: Kubernetes集群

创建 `k8s-deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cligool
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cligool
  template:
    metadata:
      labels:
        app: cligool
    spec:
      containers:
      - name: relay-server
        image: cligool:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          value: "postgres://..."
        - name: REDIS_URL
          value: "redis://..."
---
apiVersion: v1
kind: Service
metadata:
  name: cligool-service
spec:
  selector:
    app: cligool
  ports:
  - port: 80
    targetPort: 8080
```

部署到Kubernetes:
```bash
kubectl apply -f k8s-deployment.yaml
```

### 场景4: Docker Swarm

```bash
# 初始化Swarm
docker swarm init

# 部署stack
docker stack deploy -c docker-compose.swarm.yml cligool
```

## 🔧 配置选项

### 环境变量配置

编辑 `.env` 文件:
```bash
# 数据库配置
DATABASE_URL=postgres://cligool:your-password@postgres:5432/cligool
REDIS_URL=redis://redis:6379

# 服务配置
RELAY_HOST=0.0.0.0
RELAY_PORT=8080

# JWT密钥 (必须修改！)
JWT_SECRET=your-super-secret-key-here

# 日志级别
LOG_LEVEL=info
```

### 性能调优

```bash
# docker-compose.yml
services:
  relay-server:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 安全加固

1. **使用防火墙**
```bash
# 只允许本地访问
ufw deny 8080
```

2. **配置Cloudflare Access**
```bash
# 在Cloudflare Zero Trust控制台
# Settings -> Zero Trust -> Access
# 配置认证策略
```

3. **定期更新**
```bash
# 更新镜像
docker-compose pull
docker-compose up -d
```

## 📊 监控和维护

### 健康检查
```bash
# 检查服务状态
curl http://localhost:8080/api/health

# 检查Docker容器
docker-compose ps
```

### 日志管理
```bash
# 查看实时日志
docker-compose logs -f

# 查看最近日志
docker-compose logs --tail=100

# 导出日志
docker-compose logs > app.log
```

### 数据备份
```bash
# 备份数据库
docker exec cligool-postgres pg_dump -U cligool cligool > backup.sql

# 备份到云存储
rclone copy backup.sql your-backup:cligool/
```

### 服务更新
```bash
# 拉取最新代码
git pull

# 重新构建和部署
docker-compose down
docker-compose build
docker-compose up -d
```

## 🐛 故障排除

### 常见问题

1. **服务无法启动**
```bash
# 检查端口占用
netstat -tuln | grep 8080

# 检查日志
docker-compose logs relay-server

# 重启服务
docker-compose restart
```

2. **Cloudflare Tunnel连接失败**
```bash
# 检查tunnel状态
cloudflared tunnel info <tunnel-id>

# 测试连接
curl -v http://localhost:8080/api/health

# 重新配置tunnel
cloudflared tunnel delete <tunnel-id>
./scripts/cloudflare-tunnel.sh
```

3. **数据库连接问题**
```bash
# 检查数据库状态
docker-compose ps postgres

# 进入数据库
docker exec -it cligool-postgres psql -U cligool cligool
```

### 日志位置
- **应用日志**: `docker-compose logs`
- **数据库日志**: `docker-compose logs postgres`
- **Cloudflare日志**: `/var/log/cloudflared.log`

## 🚀 生产环境建议

### 高可用部署
1. 使用多个relay实例
2. 配置负载均衡
3. 设置健康检查
4. 配置自动重启

### 安全建议
1. 启用Cloudflare Access
2. 配置速率限制
3. 定期备份数据
4. 监控异常访问
5. 使用强密码

### 性能优化
1. 调整数据库连接池
2. 配置Redis持久化
3. 使用CDN缓存静态资源
4. 监控资源使用
5. 定期清理日志

## 📞 获取帮助

- **文档**: [docs/](docs/)
- **问题报告**: [GitHub Issues](https://github.com/cligool/cligool/issues)
- **讨论**: [GitHub Discussions](https://github.com/cligool/cligool/discussions)

---

*部署成功后，你可以通过Cloudflare配置的域名访问服务！*