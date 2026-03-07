# CliGool 一键部署指南

本指南介绍如何将CliGool relay server一键部署到各大免费云平台。

## 🚀 推荐的免费部署平台

### 1. Render (最推荐)

**优点**：
- ✅ 完全免费
- ✅ 支持WebSocket
- ✅ 自动HTTPS
- ✅ 一键部署

**缺点**：
- ❌ 冷启动（15-30秒）
- ❌ 免费套餐有休眠

**部署步骤**：

#### 1.1 准备工作
```bash
# 1. Fork项目到你的GitHub账号
# 访问 https://github.com/topcheer/cligool
# 点击右上角"Fork"按钮

# 2. 注册Render账号
# 访问 https://dashboard.render.com/register
```

#### 1.2 部署Web服务
```bash
# 1. 登录Render Dashboard
# 2. 点击"New +" -> "Web Service"
# 3. 连接你的GitHub账号
# 4. 选择fork的项目
# 5. 配置如下：
```

**配置**：
- **Name**: `cligool-relay`
- **Environment**: `Docker`
- **Dockerfile Path**: `./Dockerfile.multiarch`
- **Branch**: `main`
- **Plan**: `Free`

**环境变量**：
```
DATABASE_URL  # 稍后在数据库创建后配置
REDIS_URL     # 稍后在Redis创建后配置
RELAY_HOST    = 0.0.0.0
RELAY_PORT    = 8080
ENABLE_AUTO_HTTPS = true
```

#### 1.3 创建PostgreSQL数据库
```bash
# 1. 在Render Dashboard点击"New +"
# 2. 选择"PostgreSQL"
# 3. 配置如下：
```

**配置**：
- **Name**: `cligool-postgres`
- **Database**: `cligool`
- **User**: `cligool`
- **Plan**: `Free`
- **Region**: `Oregon (US West)`

**获取数据库URL**：
```bash
# 复制"Database Connection"的"Internal URL"
# 格式：postgresql://user:pass@host/dbname
```

#### 1.4 创建Redis缓存
```bash
# 1. 在Render Dashboard点击"New +"
# 2. 选择"Redis"
# 3. 配置如下：
```

**配置**：
- **Name**: `cligool-redis`
- **Plan**: `Free`
- **Maxmemory Policy**: `allkeys-lru`

#### 1.5 连接服务
```bash
# 1. 回到cligool-relay web服务
# 2. 点击"Environment"标签
# 3. 添加环境变量：
```

**添加数据库连接**：
```
DATABASE_URL = postgresql://cligool:password@cligool-postgres-xxx:5432/cligool
REDIS_URL = redis://cligool-redis-xxx:6379
```

#### 1.6 部署完成
```bash
# 1. 点击"Deploy Web Service"
# 2. 等待部署完成（约2-3分钟）
# 3. 部署完成后会获得一个URL，如：
#    https://cligool-relay.onrender.com
```

**验证部署**：
```bash
# 检查健康状态
curl https://cligool-relay.onrender.com/api/health
```

---

### 2. Fly.io (推荐)

**优点**：
- ✅ 免费套餐慷慨（3个VM）
- ✅ 无冷启动
- ✅ 全球部署
- ✅ 高性能

**缺点**：
- ❌ 需要安装CLI工具

**部署步骤**：

#### 2.1 安装Flyctl
```bash
# macOS/Linux
curl -L https://fly.io/install.sh | sh

# Windows
powershell -c "iwr https://fly.io/install.sh | iex"

# 验证安装
flyctl version
```

#### 2.2 注册账号
```bash
# 登录或注册
flyctl auth signup
flyctl auth login

# 这会打开浏览器进行授权
```

#### 2.3 部署应用
```bash
# 克隆你的项目
git clone https://github.com/YOUR_USERNAME/cligool.git
cd cligool

# 启动应用（自动创建配置）
flyctl launch

# 按提示操作：
# 1. 选择region（推荐: San Jose sjc）
# 2. 确认配置
# 3. 部署应用
```

