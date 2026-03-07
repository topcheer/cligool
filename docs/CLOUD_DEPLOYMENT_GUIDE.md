# CliGool 云平台一键部署指南

## 📋 目录

- [平台对比](#平台对比)
- [部署前准备](#部署前准备)
- [方式1：Render 部署](#方式1render-部署)
- [方式2：Koyeb 部署](#方式2koyeb-部署)
- [方式3：Railway 部署](#方式3railway-部署)
- [验证部署](#验证部署)

---

## 平台对比

| 特性 | Render | Koyeb | Railway |
|------|--------|-------|---------|
| **免费额度** | 完全免费 | $5.50/月 | $5/月 |
| **冷启动** | 15-30秒 | 无 | 无 |
| **部署难度** | ⭐ 最简单 | ⭐⭐ 中等 | ⭐⭐ 简单 |
| **性能** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **推荐指数** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## 部署前准备

### Fork 项目到你的 GitHub

1. 访问 https://github.com/topcheer/cligool
2. 点击右上角 **"Fork"** 按钮
3. 现在你有了自己的 `cligool` 仓库

### 准备账号

- **GitHub 账号**：用于 Fork 项目
- **云平台账号**：注册对应平台的账号（推荐使用 GitHub 登录）

---

## 方式1：Render 部署 ⭐ 推荐

### 特点
- ✅ **完全免费**，不需要信用卡
- ✅ **Blueprint 一键部署**，配置最简单
- ✅ **单个 Docker 容器**，无需数据库
- ⚠️ **15分钟后无流量会休眠**
- ⚠️ **休眠后首次访问需要15-30秒启动**

### 部署步骤（3分钟）

#### 步骤 1：访问 Render Blueprint

```bash
open https://dashboard.render.com/blueprints/new
```

#### 步骤 2：连接 GitHub

1. 点击 **"Connect GitHub account"** 授权
2. 选择你 Fork 的 `cligool` 仓库
3. Render 会自动检测到 `render.yaml` 文件

#### 步骤 3：确认配置

Render 会显示以下配置：

```yaml
services:
  - type: web
    name: cligool-relay
    plan: free
    env: docker
    region: singapore
    dockerfilePath: ./Dockerfile.multiarch
    healthCheckPath: /api/health
```

#### 步骤 4：部署

1. 检查配置是否正确
2. 选择区域：**Singapore**（或其他区域）
3. 点击 **"Apply Blueprint"**
4. 等待部署完成（约3-5分钟）

#### 步骤 5：获取部署 URL

部署完成后，Render 会提供一个 URL：
```
https://cligool-relay.onrender.com
```

### 配置文件说明

**`render.yaml`** - Render Blueprint 配置
```yaml
services:
  - type: web
    name: cligool-relay
    plan: free
    env: docker
    region: singapore
    dockerfilePath: ./Dockerfile.multiarch
    healthCheckPath: /api/health
    envVars:
      - key: RELAY_HOST
        value: 0.0.0.0
      - key: RELAY_PORT
        value: "8080"
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
- ✅ **无冷启动**，性能最好
- ✅ **全球 CDN**，延迟低
- ✅ **单个 Docker 容器**，无需数据库
- ⚠️ 需要信用卡验证

### 部署步骤（5分钟）

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

```
RELAY_HOST = 0.0.0.0
RELAY_PORT = 8080
```

#### 步骤 6：配置健康检查

- **路径**：`/api/health`
- **端口**：`8080`

#### 步骤 7：部署

1. 检查所有配置
2. 点击 **"Deploy"**
3. 等待部署完成（约3-5分钟）

### 配置文件说明

**`Dockerfile.koyeb`** - Koyeb 专用 Dockerfile
```dockerfile
FROM koyeb/docker-compose
COPY . /app
```

这个文件使用 Koyeb 的 Docker Compose 镜像，但 CliGool 不需要 Docker Compose 功能，直接使用 Dockerfile.multiarch 也可以。

### 常见问题

**Q: 为什么要启用 Privileged 模式？**
- Koyeb 的 `koyeb/docker-compose` 镜像需要特权模式来运行 Docker daemon

**Q: 如何查看日志？**
- 在 Koyeb 控制台点击你的服务
- 选择 **"Logs"** 标签查看实时日志

---

## 方式3：Railway 部署

### 特点
- ✅ **$5 免费额度**
- ✅ **自动环境变量注入**
- ✅ **最佳开发体验**
- ✅ **单个 Docker 容器**，无需数据库
- ⚠️ 免费额度用完后需要付费

### 部署步骤（3分钟）

#### 步骤 1：访问 Railway

```bash
open https://railway.com/new
```

#### 步骤 2：连接 GitHub

1. 点击 **"Deploy from GitHub repo"**
2. 授权 Railway 访问你的 GitHub
3. 选择你 Fork 的 `cligool` 仓库

#### 步骤 3：部署服务

1. Railway 会自动检测 `railway.toml` 配置
2. 检查配置：
   - **Builder**: DOCKERFILE
   - **Dockerfile**: Dockerfile.multiarch
   - **Healthcheck**: /api/health
3. 点击 **"Deploy"**
4. 等待部署完成（约3-5分钟）

#### 步骤 4：获取部署 URL

部署完成后，Railway 会提供一个 URL：
```
https://cligool-relay.up.railway.app
```

### 配置文件说明

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

**Q: 如何查看环境变量？**
- 点击你的服务
- 选择 **"Variables"** 标签
- 可以查看所有环境变量

**Q: 免费额度用完怎么办？**
- Railway 提供 $5 免费额度
- 用完后会自动暂停服务
- 可以升级到付费套餐或添加支付方式

---

## 验证部署

### 1. 健康检查

部署完成后，访问健康检查端点：
```bash
curl https://你的域名/api/health
```

应该返回：
```json
{"status":"ok","time":1234567890}
```

### 2. 下载并运行客户端

1. 访问你的部署 URL，例如：`https://cligool-relay.onrender.com`
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

### 3. 测试 WebSocket 连接

在浏览器中访问：
```
https://你的域名/session/test-session
```

应该看到终端界面。然后运行客户端连接到 `test-session`：
```bash
./cligool-darwin-arm64 -server https://你的域名 -session test-session
```

---

## 平台选择建议

### 个人使用/测试
**推荐：Render**
- 完全免费，不需要信用卡
- Blueprint 一键部署，配置最简单
- 适合快速测试和演示

### 生产环境
**推荐：Koyeb**
- 无冷启动，性能稳定
- 全球 CDN，用户体验好
- 适合需要稳定性能的应用

### 开发调试
**推荐：Railway**
- 开发体验最好
- 自动配置，零手动操作
- 适合快速迭代和开发

---

## 架构说明

### 简化后的架构

**CliGool Relay Server** 现在是一个**无状态**的 WebSocket 中继服务：

```
CLI客户端 ──WebSocket──> Relay Server <──WebSocket─── Web浏览器
      │                        │
      └────────本地PTY────────┘        内存会话管理
```

**特点**：
- ✅ **无需数据库**：所有会话存储在内存中
- ✅ **无需 Redis**：没有缓存依赖
- ✅ **单个容器**：部署更简单
- ✅ **自动清理**：客户端断开后自动释放资源
- ✅ **多客户端支持**：每个会话可以连接多个 Web 客户端

### 会话管理

- 会话 ID 由客户端指定（通过 `-session` 参数）
- 多个 Web 客户端可以连接到同一个会话
- CLI 客户端断开后，Web 客户端也会断开
- 服务重启后，所有会话会丢失（这是正常的）

---

## 故障排除

### 问题 1：部署失败

**症状**：构建失败，服务无法启动

**原因**：
- `Dockerfile.multiarch` 构建错误
- 端口配置错误

**解决**：
1. 查看平台日志获取详细错误信息
2. 确保 `Dockerfile.multiarch` 已推送到 GitHub
3. 检查端口是否为 `8080`

### 问题 2：无法访问 Web 界面

**症状**：访问 URL 时显示错误或超时

**原因**：
- 服务还在启动中
- 健康检查失败

**解决**：
1. 等待1-2分钟，让服务完全启动
2. 查看服务日志，检查是否有错误
3. 确认健康检查路径为 `/api/health`

### 问题 3：WebSocket 连接失败

**症状**：客户端无法连接到服务器

**原因**：
- HTTPS 配置问题
- WebSocket 路由错误

**解决**：
1. 使用 `wss://` 而不是 `ws://`
2. 检查 WebSocket URL：`wss://你的域名/api/terminal/{session_id}?type=web&user_id=web-{timestamp}`

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
