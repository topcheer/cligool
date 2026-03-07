# CliGool 云平台部署指南

本文档详细说明如何在三个免费云平台上部署 CliGool Relay Server。

## 📋 目录

- [平台对比](#平台对比)
- [部署前准备](#部署前准备)
- [方式1：Render 部署](#方式1render-部署)
- [方式2：Koyeb 部署](#方式2koyeb-部署)
- [方式3：Railway 部署](#方式3railway-部署)
- [故障排除](#故障排除)
- [验证部署](#验证部署)

---

## 平台对比

| 特性 | Render | Koyeb | Railway |
|------|--------|-------|---------|
| **免费额度** | 完全免费 | $5.50/月 | $5/月 |
| **冷启动** | 15-30秒 | 无 | 无 |
| **数据库** | ✅ 免费 | ✅ 免费 | ✅ 免费 |
| **Redis** | ✅ 免费 | ✅ 免费 | ✅ 免费 |
| **自动HTTPS** | ✅ | ✅ | ✅ |
| **部署方式** | Blueprint | Docker Compose | Config as Code |
| **推荐指数** | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## 部署前准备

### 1. Fork 项目到你的 GitHub

1. 访问 https://github.com/topcheer/cligool
2. 点击右上角 **"Fork"** 按钮
3. 现在你有了自己的 `cligool` 仓库

### 2. 准备账号

- **GitHub 账号**：用于 Fork 项目
- **云平台账号**：注册对应平台的账号（推荐使用 GitHub 登录）

---

## 方式1：Render 部署

### 特点
- ✅ **完全免费**，不需要信用卡
- ✅ **Blueprint 一键部署**，自动创建所有服务
- ⚠️ **15分钟后无流量会休眠**
- ⚠️ **休眠后首次访问需要15-30秒启动**

### 部署步骤

#### 步骤 1：访问 Render Dashboard

```bash
open https://dashboard.render.com
```

#### 步骤 2：使用 Blueprint 创建服务

1. 点击左上角 **"New +"** 按钮
2. 选择 **"Blueprint"** 选项
3. 点击 **"Connect GitHub account"** 授权
4. 选择你 Fork 的 `cligool` 仓库
5. Render 会自动检测到 `render.yaml` 文件

#### 步骤 3：确认 Blueprint 配置

Render 会显示以下配置：

```yaml
databases:
  - name: cligool-postgres
    databaseName: cligool
    user: cligool
    plan: free
    region: singapore

services:
  - type: keyvalue
    name: cligool-redis
    plan: free
    region: singapore

  - type: web
    name: cligool-relay
    plan: free
    env: docker
    region: singapore
    dockerfilePath: ./Dockerfile.multiarch
    healthCheckPath: /api/health
```

#### 步骤 4：部署

1. 检查所有配置是否正确
2. 选择区域：**Singapore**（或其他区域）
3. 点击 **"Apply Blueprint"**
4. 等待部署完成（约5-10分钟）

#### 步骤 5：获取部署 URL

部署完成后，Render 会提供一个 URL：
```
https://cligool-relay.onrender.com
```

### 文件说明

**`render.yaml`** - Render Blueprint 配置文件
```yaml
databases:
  - name: cligool-postgres
    databaseName: cligool
    user: cligool

services:
  - type: keyvalue
    name: cligool-redis

  - type: web
    name: cligool-relay
    env: docker
    dockerfilePath: ./Dockerfile.multiarch
    healthCheckPath: /api/health
```

### 常见问题

**Q: 部署失败怎么办？**
- 检查 `render.yaml` 语法是否正确
- 确保 `Dockerfile.multiarch` 已推送到 GitHub
- 查看 Render 日志获取详细错误信息

**Q: 如何避免休眠？**
- 设置外部监控，每5分钟 ping 一次
- 或者升级到付费套餐

**Q: 如何更换区域？**
- 在 `render.yaml` 中修改 `region` 字段
- 支持的区域：`oregon`, `singapore`, `frankfurt`, `ohio`

---

## 方式2：Koyeb 部署

### 特点
- ✅ **$5.50 免费额度**
- ✅ **无冷启动**，性能好
- ✅ **全球 CDN**，延迟低
- ✅ **Docker Compose 支持**
- ⚠️ 需要信用卡验证

### 部署步骤

#### 步骤 1：访问 Koyeb

```bash
open https://app.koyeb.com/apps/new
```

#### 步骤 2：连接 GitHub

1. 点击 **"GitHub"** 按钮
2. 授权 Koyeb 访问你的 GitHub 账号
3. 选择你 Fork 的 `cligool` 仓库

#### 步骤 3：配置 Docker 构建

1. **构建方式**：选择 **"Dockerfile"**
2. **Dockerfile 位置**：点击 **"Override"** 切换
3. **填写**：`Dockerfile.koyeb`
4. **启用 Privileged 模式**：打开切换开关

#### 步骤 4：配置服务

- **名称**：`cligool-relay`
- **区域**：Singapore (sin)
- **实例类型**：nano（免费）

#### 步骤 5：配置环境变量

点击 **"Add environment variable"** 添加：

```
RELAY_HOST = 0.0.0.0
RELAY_PORT = 8080
ENABLE_AUTO_HTTPS = true
```

#### 步骤 6：配置健康检查

- **路径**：`/api/health`
- **端口**：`8080`

#### 步骤 7：部署

1. 检查所有配置
2. 点击 **"Deploy"**
3. 等待部署完成（约3-5分钟）

#### 步骤 8：添加数据库（需要手动）

**PostgreSQL**：
1. 在 Koyeb 控制台点击 **"New Service"**
2. 选择 **"Database"** → **"PostgreSQL"**
3. 配置：
   - 名称：`cligool-postgres`
   - 版本：`16`
   - 区域：Singapore
4. 获取连接字符串并添加到环境变量

**Redis**：
1. 点击 **"New Service"**
2. 选择 **"Database"** → **"Redis"**
3. 配置：
   - 名称：`cligool-redis`
   - 版本：`7`
   - 区域：Singapore
4. 获取连接字符串并添加到环境变量

#### 步骤 9：更新环境变量

在 Web 服务设置中添加：
```
DATABASE_URL = <PostgreSQL连接字符串>
REDIS_URL = <Redis连接字符串>
```

### 文件说明

**`Dockerfile.koyeb`** - Koyeb 专用 Dockerfile
```dockerfile
FROM koyeb/docker-compose
COPY . /app
```

这个文件使用 Koyeb 的 Docker Compose 镜像，会自动读取 `docker-compose.yml` 并启动所有服务。

### 常见问题

**Q: 为什么使用 Dockerfile.koyeb 而不是 Dockerfile.multiarch？**
- Koyeb 的 Docker Compose 模式需要特殊的 Dockerfile
- `Dockerfile.koyeb` 使用 Koyeb 官方镜像，支持完整的 Docker Compose 功能

**Q: 如何获取数据库连接字符串？**
- 在 Koyeb 控制台点击数据库服务
- 查看 **"Connection details"** 部分
- 复制 **"Internal connection string"**

**Q: 部署后如何查看日志？**
1. 点击你的服务
2. 选择 **"Logs"** 标签
3. 可以查看实时日志和历史日志

---

## 方式3：Railway 部署

### 特点
- ✅ **$5 免费额度**
- ✅ **自动环境变量注入**
- ✅ **零配置数据库**
- ✅ **最佳开发体验**
- ⚠️ 免费额度用完后需要付费

### 部署步骤

#### 步骤 1：访问 Railway

```bash
open https://railway.com/new
```

#### 步骤 2：连接 GitHub

1. 点击 **"Deploy from GitHub repo"**
2. 授权 Railway 访问你的 GitHub
3. 选择你 Fork 的 `cligool` 仓库

#### 步骤 3：部署 Web 服务

1. Railway 会自动检测 `railway.toml` 配置
2. 检查配置：
   - **Builder**: DOCKERFILE
   - **Dockerfile**: Dockerfile.multiarch
   - **Healthcheck**: /api/health
3. 点击 **"Deploy"**
4. 等待部署完成（约3-5分钟）

#### 步骤 4：添加 PostgreSQL 数据库

1. 在 Railway 项目中点击 **"New Service"**
2. 选择 **"Database"** → **"Add PostgreSQL"**
3. Railway 会自动：
   - 创建 PostgreSQL 数据库
   - 注入 `DATABASE_URL` 环境变量到所有服务
   - 自动重新部署 Web 服务

#### 步骤 5：添加 Redis

1. 点击 **"New Service"**
2. 选择 **"Database"** → **"Add Redis"**
3. Railway 会自动：
   - 创建 Redis 实例
   - 注入 `REDIS_URL` 环境变量到所有服务
   - 自动重新部署 Web 服务

#### 步骤 6：获取部署 URL

部署完成后，Railway 会提供一个 URL：
```
https://cligool-relay.up.railway.app
```

### 文件说明

**`railway.toml`** - Railway 配置文件
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "Dockerfile.multiarch"

[deploy]
healthcheckPath = "/api/health"
healthcheckTimeout = 300
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10
```

### 常见问题

**Q: Railway 如何自动注入环境变量？**
- 当你添加数据库或 Redis 时，Railway 会自动创建 `DATABASE_URL` 和 `REDIS_URL` 环境变量
- 这些变量会自动注入到项目中的所有服务
- 无需手动配置

**Q: 如何查看环境变量？**
1. 点击你的服务
2. 选择 **"Variables"** 标签
3. 可以查看所有环境变量（包括 Railway 自动注入的）

**Q: 如何连接到数据库？**
1. 点击数据库服务
2. 选择 **"Connect"** 标签
3. 可以查看连接字符串和连接示例

**Q: 免费额度用完怎么办？**
- Railway 提供 $5 免费额度
- 用完后会自动暂停服务
- 可以升级到付费套餐或添加支付方式

---

## 故障排除

### 问题 1：部署失败

**症状**：构建失败，服务无法启动

**原因**：
- `Dockerfile.multiarch` 构建错误
- 依赖的服务（数据库、Redis）未启动
- 环境变量配置错误

**解决**：
1. 查看平台日志获取详细错误信息
2. 确保 `Dockerfile.multiarch` 已推送到 GitHub
3. 检查环境变量是否正确配置
4. 确保数据库和 Redis 服务已启动

### 问题 2：无法访问 Web 界面

**症状**：访问 URL 时显示错误或超时

**原因**：
- 服务还在启动中
- 健康检查失败
- 端口配置错误

**解决**：
1. 等待1-2分钟，让服务完全启动
2. 查看服务日志，检查是否有错误
3. 确认健康检查路径为 `/api/health`
4. 检查端口是否为 `8080`

### 问题 3：WebSocket 连接失败

**症状**：客户端无法连接到服务器

**原因**：
- HTTPS 配置问题
- WebSocket 路由错误
- 防火墙阻止

**解决**：
1. 确保 `ENABLE_AUTO_HTTPS=true`
2. 使用 `wss://` 而不是 `ws://`
3. 检查 WebSocket URL：`wss://你的域名/api/terminal/{session_id}?type=web&user_id=web-{timestamp}`

### 问题 4：数据库连接失败

**症状**：服务启动后立即崩溃，日志显示数据库连接错误

**原因**：
- `DATABASE_URL` 环境变量未设置
- 数据库服务未启动
- 网络连接问题

**解决**：
1. 确保数据库服务已启动
2. 检查 `DATABASE_URL` 环境变量是否正确
3. 在 Railway/Koyeb 中使用内部连接字符串（不是公网地址）
4. 重新部署服务

### 问题 5：Redis 连接失败

**症状**：服务无法连接到 Redis

**原因**：
- `REDIS_URL` 环境变量未设置
- Redis 服务未启动

**解决**：
1. 确保 Redis 服务已启动
2. 检查 `REDIS_URL` 环境变量
3. 使用内部连接字符串

---

## 验证部署

### 1. 健康检查

部署完成后，访问健康检查端点：
```bash
curl https://你的域名/api/health
```

应该返回：
```json
{"status":"ok"}
```

### 2. 使用测试脚本

项目提供了测试脚本：
```bash
./test-deployment.sh https://你的域名
```

这会测试：
- ✅ 健康检查
- ✅ WebSocket 端点可访问性
- ✅ HTTP 响应头
- ✅ 数据库连接
- ✅ 网络延迟

### 3. 下载并运行客户端

1. 访问你的部署 URL
2. 下载适合你平台的客户端
3. 运行客户端：
   ```bash
   # macOS ARM
   ./cligool-darwin-arm64 -server https://你的域名

   # Linux AMD64
   ./cligool-linux-amd64 -server https://你的域名

   # Windows
   cligool-windows-amd64.exe -server https://你的域名
   ```
4. 客户端会显示一个 session ID 和 Web URL
5. 在浏览器中打开 Web URL，开始使用！

---

## 平台特定提示

### Render

**优点**：
- 完全免费，不需要信用卡
- Blueprint 一键部署，配置简单
- 支持多区域部署

**注意事项**：
- 15分钟后无流量会休眠
- 休眠后首次访问需要15-30秒启动
- 免费套餐限制：512MB RAM，0.1 CPU

**最佳实践**：
- 使用 Blueprint 自动创建所有服务
- 定期访问服务避免休眠
- 监控使用量，避免超出免费额度

### Koyeb

**优点**：
- 无冷启动，性能好
- 全球 CDN，延迟低
- 支持 Docker Compose

**注意事项**：
- 需要信用卡验证
- 需要手动配置数据库和 Redis
- 免费额度：$5.50/月

**最佳实践**：
- 使用 Docker Compose 简化配置
- 配置健康检查确保服务可用
- 使用私有网络连接数据库

### Railway

**优点**：
- 自动环境变量注入
- 零配置数据库
- 最佳开发体验

**注意事项**：
- 免费额度：$5
- 额度用完后需要付费
- 相对较新的平台

**最佳实践**：
- 利用自动环境变量功能
- 使用 Railway CLI 管理项目
- 监控使用量和费用

---

## 推荐部署方案

### 个人使用/测试
**推荐：Railway**
- 开发体验最好
- 自动配置，零手动操作
- 快速部署和迭代

### 生产环境
**推荐：Koyeb**
- 无冷启动，性能稳定
- 全球 CDN，用户体验好
- 支持 Docker Compose，易于扩展

### 预算有限
**推荐：Render**
- 完全免费
- 功能完整
- 适合小型项目

---

## 相关文档

- **快速开始**: `QUICKSTART.md`
- **使用指南**: `USAGE_GUIDE.md`
- **开发指南**: `docs/DEVELOPMENT.md`
- **Windows支持**: `docs/WINDOWS_SUPPORT.md`
- **配置说明**: `docs/CONFIG.md`

---

## 获取帮助

如果遇到问题：
1. 查看平台的日志和错误信息
2. 参考本文档的故障排除部分
3. 访问 GitHub Issues：https://github.com/topcheer/cligool/issues
4. 查看平台官方文档：
   - [Render Documentation](https://render.com/docs)
   - [Koyeb Documentation](https://www.koyeb.com/docs)
   - [Railway Documentation](https://docs.railway.com)

---

**祝部署顺利！🎉**
