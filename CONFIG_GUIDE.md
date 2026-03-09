# CliGool 配置文件使用指南

CliGool 支持通过配置文件来设置常用参数，避免每次启动时都需要输入相同的参数。

## 配置文件位置

CliGool 会按以下优先级查找配置文件：

1. `./cligool.json` - 当前工作目录
2. `~/.cligool.json` - 用户主目录

**优先级**：当前目录的配置文件优先级高于用户主目录。

## 配置文件格式

配置文件使用 JSON 格式：

```json
{
  "server": "https://cligool.ystone.us",
  "cols": 120,
  "rows": 80
}
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `server` | string | `https://cligool.ystone.us` | 中继服务器URL |
| `cols` | int | `0` | 终端列数（0=自动检测） |
| `rows` | int | `0` | 终端行数（0=自动检测） |

## 自动创建配置文件

如果配置文件不存在，CliGool 会自动在用户主目录创建 `~/.cligool.json` 并设置默认值：

```json
{
  "server": "https://cligool.ystone.us",
  "cols": 0,
  "rows": 0
}
```

**注意**：`cols` 和 `rows` 设置为 `0` 表示自动检测终端大小。

## 命令行参数优先级

命令行参数的优先级高于配置文件。如果同时提供了配置文件和命令行参数，则使用命令行参数的值。

### 示例

**配置文件** (`~/.cligool.json`):
```json
{
  "server": "https://cligool.ystone.us",
  "cols": 120,
  "rows": 80
}
```

**命令**：
```bash
./cligool-darwin-arm64
```
**结果**：使用配置文件的值（server=https://cligool.ystone.us, cols=120, rows=80）

**命令**：
```bash
./cligool-darwin-arm64 -server https://other-server.com -cols 100
```
**结果**：
- `server`: 使用命令行参数 `https://other-server.com`
- `cols`: 使用命令行参数 `100`
- `rows`: 使用配置文件值 `80`

## 使用场景

### 场景1：固定服务器地址

如果你总是连接到同一个服务器，可以创建配置文件：

**~/.cligool.json**:
```json
{
  "server": "https://my-server.com",
  "cols": 0,
  "rows": 0
}
```

之后只需运行：
```bash
./cligool-darwin-arm64
```

### 场景2：自定义终端大小

如果你喜欢特定的终端大小：

**~/.cligool.json**:
```json
{
  "server": "https://cligool.ystone.us",
  "cols": 140,
  "rows": 40
}
```

### 场景3：项目特定配置

在不同项目中使用不同的配置：

**项目A/cligool.json**:
```json
{
  "server": "https://project-a-server.com",
  "cols": 120,
  "rows": 36
}
```

**项目B/cligool.json**:
```json
{
  "server": "https://project-b-server.com",
  "cols": 100,
  "rows": 30
}
```

在项目A目录运行：
```bash
cd project-a
./cligool-darwin-arm64  # 使用项目A的配置
```

在项目B目录运行：
```bash
cd project-b
./cligool-darwin-arm64  # 使用项目B的配置
```

## 配置文件管理

### 查看当前使用的配置文件

启动客户端时会显示使用的配置文件路径：

```
✅ 已加载配置文件: /Users/username/.cligool.json
```

### 编辑配置文件

使用任何文本编辑器编辑配置文件：

```bash
# macOS/Linux
nano ~/.cligool.json
vim ~/.cligool.json

# Windows
notepad %USERPROFILE%\.cligool.json
```

### 删除配置文件

如果不再需要配置文件，直接删除即可：

```bash
# macOS/Linux
rm ~/.cligool.json

# Windows
del %USERPROFILE%\.cligool.json
```

删除后，客户端会重新创建默认配置文件。

## Windows 特殊说明

Windows 上的配置文件位置：

- **当前目录**: `cligool.json`
- **用户主目录**: `C:\Users\YourUsername\.cligool.json`

创建配置文件示例（PowerShell）：

```powershell
# 创建默认配置
echo '{"server":"https://cligool.ystone.us","cols":120,"rows":80}' | Out-File -FilePath $env:USERPROFILE\.cligool.json -Encoding utf8
```

## 故障排除

### 问题：配置文件不生效

**检查**：
1. 确认配置文件路径是否正确
2. 确认JSON格式是否正确（可以使用在线JSON验证工具）
3. 查看启动日志，确认加载的配置文件路径

### 问题：配置文件格式错误

**错误示例**：
```json
{
  "server": "https://cligool.ystone.us",
  "cols": 120,
  "rows": 80,
}
```

**问题**：最后一项后面有逗号（这是无效的JSON）

**正确格式**：
```json
{
  "server": "https://cligool.ystone.us",
  "cols": 120,
  "rows": 80
}
```

### 问题：无法创建配置文件

**原因**：用户主目录没有写权限

**解决方案**：
```bash
# 检查权限
ls -la ~ | grep .cligool.json

# 手动创建
touch ~/.cligool.json
echo '{"server":"https://cligool.ystone.us","cols":0,"rows":0}' > ~/.cligool.json
```

## 高级用法

### 结合命令行参数使用

```bash
# 使用配置文件的服务器，但覆盖终端大小
./cligool-darwin-arm64 -cols 100 -rows 25

# 使用配置文件的终端大小，但连接到其他服务器
./cligool-darwin-arm64 -server https://temp-server.com

# 运行AI CLI工具，使用配置文件的终端大小
./cligool-darwin-arm64 -cmd claude
```

### 多环境配置

为不同环境创建不同的配置文件：

**~/.cligool-production.json**:
```json
{
  "server": "https://production.example.com",
  "cols": 120,
  "rows": 36
}
```

**~/.cligool-development.json**:
```json
{
  "server": "https://development.example.com",
  "cols": 140,
  "rows": 40
}
```

使用时复制为配置文件：

```bash
# 生产环境
cp ~/.cligool-production.json ~/.cligool.json
./cligool-darwin-arm64

# 开发环境
cp ~/.cligool-development.json ~/.cligool.json
./cligool-darwin-arm64
```

## 最佳实践

1. **使用版本控制**：将项目特定的 `cligool.json` 加入版本控制
2. **忽略个人配置**：将 `~/.cligool.json` 加入 `.gitignore`
3. **文档化配置**：在项目的 README 中说明配置文件的用途
4. **合理设置终端大小**：根据实际需求设置，不要过大或过小
5. **使用自动检测**：如果不确定，保持 `cols` 和 `rows` 为 `0`（自动检测）

## 相关文档

- [命令行参数使用](CMD_ARGS_USAGE.md)
- [快速开始](QUICKSTART.md)
- [Windows支持说明](docs/WINDOWS_SUPPORT.md)
