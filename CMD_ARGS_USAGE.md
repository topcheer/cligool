# CliGool 命令行参数支持

CliGool 现在支持向指定的命令传递命令行参数。

## 基本用法

```bash
cligool [选项] -cmd <命令> -args <参数>
```

### 选项说明

- `-cmd <命令>`: 要执行的命令（如 `claude`, `gemini`, `git`, `npm` 等）
- `-args <参数>`: 传递给命令的参数，用空格分隔多个参数
- `-server <URL>`: 中继服务器地址（默认：https://cligool.zty8.cn）
- `-session <ID>`: 会话ID（可选，默认自动生成）
- `-cols <数字>`: 终端列数（0=自动检测）
- `-rows <数字>`: 终端行数（0=自动检测）

## 使用示例

### Windows

```powershell
# 执行 claude 并传递 commit 命令
.\cligool-windows-amd64.exe -cmd claude -args "commit -m 'fix bug'"

# 执行 git status
.\cligool-windows-amd64.exe -cmd git -args "status"

# 执行带多个参数的命令
.\cligool-windows-amd64.exe -cmd npm -args "install --save-dev"

# 执行 ls 命令查看目录
.\cligool-windows-amd64.exe -cmd ls -args "-la /tmp"

# 执行自定义脚本
.\cligool-windows-amd64.exe -cmd "C:\path\to\script.bat" -args "param1 param2"
```

### macOS/Linux

```bash
# 执行 gemini 并传递参数
./cligool-darwin-arm64 -cmd gemini -args "chat --model gemini-pro"

# 执行带多个参数的命令
./cligool-linux-amd64 -cmd npm -args "install --save-dev"

# 执行带短横线参数的命令
./cligool-linux-amd64 -cmd ls -args "-la /tmp"

# 执行 Python 脚本
./cligool-darwin-arm64 -cmd python -args "script.py arg1 arg2"

# 使用默认 shell（不指定 -cmd）
./cligool-darwin-arm64
```

## 参数解析规则

1. **空格分隔**：多个参数用空格分隔
   ```bash
   -args "param1 param2 param3"
   ```

2. **引号处理**：可以使用引号包含带空格的参数
   ```bash
   -args "commit -m 'fix bug'"  # 'fix bug' 作为一个参数
   ```

3. **特殊字符**：特殊字符会被正确处理
   ```bash
   -args "--option=value"       # --option=value
   -args "-la"                  # -la
   -args "--verbose"            # --verbose
   ```

## 平台差异

### Windows (ConPTY)
- 使用命令行字符串方式传递参数
- 自动处理 `.cmd`, `.bat`, `.ps1` 脚本
- 参数会追加到完整命令行字符串中

### Unix/macOS (PTY)
- 使用可变参数列表方式传递参数
- 参数作为独立的数组元素传递
- 更精确的参数控制

## 注意事项

1. **参数顺序**：参数按照 `-args` 中指定的顺序传递
2. **空格处理**：参数中的空格需要用引号包围
3. **默认行为**：不指定 `-cmd` 时使用默认 shell（Windows: cmd.exe, Unix: $SHELL 或 /bin/bash）
4. **命令查找**：如果只提供命令名（无路径），会自动在 PATH 环境变量中查找

## 常见用例

### AI CLI 工具

```bash
# Claude Code
cligool -cmd claude -args "commit -m 'Add new feature'"

# Gemini
cligool -cmd gemini -args "chat --model gemini-pro"

# Cursor
cligool -cmd cursor -args "chat --help"
```

### 开发工具

```bash
# Git 操作
cligool -cmd git -args "status"
cligool -cmd git -args "log --oneline -10"

# NPM/Yarn
cligool -cmd npm -args "install"
cligool -cmd yarn -args "add lodash"

# Docker
cligool -cmd docker -args "ps -a"
```

### 系统管理

```bash
# 文件操作
cligool -cmd ls -args "-lah"
cligool -cmd find -args "/tmp -name '*.log'"

# 网络工具
cligool -cmd curl -args "-I https://example.com"
cligool -cmd ping -args "-c 4 google.com"
```

## 故障排除

### 参数未生效

1. 检查参数是否用引号正确包围
2. 确认命令路径正确
3. 查看日志输出，确认参数解析结果

### 特殊字符问题

某些特殊字符可能需要转义：
```bash
# 错误
-claude -args "message: Hello! How are you?"

# 正确
-claude -args "message: 'Hello! How are you?'"
```

## 技术细节

- **Windows**: 使用 `strings.Join()` 构建完整命令行字符串
- **Unix/macOS**: 使用 `exec.Command(fullPath, cmdArgs...)` 可变参数
- **参数解析**: 使用 `strings.Fields()` 分隔空格
- **向后兼容**: 不提供 `-args` 时行为与之前一致

## 更新日志

**2026-03-07**: 新增命令行参数支持功能
- 添加 `-args` 参数
- 支持所有平台的参数传递
- 保持向后兼容性
