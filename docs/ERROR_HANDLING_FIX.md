# 错误处理和 WebSocket 清理修复

## 问题描述

**之前的代码存在以下问题**：

1. **错误退出时没有清理 WebSocket**：当命令执行出错时（如命令不存在），客户端直接退出，但没有通知 relay 服务器
2. **Web 客户端成为"僵尸连接"**：WebSocket 连接仍然存在，但 CLI 客户端已经退出
3. **用户体验差**：用户在 Web 端不知道发生了什么错误

## 具体场景

### 场景 1：命令不存在

```bash
./bin/cligool-darwin-arm64 -cmd nonexistent-command
```

**之前的行为**：
1. WebSocket 连接成功
2. 尝试启动 PTY，失败（命令不存在）
3. 直接返回错误并退出
4. ❌ WebSocket 连接没有关闭
5. ❌ Web 客户端不知道发生了什么

**修复后的行为**：
1. WebSocket 连接成功
2. 尝试启动 PTY，失败（命令不存在）
3. ✅ 发送 "close" 消息给 relay（包含错误信息）
4. ✅ 关闭 WebSocket 连接
5. ✅ Web 客户端收到错误通知

### 场景 2：PTY 运行时错误

**之前的行为**：
- PTY 读取失败时直接退出
- ❌ 没有通知 Web 客户端

**修复后的行为**：
- ✅ 发送关闭消息并说明错误
- ✅ Web 客户端显示具体错误

## 实现细节

### Unix 客户端 (main_unix.go)

#### 1. 使用 defer 确保清理

```go
func runTerminalSession(...) error {
    // ... WebSocket 连接 ...

    var sessionError error
    defer func() {
        // 发送关闭消息
        closeMsg := TerminalMessage{
            Type:   "close",
            Session: sessionID,
            UserID:  "client",
        }

        if sessionError != nil {
            closeMsg.Data = fmt.Sprintf("客户端错误: %v", sessionError)
        } else {
            closeMsg.Data = "客户端正常退出"
        }

        jsonData, _ := json.Marshal(closeMsg)
        conn.WriteMessage(websocket.TextMessage, jsonData)
        conn.Close()
    }()

    // ... PTY 启动 ...
    if err != nil {
        sessionError = fmt.Errorf("启动PTY失败: %w", err)
        return sessionError
    }

    // ... PTY 读取循环 ...
    if err != nil {
        if err == io.EOF {
            return nil  // 正常退出
        }
        sessionError = fmt.Errorf("PTY读取失败: %w", err)
        return sessionError
    }
}
```

#### 2. 避免 log.Fatalf

**之前**：
```go
if err := runTerminalSession(...); err != nil {
    log.Fatalf("终端会话失败: %v", err)  // ❌ 不会执行 defer
}
```

**修复后**：
```go
if err := runTerminalSession(...); err != nil {
    log.Printf("❌ 终端会话失败: %v", err)
    fmt.Printf("\n❌ 错误: %v\n", err)
    os.Exit(1)  // ✅ defer 已经执行
}
```

### Windows 客户端 (main_windows.go)

应用了相同的修复模式：

1. 添加 `defer` 函数发送关闭消息
2. 使用 `sessionError` 变量跟踪错误
3. 在所有错误返回点设置 `sessionError`

### Web 客户端 (terminal.html)

Web 客户端已经正确处理 "close" 消息：

```javascript
case 'close':
    terminal.writeln('');
    if (msg.data && msg.data.includes('CLI客户端已断开')) {
        terminal.writeln('\x1b[31m❌ CLI客户端已断开连接\x1b[0m');
        terminal.writeln('');
        terminal.writeln('\x1b[33m⚠️  远程终端已关闭\x1b[0m');
        terminal.writeln('\x1b[36m💡 提示: 请确保CLI客户端正在运行\x1b[0m');
    } else if (msg.data && msg.data.includes('客户端错误')) {
        terminal.writeln('\x1b[31m❌ 客户端发生错误\x1b[0m');
        terminal.writeln(`\x1b[90m${msg.data}\x1b[0m`);
        terminal.writeln('');
        terminal.writeln('\x1b[33m⚠️  远程终端已关闭\x1b[0m');
    } else {
        terminal.writeln('\x1b[31m❌ 会话已关闭\x1b[0m');
    }
    terminal.writeln('');
    disconnect();
    break;
```

## 测试方法

### 测试命令不存在

```bash
# 终端1：启动客户端（命令不存在）
./bin/cligool-darwin-arm64 -cmd nonexistent-command -server http://localhost:8081 -session test-error

# 预期输出：
# ❌ 启动PTY失败: exec: "nonexistent-command": executable file not found
# ❌ 错误: 启动PTY失败: exec: "nonexistent-command": executable file not found

# 终端2：打开浏览器
open http://localhost:8081/session/test-error

# 预期结果：
# ✅ Web 端显示错误消息
# ✅ WebSocket 连接关闭
```

### 测试正常退出

```bash
# 终端1：启动客户端
./bin/cligool-darwin-arm64 -server http://localhost:8081 -session test-normal

# 终端2：打开浏览器
open http://localhost:8081/session/test-normal

# 终端1：按 Ctrl+C 退出

# 预期结果：
# ✅ Web 端显示 "CLI客户端已断开连接"
# ✅ WebSocket 连接关闭
```

## 改进效果

### 用户体验

**之前**：
- ❌ Web 端显示"僵尸连接"
- ❌ 用户不知道发生了什么
- ❌ 需要刷新页面才能发现问题

**修复后**：
- ✅ Web 端立即显示错误信息
- ✅ 用户清楚知道发生了什么
- ✅ WebSocket 正确关闭，可以重新连接

### 服务器资源

**之前**：
- ❌ 僵尸 WebSocket 连接占用资源
- ❌ 会话永不清理

**修复后**：
- ✅ 错误时立即关闭连接
- ✅ 会话及时清理
- ✅ 资源正确释放

## 日志示例

### 命令不存在

```
✅ WebSocket已连接
✅ 已连接到中继服务器
直接执行命令: nonexistent-command
❌ 启动PTY失败: exec: "nonexistent-command": executable file not found
❌ 发送错误关闭消息: 启动PTY失败: exec: "nonexistent-command": executable file not found
```

### 正常退出

```
✅ WebSocket已连接
...
PTY读取失败: read unix /dev/ptmx: EOF
✅ 发送正常关闭消息
```

## 相关文件

- `cmd/client/main_unix.go` - Unix 客户端错误处理修复
- `cmd/client/main_windows.go` - Windows 客户端错误处理修复
- `web/terminal.html` - Web 客户端关闭消息处理
- `internal/relay/relay.go` - 中继服务器清理逻辑
- `test-error-handling.sh` - 测试脚本
- `docs/ERROR_HANDLING_FIX.md` - 本文档
