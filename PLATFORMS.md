# CliGool 平台支持

本版本支持18个操作系统和架构组合。

## 🌟 已构建的平台

### macOS (2个)
- ✅ **darwin-amd64** - Intel处理器（Mac mini 2018, MacBook Pro 2019等）
- ✅ **darwin-arm64** - Apple Silicon（M1/M2/M3）

### Linux (8个)
- ✅ **linux-amd64** - x86_64（Ubuntu, Debian, CentOS等）
- ✅ **linux-arm64** - ARM 64位（树莓派4/5, ARM服务器）
- ✅ **linux-386** - x86 32位（老式PC）
- ✅ **linux-arm** - ARM 32位（树莓派Zero/2B/3B）
- ✅ **linux-ppc64le** - PowerPC Little-Endian
- ✅ **linux-riscv64** - RISC-V 64位
- ✅ **linux-s390x** - IBM System z
- ✅ **linux-mips64le** - MIPS 64位 Little-Endian

### *BSD系统 (6个)
- ✅ **freebsd-amd64** - FreeBSD 64位
- ✅ **freebsd-arm64** - FreeBSD ARM 64位
- ✅ **openbsd-amd64** - OpenBSD 64位
- ✅ **openbsd-arm64** - OpenBSD ARM 64位
- ✅ **netbsd-amd64** - NetBSD 64位
- ✅ **dragonfly-amd64** - DragonFlyBSD 64位

### Windows (2个)
- ✅ **windows-amd64** - Windows 64位（ZIP格式）
- ✅ **windows-arm64** - Windows ARM 64位（Surface Pro X等）

## 🚀 快速开始

### 1. 构建所有平台
```bash
./build-all.sh
```

### 2. 测试本地客户端
```bash
# 自动检测当前平台
./test-local.sh [会话ID] [服务器URL]

# 或者手动指定
./bin/cligool-darwin-arm64 -server http://localhost:8081
```

### 3. 下载页面
所有平台的二进制文件都已打包到Docker镜像中，可以通过以下方式获取：
- 访问 http://localhost:8081/ 下载页面
- 或从Docker容器中复制：`docker cp cligool-relay:/app/web/downloads ./`

## 📦 二进制文件列表

| 平台 | 文件名 | 大小 |
|------|--------|------|
| macOS Intel | `cligool-darwin-amd64` | ~8.2 MB |
| macOS ARM | `cligool-darwin-arm64` | ~7.7 MB |
| Linux amd64 | `cligool-linux-amd64` | ~8.1 MB |
| Linux arm64 | `cligool-linux-arm64` | ~7.5 MB |
| Linux 386 | `cligool-linux-386` | ~7.8 MB |
| Linux arm | `cligool-linux-arm` | ~7.8 MB |
| Linux ppc64le | `cligool-linux-ppc64le` | ~7.9 MB |
| Linux riscv64 | `cligool-linux-riscv64` | ~7.4 MB |
| Linux s390x | `cligool-linux-s390x` | ~8.3 MB |
| Linux mips64le | `cligool-linux-mips64le` | ~8.4 MB |
| FreeBSD amd64 | `cligool-freebsd-amd64` | ~8.0 MB |
| FreeBSD arm64 | `cligool-freebsd-arm64` | ~7.4 MB |
| OpenBSD amd64 | `cligool-openbsd-amd64` | ~8.0 MB |
| OpenBSD arm64 | `cligool-openbsd-arm64` | ~7.4 MB |
| NetBSD amd64 | `cligool-netbsd-amd64` | ~7.9 MB |
| DragonFlyBSD amd64 | `cligool-dragonfly-amd64` | ~8.0 MB |
| Windows amd64 | `cligool-windows-amd64.zip` | ~4.1 MB (ZIP) |
| Windows arm64 | `cligool-windows-arm64.zip` | ~3.8 MB (ZIP) |

## 🎯 新功能：双端同步

本版本实现了**本地终端和Web终端双端同步**：

### 特性
- ✅ 本地终端实时显示所有输出
- ✅ 本地终端可以直接输入命令
- ✅ Web终端同步显示所有内容
- ✅ Web终端也可以输入命令
- ✅ 两边完全同步，互不干扰

### 使用场景
1. **本地工作**：在本地终端操作，Web端查看演示
2. **远程协作**：在Web端操作，本地终端监控
3. **教学演示**：一边操作一边给其他人看
4. **技术支持**：本地查看，远程操作

## 🔧 技术细节

### 构建方式
- **交叉编译**：使用Go的交叉编译功能
- **无依赖部署**：所有二进制都是静态编译
- **小体积**：优化后每个文件约7-8MB

### 系统要求
- Go 1.21+（用于构建）
- 对应的操作系统（用于运行）

### 编码支持
- **Unix/Linux/macOS**：UTF-8原生支持
- **Windows**：GBK自动转换为UTF-8

## 📝 许可证
MIT License - 详见 [LICENSE](LICENSE) 文件
