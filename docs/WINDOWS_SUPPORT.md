# CliGool Windows 支持说明

## ✅ 完整的Windows支持！

CliGool 客户端现在完全支持 Windows 平台，使用 Windows ConPTY (Console Pseudo Terminal) 技术。

## 📦 编译说明

### Windows 编译
```bash
# 直接在 Windows 上编译
go build -o cligool.exe ./cmd/client

# 或者在 Unix/Linux 上交叉编译
GOOS=windows GOARCH=amd64 go build -o cligool.exe ./cmd/client
```

### Linux/macOS 编译
```bash
go build -o cligool ./cmd/client
```

## 🔧 技术实现

### 平台特定实现

**Unix/Linux/macOS** (`main_unix.go`):
- 使用真正的 PTY (伪终端)
- 完整的终端功能支持
- 支持窗口大小调整
- 信号处理 (SIGWINCH)

**Windows** (`main_windows.go`):
- 使用 **ConPTY** (Windows Console Pseudo Terminal)
- 完整的终端功能支持（与Unix版相当）
- 支持窗口大小动态调整
- 自动检测控制台编码并转换为UTF-8
- 支持所有终端特性（颜色、光标控制等）

### 架构对比

```
Unix/Linux/macOS:
用户输入 → PTY → Shell → PTY → 输出
          ↑                    ↓
          └──── WebSocket ──────┘

Windows (ConPTY):
用户输入 → ConPTY → cmd.exe → ConPTY → 输出
          ↑                         ↓
          └──────── WebSocket ───────┘
```

**重要改进**：Windows版本现在使用ConPTY，功能与Unix版本完全对等！

## 🚀 使用方法

### Windows
```cmd
REM 启动默认客户端（cmd.exe）
cligool-windows-amd64.exe -server https://cligool.ystone.us

REM 使用自定义会话ID
cligool-windows-amd64.exe -server https://cligool.ystone.us -session your-session-id

REM 运行AI CLI工具
cligool-windows-amd64.exe -cmd claude -server https://cligool.ystone.us

REM 运行带参数的命令
cligool-windows-amd64.exe -cmd git -args "status" -server https://cligool.ystone.us

REM 使用PowerShell
cligool-windows-amd64.exe -cmd powershell.exe -server https://cligool.ystone.us
```

### Linux/macOS
```bash
# 启动客户端
./cligool-darwin-arm64 -server https://cligool.ystone.us

# 使用自定义会话ID
./cligool-darwin-arm64 -server https://cligool.ystone.us -session your-session-id

# 运行AI CLI工具
./cligool-darwin-arm64 -cmd claude -server https://cligool.ystone.us
```

## 🎯 Windows ConPTY特性

### 完整终端支持
- ✅ **真正的PTY**：Windows ConPTY提供完整的伪终端支持
- ✅ **终端功能**：支持所有终端特性（颜色、光标控制、屏幕清除等）
- ✅ **窗口大小**：动态窗口大小调整（自动检测和监控）
- ✅ **编码转换**：自动检测控制台编码并转换为UTF-8
- ✅ **多Shell支持**：cmd.exe、PowerShell、.cmd/.bat/.ps1脚本

### 编码自动检测
客户端自动检测Windows控制台编码：
- 936: GBK（简体中文）
- 932: Shift-JIS（日文）
- 949: EUC-KR（韩文）
- 950: Big5（繁体中文）
- 1252: Windows-1252（西欧）
- 437: CP437（英文）
- 其他：默认使用Windows-1252

### 动态窗口大小
- ✅ 启动时自动检测终端大小
- ✅ 运行时动态调整（每500ms监控）
- ✅ 零日志输出（不破坏终端布局）

## 💡 Windows 最佳实践

### 1. 选择合适的Shell
```cmd
REM 使用cmd.exe（默认）
cligool-windows-amd64.exe

REM 使用PowerShell
cligool-windows-amd64.exe -cmd powershell.exe

REM 使用PowerShell 7（如果已安装）
cligool-windows-amd64.exe -cmd pwsh.exe
```

### 2. 使用Windows Terminal
配合 Windows Terminal 获得最佳终端体验。

### 3. 运行脚本文件
```cmd
REM 运行.cmd脚本
cligool-windows-amd64.exe -cmd "C:\scripts\build.cmd" -args "-release"

REM 运行PowerShell脚本
cligool-windows-amd64.exe -cmd "C:\scripts\setup.ps1"
```

## 📊 功能对比

| 功能 | Unix/Linux/macOS | Windows (ConPTY) |
|------|------------------|------------------|
| 基本 WebSocket 连接 | ✅ | ✅ |
| 终端输入输出 | ✅ | ✅ |
| PTY 支持 | ✅ (PTY) | ✅ (ConPTY) |
| 窗口大小调整 | ✅ (SIGWINCH) | ✅ (监控) |
| 信号处理 | ✅ | ✅ (Windows API) |
| 完整终端功能 | ✅ | ✅ |
| 颜色支持 | ✅ | ✅ |
| 光标控制 | ✅ | ✅ |
| 多用户会话 | ✅ | ✅ |
| 数据双向转发 | ✅ | ✅ |
| 编码自动转换 | ❌ (UTF-8原生) | ✅ (检测并转换) |
| AI CLI工具支持 | ✅ | ✅ |

**结论**：Windows版本功能与Unix版本完全对等！

## 🛠️ 故障排除

### Windows 常见问题

**问题**: 编译失败
```bash
# 解决方案：确保安装了 Go 和必要的依赖
go version
go mod tidy

# Windows交叉编译
GOOS=windows GOARCH=amd64 go build -o cligool-windows-amd64.exe ./cmd/client
```

**问题**: 中文显示乱码
```cmd
# 解决方案：已自动处理 - ConPTY输出自动转换为UTF-8
# 如果仍有问题，检查控制台字体是否支持中文字符
```

**问题**: AI CLI工具无响应
```cmd
# 解决方案：确保使用 -cmd 参数直接运行工具
cligool-windows-amd64.exe -cmd claude -server https://your-domain.com

# 不要在cmd.exe中手动执行工具
```

**问题**: 窗口大小不正确
```cmd
# 解决方案：客户端会自动检测，如需手动指定
cligool-windows-amd64.exe -cols 120 -rows 36 -server https://your-domain.com
```

## 🎉 总结

✅ **完全支持 Windows** - ConPTY技术提供与Unix版本完全对等的功能
✅ **完整终端特性** - 颜色、光标控制、窗口大小调整全部支持
✅ **AI CLI工具完美运行** - Claude、Gemini、Aider等工具完美支持
✅ **自动编码转换** - 多语言字符完美显示

现在 CliGool 是真正的跨平台远程终端解决方案！

## 📚 相关文档

- **命令行参数使用**: [CMD_ARGS_USAGE.md](../CMD_ARGS_USAGE.md)
- **AI CLI工具指南**: [AI_CLI_GUIDE.md](AI_CLI_GUIDE.md)
- **快速开始**: [../QUICKSTART.md](../QUICKSTART.md)
- **主README**: [../README.md](../README.md)
