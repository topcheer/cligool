// +build !windows

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
	"time"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type TerminalMessage struct {
	Type        string `json:"type"`
	Data        string `json:"data"`
	Session     string `json:"session"`
	UserID      string `json:"user_id"`
	Source      string `json:"source,omitempty"` // "local" or "web"
	WorkingDir  string `json:"working_dir,omitempty"`
	OSInfo      string `json:"os_info,omitempty"`
	Rows        int    `json:"rows,omitempty"`    // 终端行数
	Cols        int    `json:"cols,omitempty"`    // 终端列数
}

func main() {
	serverURL := flag.String("server", "https://cligool.zty8.cn", "中继服务器URL")
	sessionID := flag.String("session", "", "会话ID")
	cols := flag.Int("cols", 120, "终端列数")
	rows := flag.Int("rows", 80, "终端行数")
	execCmd := flag.String("cmd", "", "直接执行的命令（如 claude, gemini 等）")
	flag.Parse()

	sid := *sessionID
	if sid == "" {
		sid = uuid.New().String()
	}

	// 显示连接信息
	printHeader(sid, *serverURL)

	// 启动WebSocket并运行PTY
	if err := runTerminalSession(*serverURL, sid, *cols, *rows, *execCmd); err != nil {
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

func runTerminalSession(serverURL, sessionID string, cols, rows int, execCmd string) error {
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

	// 设置心跳机制
	// 注意：心跳是控制帧，直接发送，不通过channel
	setupHeartbeat(conn)

	// 发送初始化消息（工作目录和系统信息）
	wd, _ := os.Getwd()
	initMsg := TerminalMessage{
		Type:        "init",
		Session:     sessionID,
		UserID:      "client",
		WorkingDir:  wd,
		OSInfo:      "unix",
		Rows:        rows,
		Cols:        cols,
	}
	jsonData, _ := json.Marshal(initMsg)
	wsWriteChan <- jsonData
	log.Printf("✅ 已发送初始化消息: 工作目录=%s, 大小=%dx%d", wd, cols, rows)

	// 创建 PTY
	var command *exec.Cmd
	if execCmd != "" {
		// 直接执行指定的命令
		log.Printf("直接执行命令: %s", execCmd)
		command = exec.Command(execCmd)
	} else {
		// 使用默认 shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		log.Printf("使用默认shell: %s", shell)
		command = exec.Command(shell, "-i", "-l")
	}

	// 设置环境变量以确保终端工具正确工作
	env := append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",  // 启用真色支持
		"FORCE_COLOR=1",         // 强制启用颜色
	)

	// 确保 LANG/LC_ALL 设置为 UTF-8
	hasUTF8Locale := false
	for _, e := range env {
		if len(e) > 9 && (e[:9] == "LANG=" || e[:9] == "LC_ALL=") {
			if len(e) > 10 && e[len(e)-5:] == ".UTF-8" {
				hasUTF8Locale = true
				break
			}
		}
	}
	if !hasUTF8Locale {
		env = append(env, "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")
	}

	command.Env = env

	// 启动PTY
	ptmx, err := pty.Start(command)
	if err != nil {
		return fmt.Errorf("启动PTY失败: %w", err)
	}
	defer ptmx.Close()

	// 设置初始PTY窗口大小
	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}); err != nil {
		log.Printf("设置PTY窗口大小失败: %v", err)
	}

	// 处理窗口大小变化
	handleResize := func() {
		// 获取当前终端窗口大小
		if size, err := pty.GetsizeFull(ptmx); err == nil {
			if err := pty.Setsize(ptmx, size); err != nil {
				log.Printf("更新PTY窗口大小失败: %v", err)
			}
		}
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			handleResize()
		}
	}()

	// 本地终端输入 -> PTY
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

			// 写入PTY
			if _, err := ptmx.Write(buf[:n]); err != nil {
				log.Printf("PTY写入失败（本地输入）: %v", err)
				continue
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
					log.Printf("PTY写入失败（Web输入）: %v", err)
					continue
				}
			} else if msg.Type == "resize" {
				// 处理窗口大小调整
				// TODO: 实现PTY窗口大小调整
			}
		}
	}()

	// PTY -> 本地stdout + WebSocket (PTY输出同时显示在两端)
	// 使用更小的缓冲区以减少延迟
	buf := make([]byte, 1024)
	emulator := NewTerminalEmulator()
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("PTY读取失败: %w", err)
		}

		data := buf[:n]

		// 使用终端仿真器处理数据，拦截查询并生成响应
		output, responses := emulator.Process(data)

		// 输出正常数据
		if len(output) > 0 {
			// 1. 显示到本地终端
			os.Stdout.Write(output)

			// 2. 同时发送到WebSocket
			msg := TerminalMessage{
				Type:    "output",
				Data:    string(output),
				Session: sessionID,
				UserID:  "client",
			}

			jsonData, _ := json.Marshal(msg)
			wsWriteChan <- jsonData
		}

		// 发送查询响应到 PTY
		if len(responses) > 0 {
			ptmx.Write(responses)
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
