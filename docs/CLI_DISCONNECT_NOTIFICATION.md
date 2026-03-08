# CLI 断开通知功能

## 功能概述

当 CLI 客户端意外断开连接时，中继服务器会立即通知所有连接的 Web 客户端，并自动关闭 WebSocket 连接。

## 问题场景

**之前的行为**：
- CLI 客户端断开后，Web 客户端不知道
- Web 客户端继续等待输入，用户体验差
- 没有明确的错误提示

**改进后的行为**：
- CLI 客户端断开时，Web 客户端立即收到通知
- 显示友好的错误消息
- 自动关闭 WebSocket 连接
- 提供明确的故障排查建议

## 实现细节

### 中继服务器端

**文件**：`internal/relay/relay.go`

#### 1. 检测 CLI 断开

在 `handleClientConnection` 函数中：

```go
for {
    var msg TerminalMessage
    err := conn.ReadJSON(&msg)
    if err != nil {
        log.Printf("❌ CLI客户端连接断开: %v", err)
        break
    }
    // ... 处理消息
}

// CLI客户端断开，通知所有Web客户端
log.Printf("🔔 通知所有Web客户端: CLI已断开")
s.notifyWebClientsClientDisconnected(session)
```

#### 2. 通知 Web 客户端

新增函数 `notifyWebClientsClientDisconnected`：

```go
func (s *Service) notifyWebClientsClientDisconnected(session *Session) {
    // 创建关闭消息
    closeMsg := TerminalMessage{
        Type:   "close",
        Data:   "CLI客户端已断开连接",
        Session: session.ID,
    }

    // 向所有Web客户端发送关闭消息
    for userID, conn := range session.Clients {
        conn.WriteMessage(websocket.TextMessage, jsonData)
        conn.Close()  // 关闭连接
    }

    // 清空Web客户端列表
    session.Clients = make(map[string]*websocket.Conn)
}
```

### Web 客户端端

**文件**：`web/terminal.html`

#### 处理关闭消息

```javascript
case 'close':
    terminal.writeln('');
    if (msg.data && msg.data.includes('CLI客户端已断开')) {
        terminal.writeln('\x1b[31m❌ CLI客户端已断开连接\x1b[0m');
        terminal.writeln('');
        terminal.writeln('\x1b[33m⚠️  远程终端已关闭\x1b[0m');
        terminal.writeln('\x1b[36m💡 提示: 请确保CLI客户端正在运行\x1b[0m');
    } else {
        terminal.writeln('\x1b[31m❌ 会话已关闭\x1b[0m');
    }
    terminal.writeln('');
    disconnect();
    break;
```

## 测试方法

### 手动测试

使用提供的测试脚本：

```bash
./test-cli-disconnect.sh
```

**测试步骤**：

1. **启动 CLI 客户端**
   ```bash
   ./bin/cligool-darwin-arm64 -server http://localhost:8081 -session test-disconnect -no-browser
   ```

2. **打开浏览器访问**
   ```bash
   open http://localhost:8081/session/test-disconnect
   ```

3. **验证连接正常**
   - Web 终端显示 CLI 端的输出
   - 可以看到提示符

4. **强制关闭 CLI 客户端**
   - 按 `Ctrl+C` 关闭 CLI 客户端
   - 或者直接关闭终端窗口

5. **验证 Web 端响应**
   - ✅ 立即显示错误消息
   - ✅ 显示红色的错误提示
   - ✅ 显示故障排查建议
   - ✅ WebSocket 连接自动关闭

### 自动化测试

```bash
# 终端1：启动 CLI 客户端
./bin/cligool-darwin-arm64 -server http://localhost:8081 -session test-auto -no-browser &
CLI_PID=$!

# 等待连接
sleep 2

# 在浏览器中打开会话
open http://localhost:8081/session/test-auto

# 等待观察
sleep 5

# 强制关闭 CLI 客户端
kill -INT $CLI_PID

# 观察浏览器中的反应
```

## 服务器日志

**正常的断开流程**：

```
❌ CLI客户端连接断开: read tcp 192.168.x.x:xxxxx->192.168.x.x:8080: use of closed network connection
🔔 通知所有Web客户端: CLI已断开
📡 通知 2 个Web客户端: CLI已断开
✅ 已通知Web客户端: web-1741234567890
✅ 已通知Web客户端: web-1741234567891
🧹 已清理所有Web客户端连接
```

**无 Web 客户端时**：

```
❌ CLI客户端连接断开: ...
🔔 通知所有Web客户端: CLI已断开
📭 没有Web客户端需要通知
```

## 用户体验

### Web 端显示

**CLI 断开时**：

```
❌ CLI客户端已断开连接

⚠️  远程终端已关闭
💡 提示: 请确保CLI客户端正在运行
```

**其他类型关闭时**：

```
❌ 会话已关闭
[附加消息]
```

## 优势

1. **实时通知**：Web 客户端立即知道 CLI 端断开
2. **清晰提示**：友好的错误消息，帮助用户理解问题
3. **自动清理**：自动关闭 WebSocket 连接，避免资源泄漏
4. **故障排查**：提供明确的建议，帮助用户解决问题
5. **多客户端支持**：同时通知所有连接的 Web 客户端

## 相关场景

### 1. 正常关闭

- 用户按 `Ctrl+C` 关闭 CLI 客户端
- ✅ Web 客户端收到通知

### 2. 异常断开

- 网络中断
- 进程崩溃
- 系统关机
- ✅ Web 客户端收到通知

### 3. 主动重连

- Web 客户端可以显示"重连"按钮
- 提示用户重新启动 CLI 客户端
- 自动刷新页面或重新连接

## 未来改进

1. **自动重连**：Web 客户端尝试自动重连
2. **重连按钮**：提供手动重连按钮
3. **状态显示**：显示 CLI 连接状态（已连接/已断开）
4. **历史记录**：保留断开前的终端输出
5. **通知选项**：桌面通知或声音提示

## 相关文件

- `internal/relay/relay.go` - 中继服务器实现
- `web/terminal.html` - Web 客户端实现
- `test-cli-disconnect.sh` - 测试脚本
- `docs/CLI_DISCONNECT_NOTIFICATION.md` - 本文档
