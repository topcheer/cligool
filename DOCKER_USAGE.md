# Docker 部署说明

## 📦 两种部署方式

### 方式 1：生产环境部署（推荐）

**使用预构建的镜像**，快速启动服务。

```bash
# 使用预构建的镜像（从 GitHub Container Registry）
docker-compose up -d
```

**特点**：
- ✅ 快速启动，无需等待构建
- ✅ 使用最新的稳定版本
- ✅ 适合生产环境和快速测试

**⚠️ 注意**：预构建镜像还在 CI/CD 流程中，暂时请使用方式 2（开发环境）

**配置文件**：`docker-compose.yml`

---

### 方式 2：本地开发环境

**自动构建镜像**，使用本地最新代码。

```bash
# 本地开发环境（自动构建）
docker-compose -f docker-compose.dev.yml up -d --build
```

**特点**：
- ✅ 使用本地最新代码
- ✅ 自动构建镜像
- ✅ 适合开发和调试

**配置文件**：`docker-compose.dev.yml`

**可选：挂载源代码**（热重载）

如果需要实时修改代码并测试，取消注释 `docker-compose.dev.yml` 中的卷挂载：

```yaml
volumes:
  - ./cmd:/app/cmd
  - ./internal:/app/internal
  - ./web:/app/web
```

然后：
```bash
docker-compose -f docker-compose.dev.yml up --build
```

---

## 🔧 常用命令

### 生产环境

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看状态
docker-compose ps
```

### 开发环境

```bash
# 构建并启动
docker-compose -f docker-compose.dev.yml up -d --build

# 查看日志
docker-compose -f docker-compose.dev.yml logs -f

# 停止服务
docker-compose -f docker-compose.dev.yml down

# 重新构建
docker-compose -f docker-compose.dev.yml build --no-cache
```

---

## 📝 配置说明

### 环境变量

两个配置文件都使用相同的环境变量：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `RELAY_HOST` | 监听地址 | 0.0.0.0 |
| `RELAY_PORT` | 监听端口 | 8080 |

### 端口映射

- **主机端口**：8081
- **容器端口**：8080

访问地址：`http://localhost:8081`

---

## 🚀 快速开始

### 生产环境（最简单）

```bash
# 1. 克隆项目
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 启动服务
docker-compose up -d

# 3. 访问服务
open http://localhost:8081
```

### 开发环境

```bash
# 1. 克隆项目
git clone https://github.com/topcheer/cligool.git
cd cligool

# 2. 构建并启动
docker-compose -f docker-compose.dev.yml up -d --build

# 3. 查看日志
docker-compose -f docker-compose.dev.yml logs -f

# 4. 访问服务
open http://localhost:8081
```

---

## 🔍 故障排除

### 问题 1：端口被占用

```bash
# 检查端口占用
lsof -i :8081

# 修改端口
# 编辑 docker-compose.yml 或 docker-compose.dev.yml
# 将 "8081:8080" 改为 "8082:8080"
```

### 问题 2：镜像拉取失败

```bash
# 查看可用的镜像
docker images | grep cligool

# 手动拉取镜像
docker pull ghcr.io/topcheer/cligool:latest

# 或者使用开发环境自动构建
docker-compose -f docker-compose.dev.yml up -d --build
```

### 问题 3：容器启动失败

```bash
# 查看详细日志
docker-compose logs -f relay-server

# 或
docker-compose -f docker-compose.dev.yml logs -f relay-server

# 检查容器状态
docker-compose ps

# 进入容器调试
docker-compose exec relay-server sh
```

---

## 📚 相关文档

- **云平台部署指南**: `docs/CLOUD_DEPLOYMENT_GUIDE.md`
- **配置说明**: `docs/CONFIG.md`
- **使用指南**: `USAGE_GUIDE.md`

---

**选择合适的方式开始使用吧！🎉**
