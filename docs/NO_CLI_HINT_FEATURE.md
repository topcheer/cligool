# Session 存在性检查和友好提示功能

## 功能概述

当 Web 客户端尝试连接到一个没有 CLI 客户端的 session 时，系统会显示一个友好的提示界面，包含：
1. 清晰的错误说明
2. 适合当前操作系统的启动命令示例
3. 下载客户端的链接
4. 使用提示

## 问题场景

**之前的行为**：
- Web 客户端连接到不存在的 session
- 终端一片空白，没有任何输出
- 用户不知道发生了什么
- ❌ 用户体验很差

**改进后的行为**：
- Web 客户端检测到没有 CLI 客户端
- 显示友好的提示界面
- 提供具体的命令示例
- 包含下载链接和使用提示
- ✅ 用户体验清晰友好

## 实现细节

### 中继服务器端

**文件**：`internal/relay/relay.go`

#### 1. 检查 CLI 客户端状态

```go
// 检查是否有CLI客户端连接
session.Mutex.RLock()
cliConnected := session.ClientCon != nil
session.Mutex.RUnlock()

if !cliConnected {
    // 没有CLI客户端连接，发送提示消息
    log.Printf("⚠️  Web客户端连接但无CLI客户端，发送提示消息")
    s.sendNoCliClientMessage(session, conn)
    return
}
```

#### 2. 生成平台特定的命令示例

```go
func buildNoCliHintMessage(sessionID string) string {
    ostype := runtime.GOOS
    arch := runtime.GOARCH

    var cmdExample string
    switch ostype {
    case "darwin":
        if arch == "arm64" {
            cmdExample = `./cligool-darwin-arm64 -server http://localhost:8081 -session %s`
        } else {
            cmdExample = `./cligool-darwin-amd64 -server http://localhost:8081 -session %s`
        }
    // ... 其他平台
    }

    // 返回格式化的提示消息
    return fmt.Sprintf(`⚠️  CLI客户端未连接
...
启动命令示例：
%s
...`, cmdExample, sessionID)
}
```

#### 3. 发送特殊消息类型

```go
hintMsg := TerminalMessage{
    Type:   "no_cli",  // 新的消息类型
    Data:   buildNoCliHintMessage(sessionID),
    Session: session.ID,
}
conn.WriteMessage(websocket.TextMessage, jsonData)
```

### Web 客户端端

**文件**：`web/terminal.html`

#### 1. 处理 "no_cli" 消息

```javascript
case 'no_cli':
    // 没有CLI客户端连接，显示帮助信息
    terminal.reset(); // 清空终端
    displayNoCliMessage(msg.data);
    break;
```

#### 2. 显示友好的界面

```javascript
function displayNoCliMessage(message) {
    // 清空终端并显示友好的帮助信息
    terminal.writeln('\x1b[2J\x1b[H'); // 清屏并移动到首页

    // 显示标题框
    terminal.writeln('╔═══════════════════════════════════════════════════════════════╗');
    terminal.writeln('║               ⚠️  CLI客户端未连接                              ║');
    terminal.writeln('╠═══════════════════════════════════════════════════════════════╣');
    terminal.writeln('║ 请先启动CLI客户端，然后再刷新此页面                        ║');
    terminal.writeln('╚═══════════════════════════════════════════════════════════════╝');

    // 解析并显示消息内容，带颜色高亮
    // ...
}
```

## 显示效果

### 终端界面

```
╔═══════════════════════════════════════════════════════════════╗
║               ⚠️  CLI客户端未连接                              ║
╠═══════════════════════════════════════════════════════════════╣
║ 请先启动CLI客户端，然后再刷新此页面                        ║
╚═══════════════════════════════════════════════════════════════╝

⚠️  CLI客户端未连接

请先启动CLI客户端，然后再刷新此页面。

启动命令示例：
  ./cligool-darwin-arm64 -server http://localhost:8081 -session test-123

或者使用配置文件：
1. 编辑 ~/.cligool.json 设置服务器地址
2. 运行: ./cligool -session test-123

💡 提示：
- CLI客户端必须在Web客户端之前启动
- 确保使用相同的session ID
- 检查防火墙设置

📥 下载客户端：
- https://cligool.ystone.us/

