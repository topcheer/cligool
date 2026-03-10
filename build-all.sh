#!/bin/bash
# 构建所有受支持平台的客户端 + relay

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🔨 开始构建所有平台的客户端..."
echo "📊 目标平台：29 个客户端平台 + 1 个中继服务器"
echo ""

mkdir -p bin

CLIENT_TARGETS=(
  "Windows amd64|windows|amd64|bin/cligool-windows-amd64.exe||"
  "Windows arm64|windows|arm64|bin/cligool-windows-arm64.exe||"
  "macOS amd64|darwin|amd64|bin/cligool-darwin-amd64||"
  "macOS arm64|darwin|arm64|bin/cligool-darwin-arm64||"
  "Linux amd64|linux|amd64|bin/cligool-linux-amd64||"
  "Linux arm64|linux|arm64|bin/cligool-linux-arm64||"
  "Linux 386|linux|386|bin/cligool-linux-386||"
  "Linux arm|linux|arm|bin/cligool-linux-arm|GOARM|6"
  "Linux armbe (ARM64 Big-Endian)|linux|arm64|bin/cligool-linux-armbe|GOARM|7"
  "Linux ppc64le|linux|ppc64le|bin/cligool-linux-ppc64le||"
  "Linux ppc64 (Big-Endian)|linux|ppc64|bin/cligool-linux-ppc64|GOBIGENDIAN|true"
  "Linux riscv64|linux|riscv64|bin/cligool-linux-riscv64||"
  "Linux s390x|linux|s390x|bin/cligool-linux-s390x||"
  "Linux mips|linux|mips|bin/cligool-linux-mips||"
  "Linux mips64le|linux|mips64le|bin/cligool-linux-mips64le||"
  "Linux mips64|linux|mips64|bin/cligool-linux-mips64|GOMIPS|hardfloat"
  "Linux loong64|linux|loong64|bin/cligool-linux-loong64||"
  "FreeBSD amd64|freebsd|amd64|bin/cligool-freebsd-amd64||"
  "FreeBSD arm64|freebsd|arm64|bin/cligool-freebsd-arm64||"
  "FreeBSD 386|freebsd|386|bin/cligool-freebsd-386||"
  "FreeBSD arm|freebsd|arm|bin/cligool-freebsd-arm|GOARM|6"
  "FreeBSD riscv64|freebsd|riscv64|bin/cligool-freebsd-riscv64||"
  "OpenBSD amd64|openbsd|amd64|bin/cligool-openbsd-amd64||"
  "OpenBSD arm64|openbsd|arm64|bin/cligool-openbsd-arm64||"
  "NetBSD amd64|netbsd|amd64|bin/cligool-netbsd-amd64||"
  "NetBSD arm64|netbsd|arm64|bin/cligool-netbsd-arm64||"
  "NetBSD arm|netbsd|arm|bin/cligool-netbsd-arm|GOARM|6"
  "NetBSD 386|netbsd|386|bin/cligool-netbsd-386||"
  "DragonFlyBSD amd64|dragonfly|amd64|bin/cligool-dragonfly-amd64||"
)

TOTAL_STEPS=$((${#CLIENT_TARGETS[@]} + 1))
STEP=1
GENERATED_OUTPUTS=()

build_client_target() {
  local label="$1"
  local goos="$2"
  local goarch="$3"
  local output="$4"
  local extra_key="$5"
  local extra_value="$6"

  echo "📦 [${STEP}/${TOTAL_STEPS}] 构建 ${label}..."
  if [[ -n "$extra_key" ]]; then
    env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" "${extra_key}=${extra_value}" \
      go build -o "$output" ./cmd/client
  else
    env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -o "$output" ./cmd/client
  fi
  GENERATED_OUTPUTS+=("$output")
  STEP=$((STEP + 1))
}

for target in "${CLIENT_TARGETS[@]}"; do
  IFS='|' read -r label goos goarch output extra_key extra_value <<<"$target"
  build_client_target "$label" "$goos" "$goarch" "$output" "$extra_key" "$extra_value"
done

echo "📦 [${STEP}/${TOTAL_STEPS}] 构建中继服务器 (Linux amd64)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/relay-server ./cmd/relay
GENERATED_OUTPUTS+=("bin/relay-server")

echo ""
echo "✅ 构建完成！"
echo ""
echo "📊 构建统计："
echo "   - Windows: 2 个平台"
echo "   - macOS: 2 个平台"
echo "   - Linux: 13 个平台"
echo "   - FreeBSD: 5 个平台"
echo "   - OpenBSD: 2 个平台 (仅 amd64 和 arm64)"
echo "   - NetBSD: 4 个平台"
echo "   - DragonFlyBSD: 1 个平台 (仅 amd64)"
echo "   - 中继服务器: 1 个"
echo "   总计: 30 个二进制文件"
echo ""
echo "⚠️  平台限制说明："
echo "   - OpenBSD 386/arm/riscv64: pty 库限制"
echo "   - DragonFlyBSD arm64: Go 不支持"
echo ""
echo "📦 构建的文件："
for output in "${GENERATED_OUTPUTS[@]}"; do
  ls -lh "$output"
done
echo ""
echo "💡 提示："
echo "   - build-all.sh 会直接交叉编译 Windows / macOS / Linux / *BSD 版本"
echo "   - 如只需 Windows 产物，可单独运行 ./build-windows.sh"
echo "   - 所有二进制文件都是静态编译，无需额外依赖"
echo "   - 可以直接复制到目标系统运行"
