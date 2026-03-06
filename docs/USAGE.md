# CliGool 使用指南

## 📖 简介

CliGool 是一个可以让用户通过Web界面远程操作终端的CLI wrapper应用，支持实时终端访问、完整终端特性和多用户协作功能。

## 🎯 主要功能

- ✅ **实时终端访问** - 通过Web界面实时控制远程终端
- ✅ **完整终端特性** - 支持颜色、光标控制、屏幕操作等
- ✅ **多用户协作** - 多个用户可以同时访问同一终端会话
- ✅ **安全加密** - WebSocket通信采用TLS加密
- ✅ **会话管理** - 支持创建、删除、列出终端会话
- ✅ **跨平台** - 支持 Linux、macOS、Windows

## 🚀 快速开始

### 1. 启动中继服务器

```bash
# 使用Docker启动
docker-compose up -d

# 或者直接运行
make run-relay
```

### 2. 启动CLI客户端

```bash
# 连接到中继服务器
./bin/cligool-client -server https://your-relay-server.com -session my-session

# 或者使用环境变量
export CLIGOOL_SERVER=https://your-relay-server.com
./bin/cligool-client -session my-session
```

### 3. 访问Web界面

在浏览器中打开：
```
https://your-relay-server.com/?session=my-session
```

## 💡 使用场景

### 场景1: 个人远程访问

**需求**: 在外出时需要访问家里的电脑终端

**操作步骤**:
1. 在家里电脑上启动CLI客户端
   ```bash
   ./cligool-client -server https://relay.example.com -session home-pc
   ```

2. 在任何地方打开浏览器访问
   ```
   https://relay.example.com/?session=home-pc
   ```

### 场景2: 团队协作调试

**需求**: 多个团队成员需要同时查看和调试服务器问题

**操作步骤**:
1. 在服务器上启动CLI客户端
   ```bash
   ./cligool-client -server https://relay.example.com -session debug-session
   ```

2. 团队成员分别访问同一个会话
   ```
   https://relay.example.com/?session=debug-session
   ```

3. 协作进行调试（注意：同一时间只能有一个用户输入）

### 场景3: 远程教育演示

**需求**: 老师需要向学生演示命令行操作

**操作步骤**:
1. 老师启动CLI客户端
   ```bash
   ./cligool-client -server https://relay.example.com -session demo-session
   ```

2. 学生访问同一个会话观看演示
   ```
   https://relay.example.com/?session=demo-session
   ```

## 🔧 命令行参数

### CLI客户端参数

```bash
Usage: cligool-client [options]

Options:
  -server string
        中继服务器URL (环境变量: CLIGOOL_SERVER)
  -session string
        会话ID (可选，自动生成如果未提供)
  -shell string
        Shell路径 (默认使用系统默认)
  -debug
        启用调试模式
```

### 使用示例

```bash
# 基础使用
./cligool-client -server https://relay.example.com

# 指定会话ID
./cligool-client -server https://relay.example.com -session my-session

# 使用特定Shell
./cligool-client -server https://relay.example.com -shell /bin/zsh

# 调试模式
./cligool-client -server https://relay.example.com -debug
```

## 🖥️ Web界面功能

### 界面说明

```
┌─────────────────────────────────────────────┐
│ CliGool          [状态] [会话ID] [连接/断开] │
├─────────────────────────────────────────────┤
│                                             │
│              终端显示区域                     │
│                                             │
│              (实时终端输出)                   │
│                                             │
├─────────────────────────────────────────────┤
│ 快捷键 | 连接信息                           │
└─────────────────────────────────────────────┘
```

### 功能按钮

- **连接按钮**: 连接到指定的终端会话
- **断开按钮**: 断开当前连接
- **状态指示**: 显示当前连接状态

### 快捷键

- `Ctrl+C` - 复制选中文本
- `Ctrl+V` - 粘贴文本
- `Ctrl+Shift+F` - 搜索终端内容

### 连接状态

- **未连接** - 灰色，没有活动连接
- **连接中** - 黄色，正在建立连接
- **已连接** - 绿色，连接正常

## 📋 API接口

### 健康检查

```bash
GET /api/health
```

响应：
```json
{
  "status": "ok",
  "time": 1699123456
}
```

### 创建会话

```bash
POST /api/sessions
Content-Type: application/json

{
  "owner": "user123"
}
```

响应：
```json
{
  "id": "session-id",
  "owner": "user123",
  "created_at": "2023-11-04T12:34:56Z",
  "updated_at": "2023-11-04T12:34:56Z",
  "active": false
}
```

### 获取会话信息

```bash
GET /api/sessions/{session_id}
```

### 删除会话

```bash
DELETE /api/sessions/{session_id}
```

### 列出所有会话

```bash
GET /api/sessions
```

### WebSocket连接

```bash
WS /api/terminal/{session_id}?type=web&user_id={user_id}
```

## 🔐 安全注意事项

### 访问控制

1. **会话ID保密**: 会话ID相当于访问凭证，不要分享给不信任的人
2. **使用HTTPS**: 生产环境务必使用HTTPS加密通信
3. **定期更换会话ID**: 定期更换会话ID提高安全性
4. **监控连接**: 定期检查活动连接，及时发现异常访问

### 最佳实践

1. **生产环境部署**
   - 使用域名和HTTPS证书
   - 配置防火墙规则
   - 启用日志监控
   - 定期更新系统

2. **会话管理**
   - 为不同用途使用不同会话ID
   - 及时删除不需要的会话
   - 限制每个用户的会话数量

3. **用户协作**
   - 明确会话的所有者和参与者角色
   - 避免同时多个用户输入
   - 尊重他人的会话隐私

## 🐛 故障排除

### 连接问题

**问题**: 无法连接到中继服务器
```
解决方案:
1. 检查网络连接
2. 确认服务器地址正确
3. 检查防火墙设置
4. 验证服务器状态
```

**问题**: WebSocket连接断开
```
解决方案:
1. 检查网络稳定性
2. 确认服务器运行正常
3. 验证会话ID正确
4. 查看浏览器控制台错误信息
```

### 终端问题

**问题**: 终端显示异常
```
解决方案:
1. 检查PTY权限
2. 验证Shell路径
3. 尝试不同的Shell
4. 查看客户端日志
```

**问题**: 特殊按键不工作
```
解决方案:
1. 检查终端类型设置
2. 验证TERM环境变量
3. 尝试不同的终端类型
```

### 性能问题

**问题**: 终端响应缓慢
```
解决方案:
1. 检查网络延迟
2. 减少终端输出
3. 调整缓冲区大小
4. 优化网络配置
```

## 📞 获取帮助

- **文档**: [https://github.com/cligool/cligool/docs](https://github.com/cligool/cligool/docs)
- **问题反馈**: [https://github.com/cligool/cligool/issues](https://github.com/cligool/cligool/issues)
- **讨论区**: [https://github.com/cligool/cligool/discussions](https://github.com/cligool/cligool/discussions)

## 🔄 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史。