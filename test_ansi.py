#!/usr/bin/env python3
"""
测试 ANSI 转义序列的 Python 脚本
用于验证 PTY 和 WebSocket 能否正确处理复杂的终端输出
"""

import sys
import time

def test_ansi_sequences():
    """测试各种 ANSI 转义序列"""

    # 清屏
    print("\x1b[2J\x1b[H", end="", flush=True)

    # 测试基本颜色
    print("\x1b[31m红色\x1b[0m", flush=True)
    print("\x1b[32m绿色\x1b[0m", flush=True)
    print("\x1b[33m黄色\x1b[0m", flush=True)
    print("\x1b[34m蓝色\x1b[0m", flush=True)
    print("\x1b[1m粗体\x1b[0m", flush=True)
    print("")

    # 测试光标移动
    print("测试光标移动...")
    for i in range(5):
        print(f"\r计数: {i}/4", end="", flush=True)
        time.sleep(0.5)
    print("\r完成！     ", flush=True)
    print("")

    # 测试光标隐藏/显示
    print("测试光标隐藏/显示...")
    sys.stdout.write("\x1b[?25l")  # 隐藏光标
    sys.stdout.flush()
    time.sleep(1)
    sys.stdout.write("\x1b[?25h")  # 显示光标
    sys.stdout.flush()
    print(" 完成")
    print("")

    # 测试清除行
    print("测试清除行...")
    print("这行会被删除...")
    time.sleep(1)
    print("\x1b[1F\x1b[K", end="")  # 移到上一行并清除
    print("已经删除！")
    print("")

    # 测试屏幕滚动
    print("测试屏幕滚动区域...")
    print("\x1b[?25l", end="", flush=True)  # 隐藏光标以避免闪烁

    for i in range(10):
        print(f"行 {i+1}/10")
        time.sleep(0.2)

    time.sleep(1)
    print("\x1b[?25h", end="", flush=True)  # 显示光标
    print("")

    # 测试复杂格式
    print("测试复杂格式组合:")
    print("\x1b[1;31;40m 粗体红色黑底 \x1b[0m", flush=True)
    print("\x1b[4;32;40m 下划线绿色黑底 \x1b[0m", flush=True)
    print("\x1b[5;33;40m 闪烁黄色黑底 \x1b[0m", flush=True)
    print("\x1b[7;34;47m 反色蓝色白底 \x1b[0m", flush=True)
    print("")

    # 测试光标位置保存和恢复
    print("测试光标位置保存/恢复:")
    print("第一行")
    print("第二行")
    print("\x1b[s", end="", flush=True)  # 保存光标位置
    print("第三行")
    print("第四行")
    time.sleep(1)
    print("\x1b[u", end="", flush=True)  # 恢复光标位置
    print("<-- 这里恢复了光标位置")
    print("")

    # 测试终端查询（如果终端支持，应该有响应）
    print("测试终端能力查询 (DA1):")
    sys.stdout.write("\x1b[c")
    sys.stdout.flush()
    time.sleep(0.5)
    print("")
    print("(如果终端支持，上面应该返回能力代码)")
    print("")

    print("=== 所有测试完成 ===")

if __name__ == "__main__":
    try:
        test_ansi_sequences()
    except KeyboardInterrupt:
        print("\n测试被中断")
        sys.exit(0)
