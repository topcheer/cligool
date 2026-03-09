# 禁止自动打开浏览器

## 功能说明

默认情况下，CliGool 客户端连接到中继服务器后会自动打开浏览器。如果你不希望自动打开浏览器，可以使用 `-no-browser` 参数。

## 使用方法

### 命令行参数

```bash
# Unix/Linux/macOS
./cligool-darwin-arm64 -server http://localhost:8081 -no-browser

# Windows
cligool-windows-amd64.exe -server http://localhost:8081 -no-browser
```

### 配置文件

你也可以在配置文件中设置默认行为：

```json
{
  "server": "https://cligool.ystone.us",
  "proxy": "",
  "cols": 0,
  "rows": 0,
  "no_browser": true
}
```

配置文件位置：
- `./cligool.json` (当前目录)
- `~/.cligool.json` (用户主目录)

## 使用场景

- **服务器环境**：在没有图形界面的服务器上运行时
- **远程连接**：通过 SSH 连接到远程机器时
- **自动化脚本**：在脚本中使用 CliGool 时
- **多窗口开发**：自己打开浏览器并手动访问 URL 时

## 行为对比

### 默认行为（不使用 `-no-browser`）

```
✅ 已连接到中继服务器
✅ 已在浏览器中打开: http://localhost:8081/session/xxx
```

### 使用 `-no-browser` 参数

```
✅ 已连接到中继服务器
```

（浏览器不会自动打开）

## 手动访问

即使使用了 `-no-browser` 参数，你仍然可以手动访问 Web 终端：

1. 复制连接信息中显示的 Web 访问地址
2. 在浏览器中打开该地址

示例地址：`http://localhost:8081/session/your-session-id`

## 平台差异

### Unix/Linux/macOS
- 使用系统通知（terminal-notifier）
- 自动打开默认浏览器

### Windows
- 使用 PowerShell BurntToast 模块（如果可用）
- 使用 rundll32 或 cmd start 命令打开浏览器

使用 `-no-browser` 参数后，这些平台特定的行为都会被禁用。
