# CliGool Windows 快速开始指南

## 🎉 现在支持 Windows！

CliGool 客户端现在完全支持 Windows 平台。

## 📥 下载和安装

### 方法1: 下载预编译版本
从 releases 页面下载 Windows 版本：
- `cligool-windows-amd64.exe` - 64位 Windows

### 方法2: 自己编译
```cmd
# 需要安装 Go 1.21+
git clone https://github.com/your-repo/cligool.git
cd cligool
go build -o cligool.exe ./cmd/client
```

## 🚀 快速开始

### 1. 启动客户端
```cmd
cligool.exe -server https://cligool.ystone.us
```

### 2. 查看输出
客户端会显示：
```
╔═══════════════════════════════════════════════════════════╗
║                    🚀 CliGool 远程终端                      ║
╠═══════════════════════════════════════════════════════════╣
║ 📋 会话ID: your-uuid-here                                ║
║ 🌐 Web访问: https://cligool.ystone.us/session/your-uuid    ║
║ 🔗 连接状态: 🟡 连接中...                                   ║
╚═══════════════════════════════════════════════════════════╝

✅ WebSocket已连接
✅ 已连接到中继服务器
💡 现在可以在Web终端中输入命令了
```

### 3. 打开Web终端
在浏览器中访问显示的URL，就可以开始使用远程终端了！

## 🛠️ 配置选项

### 指定服务器
```cmd
cligool.exe -server https://your-server.com
```

### 使用自定义会话ID
```cmd
cligool.exe -server https://cligool.ystone.us -session my-custom-session-id
```

## ⚠️ Windows 限制和注意事项

### 功能限制
相比 Linux/macOS 版本，Windows 版本有以下限制：

1. **终端模拟**: Windows 使用标准管道而不是真正的 PTY
2. **窗口大小**: 不支持动态调整终端窗口大小
3. **Shell选项**: 默认使用 `cmd.exe`，某些Unix命令不可用
4. **信号处理**: 不支持Unix信号（如SIGWINCH）

### 推荐使用环境
- ✅ Windows 10/11
- ✅ Windows Terminal（推荐）
- ✅ PowerShell 5.1+
- ⚠️ 传统的 CMD 窗口（功能有限）

### 更好的替代方案

如果您需要完整的Unix终端体验，推荐：

#### 使用 WSL (Windows Subsystem for Linux)
```cmd
# 在WSL中使用Linux版本的CliGool
wsl
cd /path/to/cligool
./cligool -server https://cligool.ystone.us
```

#### 使用 Git Bash
```bash
# 在Git Bash中使用
./cligool.exe -server https://cligool.ystone.us
```

## 📊 功能对比

| 功能 | Windows | Linux/macOS |
|------|---------|-------------|
| 基本终端 | ✅ | ✅ |
| 命令执行 | ✅ | ✅ |
| 文本编辑 | ⚠️ 有限 | ✅ |
| 彩色输出 | ✅ | ✅ |
| 特殊字符 | ⚠️ 有限 | ✅ |
| 窗口调整 | ❌ | ✅ |
| 完整PTY | ❌ | ✅ |

## 🐛 故障排除

### 问题1: 终端显示乱码
**解决方案**:
```cmd
chcp 65001
cligool.exe -server https://cligool.ystone.us
```

### 问题2: 某些命令不工作
**解决方案**: 使用 PowerShell 或 WSL
```cmd
# 使用PowerShell
powershell.exe -Command "& { ./cligool.exe -server https://cligool.ystone.us }"
```

### 问题3: 编译失败
**解决方案**:
```cmd
# 确保Go版本正确
go version

# 清理依赖
go mod tidy

# 重新编译
go build -o cligool.exe ./cmd/client
```

## 📚 进一步阅读

- [完整功能对比](./WINDOWS_SUPPORT.md)
- [开发文档](./DEVELOPMENT.md)
- [API文档](../README.md)

## 💡 最佳实践

1. **使用Windows Terminal**: 获得更好的终端体验
2. **考虑WSL**: 如需完整Unix功能，使用WSL
3. **PowerShell集成**: 可以与PowerShell脚本集成
4. **防火墙设置**: 确保允许WebSocket连接

## 🎯 总结

✅ **完全支持Windows** - 现在可以在Windows上使用CliGool
⚠️ **功能有限** - 某些高级功能在Windows上受限
💡 **推荐WSL** - 如需完整功能，建议使用WSL

**开始使用**: 只需下载 `cligool-windows-amd64.exe` 并运行！

如有问题，请查看 [故障排除](#-故障排除) 或提交issue。
