# CliGool 安装包设计方案

## 目标

为所有支持的操作系统提供原生安装包，提升用户体验。

## 支持的平台

### 1. Windows
- **格式**: Inno Setup (.exe)
- **特权用户模式**:
  - 安装位置: `C:\Program Files\CliGool`
  - PATH: 系统级 PATH
  - 快捷方式: 所有用户开始菜单 + 桌面（可选）
  - 右键菜单: "在 CliGool 中打开"（可选）
- **普通用户模式**:
  - 安装位置: `%LOCALAPPDATA%\CliGool`
  - PATH: 用户级 PATH
  - 快捷方式: 当前用户开始菜单 + 桌面（可选）
  - 限制: 无法创建系统级右键菜单

### 2. macOS
- **格式**: .pkg (Installer) + .dmg (拖拽安装)
- **安装位置**: `/usr/local/bin/` 或 `/opt/cligool/`
- **Homebrew**: Tap + Formula
- **功能**:
  - 添加到 PATH
  - 创建 LaunchAgent（可选，用于后台服务）
  - 创建 alias (可选)

### 3. Linux
- **格式**: .deb, .rpm, .sh (通用)
- **安装位置**: `/usr/local/bin/` 或 `/opt/cligool/`
- **功能**:
  - 添加到 PATH
  - 创建桌面快捷方式
  - 创建 systemd service（可选）

## 安装包结构

```
installers/
├── windows/
│   ├── cligool.iss            # Inno Setup 脚本
│   ├── cligool-admin.ico      # 管理员模式图标
│   ├── cligool-user.ico       # 普通用户模式图标
│   └── output/
│       ├── cligool-setup.exe   # 生成的安装程序
│       └── cligool-portable.exe # 便携版（可选）
├── macos/
│   ├── cligool.pkgproj        # Packages 项目文件
│   ├── distribution.xml       # 安装配置
│   ├── scripts/
│   │   ├── postinstall        # 安装后脚本
│   │   └── preinstall         # 安装前脚本
│   ├── dmg/
│   │   └── cligool.dmg        # 拖拽安装镜像
│   └── homebrew/
│       └── cligool.rb         # Homebrew Formula
├── linux/
│   ├── deb/
│   │   ├── DEBIAN/
│   │   │   ├── control        # 包控制信息
│   │   │   ├── postinst       # 安装后脚本
│   │   │   └── prerm          # 卸载前脚本
│   │   └── usr/
│   │       └── local/
│   │           └── bin/
│   │               └── cligool
│   ├── rpm/
│   │   └── cligool.spec       # RPM spec 文件
│   └── sh/
│       └── install.sh         # 通用安装脚本
└── build-all.sh               # 构建所有安装包
```

## Windows Inno Setup 设计

### 特性
1. **自动权限检测**:
   ```pascal
   [Setup]
   PrivilegesRequired=admin      # 请求管理员权限
   Uninstallable=yes
   CreateAppDir=yes
   ```

2. **动态安装路径**:
   ```pascal
   [Setup]
   DefaultDirName={code:GetInstallDir}
   ```

3. **权限提示**:
   ```pascal
   [Messages]
   WelcomeLabel2=本程序将在您的计算机上安装 CliGool 远程终端工具。
   AdminPermissions=如需安装到系统目录，请以管理员身份运行。

   [Code]
   function GetInstallDir(Param: String): String;
   begin
     if IsAdminLoggedOn then
       Result := '{pf}\CliGool'  // Program Files
     else
       Result := '{localappdata}\CliGool';  // 用户目录
   end;
   ```

4. **PATH 环境变量**:
   ```pascal
   [Registry]
   ; 系统级 PATH（管理员）
   Root: HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment; \
       ValueType: expandsz; ValueName: Path; ValueData: "{app}"; \
       Check: IsAdminLoggedOn

   ; 用户级 PATH（普通用户）
   Root: HKCU\Environment; \
       ValueType: expandsz; ValueName: Path; ValueData: "{app}"; \
       Check: not IsAdminLoggedOn
   ```

