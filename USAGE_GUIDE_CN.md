# CliGool - 完整的远程终端解决方案

## 🎉 架构修复完成！

经过重大修复，现在 CliGool 采用了正确的架构设计，完全解决了之前的 JSON 格式不匹配问题。

## 🏗️ 最终架构设计

```
┌──────────────┐              ┌──────────┐              ┌─────────────┐
│ 用户A的电脑   │──WebSocket──▶│ 中继服务器│◀──WebSocket───│  用户B的浏览器│
│ (真实PTY)    │              │ 消息转发器│              │  (本地HTML)   │
│ CLI客户端    │              │          │              │  Web界面     │
└──────────────┘              └──────────┘              └─────────────┘
```

**关键特点**：
- ✅ **独立组件**：CLI客户端、中继服务器、Web界面完全分离
- ✅ **Base64编码**：解决了消息格式兼容性问题
- ✅ **真实PTY**：只有CLI客户端提供真实的终端环境
- ✅ **本地Web界面**：HTML文件可以在任何地方打开使用

## 🚀 快速开始

### 第一步：启动中继服务器

```bash
# 启动所有服务（包括数据库和缓存）
docker-compose up -d

# 检查服务状态
docker-compose ps
```

**预期输出**：
```
NAME                STATUS
cligool-relay       Up (healthy)
cligool-postgres    Up
cligool-redis       Up
```

### 第二步：启动CLI客户端

**在有真实终端的机器上**（如你的Mac/PC）：

```bash
# 连接到本地中继服务器
./bin/cligool-simple -server http://localhost:8081 -connect-only

# 或连接到远程服务器
./bin/cligool-simple -server https://cligool.zty8.cn -connect-only
```

**输出示例**：
```
🚀 连接到中继服务器: http://localhost:8081
📋 会话ID: abc123-def456-7890-abcd-ef1234567890
🌐 Web访问地址: http://localhost:8081/?session=abc123-def456-7890-abcd-ef1234567890
✅ WebSocket连接成功！
📡 连接模式：仅保持WebSocket连接，不启动本地shell
💡 现在你可以在Web界面中使用这个会话：
   http://localhost:8081/?session=abc123-def456-7890-abcd-ef1234567890
```

### 第三步：打开Web界面

```bash
# 直接打开项目中的web-client.html文件
open web-client.html

# 或者复制到其他位置使用
cp web-client.html ~/Desktop/cligool.html
open ~/Desktop/cligool.html
```

### 第四步：连接Web界面

1. 在浏览器中打开 `web-client.html`
2. 输入中继服务器地址：`ws://localhost:8081`（或远程地址）
3. 输入会话ID：复制第二步中的会话ID
4. 点击"连接"按钮
5. 开始远程控制终端！

## 🔧 技术细节

### 消息格式修复

之前的问题：
- JavaScript `Uint8Array` → JSON → Go `[]byte` 导致格式不匹配

现在的解决方案：
- JavaScript字符串 → Base64编码 → JSON → Go字符串
- 所有组件使用统一的Base64字符串格式

**示例消息**：
```json
{
  "type": "input",
  "data": "SGVsbG8gZnJvbSB3ZWIgY2xpZW50IQ==",
  "session": "abc123-def456-7890-abcd-ef1234567890"
}
```

## 📋 文件说明

### 核心组件

- **`cmd/relay/`** - 中继服务器代码
  - `main.go` - 服务器主程序
  - 使用Gin框架处理WebSocket连接
  - 支持多用户会话管理

- **`cmd/client/simple.go`** - 简化版CLI客户端
  - 只建立WebSocket连接，不启动本地shell
  - 自动生成UUID会话ID
  - 支持心跳保持连接

- **`web-client.html`** - 独立Web界面
  - 单一HTML文件，包含所有CSS和JavaScript
  - 使用xterm.js提供终端仿真
  - 可以在任何地方本地打开

### 配置文件

- **`Dockerfile`** - 中继服务器容器镜像
  - 多阶段构建，最终镜像基于Alpine Linux
  - 包含健康检查机制
  - 暴露8080端口

- **`docker-compose.yml`** - 服务编排
  - 定义中继服务器、PostgreSQL、Redis服务
  - 配置网络和数据卷
  - 设置环境变量

- **`cloudflare-tunnel.yml.example`** - Cloudflare Tunnel配置示例

## 🌟 使用场景

### 场景1：远程访问你的Mac

```bash
# 在家里的Mac上
./bin/cligool-simple -server http://localhost:8081 -connect-only

# 在办公室电脑的浏览器中打开web-client.html
# 输入会话ID，连接！
```

### 场景2：技术支持

```bash
# 朋友的电脑出现问题
# 让朋友运行: ./bin/cligool-simple -connect-only
# 你得到会话ID后，在浏览器中连接
# 远程查看和操作朋友的终端
```

### 场景3：团队协作

```bash
# 服务器上运行: ./bin/cligool-simple -connect-only
# 团队成员各自在浏览器中打开web-client.html
# 输入同一个会话ID，多人同时查看
```

## 💡 关键优势

- ✅ **真正安全** - Web界面是本地文件，可托管在任何地方
- ✅ **灵活部署** - 中继服务器可以在任何Docker机器上运行
- ✅ **协作友好** - 多人可同时查看同一会话
- ✅ **无复杂权限** - Web端不需要PTY权限
- ✅ **格式兼容** - Base64编码确保消息正确传递

## 🔍 故障排除

### 连接问题

1. **检查服务器状态**
   ```bash
   docker-compose ps
   docker logs cligool-relay
   ```

2. **验证WebSocket连接**
   ```bash
   curl -I http://localhost:8081/api/health
   ```

3. **查看客户端日志**
   ```bash
   tail -f /tmp/cligool-client.log
   ```

### 常见错误

**"WebSocket连接失败"**
- 检查服务器地址是否正确
- 确认中继服务器正在运行
- 验证端口没有被占用

**"会话不存在"**
- 确保CLI客户端正在运行
- 检查会话ID是否正确复制
- 查看中继服务器日志

## 🎯 验证成功标准

当系统正常工作时，你应该能够：

- ✅ 看到CLI客户端成功连接并显示会话ID
- ✅ 在Web界面中成功连接到同一会话
- ✅ 在两个终端中看到实时消息同步
- ✅ 在Web界面中输入命令并看到响应
- ✅ 多个Web客户端同时连接同一会话

## 🚀 部署到生产环境

### 使用Cloudflare Tunnel

1. **安装Cloudflare Tunnel**
   ```bash
   # 按照Cloudflare官方文档安装cloudflared
   ```

2. **配置隧道**
   ```bash
   # 复制配置文件
   cp cloudflare-tunnel.yml.example ~/.cloudflared/config.yml

   # 修改配置中的域名和服务地址
   ```

3. **启动隧道**
   ```bash
   cloudflared tunnel run
   ```

### Docker部署

```bash
# 构建镜像
docker build -t cligool-relay-server .

# 运行容器
docker run -d -p 8080:8080 --name cligool-relay cligool-relay-server
```

## 🎉 总结

现在 CliGool 已经完全修复并可以正常使用了！

**修复内容**：
- ✅ 修复了JSON消息格式不匹配问题
- ✅ 采用Base64编码确保数据正确传递
- ✅ 分离了Web界面和CLI客户端
- ✅ 实现了真正的消息中继架构

**下一步**：
1. 启动中继服务器
2. 运行CLI客户端
3. 打开Web界面
4. 开始使用远程终端功能！

祝你使用愉快！ 🚀