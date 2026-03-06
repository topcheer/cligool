// +build !windows

package main

import (
	"log"
	"os/exec"
	"runtime"
)

// SystemNotifier 系统通知器
type SystemNotifier struct{}

// NewSystemNotifier 创建系统通知器
func NewSystemNotifier() *SystemNotifier {
	return &SystemNotifier{}
}

// SendWebTerminalNotification 发送 Web 终端 URL 通知
func (sn *SystemNotifier) SendWebTerminalNotification(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return sn.sendMacOSNotification(url)
	case "linux":
		return sn.sendLinuxNotification(url)
	default:
		log.Printf("系统通知不支持此平台: %s", runtime.GOOS)
		return nil
	}
}

// sendMacOSNotification 发送 macOS 通知
func (sn *SystemNotifier) sendMacOSNotification(url string) error {
	// 发送通知，显示 URL
	script := `display notification "Web 终端已就绪，URL 已复制到剪贴板" with title "🌐 CliGool Web 终端" sound name "Glass" subtitle "` + url + `"`
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Printf("发送 macOS 通知失败: %v", err)
		return err
	}

	// 同时自动打开浏览器
	return sn.openBrowser(url)
}

// sendLinuxNotification 发送 Linux 通知
func (sn *SystemNotifier) sendLinuxNotification(url string) error {
	// 检查是否安装了 notify-send
	cmd := exec.Command("which", "notify-send")
	if err := cmd.Run(); err != nil {
		log.Printf("notify-send 不可用，跳过系统通知")
		// 直接打开浏览器
		return sn.openBrowser(url)
	}

	notifyCmd := exec.Command("notify-send",
		"🌐 CliGool Web 终端",
		url+" (正在打开浏览器...)",
		"-u", "normal",
		"-i", "web-browser",
		"-t", "5000") // 显示 5 秒

	// 发送通知（不等待）
	go notifyCmd.Run()

	// 同时自动打开浏览器
	return sn.openBrowser(url)
}

// openBrowser 直接打开浏览器
func (sn *SystemNotifier) openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS: 使用 open 命令
		cmd = exec.Command("open", url)
	case "linux":
		// Linux: 使用 xdg-open
		cmd = exec.Command("xdg-open", url)
	default:
		log.Printf("自动打开浏览器不支持此平台: %s", runtime.GOOS)
		return nil
	}

	if err := cmd.Run(); err != nil {
		log.Printf("打开浏览器失败: %v", err)
		return err
	}

	log.Printf("✅ 已在浏览器中打开: %s", url)
	return nil
}
