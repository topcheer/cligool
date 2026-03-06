package terminal

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// TerminalSize 终端大小
type TerminalSize struct {
	Rows uint16
	Cols uint16
}

// ShellConfig Shell配置
type ShellConfig struct {
	Path    string
	Env     []string
	Timeout int // 超时时间（秒）
}

// DefaultShellConfig 默认Shell配置
func DefaultShellConfig() *ShellConfig {
	return &ShellConfig{
		Path:    getDefaultShell(),
		Env:     defaultEnvVars(),
		Timeout: 0, // 无超时
	}
}

// CreateShellCommand 创建shell命令
func CreateShellCommand(shellPath string) *exec.Cmd {
	config := DefaultShellConfig()
	if shellPath != "" {
		config.Path = shellPath
	}

	cmd := exec.Command(config.Path)
	cmd.Env = append(os.Environ(), config.Env...)
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")

	// 设置标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// getDefaultShell 获取默认shell路径
func getDefaultShell() string {
	// 检查SHELL环境变量
	if shell := os.Getenv("SHELL"); shell != "" {
		if _, err := exec.LookPath(shell); err == nil {
			return shell
		}
	}

	// 根据操作系统返回默认shell
	switch runtime.GOOS {
	case "windows":
		if _, err := exec.LookPath("powershell.exe"); err == nil {
			return "powershell.exe"
		}
		return "cmd.exe"
	case "darwin":
		// macOS优先使用zsh
		if _, err := exec.LookPath("/bin/zsh"); err == nil {
			return "/bin/zsh"
		}
		return "/bin/bash"
	case "linux":
		// Linux常见shells
		shells := []string{"/bin/bash", "/bin/zsh", "/bin/sh", "/usr/bin/zsh"}
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

// defaultEnvVars 默认环境变量
func defaultEnvVars() []string {
	envVars := []string{
		"TERM=xterm-256color",
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	}

	// 添加常用工具路径
	path := os.Getenv("PATH")
	if path == "" {
		path = "/usr/local/bin:/usr/bin:/bin:/usr/local/sbin:/usr/sbin:/sbin"
	}
	envVars = append(envVars, "PATH="+path)

	return envVars
}

// GetTerminalSize 获取终端大小
func GetTerminalSize() (*TerminalSize, error) {
	// 这里需要使用系统调用获取终端大小
	// 具体实现依赖于操作系统
	return nil, nil
}

// IsSupportedShell 检查是否为支持的shell
func IsSupportedShell(shell string) bool {
	supportedShells := []string{
		"sh", "bash", "zsh", "fish", "powershell", "cmd",
		"/bin/sh", "/bin/bash", "/bin/zsh", "/usr/bin/zsh",
		"powershell.exe", "cmd.exe",
	}

	shell = strings.ToLower(shell)
	shell = strings.TrimSuffix(shell, ".exe")

	for _, supported := range supportedShells {
		if strings.HasSuffix(shell, supported) {
			return true
		}
	}

	return false
}

// IsWindowsOS 判断是否为Windows系统
func IsWindowsOS() bool {
	return runtime.GOOS == "windows"
}

// IsUnixOS 判断是否为Unix系统
func IsUnixOS() bool {
	return strings.Contains(runtime.GOOS, "darwin") ||
		strings.Contains(runtime.GOOS, "linux") ||
		strings.Contains(runtime.GOOS, "bsd") ||
		strings.Contains(runtime.GOOS, "unix")
}