#### 2.4 创建数据库
```bash
# 创建PostgreSQL
flyctl postgres create --name cligool-postgres

# 创建Redis
flyctl redis create --name cligool-redis
```

#### 2.5 连接数据库
```bash
# 附加数据库到应用
flyctl postgres attach cligool-postgres --app cligool-relay

# 获取Redis连接URL
flyctl redis status --app cligool-redis

# 设置环境变量
flyctl secrets set DATABASE_URL="postgresql://..." --app cligool-relay
flyctl secrets set REDIS_URL="redis://..." --app cligool-relay
```

#### 2.6 重新部署
```bash
# 部署应用
flyctl deploy

# 获取应用URL
flyctl info --app cligool-relay
```

**验证部署**：
```bash
# 检查健康状态
curl https://cligool-relay.fly.dev/api/health
```

---

### 3. Koyeb

**优点**：
- ✅ $5.50免费额度/月
- ✅ 支持WebSocket
- ✅ 全球CDN
- ✅ 自动HTTPS

**部署步骤**：

#### 3.1 注册账号
```bash
# 访问 https://www.koyeb.com
# 注册账号（支持GitHub登录）
```

#### 3.2 部署应用
```bash
# 1. 登录Koyeb Dashboard
# 2. 点击"Create App"
# 3. 选择"Dockerfile"
# 4. 输入你的GitHub仓库信息：
```

**配置**：
- **GitHub Repository**: `YOUR_USERNAME/cligool`
- **Dockerfile Path**: `./Dockerfile.multiarch`
- **Region**: `Washington D.C.` (或其他)
- **Instance Type**: `Nano` (免费)

#### 3.3 创建数据库
```bash
# 创建PostgreSQL
# 1. 点击"Create Service"
# 2. 选择"PostgreSQL"
# 3. 配置并部署

# 创建Redis
# 1. 点击"Create Service"
# 2. 选择"Redis"
# 3. 配置并部署
```

#### 3.4 配置环境变量
```bash
# 在应用设置中添加环境变量：
DATABASE_URL = postgresql://...
REDIS_URL = redis://...
RELAY_HOST = 0.0.0.0
RELAY_PORT = 8080
```

**验证部署**：
```bash
# 检查健康状态
curl https://YOUR_APP.koyeb.com/api/health
```

---

## 📊 平台对比

| 特性 | Render | Fly.io | Koyeb | Railway |
|------|--------|--------|-------|---------|
| 免费额度 | 完全免费 | 3个VM | $5.50/月 | $5/月 |
| 冷启动 | 15-30秒 | 无 | 无 | 无 |
| WebSocket | ✅ | ✅ | ✅ | ✅ |
| 自动HTTPS | ✅ | ✅ | ✅ | ✅ |
| 数据库 | ✅ | ✅ | ✅ | ✅ |
| Redis | ✅ | ✅ | ✅ | ✅ |
| 易用性 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## 🎯 推荐选择

### 个人使用/开发
1. **Fly.io** - 性能最好，无冷启动
2. **Render** - 最简单，完全免费

### 生产环境
1. **Railway** - 功能最全
2. **Render** - 稳定性好

## ⚠️ 注意事项

1. **免费套餐限制**：
   - 流量限制
   - 内存限制
   - CPU限制

2. **休眠策略**：
   - Render免费版会休眠
   - 首次访问需要等待启动

3. **域名配置**：
   - 所有平台都支持自定义域名
   - 需要配置DNS记录

## 🔧 通用配置

所有平台都需要配置以下环境变量：

```bash
# 数据库连接
DATABASE_URL=postgres://user:pass@host:5432/dbname

# Redis连接
REDIS_URL=redis://host:6379

# 服务器配置
RELAY_HOST=0.0.0.0
RELAY_PORT=8080
ENABLE_AUTO_HTTPS=true
```

## 📞 需要帮助？

如果遇到问题，请访问：
- [Render文档](https://render.com/docs)
- [Fly.io文档](https://fly.io/docs/)
- [Koyeb文档](https://www.koyeb.com/docs)
- [项目Issues](https://github.com/topcheer/cligool/issues)
