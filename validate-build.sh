#!/bin/bash
# 测试构建脚本 - 快速测试几个关键平台

echo "🧪 快速构建测试"
echo "=================="
echo ""

# 创建测试目录
TEST_DIR="test-builds"
rm -rf $TEST_DIR
mkdir -p $TEST_DIR

# 测试平台列表
declare -a platforms=(
    "darwin:amd64:macOS Intel"
    "darwin:arm64:macOS ARM"
    "linux:amd64:Linux x86_64"
    "linux:arm64:Linux ARM64"
    "linux:loong64:Linux LoongArch"
    "freebsd:amd64:FreeBSD x86_64"
    "freebsd:riscv64:FreeBSD RISC-V"
)

success_count=0
fail_count=0

for platform in "${platforms[@]}"; do
    IFS=':' read -r goos goarch description <<< "$platform"

    echo "📦 测试构建: $description ($goos/$goarch)"

    output_file="$TEST_DIR/cligool-$goos-$goarch"

    if GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build -o "$output_file" ./cmd/client 2>/dev/null; then
        if [ -f "$output_file" ]; then
            size=$(ls -lh "$output_file" | awk '{print $5}')
            echo "   ✅ 成功 ($size)"
            ((success_count++))
        else
            echo "   ❌ 失败 - 文件未生成"
            ((fail_count++))
        fi
    else
        echo "   ❌ 失败 - 编译错误"
        ((fail_count++))
    fi
done

echo ""
echo "=================="
echo "📊 测试结果:"
echo "   成功: $success_count/${#platforms[@]}"
echo "   失败: $fail_count/${#platforms[@]}"

if [ $fail_count -eq 0 ]; then
    echo "   🎉 所有测试通过！"
    echo ""
    echo "📦 架构验证："
    for file in $TEST_DIR/*; do
        if [ -f "$file" ]; then
            echo ""
            echo "文件: $(basename $file)"
            file "$file" | head -1
        fi
    done
    exit 0
else
    echo "   ⚠️  存在失败的构建"
    exit 1
fi
