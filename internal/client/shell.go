package client

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Cmd 表示外部命令
type Cmd struct {
	*exec.Cmd
}

// CreateShellCommand 创建shell命令
func CreateShellCommand(shellPath string) *Cmd {
	if shellPath == "" {
		shellPath = getDefaultShell()
	}

	// 检查shell是否存在
	if _, err := exec.LookPath(shellPath); err != nil {
		// 如果指定的shell不存在，使用默认shell
		shellPath = getDefaultShell()
	}

	cmd := exec.Command(shellPath)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// 设置标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Cmd{Cmd: cmd}
}

// getDefaultShell 获取系统默认shell
func getDefaultShell() string {
	// 首先检查SHELL环境变量
	if shell := os.Getenv("SHELL"); shell != "" {
		if _, err := exec.LookPath(shell); err == nil {
			return shell
		}
	}

	// 根据操作系统选择默认shell
	switch runtime.GOOS {
	case "windows":
		return "cmd.exe"
	case "darwin", "linux":
		// 尝试常见的Unix shells
		shells := []string{"/bin/zsh", "/bin/bash", "/bin/sh"}
		for _, shell := range shells {
			if _, err := exec.LookPath(shell); err == nil {
				return shell
			}
		}
		return "/bin/sh"
	default:
		return "/bin/sh"
	}
}

// IsWindows 判断是否为Windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsUnix 判断是否为Unix系统
func IsUnix() bool {
	return strings.Contains(runtime.GOOS, "darwin") ||
		strings.Contains(runtime.GOOS, "linux") ||
		strings.Contains(runtime.GOOS, "bsd")
}