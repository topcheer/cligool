# CliGool - 正确的远程终端架构

## 🏗️ 正确的系统架构

```
用户A的电脑(有真实终端)          中继服务器               用户B的浏览器
┌──────────────┐              ┌──────────┐              ┌─────────────┐
│ CLI客户端     │──WebSocket──▶│ 消息转发器│◀──WebSocket───│  独立HTML     │
│ (真实PTY)     │              │          │              │  (本地文件)   │
└──────────────┘              └──────────┘              └─────────────┘
```

## 🚀 正确的使用流程

### 第一步：启动CLI客户端
```bash
# 在有真实终端的机器上（如你的Mac/PC）
./bin/cligool-simple -connect-only
```

输出示例：
```
🚀 连接到中继服务器: https://cligool.zty8.cn
📋 会话ID: abc123-def456-7890-abcd-ef1234567890
🌐 Web访问地址: https://cligool.zty8.cn/?session=abc123-def456-7890-abcd-ef1234567890
✅ WebSocket连接成功！
```

### 第二步：打开Web界面
```bash
# 下载独立的Web界面文件到本地
curl -o ~/Desktop/cligool.html https://cligool.zty8.cn/web-client.html

# 或者直接复制以下内容保存为HTML文件...
```

或者将 `web-client.html` 文件传给你的用户，让他们在浏览器中打开。

### 第三步：连接Web界面
1. 在浏览器中打开 `web-client.html`
2. 输入会话ID：`abc123-def456-7890-abcd-ef1234567890`
3. 点击"连接"按钮
4. 开始远程控制终端！

## 📋 文件说明

- **web-client.html** - 独立的Web界面，用户本地打开
- **cligool-simple** - 简化版客户端，只建立连接
- **cligool-client** - 完整客户端（需要PTY权限）

## 🌟 完整场景示例

### 场景1：远程访问你的Mac
```bash
# 在你的Mac上
./bin/cligool-simple -connect-only

# 得到会话ID后，在办公室电脑的浏览器中打开web-client.html
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

## 🔧 Web界面使用方式

**方式1：托管版本** (暂时不可用，需要重新设计)
```bash
# 访问 https://cligool.zty8.cn (当前版本不正确)
```

**方式2：本地文件** (推荐)
```bash
# 下载独立HTML文件
curl -o cligool.html https://cligool.zty8.cn/web-client.html

# 在浏览器中打开
open cligool.html
```

## 💡 关键点

1. **CLI客户端必须先启动** - 提供真实的PTY
2. **Web界面后连接** - 连接到已建立的会话
3. **中继服务器只转发** - 不启动任何终端

## 🎯 优势

- ✅ **真正安全** - Web界面是本地文件，无中间人
- ✅ **灵活部署** - Web界面可以托管在任何地方
- ✅ **协作友好** - 多人可同时查看同一会话
- ✅ **无需复杂权限** - Web端不需要PTY权限

## ⚠️ 当前状态

- ✅ **中继服务器**: 正常运行，支持WebSocket转发
- ✅ **CLI客户端**: 可以建立连接，部分环境下PTY有限制
- ✅ **Web客户端**: 独立HTML文件，可本地打开使用
- ⚠️ **完整实现**: 需要Web界面能正确显示和连接

你现在可以试用正确的架构了！
