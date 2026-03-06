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
	script := `display notification "点击打开 Web 终端" with title "🌐 CliGool Web 终端已就绪" sound name "Glass" subtitle "` + url + `"`
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Printf("发送 macOS 通知失败: %v", err)
		return err
	}

	log.Printf("✅ 已发送系统通知: %s", url)
	return nil
}

// sendLinuxNotification 发送 Linux 通知
func (sn *SystemNotifier) sendLinuxNotification(url string) error {
	// 检查是否安装了 notify-send
	cmd := exec.Command("which", "notify-send")
	if err := cmd.Run(); err != nil {
		log.Printf("notify-send 不可用，跳过系统通知")
		return nil
	}

	notifyCmd := exec.Command("notify-send",
		"🌐 CliGool Web 终端已就绪",
		url,
		"-u", "normal",
		"-i", "web-browser",
		"-t", "10000") // 显示 10 秒

	if err := notifyCmd.Run(); err != nil {
		log.Printf("发送 Linux 通知失败: %v", err)
		return err
	}

	log.Printf("✅ 已发送系统通知: %s", url)
	return nil
}
