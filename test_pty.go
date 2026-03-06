package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
)

func main() {
	// 创建一个简单的命令
	cmd := exec.Command("echo", "test")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// 尝试启动PTY
	ptyMaster, err := pty.Start(cmd)
	if err != nil {
		log.Printf("PTY启动失败: %v", err)

		// 检查具体错误类型
		if pathErr, ok := err.(*os.PathError); ok {
			log.Printf("路径错误: %s, Op: %s", pathErr.Path, pathErr.Op)
			if pathErr.Err == syscall.EPERM {
				log.Println("错误: 权限不足 (EPERM)")
				log.Println("可能的原因:")
				log.Println("1. 在沙盒环境中运行")
				log.Println("2. 终端权限不足")
				log.Println("3. 需要在真实终端中运行")
			}
		}

		// 尝试不用PTY的方式
		log.Println("尝试不使用PTY运行命令...")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("命令执行失败: %v", err)
		}
		fmt.Printf("输出: %s\n", output)
		return
	}
	defer ptyMaster.Close()

	log.Println("PTY启动成功！")

	// 等待命令完成
	cmd.Wait()
}