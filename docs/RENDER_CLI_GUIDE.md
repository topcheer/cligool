# Render CLI 使用指南

## 🛠️ Render CLI (renderctl)

Render提供了CLI工具 `renderctl`，但主要用于高级功能。**对于基本部署，Web Dashboard更简单**。

---

## 📥 安装 Render CLI

### macOS
```bash
# 使用Homebrew安装
brew install render-ci

# 验证安装
renderctl version
```

### Windows
```bash
# 使用Scoop安装
scoop bucket add render-ci
scoop install render-ci

# 验证安装
renderctl version
```

### Linux
```bash
# 下载最新版本
curl -O https://github.com/render-oss/render-cli/releases/download/v0.1.10/renderctl-linux-amd64
chmod +x renderctl-linux-amd64
sudo mv renderctl-linux-amd64 /usr/local/bin/renderctl

# 验证安装
renderctl version
```

---

## 🎯 Render CLI 主要功能

### 1. 蓝绿部署
```bash
# 预览部署
renderctl preview create

# 调整流量比例
renderctl scale --service cligool-relay --percent-production 80
```

### 2. 查看日志
```bash
# 实时查看日志
renderctl logs --service cligool-relay --tail

# 查看最近100行日志
renderctl logs --service cligool-relay --limit 100
```

### 3. 管理服务
```bash
# 列出所有服务
renderctl list

# 获取服务详情
renderctl get --service cligool-relay

# 重启服务
renderctl restart --service cligool-relay
```

### 4. 环境变量管理
```bash
# 列出环境变量
renderctl env list --service cligool-relay

# 设置环境变量
renderctl env set DATABASE_URL="..." --service cligool-relay

# 删除环境变量
renderctl env unset DATABASE_URL --service cligool-relay
```

---

## ⚠️ 重要限制

**Render CLI 不能用于初始部署！**

- ❌ 不能创建新服务
- ❌ 不能连接GitHub仓库
- ❌ 不能配置Dockerfile路径
- ✅ 只能管理已存在的服务

---

## 🎯 推荐的部署方式

### 方式1: Web Dashboard（推荐新手）

**优点**：
- ✅ 可视化界面
- ✅ 步骤清晰
- ✅ 不需要安装任何工具
- ✅ 支持所有功能

**步骤**：
1. 访问 https://dashboard.render.com
2. 点击 "New +" 创建服务
3. 配置参数
4. 部署

### 方式2: Render Blueprint（高级）

Blueprint是Render的IaC（基础设施即代码）功能。

**创建 `render.yaml` 文件**：
```yaml
services:
  # Redis服务
  - type: redis
    name: cligool-redis
    plan: free
    maxmemoryPolicy: allkeys-lru
    ipAllowList: [] # 允许所有IP访问

  # PostgreSQL数据库
  - type: pserv
    name: cligool-postgres
    plan: free
    databaseName: cligool
    user: cligool

  # Web服务
  - type: web
    name: cligool-relay
    plan: free
    env: docker
    dockerfilePath: ./Dockerfile.multiarch
    dockerContext: .
    healthCheckPath: /api/health
    envVars:
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: cligool-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: cligool-redis
          property: connectionString
      - key: RELAY_HOST
        value: 0.0.0.0
      - key: RELAY_PORT
        value: 8080
      - key: ENABLE_AUTO_HTTPS
        value: true
```

**使用Blueprint部署**：
```bash
# 1. 安装Render CLI
brew install render-ci

# 2. 登录
renderctl login

# 3. 使用Blueprint一键部署
renderctl blueprint apply ./render.yaml

# 这个命令会：
# - 自动创建所有服务
# - 自动配置环境变量
# - 自动建立服务间连接
# - 自动部署
```

---

## 🚀 实际使用建议

### 对于CliGool项目

**推荐方案：使用Blueprint**

1. 创建正确的 `render.yaml` 文件
2. 一行命令部署所有服务

让我为你创建正确的Blueprint配置：