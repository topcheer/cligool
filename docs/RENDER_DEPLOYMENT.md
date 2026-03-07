# Render 手动部署详细步骤

## 🚀 部署 CliGool 到 Render（完全免费）

### 准备工作

1. **Fork项目到你的GitHub**
   - 访问 https://github.com/topcheer/cligool
   - 点击右上角"Fork"按钮

2. **注册Render账号**
   - 访问 https://dashboard.render.com/register
   - 使用GitHub账号授权登录

---

## 步骤1: 创建PostgreSQL数据库

1. 登录Render Dashboard后，点击左上角 **"New +"** 按钮
2. 选择 **"PostgreSQL"**
3. 配置数据库：
   ```
   Name: cligool-postgres
   Database: cligool
   User: cligool
   Region: Oregon (US West)  # 或选择离你最近的
   Plan: Free
   ```
4. 点击 **"Create Database"**
5. 等待创建完成（约1-2分钟）
6. 创建后，会显示 **"Internal Database URL"**，格式类似：
   ```
   postgresql://cligool:xxxxxx@cligool-postgres-xxx:5432/cligool
   ```
7. **复制这个URL保存起来**，稍后需要用到

---

## 步骤2: 创建Redis

1. 再次点击 **"New +"**
2. 选择 **"Redis"**
3. 配置Redis：
   ```
   Name: cligool-redis
   Region: Oregon (US West)  # 与数据库相同区域
   Plan: Free
   Maxmemory Policy: allkeys-lru
   ```
4. 点击 **"Create Redis"**
5. 等待创建完成
6. 会显示 **"Connection Info"**，找到 **"Internal URL"**：
   ```
   redis://cligool-redis-xxx:6379
   ```
7. **复制这个URL保存起来**

---

## 步骤3: 部署Web服务（Relay Server）

1. 点击 **"New +"**
2. 选择 **"Web Service"**
3. 会看到你的GitHub仓库列表，选择 **"cligool"**（你fork的仓库）
4. 配置Web服务：

   **Basic**:
   ```
   Name: cligool-relay
   Environment: Docker
   Region: Oregon (US West)
   Plan: Free
   ```

   **Dockerfile**:
   ```
   Dockerfile Path: ./Dockerfile.multiarch
   Docker Context: .
   ```

   **Advanced**:
   ```
   Command: 空着不填（使用Dockerfile默认命令）
   Working Directory: 空着不填
   ```

5. 点击 **"Advanced"** 展开更多选项

6. 配置环境变量（在"Environment Variables"部分）：

   点击 **"Add Environment Variable"**，逐个添加：

   ```
   Variable: RELAY_HOST
   Value: 0.0.0.0

   Variable: RELAY_PORT
   Value: 8080

   Variable: ENABLE_AUTO_HTTPS
   Value: true

   Variable: DATABASE_URL
   Value: [粘贴步骤1保存的PostgreSQL URL]

   Variable: REDIS_URL
   Value: [粘贴步骤2保存的Redis URL]
   ```

7. 检查所有配置，确认无误
8. 点击 **"Create Web Service"**
9. 等待部署完成（约3-5分钟）

---

## 步骤4: 验证部署

1. 部署完成后，会看到一个URL，类似：
   ```
   https://cligool-relay.onrender.com
   ```

2. 在浏览器中访问：
   ```
   https://cligool-relay.onrender.com/api/health
   ```

   应该看到：
   ```json
   {"status":"ok"}
   ```

3. 测试WebSocket连接：
   ```
   wss://cligool-relay.onrender.com/api/terminal/test-session?type=web&user_id=test
   ```

---

## 步骤5: 下载并运行客户端

1. 访问你的部署URL：
   ```
   https://cligool-relay.onrender.com
   ```

2. 点击下载适合你平台的客户端

3. 运行客户端（以macOS ARM为例）：
   ```bash
   # 解压下载的文件
   tar -xzf cligool-darwin-arm64.tar.gz

   # 运行客户端
   ./cligool-darwin-arm64 -server https://cligool-relay.onrender.com
   ```

4. 客户端会显示一个session ID，例如：
   ```
   Session ID: abc-123-def
   Web URL: https://cligool-relay.onrender.com/terminal/abc-123-def
   ```

5. 在浏览器中打开Web URL，开始使用！

---

## 🛠️ 故障排除

### 问题1: 部署失败

**原因**：数据库或Redis连接失败

**解决**：
1. 检查DATABASE_URL和REDIS_URL是否正确
2. 确保数据库和Redis已经创建成功
3. 确保它们都在同一个Region

### 问题2: 无法访问Web界面

**原因**：服务还在启动中

**解决**：
1. 等待1-2分钟，让服务完全启动
2. 查看Render日志：点击服务 -> "Logs"标签
3. 检查是否有错误信息

### 问题3: WebSocket连接失败

**原因**：Render的免费套餐需要启用HTTPS

**解决**：
1. 确保`ENABLE_AUTO_HTTPS=true`
2. Render会自动配置SSL证书
3. 使用`wss://`而不是`ws://`

---

## 📊 Render免费套餐限制

- ✅ 750小时/月运行时间
- ✅ 512MB RAM
- ✅ 0.1 CPU
- ⚠️  15分钟后无流量会休眠
- ⚠️ 休眠后首次访问需要15-30秒启动

---

## 💡 提示

1. **避免休眠**：
   - 可以设置外部监控，每5分钟ping一次
   - 或者升级到付费套餐

2. **自定义域名**：
   - 在服务设置中添加自定义域名
   - Render会自动配置SSL证书

3. **查看日志**：
   - 随时可以在Render Dashboard查看日志
   - 帮助调试问题

---

## ✅ 部署完成后

你的relay server现在运行在：
```
https://cligool-relay.onrender.com
```

可以：
- ✅ 从任何地方访问
- ✅ 使用多个Web客户端连接
- ✅ 完全免费运行
- ✅ 自动HTTPS加密

有问题？查看日志或访问：
https://github.com/topcheer/cligool/issues
