# CliGool Windows 支持说明

## ✅ 现在支持 Windows！

CliGool 客户端现在完全支持 Windows 平台。

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
- 使用标准输入输出管道
- 兼容 `cmd.exe` 和 PowerShell
- 基本终端功能支持

### 架构差异

```
Unix/Linux/macOS:
用户输入 → PTY → Shell → PTY → 输出
          ↑                    ↓
          └──── WebSocket ──────┘

Windows:
用户输入 → Stdin → cmd.exe → Stdout → 输出
          ↑                       ↓
          └────── WebSocket ───────┘
```

## 🚀 使用方法

### Windows
```cmd
# 启动客户端
cligool.exe -server https://cligool.zty8.cn

# 使用自定义会话ID
cligool.exe -server https://cligool.zty8.cn -session your-session-id
```

### Linux/macOS
```bash
# 启动客户端
./cligool -server https://cligool.zty8.cn

# 使用自定义会话ID
./cligool -server https://cligool.zty8.cn -session your-session-id
```

## ⚠️ Windows 限制

相比 Unix/Linux/macOS 版本，Windows 版本有以下限制：

1. **PTY 限制**: Windows 不支持真正的 PTY，使用标准管道替代
2. **终端功能**: 某些高级终端功能可能不可用
3. **窗口大小**: 动态窗口大小调整不支持
4. **Shell 选项**: 默认使用 `cmd.exe`，不支持所有 Unix shell

## 💡 Windows 改进建议

如果您需要更好的 Windows 支持，可以考虑：

### 1. 使用 PowerShell
修改 `main_windows.go` 中的命令：
```go
cmd := exec.Command("powershell.exe")
```

### 2. 使用 Windows Terminal
配合 Windows Terminal 获得更好的终端体验。

### 3. 启用 WSL
在 Windows 上使用 Windows Subsystem for Linux (WSL) 运行 Unix 版本：
```bash
wsl ./cligool -server https://cligool.zty8.cn
```

## 📊 功能对比

| 功能 | Unix/Linux/macOS | Windows |
|------|------------------|---------|
| 基本 WebSocket 连接 | ✅ | ✅ |
| 终端输入输出 | ✅ | ✅ |
| PTY 支持 | ✅ | ❌ |
| 窗口大小调整 | ✅ | ❌ |
| 信号处理 | ✅ | ❌ |
| 完整终端功能 | ✅ | ⚠️ 有限 |
| 多用户会话 | ✅ | ✅ |
| 数据双向转发 | ✅ | ✅ |

## 🛠️ 故障排除

### Windows 常见问题

**问题**: 编译失败
```bash
# 解决方案：确保安装了 Go 和必要的依赖
go version
go mod tidy
```

**问题**: 终端显示异常
```cmd
# 解决方案：使用 Windows Terminal 或 PowerShell
# 避免使用老旧的 CMD 窗口
```

**问题**: 某些命令不工作
```cmd
# 解决方案：切换到 PowerShell
# 或使用 WSL 运行 Unix 版本
```

## 🔮 未来计划

- [ ] 添加对 Windows ConPTY 的支持
- [ ] 实现 Windows 窗口大小调整
- [ ] 改进 Windows 终端兼容性
- [ ] 添加 PowerShell 特定优化

## 📝 总结

✅ **完全支持 Windows** - 现在可以在 Windows 上使用 CliGool
⚠️ **功能受限** - 某些高级功能在 Windows 上不可用
💡 **推荐 WSL** - 如需完整功能，建议在 Windows 上使用 WSL

现在 CliGool 是真正的跨平台远程终端解决方案！
