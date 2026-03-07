# CliGool Docker 部署指南

本文档介绍如何使用Docker部署CliGool中继服务器。

## 🐳 镜像信息

### 镜像仓库

- **GitHub Container Registry**: `ghcr.io/topcheer/cligool`
- **支持架构**: linux/amd64, linux/arm64
- **包含内容**: 中继服务器 + 所有33个平台的客户端二进制文件

### 可用标签

- `latest` - 最新稳定版本
- `v1.1.0`, `v1.2.0`, 等 - 特定版本
- `v1` - 主版本latest
- `v1.1` - 次版本latest

## 🚀 快速开始

### 方法1：使用Docker Compose（推荐）

```bash
# 克隆仓库
git clone https://github.com/topcheer/cligool.git
cd cligool

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f relay-server

# 访问Web界面
open http://localhost:8081
```

### 方法2：使用Docker命令

```bash
# 拉取镜像
docker pull ghcr.io/topcheer/cligool:latest

# 运行容器（需要先启动PostgreSQL和Redis）
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://cligool:cligool123@postgres:5432/cligool?sslmode=disable \
  -e REDIS_URL=redis://redis:6379 \
  ghcr.io/topcheer/cligool:latest
```

### 方法3：完整环境（包含数据库）

```bash
# 创建网络
docker network create cligool-network

# 启动PostgreSQL
docker run -d \
  --name cligool-postgres \
  --network cligool-network \
  -e POSTGRES_DB=cligool \
  -e POSTGRES_USER=cligool \
  -e POSTGRES_PASSWORD=cligool123 \
  postgres:15-alpine

# 启动Redis
docker run -d \
  --name cligool-redis \
  --network cligool-network \
  redis:7-alpine

# 启动CliGool中继服务器
docker run -d \
  --name cligool-relay \
  --network cligool-network \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://cligool:cligool123@cligool-postgres:5432/cligool?sslmode=disable \
  -e REDIS_URL=redis://cligool-redis:6379 \
  ghcr.io/topcheer/cligool:latest
```

## ⚙️ 环境变量

### 数据库配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DATABASE_URL` | PostgreSQL连接字符串 | - |
| `REDIS_URL` | Redis连接字符串 | - |

### 服务器配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `RELAY_HOST` | 监听地址 | 0.0.0.0 |
| `RELAY_PORT` | 监听端口 | 8080 |
| `ENABLE_AUTO_HTTPS` | 自动启用HTTPS | false |

### 示例

```bash
docker run -d \
  --name cligool-relay \
  -p 8081:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable \
  -e REDIS_URL=redis://host:6379 \
  -e RELAY_HOST=0.0.0.0 \
  -e RELAY_PORT=8080 \
  ghcr.io/topcheer/cligool:latest
```

## 📦 多架构支持

Docker镜像支持以下架构：

- **linux/amd64** - x86_64服务器（Intel/AMD）
- **linux/arm64** - ARM 64位服务器（AWS Graviton、Apple Silicon等）

Docker会自动拉取适合您系统架构的镜像：

```bash
# 自动选择正确的架构
docker pull ghcr.io/topcheer/cligool:latest

# 查看镜像架构
docker inspect ghcr.io/topcheer/cligool:latest | grep Architecture
```

## 🔧 构建自定义镜像

### 本地构建（单架构）

```bash
# 使用标准Dockerfile
docker build -t cligool:local .

# 或使用多架构Dockerfile
docker build -f Dockerfile.multiarch -t cligool:local .
```

### 本地构建（多架构）

```bash
# 启用buildx
docker buildx create --use

# 构建多架构镜像
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t cligool:multiarch \
  -f Dockerfile.multiarch \
  --push \
  .
```

## 📊 容器管理

### 查看日志

```bash
# 查看所有日志
docker logs cligool-relay

# 实时查看日志
docker logs -f cligool-relay

# 查看最近100行日志
docker logs --tail 100 cligool-relay
```

### 进入容器

```bash
# 使用shell进入容器
docker exec -it cligool-relay sh

# 查看服务器状态
docker exec cligool-relay wget -qO- http://localhost:8080/api/health
```

### 重启容器

```bash
# 重启容器
docker restart cligool-relay

# 停止容器
docker stop cligool-relay

# 启动容器
docker start cligool-relay
```

### 清理

```bash
# 停止并删除容器
docker stop cligool-relay && docker rm cligool-relay

# 删除镜像
docker rmi ghcr.io/topcheer/cligool:latest

# 清理所有相关资源（包括数据卷）
docker-compose down -v
```

## 🌐 反向代理配置

### Nginx示例

```nginx
server {
    listen 80;
    server_name cligool.example.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket支持
        proxy_read_timeout 86400;
    }
}
```

### Caddy示例

```
cligool.example.com {
    reverse_proxy localhost:8081
}
```

## 🔒 安全建议

1. **不要在生产环境使用默认密码**
2. **使用HTTPS**：配置反向代理并启用SSL
3. **限制网络访问**：使用防火墙限制数据库访问
4. **定期更新镜像**：`docker pull ghcr.io/topcheer/cligool:latest`
5. **备份PostgreSQL数据**：定期备份PostgreSQL数据卷

## 🐛 故障排除

### 容器无法启动

```bash
# 检查日志
docker logs cligool-relay

# 检查数据库连接
docker exec cligool-relay ping -c 3 cligool-postgres
```

### 无法访问Web界面

```bash
# 检查端口映射
docker ps | grep cligool-relay

# 检查防火墙
sudo ufw status
```

### 数据库连接失败

```bash
# 检查数据库容器
docker ps | grep postgres

# 测试数据库连接
docker exec cligool-relay sh -c 'apk add pgadmin && psql $DATABASE_URL'
```

## 📚 更多信息

- [项目主页](https://github.com/topcheer/cligool)
- [快速开始指南](https://github.com/topcheer/cligool/blob/main/QUICKSTART.md)
- [使用指南](https://github.com/topcheer/cligool/blob/main/USAGE_GUIDE_CN.md)
- [配置文件指南](https://github.com/topcheer/cligool/blob/main/CONFIG_GUIDE.md)