5. **快捷方式**:
   ```pascal
   [Icons]
   Name: {group}\CliGool; Filename: "{app}\cligool-windows-amd64.exe"
   Name: {group}\卸载 CliGool; Filename: "{uninstallexe}"
   Name: "{autodesktop}\CliGool"; Filename: "{app}\cligool-windows-amd64.exe"; \
       Tasks: desktopicon
   ```

6. **右键菜单集成**（仅管理员）:
   ```pascal
   [Registry]
   ; 添加"在 CliGool 中打开"右键菜单
   Root: HKCR\Directory\shell\CliGool; \
       ValueType: string; ValueData: "在 CliGool 中打开"; \
       Flags: deletekey; Check: IsAdminLoggedOn
   Root: HKCR\Directory\shell\CliGool\command; \
       ValueType: string; ValueData: """{app}\cligool-windows-amd64.exe"" -cmd ""%1"""; \
       Check: IsAdminLoggedOn
   ```

### 安装界面
- 欢迎页面（权限提示）
- 许可协议
- 安装位置选择
- 组件选择：
  - [x] 主程序
  - [ ] 创建桌面快捷方式
  - [ ] 添加到右键菜单（仅管理员）
  - [ ] 配置文件初始化
- 安装进度
- 完成页面（运行选项）

### 卸载支持
- 完整卸载（包括注册表、快捷方式、PATH）
- 保留配置文件选项
- 卸载调查（可选）

## macOS .pkg 设计

### Packages 项目配置
1. **项目设置**:
   - 标题: CliGool
   - 版本: [自动从 Git 获取]
   - 安装位置: /usr/local/bin/

2. **组件**:
   - 主程序包
   - 配置文件示例（可选）
   - 文档（可选）

3. **脚本**:
   ```bash
   # postinstall
   #!/bin/bash
   # 添加到 PATH（如果需要）
   if ! grep -q "/usr/local/bin" /etc/paths; then
       echo "/usr/local/bin" >> /etc/paths
   fi

   # 设置权限
   chmod +x /usr/local/bin/cligool

   # 创建配置文件（如果不存在）
   if [ ! -f "$HOME/.cligool.json" ]; then
       cat > "$HOME/.cligool.json" << EOF
   {
     "server": "https://cligool.zty8.cn",
     "cols": 0,
     "rows": 0
   }
   EOF
   fi
   ```

### DMG 设计
- 简单的拖拽安装
- 包含应用程序别名
- README 快捷方式
- 背景图片和样式

### Homebrew Formula
```ruby
class Cligool < Formula
  desc "Cross-platform remote terminal solution"
  homepage "https://github.com/topcheer/cligool"
  url "https://github.com/topcheer/cligool/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "..."
  license "MIT"

  def install
    bin.install "cligool-darwin-#{Hardware::CPU.arch}" => "cligool"
  end

  test do
    system "#{bin}/cligool", "-help"
  end
end
```

## Linux 安装包设计

### .deb 包 (Debian/Ubuntu)

1. **control 文件**:
   ```
   Package: cligool
   Version: 1.0.0
   Architecture: amd64
   Maintainer: Your Name <you@example.com>
   Description: Cross-platform remote terminal solution
    CliGool is a WebSocket-based remote terminal solution
    supporting 18 OS architectures.

   Depends: libc6
   Section: utils
   Priority: optional
   ```

2. **postinst 脚本**:
   ```bash
   #!/bin/bash
   set -e

   case "$1" in
       configure)
           # 创建符号链接
           ln -sf /opt/cligool/bin/cligool /usr/local/bin/cligool

           # 添加到桌面菜单
           if [ -d /usr/share/applications ]; then
               cat > /usr/share/applications/cligool.desktop << EOF
   [Desktop Entry]
   Name=CliGool
   Comment=Cross-platform remote terminal
   Exec=/usr/local/bin/cligool
   Terminal=true
   Type=Application
   Categories=Network;RemoteAccess;
   EOF
           fi
           ;;
   esac

   exit 0
   ```

### .rpm 包 (RedHat/CentOS/Fedora)

