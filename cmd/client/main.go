package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unicode/utf8"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type TerminalMessage struct {
	Type    string `json:"type"`
	Data    string `json:"data"`
	Session string `json:"session"`
	UserID  string `json:"user_id"`
}

func main() {
	serverURL := flag.String("server", "https://cligool.zty8.cn", "中继服务器URL")
	sessionID := flag.String("session", "", "会话ID")
	flag.Parse()

	sid := *sessionID
	if sid == "" {
		sid = uuid.New().String()
	}

	// 显示连接信息
	printHeader(sid, *serverURL)

	// 启动WebSocket并运行PTY
	if err := runTerminalSession(*serverURL, sid); err != nil {
		log.Fatalf("终端会话失败: %v", err)
	}
}

func printHeader(sessionID, serverURL string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🚀 CliGool 远程终端                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Printf("║ 📋 会话ID: %-43s ║\n", sessionID)
	fmt.Printf("║ 🌐 Web访问: %-43s ║\n", serverURL+"/session/"+sessionID)
	fmt.Printf("║ 🔗 连接状态: %-43s ║\n", "🟡 连接中...")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func runTerminalSession(serverURL, sessionID string) error {
	// 建立 WebSocket 连接
	wsURL, _ := buildWebSocketURL(serverURL, sessionID)
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("WebSocket连接失败: %w", err)
	}
	defer conn.Close()

	log.Println("✅ WebSocket已连接")
	fmt.Println("✅ 已连接到中继服务器")
	fmt.Println("💡 现在可以在Web终端中输入命令了")
	fmt.Println()

	// 创建 PTY
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	// 创建命令
	cmd := exec.Command(shell, "-i", "-l")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// 启动PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("启动PTY失败: %w", err)
	}
	defer ptmx.Close()

	// 处理窗口大小变化
	handleResize := func() {
		// 这里可以添加窗口大小调整逻辑
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			handleResize()
		}
	}()

	// WebSocket -> PTY (网页输入写入PTY)
	// 使用单独的 goroutine 确保输入立即处理
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket读取失败: %v", err)
				return
			}

			var msg TerminalMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			if msg.Type == "input" && msg.Data != "" {
				// 立即写入PTY，不缓冲
				if _, err := ptmx.Write([]byte(msg.Data)); err != nil {
					log.Printf("PTY写入失败: %v", err)
					continue
				}
			} else if msg.Type == "resize" {
				// 处理窗口大小调整
				// TODO: 实现PTY窗口大小调整
			}
		}
	}()

	// PTY -> WebSocket (PTY输出发送到网页)
	// 使用更小的缓冲区以减少延迟
	buf := make([]byte, 128) // 从1024减少到128字节以减少延迟
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("PTY读取失败: %w", err)
		}

		// 确保数据是有效的UTF-8
		data := buf[:n]
		if !utf8.Valid(data) {
			// 跳过无效的UTF-8序列
			continue
		}

		// 立即发送到WebSocket，不等待更多数据
		msg := TerminalMessage{
			Type:    "output",
			Data:    string(data),
			Session: sessionID,
			UserID:  "client",
		}

		// 使用WriteMessage而不是WriteJSON以提高性能
		jsonData, _ := json.Marshal(msg)
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			return fmt.Errorf("WebSocket写入失败: %w", err)
		}
	}
}

func buildWebSocketURL(serverURL, sessionID string) (string, error) {
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return "", err
	}

	scheme := "ws"
	if parsedURL.Scheme == "https" {
		scheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/api/terminal/%s?type=client&user_id=client",
		scheme, parsedURL.Host, sessionID)
	return wsURL, nil
}
