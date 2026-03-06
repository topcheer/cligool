#!/bin/bash
# 测试终端能力

echo "=== 终端能力测试 ==="
echo ""

echo "1. 检查 TERM 环境变量:"
echo "   TERM=$TERM"
echo ""

echo "2. 检查终端大小:"
echo "   COLUMNS=$COLUMNS"
echo "   LINES=$LINES"
echo ""

echo "3. 测试 ANSI 颜色:"
echo -e "\x1b[31m红色文本\x1b[0m"
echo -e "\x1b[32m绿色文本\x1b[0m"
echo -e "\x1b[33m黄色文本\x1b[0m"
echo -e "\x1b[34m蓝色文本\x1b[0m"
echo -e "\x1b[1m粗体文本\x1b[0m"
echo -e "\x1b[4m下划线文本\x1b[0m"
echo ""

echo "4. 测试光标控制:"
echo -e "清除屏幕:\x1b[2J"
echo -e "移动光标到首页:\x1b[H"
echo "   请按回车继续..."
read

echo "5. 测试终端设备查询 (DA1):"
echo -ne "\x1b[c"
sleep 0.1
echo ""
echo "   (如果终端支持，应该返回能力响应)"
echo ""

echo "6. 测试光标位置查询 (CPR):"
echo -ne "\x1b[6n"
sleep 0.1
echo ""
echo "   (如果终端支持，应该返回光标位置)"
echo ""

echo "7. 测试复杂格式:"
echo -e "\x1b[?25l隐藏光标"
sleep 1
echo -e "\x1b[?25h显示光标"
echo ""

echo "8. 测试替代屏幕模式:"
echo -e "进入替代屏幕...\x1b[?1049h"
echo "这是替代屏幕内容"
sleep 2
echo -e "退出替代屏幕...\x1b[?1049l"
echo "已返回主屏幕"
echo ""

echo "=== 测试完成 ==="