1. **spec 文件**:
   ```spec
   Name:           cligool
   Version:        1.0.0
   Release:        1%{?dist}
   Summary:        Cross-platform remote terminal solution
   License:        MIT
   URL:            https://github.com/topcheer/cligool
   Source0:        %{name}-%{version}.tar.gz

   %description
   CliGool is a WebSocket-based remote terminal solution.

   %install
   mkdir -p %{buildroot}/usr/local/bin
   install -m 755 cligool %{buildroot}/usr/local/bin/

   %files
   /usr/local/bin/cligool

   %post
   # 安装后脚本

   %postun
   # 卸载后脚本
   ```

### Shell 安装脚本
```bash
#!/bin/bash
# 通用 Linux 安装脚本

set -e

echo "正在安装 CliGool..."

# 检测架构
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) BINARY="cligool-linux-amd64" ;;
    aarch64) BINARY="cligool-linux-arm64" ;;
    i386) BINARY="cligool-linux-386" ;;
    armv7l) BINARY="cligool-linux-arm" ;;
    *) echo "不支持的架构: $ARCH"; exit 1 ;;
esac

# 下载二进制文件
DOWNLOAD_URL="https://github.com/topcheer/cligool/releases/latest/download/${BINARY}"
echo "正在从 $DOWNLOAD_URL 下载..."

# 安装
sudo mkdir -p /usr/local/bin
sudo curl -L "$DOWNLOAD_URL" -o /usr/local/bin/cligool
sudo chmod +x /usr/local/bin/cligool

echo "✅ CliGool 安装完成！"
echo "运行 'cligool' 开始使用"
```

## 构建脚本

### build-all.sh
```bash
#!/bin/bash
# 构建所有安装包

set -e

VERSION=$(git describe --tags --always)
BUILD_DIR="build/installers"

echo "🔨 开始构建安装包 v${VERSION}..."

# 清理旧的构建
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 构建 Windows 安装包
echo "📦 构建 Windows 安装包..."
./installers/windows/build.sh

# 构建 macOS 安装包
echo "📦 构建 macOS 安装包..."
./installers/macos/build.sh

# 构建 Linux 安装包
echo "📦 构建 Linux 安装包..."
./installers/linux/build.sh

echo "✅ 所有安装包构建完成！"
echo "📂 输出目录: $BUILD_DIR"
```

## 安装流程设计

### Windows
1. 用户下载 cligool-setup.exe
2. 双击运行
3. **权限检测**:
   - 如果是管理员 → 提示"管理员模式"（可安装到 Program Files）
   - 如果是普通用户 → 提示"用户模式"（安装到用户目录）
4. 选择安装选项
5. 安装完成，提示打开

### macOS
1. 用户下载 cligool.pkg 或 cligool.dmg
2. **.pkg**: 双击运行，按向导安装
3. **.dmg**: 打开镜像，拖拽到应用程序
4. 安装完成，自动添加到 PATH

### Linux
1. **Debian/Ubuntu**: `sudo dpkg -i cligool.deb`
2. **RedHat/CentOS**: `sudo rpm -i cligool.rpm`
3. **通用**: `curl -sSL https://get.cligool.dev | bash`

## 卸载支持

### Windows
- 控制面板 → 程序和功能 → CliGool → 卸载
- 或运行安装目录中的 unins000.exe

### macOS
```bash
sudo rm -rf /usr/local/bin/cligool
sudo rm -rf /Library/Receipts/cligool.pkg
```

### Linux
```bash
# Debian/Ubuntu
sudo apt remove cligool

# RedHat/CentOS
sudo yum remove cligool

# 手动
sudo rm -rf /usr/local/bin/cligool
```

## 更新机制

1. **版本检查**:
   - 启动时检查 GitHub Releases
   - 提示有新版本可用

2. **自动更新**（可选）:
   - Windows: 下载新安装程序
   - macOS: 提示下载新版本
   - Linux: 包管理器更新

## 下一步

1. 实现 Windows Inno Setup 脚本
2. 实现 macOS Packages 项目
3. 实现 Linux .deb/.rpm 包
4. 创建构建自动化脚本
5. 测试所有安装包
6. 文档更新
