# AI CLI 工具使用指南

本文档说明如何使用 CliGool 客户端运行各种 AI agent CLI 工具。

## 支持的 AI CLI 工具

改进后的 CliGool 客户端现在支持以下 AI CLI 工具：

- **Claude CLI** (Anthropic)
- **Gemini CLI** (Google)
- **Codex CLI** (OpenAI)
- **Aider CLI**
- **Cursor CLI**
- **其他基于终端的 AI 工具**

## 基本用法

### Unix/macOS/Linux

```bash
# 运行 Claude CLI
./cligool-darwin-arm64 -cmd claude

# 运行 Gemini CLI
./cligool-linux-amd64 -cmd gemini

# 运行 Aider
./cligool-darwin-arm64 -cmd aider

# 使用参数（推荐方式）
./cligool-darwin-arm64 -cmd claude -args "--model claude-3-opus-20240229"

# 多个参数
./cligool-darwin-arm64 -cmd git -args "commit -m 'fix bug'"
```

### Windows

```cmd
REM 运行 Claude CLI
cligool-windows-amd64.exe -cmd claude

REM 运行 Gemini CLI
cligool-windows-amd64.exe -cmd gemini

REM 使用参数（推荐方式）
cligool-windows-amd64.exe -cmd claude -args "--model claude-3-opus-20240229"

REM 多个参数
cligool-windows-amd64.exe -cmd git -args "commit -m 'fix bug'"
```

## 终端能力支持

CliGool 客户端现在支持以下终端特性：

### 基本特性
- ✅ **颜色支持**: 支持 256 色和真色 (TrueColor)
- ✅ **光标控制**: 支持光标移动、显示/隐藏
- ✅ **文本样式**: 支持粗体、下划线、闪烁等
- ✅ **屏幕清除**: 支持清屏和清除行
- ✅ **滚动区域**: 支持设置滚动区域

### 高级特性
- ✅ **替代屏幕**: 支持独立的屏幕缓冲区
- ✅ **状态行**: 支持设置状态行
- ✅ **鼠标支持**: 基本鼠标事件支持
- ✅ **窗口大小**: 动态窗口大小调整
- ✅ **终端查询**: 自动响应终端能力查询

### 终端能力查询

客户端自动处理以下查询序列：

- **DA1** (Device Attributes 1): `ESC c` 或 `ESC [0c`
  - 响应: `ESC [?1;2c` (VT100 with video)

- **DA2** (Device Attributes 2): `ESC [>c` 或 `ESC [>0c`
  - 响应: `ESC [>0;272;0c` (xterm version)

- **CPR** (Cursor Position Report): `ESC [6n`
  - 响应: `ESC [1;1R` (光标位置)

- **XTVERSION**: `ESC [>q`
  - 响应: `ESC [>0;272;0c`

- **XTGETTCAP**: `ESC P+q...ESC \`
  - 响应: 空响应（避免阻塞）

## 环境变量

客户端自动设置以下环境变量：

### Unix/macOS/Linux
```bash
TERM=xterm-256color
COLORTERM=truecolor
FORCE_COLOR=1
LANG=en_US.UTF-8
LC_ALL=en_US.UTF-8
```

### Windows
```cmd
TERM=xterm-256color
COLORTERM=truecolor
FORCE_COLOR=1
PYTHONIOENCODING=utf-8
```

## Web 终端特性

Web 端 xterm.js 配置：

- ✅ 启用实验性 API (`allowProposedApi`)
- ✅ 自动换行符转换 (`convertEol`)
- ✅ Alt+点击移动光标 (`altClickMovesCursor`)
- ✅ 光标闪烁动画
- ✅ 1000 行滚动缓冲

## 常见问题

### Q: AI CLI 工具显示乱码怎么办？

**A**: 这可能是编码问题。确保：
1. 系统 locale 设置为 UTF-8
2. CLI 工具本身支持 UTF-8
3. 使用 `-cmd` 参数直接运行 CLI 工具

### Q: 颜色不显示怎么办？

**A**: 检查：
1. CLI 工具是否支持颜色（查看其 `--color` 选项）
2. 环境变量 `FORCE_COLOR` 是否设置
3. 是否使用 `-cmd` 参数直接运行

### Q: 交互式输入不工作怎么办？

**A**: 确保使用 `-cmd` 参数直接运行 CLI 工具，而不是在 shell 中执行命令。

### Q: 性能问题怎么办？

**A**: 如果遇到性能问题：
1. 减少终端窗口大小（行数和列数）
2. 使用本地客户端而非 Web 端
3. 检查网络延迟

## 性能优化建议

### 客户端选择
- **最佳性能**: 使用原生客户端（Unix/macOS/Linux/Windows ConPTY）
- **Windows**: 使用ConPTY客户端（功能与Unix完全对等）
- **远程访问**: 使用 Web 端

### 终端大小
- 推荐大小: 80x24（标准）
- 大型终端: 120x36（更多空间）
- 小型终端: 40x12（最小化延迟）

### 网络优化
- 使用低延迟网络
- 启用 WebSocket 压缩（如果支持）
- 减少不必要的输出

## 示例配置

### Claude CLI
```bash
# 基本使用
./cligool-darwin-arm64 -cmd claude

# 指定模型（推荐使用 -args）
./cligool-darwin-arm64 -cmd claude -args "--model claude-3-opus-20240229"

# 会话 ID
./cligool-darwin-arm64 -cmd claude -session my-session-123

# Windows
cligool-windows-amd64.exe -cmd claude -args "--model claude-3-opus-20240229"
```

### Aider
```bash
# 基本使用
./cligool-darwin-arm64 -cmd aider

# 指定模型
./cligool-darwin-arm64 -cmd aider -args "--model gpt-4"

# Git 仓库
./cligool-darwin-arm64 -cmd aider -args "--model gpt-4" -session aider-session
```

### Gemini CLI
```bash
# 基本使用
./cligool-darwin-arm64 -cmd gemini

# 指定模型
./cligool-darwin-arm64 -cmd gemini -args "--model gemini-pro"

# Windows
cligool-windows-amd64.exe -cmd gemini -args "chat --model gemini-pro"
```

## 技术细节

### PTY 设置
- **Unix/macOS/Linux**: 使用 `github.com/creack/pty` 库
- **Windows**: 使用 ConPTY (Windows Console Pseudo Terminal)
- 终端类型: `xterm-256color`
- 窗口大小:
  - Unix: 自动检测 + SIGWINCH信号处理
  - Windows: 自动检测 + 后台监控（每500ms）
  - 默认: 120x80（可手动覆盖）

### 终端仿真
- 实现了基本的终端能力查询响应
- 支持 CSI、OSC、DCS 序列
- 自动处理转义序列

### WebSocket 通信
- 二进制安全传输
- 心跳机制（30 秒间隔）
- 自动重连（Web 端）

## 故障排查

### 调试模式
```bash
# 启用详细日志
./cligool-darwin-arm64 -cmd claude 2>&1 | tee debug.log
```

### 测试终端能力
```bash
# 运行测试脚本
python3 test_ansi.py
bash test_terminal_capabilities.sh
```

### 检查环境
```bash
# 查看环境变量
env | grep -E "TERM|LANG|LC_"

# 查看终端大小
stty size
```

## 贡献

如果发现任何问题或有改进建议，请：
1. 创建 Issue 描述问题
2. 提供复现步骤
3. 包含日志输出

## 许可证

MIT License - 详见 LICENSE 文件
