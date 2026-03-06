// +build windows

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type TerminalMessage struct {
	Type        string `json:"type"`
	Data        string `json:"data"`
	Session     string `json:"session"`
	UserID      string `json:"user_id"`
	Source      string `json:"source,omitempty"` // "local" or "web"
	WorkingDir  string `json:"working_dir,omitempty"`
	OSInfo      string `json:"os_info,omitempty"`
}

func main() {
	log.Println("🚀 CliGool Windows客户端启动...")

	serverURL := flag.String("server", "https://cligool.zty8.cn", "中继服务器URL")
	sessionID := flag.String("session", "", "会话ID")
	flag.Parse()

	sid := *sessionID
	if sid == "" {
		sid = uuid.New().String()
	}

	log.Println("📋 会话ID:", sid)

	// 显示连接信息
	printHeader(sid, *serverURL)

	log.Println("🔗 开始连接WebSocket...")

	// 启动WebSocket并运行终端会话
	if err := runTerminalSession(*serverURL, sid); err != nil {
		log.Printf("❌ 终端会话失败: %v", err)
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
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
	log.Println("🔧 开始建立WebSocket连接...")

	// 建立 WebSocket 连接
	wsURL, _ := buildWebSocketURL(serverURL, sessionID)
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Println("📡 WebSocket URL:", wsURL)

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

	// 创建WebSocket写入channel，确保串行写入
	wsWriteChan := make(chan []byte, 100)

	// 启动WebSocket写入goroutine
	go func() {
		for data := range wsWriteChan {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("❌ WebSocket写入失败: %v", err)
				return
			}
		}
	}()

	// 设置ping handler和心跳
	// 注意：心跳是控制帧，直接发送，不通过channel
	setupHeartbeat(conn)

	// 发送初始化消息（工作目录和系统信息）
	wd, _ := os.Getwd()
	initMsg := TerminalMessage{
		Type:        "init",
		Session:     sessionID,
		UserID:      "client",
		WorkingDir:  wd,
		OSInfo:      "windows",
	}
	jsonData, _ := json.Marshal(initMsg)
	wsWriteChan <- jsonData
	log.Printf("✅ 已发送初始化消息: 工作目录=%s", wd)

	log.Println("🔧 准备启动cmd.exe...")

	// Windows上使用cmd.exe而不是PTY
	cmd := exec.Command("cmd.exe")
	// 使用交互模式以获得完整的终端体验
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"CMD_QUIT=exit", // 退出命令
	)
	// 禁用命令行参数处理，保持交互模式
	cmd.Args = []string{"cmd.exe"}

	log.Println("📝 创建输入管道...")
	// 创建输入输出管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("创建stdin管道失败: %w", err)
	}

	log.Println("📝 创建输出管道...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %w", err)
	}

	log.Println("🚀 启动cmd.exe进程...")
	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("❌ cmd.exe启动失败: %v", err)
		return fmt.Errorf("启动命令失败: %w", err)
	}
	defer cmd.Process.Kill()

	log.Println("✅ cmd.exe已启动，PID:", cmd.Process.Pid)

	// 使用带缓冲的写入器并确保每次写入都flush
	stdinWriter := bufio.NewWriter(stdin)

	// 本地终端输入 -> Stdin
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Printf("本地stdin读取失败: %v", err)
				continue
			}

			// 写入cmd.exe
			_, writeErr := stdinWriter.Write(buf[:n])
			if writeErr != nil {
				log.Printf("❌ 写入stdin失败（本地输入）: %v", writeErr)
				continue
			}

			if err := stdinWriter.Flush(); err != nil {
				log.Printf("❌ Flush失败（本地输入）: %v", err)
			}

			// 同时发送到WebSocket，让Web端看到本地输入
			msg := TerminalMessage{
				Type:    "input",
				Data:    string(buf[:n]),
				Session: sessionID,
				UserID:  "client",
				Source:  "local",
			}
			jsonData, _ := json.Marshal(msg)
			wsWriteChan <- jsonData
		}
	}()

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
				log.Printf("JSON解析失败: %v", err)
				continue
			}

			if msg.Type == "input" && msg.Data != "" {
				// Windows cmd.exe需要CRLF换行符
				data := []byte(msg.Data)
				// 将单独的"\r"转换为"\r\n" (Windows标准换行符)
				if msg.Data == "\r" {
					data = []byte("\r\n")
				}

				_, err := stdinWriter.Write(data)
				if err != nil {
					log.Printf("❌ 写入stdin失败: %v", err)
				} else {
					// 立即flush确保cmd.exe收到输入
					if err := stdinWriter.Flush(); err != nil {
						log.Printf("❌ Flush失败: %v", err)
					}
				}
			}
		}
	}()

	// stderr读取（发送到WebSocket和本地stderr）
	go func() {
		// 使用更大的缓冲区
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					// log.Printf("stderr读取失败: %v", err)
				}
				return
			}

			data := buf[:n]
			// Windows cmd.exe使用GBK编码，需要转换为UTF-8
			converted, err := convertGBKToUTF8(data)
			if err != nil {
				// 如果转换失败，使用原始数据
				converted = string(data)
			}

			// 1. 显示到本地终端stderr（UTF-8数据）
			os.Stderr.Write([]byte(converted))

			// 2. 发送到WebSocket
			msg := TerminalMessage{
				Type:    "output",
				Data:    converted,
				Session: sessionID,
				UserID:  "client",
			}

			jsonData, _ := json.Marshal(msg)
			wsWriteChan <- jsonData
		}
	}()

	// Stdout -> 本地stdout + WebSocket (命令输出同时显示在两端)
	// 使用更大的缓冲区以处理cmd.exe的大量输出
	buf := make([]byte, 4096)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("stdout读取失败: %w", err)
		}

		// Windows cmd.exe使用GBK编码，需要转换为UTF-8
		data := buf[:n]

		// 转换为UTF-8
		converted, err := convertGBKToUTF8(data)
		if err != nil {
			// 如果转换失败，使用原始数据
			converted = string(data)
		}

		// 1. 显示到本地终端（UTF-8数据）
		os.Stdout.Write([]byte(converted))

		// 2. 发送UTF-8到WebSocket
		msg := TerminalMessage{
			Type:    "output",
			Data:    converted,
			Session: sessionID,
			UserID:  "client",
		}

		jsonData, _ := json.Marshal(msg)
		wsWriteChan <- jsonData
	}
}

// convertGBKToUTF8 将GBK编码的字节转换为UTF-8字符串
func convertGBKToUTF8(data []byte) (string, error) {
	// 使用简体中文GBK编码
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(converted), nil
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

// setupHeartbeat 设置心跳机制
func setupHeartbeat(conn *websocket.Conn) {
	// 设置ping handler，自动回复pong
	conn.SetPingHandler(func(appData string) error {
		// log.Printf("💓 收到服务器ping")
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	// 设置pong handler
	conn.SetPongHandler(func(appData string) error {
		// log.Printf("💓 收到服务器pong")
		return nil
	})

	// 定期发送ping
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte("heartbeat")); err != nil {
					log.Printf("❌ 发送ping失败: %v", err)
					return
				}
				// log.Printf("💓 发送ping到服务器")
			}
		}
	}()
}