─────────────────────────────────────────────────────────────

💡 提示：
  1. 复制上面的命令到终端执行
  2. 或者编辑 ~/.cligool.json 配置文件
  3. 启动CLI客户端后刷新此页面

等待CLI客户端连接...
(终端已禁用，请启动CLI客户端后刷新页面)
```

## 平台特定命令

服务器会根据运行平台自动生成适合的命令示例：

### macOS (Intel)
```bash
./cligool-darwin-amd64 -server http://localhost:8081 -session <session-id>
```

### macOS (Apple Silicon)
```bash
./cligool-darwin-arm64 -server http://localhost:8081 -session <session-id>
```

### Linux (AMD64)
```bash
./cligool-linux-amd64 -server http://localhost:8081 -session <session-id>
```

### Linux (ARM64)
```bash
./cligool-linux-arm64 -server http://localhost:8081 -session <session-id>
```

### Windows (AMD64)
```bash
cligool-windows-amd64.exe -server http://localhost:8081 -session <session-id>
```

## 测试方法

### 方法1：使用测试脚本

```bash
./test-no-cli-hint.sh
```

### 方法2：手动测试

```bash
# 1. 生成一个随机 session ID
SESSION="test-$(date +%s)"

# 2. 在浏览器中打开（不要启动CLI客户端）
open http://localhost:8081/session/$SESSION

# 3. 应该看到友好的提示界面

# 4. 在另一个终端启动CLI客户端
./bin/cligool-darwin-arm64 -server http://localhost:8081 -session $SESSION

# 5. 刷新浏览器，应该正常连接
```

### 方法3：查看服务器日志

```bash
# 启动服务器
docker-compose -f docker-compose.dev.yml up -d

# 监控日志
docker-compose -f docker-compose.dev.yml logs -f relay-server

# 在浏览器访问一个不存在的session
# 应该看到日志：
# ⚠️  Web客户端连接但无CLI客户端，发送提示消息
# ✅ 已发送无CLI提示消息给Web客户端
```

## 优势

### 1. **用户体验**

**之前**：
- ❌ 终端空白，用户困惑
- ❌ 不知道要做什么
- ❌ 可能以为服务有问题

**现在**：
- ✅ 清晰的错误说明
- ✅ 具体的命令示例
- ✅ 友好的使用提示

### 2. **自包含命令**

- ✅ 命令示例包含实际的 session ID
- ✅ 用户可以直接复制粘贴执行
- ✅ 根据平台自动选择合适的客户端

### 3. **视觉吸引力**

- ✅ 使用颜色突出重要信息
- ✅ 结构化的提示框设计
- ✅ 专业的终端界面

### 4. **教育性**

- ✅ 告诉用户正确的使用顺序
- ✅ 提供配置文件选项
- ✅ 包含下载链接

## 技术细节

### 消息流程

```
Web客户端连接 → relay检查session状态 → 无CLI客户端
                                                    ↓
                                            发送 "no_cli" 消息
                                                    ↓
                                            Web端接收消息
                                                    ↓
                                            显示友好提示界面
```

### 特殊消息类型

新的消息类型 `"no_cli"` 与现有的消息类型：
- `"init"` - 初始化消息
- `"output"` - 终端输出
- `"input"` - 用户输入
- `"close"` - 关闭连接
- **`"no_cli"`** - **新增：无CLI客户端提示**

### 兼容性

- ✅ 不影响现有功能
- ✅ 只在没有CLI客户端时触发
- ✅ Web客户端仍然可以连接（只是显示提示）
- ✅ 一旦CLI客户端启动，刷新页面即可正常使用

## 相关文件

- `internal/relay/relay.go` - 中继服务器检查和消息发送
- `web/terminal.html` - Web 客户端提示界面显示
- `test-no-cli-hint.sh` - 测试脚本
- `docs/NO_CLI_HINT_FEATURE.md` - 本文档

## 未来改进

1. **自动重试**：定期检查是否有CLI客户端连接
2. **一键启动**：提供按钮直接下载并启动CLI客户端
3. **Session管理**：显示所有活跃的session列表
4. **快速连接**：选择一个存在的session快速连接
5. **二维码**：显示命令二维码，方便手机扫描
