// +build !windows

package main

import (
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// URLNotifier 检测 URL 并发送系统通知
type URLNotifier struct {
	urlRegex  *regexp.Regexp
	lastURL   string
	enabled   bool
}

// NewURLNotifier 创建 URL 通知器
func NewURLNotifier(enabled bool) *URLNotifier {
	return &URLNotifier{
		urlRegex: regexp.MustCompile(`https?://[^\s"'<>]+`),
		enabled:  enabled,
	}
}

// Process 处理输出数据，检测 URL 并发送通知
func (un *URLNotifier) Process(data []byte) {
	if !un.enabled {
		return
	}

	// 检测 URL
	matches := un.urlRegex.FindAll(data, -1)

	for _, match := range matches {
		url := string(match)

		// 只通知公网 URL（不是 localhost 或私有 IP）
		if un.isPublicURL(url) && url != un.lastURL {
			un.lastURL = url
			un.sendNotification(url)
		}
	}
}

// isPublicURL 检查是否是公网 URL
func (un *URLNotifier) isPublicURL(url string) bool {
	// 排除 localhost 和私有 IP
	excludePatterns := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"[::1]",
	}

	lowerURL := strings.ToLower(url)
	for _, pattern := range excludePatterns {
		if strings.Contains(lowerURL, pattern) {
			return false
		}
	}

	// 简单检查：http:// 或 https:// 开头
	return strings.HasPrefix(lowerURL, "http://") || strings.HasPrefix(lowerURL, "https://")
}

// sendNotification 发送系统通知
func (un *URLNotifier) sendNotification(url string) {
	switch runtime.GOOS {
	case "darwin":
		un.sendMacOSNotification(url)
	case "linux":
		un.sendLinuxNotification(url)
	default:
		log.Printf("检测到 URL (系统通知不支持): %s", url)
	}
}

// sendMacOSNotification 发送 macOS 通知
func (un *URLNotifier) sendMacOSNotification(url string) {
	// 使用 osascript 显示通知
	script := `display notification "` + url + `" with title "🔗 CliGool 检测到链接" sound name "Glass"`
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Printf("发送 macOS 通知失败: %v", err)
		// 降级：使用终端提示
		log.Printf("🔗 [URL] %s", url)
	} else {
		log.Printf("✅ 已发送系统通知: %s", url)
	}
}

// sendLinuxNotification 发送 Linux 通知
func (un *URLNotifier) sendLinuxNotification(url string) {
	// 检查是否安装了 notify-send
	cmd := exec.Command("which", "notify-send")
	if err := cmd.Run(); err != nil {
		// notify-send 不可用，使用终端提示
		log.Printf("🔗 [URL] %s", url)
		return
	}

	// 使用 notify-send 发送通知
	notifyCmd := exec.Command("notify-send",
		"🔗 CliGool 检测到链接",
		url,
		"-u", "normal",
		"-i", "dialog-information",
		"-t", "5000")

	if err := notifyCmd.Run(); err != nil {
		log.Printf("发送 Linux 通知失败: %v", err)
		log.Printf("🔗 [URL] %s", url)
	} else {
		log.Printf("✅ 已发送系统通知: %s", url)
	}
}

// SendDirectNotification 直接发送通知（用于特殊情况）
func (un *URLNotifier) SendDirectNotification(url string) {
	if un.enabled && un.isPublicURL(url) {
		un.lastURL = url
		un.sendNotification(url)
	}
}
