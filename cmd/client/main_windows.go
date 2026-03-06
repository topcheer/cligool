// +build windows

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
	"unicode/utf8"

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

	// 启动WebSocket并运行终端会话
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
	fmt.Println("⚠️  Windows模式：功能可能受限")
	fmt.Println()

	// Windows上使用cmd.exe而不是PTY
	cmd := exec.Command("cmd.exe")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// 创建输入输出管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建stdin管道失败: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %w", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %w", err)
	}
	defer cmd.Process.Kill()

	// WebSocket -> Stdin (网页输入写入命令)
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
				stdin.Write([]byte(msg.Data))
			}
		}
	}()

	// 合并stdout和stderr
	go func() {
		io.Copy(stdin, stdout)
	}()

	// Stdout/Stderr -> WebSocket (命令输出发送到网页)
	buf := make([]byte, 128)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				// 尝试读取stderr
				n, err := stderr.Read(buf)
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return fmt.Errorf("stderr读取失败: %w", err)
				}

				data := buf[:n]
				if !utf8.Valid(data) {
					continue
				}

				msg := TerminalMessage{
					Type:    "output",
					Data:    string(data),
					Session: sessionID,
					UserID:  "client",
				}

				jsonData, _ := json.Marshal(msg)
				if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					return fmt.Errorf("WebSocket写入失败: %w", err)
				}
				continue
			}
			return fmt.Errorf("stdout读取失败: %w", err)
		}

		// 确保数据是有效的UTF-8
		data := buf[:n]
		if !utf8.Valid(data) {
			continue
		}

		// 发送到WebSocket
		msg := TerminalMessage{
			Type:    "output",
			Data:    string(data),
			Session: sessionID,
			UserID:  "client",
		}

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
