# CliGool 安装包

本目录包含 CliGool 的各种平台安装包和安装脚本。

## 📦 支持的平台

### macOS
- **安装脚本**: `macos/install.sh`
- **安装位置**: `/usr/local/bin/cligool`
- **要求**: macOS 10.15+, sudo 权限

### Linux
- **安装脚本**: `linux/install.sh`
- **安装位置**: `/usr/local/bin/cligool`
- **支持架构**: amd64, arm64, 386, arm, armbe, ppc64le, ppc64, riscv64, s390x, mips, mips64le, mips64, loong64
- **要求**: sudo 权限

### Windows
- **安装方式**: 使用zip压缩包
- **下载文件**: `cligool-windows-amd64.zip` 或 `cligool-windows-arm64.zip`
- **支持架构**: amd64, arm64
- **安装步骤**:
  1. 下载对应的zip文件
  2. 解压到任意目录（如 `C:\Program Files\CliGool`）
  3. 将目录添加到系统PATH环境变量
  4. 在命令行中运行 `cligool.exe`

## 🚀 快速安装

### macOS
```bash
# 下载并运行安装脚本
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/macos/install.sh | bash

# 或本地运行
cd installers/macos
sudo ./install.sh
```

### Linux
```bash
# 下载并运行安装脚本
curl -sSL https://raw.githubusercontent.com/topcheer/cligool/main/installers/linux/install.sh | bash

# 或本地运行
cd installers/linux
sudo ./install.sh
```

### Windows
```powershell
# 1. 下载zip文件
# 从 https://github.com/topcheer/cligool/releases 下载:
# - cligool-windows-amd64.zip (Intel/AMD 64位)
# - cligool-windows-arm64.zip (ARM 64位)

# 2. 解压文件
# 右键点击zip文件 -> "解压到..."
# 解压到您想要的目录，例如: C:\Program Files\CliGool

# 3. 添加到PATH
# 方法1：系统设置
# - 打开"系统属性" -> "高级" -> "环境变量"
# - 编辑"Path"变量，添加解压目录
#
# 方法2：PowerShell命令（需要管理员权限）
# [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\CliGool", "Machine")

# 4. 验证安装
cligool --version
```

## 🔨 构建安装包

### 构建所有平台（推荐）
```bash
cd installers
chmod +x build-all.sh
./build-all.sh
```

### 单独构建

#### macOS
```bash
cd installers/macos
chmod +x install.sh
# 脚本已就绪，无需编译
```

#### Linux
```bash
cd installers/linux
chmod +x install.sh
# 脚本已就绪，无需编译
```

#### Windows

**注意**: Windows不再提供自动生成的安装程序(.exe)，请使用zip压缩包安装。

**如需手动构建Windows二进制文件**:
```powershell
# 方法1：使用Go编译
go build -o cligool.exe ../cmd/client

# 方法2：使用GitHub Actions自动构建的zip文件
# 从GitHub Releases下载预编译的zip文件
```

**在 macOS/Linux 上**:
```bash
# 仅准备脚本，无法编译 .exe
cd installers/windows
chmod +x build.sh
./build.sh  # 会提示需要在 Windows 上编译
```

## 📋 安装包功能

### macOS/Linux
- ✅ 自动检测架构
- ✅ 下载正确的二进制文件
- ✅ 安装到系统路径
- ✅ 创建配置文件
- ✅ 创建桌面快捷方式（Linux）
- ✅ 添加到 PATH

### Windows (zip压缩包)
- ✅ 跨架构支持（amd64/arm64）
- ✅ 绿色软件，无需安装
- ✅ 可移植到任意目录
- ✅ 支持手动配置PATH
- ✅ 支持自定义配置文件
- ✅ 可多个版本并存

## 🗑️ 卸载

### macOS
```bash
sudo rm -f /usr/local/bin/cligool
rm -rf ~/.cligool
rm -f ~/.cligool.json
```

### Linux
```bash
sudo rm -f /usr/local/bin/cligool
rm -rf ~/.cligool
rm -f ~/.cligool.json
rm -f ~/.local/share/applications/cligool.desktop
```

### Windows
```powershell
# 1. 删除解压目录
Remove-Item "C:\Program Files\CliGool" -Recurse -Force

# 2. 从PATH移除
# 系统设置 -> 高级 -> 环境变量 -> 编辑Path -> 删除CliGool目录

# 3. 删除配置文件（可选）
Remove-Item $env:USERPROFILE\.cligool.json -Force
Remove-Item $env:USERPROFILE\.cligool -Recurse -Force
```

## ⚙️ 高级配置

### 自定义安装位置

**macOS/Linux**:
```bash
# 下载后手动安装到自定义位置
INSTALL_DIR=/opt/cligool bash install.sh
```

**Windows**:
- 解压zip文件到任意目录即可
- 推荐目录：`C:\Program Files\CliGool` 或 `C:\CliGool`

### 仅下载不安装

**macOS/Linux**:
```bash
# 下载到当前目录
curl -LO https://github.com/topcheer/cligool/releases/latest/download/cligool-darwin-arm64
chmod +x cligool-darwin-arm64
./cligool-darwin-arm64
```

**Windows**:
```powershell
# 下载zip文件到当前目录
Invoke-WebRequest -Uri "https://github.com/topcheer/cligool/releases/latest/download/cligool-windows-amd64.zip" -OutFile "cligool-windows-amd64.zip"

# 解压文件
Expand-Archive -Path cligool-windows-amd64.zip -DestinationPath .

# 直接运行
.\cligool-windows-amd64.exe
```

## 🐛 故障排除

### macOS/Linux

#### 问题：Permission denied
```bash
# 解决方案：添加执行权限
chmod +x installers/macos/install.sh
sudo ./installers/macos/install.sh
```

#### 问题：command not found after install
```bash
# 检查 PATH
echo $PATH

# 重新加载 shell 配置
source ~/.zshrc  # 或 source ~/.bash_profile

# 或手动指定路径
/usr/local/bin/cligool
```

### Windows

#### 问题：无法在命令行中使用cligool
- **原因**: 目录未添加到PATH
- **解决方案**:
  ```powershell
  # 临时添加（当前会话有效）
  $env:Path += ";C:\Program Files\CliGool"

  # 永久添加（需要管理员权限）
  [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\CliGool", "Machine")
  ```

#### 问题：杀毒软件警告
- **原因**: 未知软件，没有数字签名
- **解决方案**:
  1. 临时禁用杀毒软件
  2. 添加到信任列表
  3. 下载时点击"保留文件"

#### 问题：找不到vcruntime140.dll等依赖
- **原因**: 缺少Visual C++运行库
- **解决方案**: CliGool是纯Go编译，不依赖外部DLL

## 📝 许可证

MIT License - 详见项目根目录的 LICENSE 文件

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

如果发现任何问题或有改进建议：
1. 创建 Issue 描述问题
2. 提供复现步骤
3. 包含日志输出（如有）

## 📚 相关文档

- [项目主页](https://github.com/topcheer/cligool)
- [使用指南](../USAGE_GUIDE_CN.md)
- [配置文件指南](../CONFIG_GUIDE.md)
- [命令行参数使用](../CMD_ARGS_USAGE.md)
