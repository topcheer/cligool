//go:build !windows
// +build !windows

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
	"golang.org/x/term"
)

type TerminalMessage struct {
	Type       string `json:"type"`
	Data       string `json:"data"`
	Session    string `json:"session"`
	UserID     string `json:"user_id"`
	Source     string `json:"source,omitempty"` // "local" or "web"
	WorkingDir string `json:"working_dir,omitempty"`
	OSInfo     string `json:"os_info,omitempty"`
	Rows       int    `json:"rows,omitempty"` // 终端行数
	Cols       int    `json:"cols,omitempty"` // 终端列数
}

func main() {
	// 加载配置文件
	config, configPath, err := LoadConfig()
	if err != nil {
		log.Printf("警告: 加载配置文件失败: %v，使用默认配置", err)
		config = DefaultConfig()
	} else {
		log.Printf("✅ 已加载配置文件: %s", configPath)
	}

	serverURL := flag.String("server", config.Server, "中继服务器URL")
	proxyURL := flag.String("proxy", config.Proxy, "代理服务器地址（如 http://proxy.example.com:8080 或 socks5://proxy.example.com:1080）")
	sessionID := flag.String("session", "", "会话ID")
	cols := flag.Int("cols", config.Cols, "终端列数（0=自动检测）")
	rows := flag.Int("rows", config.Rows, "终端行数（0=自动检测）")
	execCmd := flag.String("cmd", "", "直接执行的命令（如 claude, gemini 等）")
	execArgs := flag.String("args", "", "传递给命令的参数（可选，用空格分隔）")
	noBrowser := flag.Bool("no-browser", config.NoBrowser, "禁止自动打开浏览器")
	flag.Parse()

	// 显示代理信息
	if *proxyURL != "" {
		log.Printf("🔧 使用代理: %s", *proxyURL)
	}

	// 构建完整的命令行
	var commandPath string
	var cmdArgs []string
	if *execCmd != "" {
		commandPath = *execCmd
		// 解析参数
		if *execArgs != "" {
			cmdArgs = strings.Fields(*execArgs)
		}
	} else {
		// 使用默认shell
		commandPath = ""
		cmdArgs = nil
	}

	// 自动检测终端大小
	if *cols == 0 || *rows == 0 {
		if size, err := pty.GetsizeFull(os.Stdout); err == nil {
			if *cols == 0 {
				*cols = int(size.Cols)
			}
			if *rows == 0 {
				*rows = int(size.Rows)
			}
		} else {
			if *cols == 0 {
				*cols = 120
			}
			if *rows == 0 {
				*rows = 80
			}
		}
	}

	sid := *sessionID
	if sid == "" {
		sid = uuid.New().String()
	}

	// 显示连接信息
	printHeader(sid, *serverURL, *proxyURL)

	// 启动WebSocket并运行PTY
	if err := runTerminalSession(*serverURL, *proxyURL, sid, *cols, *rows, commandPath, cmdArgs, *noBrowser); err != nil {
		log.Printf("❌ 终端会话失败: %v", err)
		fmt.Printf("\n❌ 错误: %v\n", err)
		os.Exit(1)
	}
}

func printHeader(sessionID, serverURL, proxyURL string) {
	// 清理URL，移除末尾的斜杠
	cleanURL := strings.TrimSuffix(serverURL, "/")
	webURL := fmt.Sprintf("%s/session/%s", cleanURL, sessionID)

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    🚀 CliGool 远程终端                      ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Printf("║ 📋 会话ID: %-43s ║\n", sessionID)
	fmt.Printf("║ 🌐 Web访问: %-43s ║\n", webURL)
	if proxyURL != "" {
		fmt.Printf("║ 🔧 代理服务器: %-39s ║\n", proxyURL)
	}
	fmt.Printf("║ 🔗 连接状态: %-43s ║\n", "🟡 连接中...")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func runTerminalSession(serverURL, proxyURL, sessionID string, cols, rows int, commandPath string, cmdArgs []string, noBrowser bool) error {
	// 建立 WebSocket 连接参数
	wsURL, err := buildWebSocketURL(serverURL, sessionID)
	if err != nil {
		return fmt.Errorf("构建WebSocket URL失败: %w", err)
	}

	// 创建拨号器
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 如果配置了代理，使用代理拨号器
	if proxyURL != "" {
		proxyDialer, err := createProxyDialer(proxyURL)
		if err != nil {
			return fmt.Errorf("创建代理拨号器失败: %w", err)
		}
		dialer.NetDial = proxyDialer
		log.Printf("✅ 已配置代理: %s", proxyURL)
	}

	dialRelay := func() (*websocket.Conn, error) {
		conn, _, err := dialer.Dial(wsURL, nil)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	// 确保在退出时总是通知 relay 服务器
	var sessionError error
	var relayClient *relayConnectionManager
	defer func() {
		if relayClient == nil {
			return
		}

		// 发送关闭消息
		closeMsg := TerminalMessage{
			Type:    "close",
			Session: sessionID,
			UserID:  "client",
		}

		if sessionError != nil {
			closeMsg.Data = fmt.Sprintf("客户端错误: %v", sessionError)
			log.Printf("❌ 发送错误关闭消息: %v", sessionError)
		} else {
			closeMsg.Data = "客户端正常退出"
			log.Printf("✅ 发送正常关闭消息")
		}

		relayClient.Shutdown(&closeMsg)
	}()

	// 创建 PTY
	var command *exec.Cmd
	if commandPath != "" {
		// 直接执行指定的命令
		log.Printf("直接执行命令: %s", commandPath)
		if len(cmdArgs) > 0 {
			log.Printf("命令参数: %v", cmdArgs)
			// 确保命令在PATH中能找到
			fullPath := commandPath
			if filepath.Base(commandPath) == commandPath {
				// 如果只给了命令名，没有路径，查找完整路径
				if path, err := exec.LookPath(commandPath); err == nil {
					fullPath = path
				}
			}
			command = exec.Command(fullPath, cmdArgs...)
		} else {
			// 确保命令在PATH中能找到
			fullPath := commandPath
			if filepath.Base(commandPath) == commandPath {
				// 如果只给了命令名，没有路径，查找完整路径
				if path, err := exec.LookPath(commandPath); err == nil {
					fullPath = path
				}
			}
			command = exec.Command(fullPath)
		}
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
		"COLORTERM=truecolor", // 启用真色支持
		"FORCE_COLOR=1",       // 强制启用颜色
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
		sessionError = fmt.Errorf("启动PTY失败: %w", err)
		return sessionError
	}
	defer ptmx.Close()

	// 设置初始PTY窗口大小
	if err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}); err != nil {
		log.Printf("设置PTY窗口大小失败: %v", err)
	}

	// 在本地终端启动后再连接 relay，这样即使 relay 暂时不可用也能继续本地运行
	wd, _ := os.Getwd()
	cleanURL := strings.TrimSuffix(serverURL, "/")
	webURL := fmt.Sprintf("%s/session/%s", cleanURL, sessionID)
	relayClient = newRelayConnectionManager(relayConnectionConfig{
		Dial: dialRelay,
		InitMessage: TerminalMessage{
			Type:       "init",
			Session:    sessionID,
			UserID:     "client",
			WorkingDir: wd,
			OSInfo:     "unix",
			Rows:       rows,
			Cols:       cols,
		},
		InboundHandler: func(msg TerminalMessage) {
			switch msg.Type {
			case "input":
				if msg.Data == "" {
					return
				}
				if _, err := ptmx.Write([]byte(msg.Data)); err != nil {
					log.Printf("PTY写入失败（Web输入）: %v", err)
				}
			case "resize":
				if msg.Rows <= 0 || msg.Cols <= 0 {
					return
				}
				if err := pty.Setsize(ptmx, &pty.Winsize{
					Rows: uint16(msg.Rows),
					Cols: uint16(msg.Cols),
				}); err != nil {
					log.Printf("PTY窗口大小调整失败（Web输入）: %v", err)
				}
			}
		},
		OnConnected: func(reconnected bool) {
			if reconnected {
				log.Println("✅ 已重新连接到中继服务器")
				fmt.Println("✅ 已重新连接到中继服务器")
				return
			}

			log.Println("✅ WebSocket已连接")
			fmt.Println("✅ 已连接到中继服务器")
			fmt.Println("💡 现在可以在Web终端中输入命令了")

			if !noBrowser {
				notifier := NewSystemNotifier()
				if err := notifier.SendWebTerminalNotification(webURL); err == nil {
					fmt.Println("📱 已发送系统通知并在浏览器中打开 Web 终端")
				}
			}
			fmt.Println()
		},
		OnDisconnected: func(hadConnectedBefore bool, err error) {
			if hadConnectedBefore {
				message := fmt.Sprintf("⚠️ 与中继服务器连接已断开（%v），正在自动重试并缓冲未发送消息", err)
				log.Println(message)
				fmt.Println(message)
				return
			}

			message := fmt.Sprintf("⚠️ 暂时无法连接到中继服务器（%v），正在自动重试并缓冲未发送消息", err)
			log.Println(message)
			fmt.Println(message)
		},
	})
	relayClient.Start()

	// 处理窗口大小变化
	handleResize := func() {
		// 从标准输入获取当前终端窗口大小
		if size, err := pty.GetsizeFull(os.Stdin); err == nil {
			// 更新 PTY 窗口大小
			if err := pty.Setsize(ptmx, size); err != nil {
				// 如果失败，尝试使用完整结构体
				fullSize := &pty.Winsize{
					Rows: size.Rows,
					Cols: size.Cols,
					X:    size.X,
					Y:    size.Y,
				}
				_ = pty.Setsize(ptmx, fullSize)
			}

			// 发送新的终端大小到 WebSocket 服务器
			resizeMsg := TerminalMessage{
				Type:    "resize",
				Rows:    int(size.Rows),
				Cols:    int(size.Cols),
				Session: sessionID,
				UserID:  "client",
			}
			if err := relayClient.Send(resizeMsg); shouldLogRelaySendError(err) {
				log.Printf("❌ 发送终端大小失败: %v", err)
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

	// 将本地终端设置为原始模式，以便立即发送每个按键
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("警告: 无法设置终端为原始模式: %v", err)
	} else {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

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

			data := buf[:n]

			// 写入PTY
			if _, err := ptmx.Write(data); err != nil {
				log.Printf("PTY写入失败（本地输入）: %v", err)
				continue
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
			sessionError = fmt.Errorf("PTY读取失败: %w", err)
			return sessionError
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

			if err := relayClient.Send(msg); shouldLogRelaySendError(err) {
				log.Printf("❌ 发送终端输出失败: %v", err)
			}
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

// createProxyDialer 创建代理拨号器
func createProxyDialer(proxyURL string) (func(network, addr string) (net.Conn, error), error) {
	// 解析代理URL
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("解析代理URL失败: %w", err)
	}

	// 根据代理类型创建拨号器
	switch proxy.Scheme {
	case "http", "https":
		// HTTP 代理 - 使用环境变量方式
		// 设置环境变量供http包使用
		// 注意：WebSocket连接需要特殊处理
		return createHTTPProxyDialer(proxyURL)

	case "socks5":
		// SOCKS5 代理
		return createSocks5Dialer(proxy.Host)

	default:
		return nil, fmt.Errorf("不支持的代理类型: %s（支持 http、https、socks5）", proxy.Scheme)
	}
}

// createHTTPProxyDialer 创建HTTP代理拨号器
func createHTTPProxyDialer(proxyURL string) (func(network, addr string) (net.Conn, error), error) {
	// 对于HTTP代理，我们通过环境变量设置
	// 但这里使用更直接的方式：先连接到代理，然后发送CONNECT请求
	// 简化实现：使用 golang.org/x/net/proxy 的 SOCKS5 拨号器
	// 因为 SOCKS5 协议也支持 HTTP 代理

	// 尝试将 HTTP 代理转换为 SOCKS5 格式
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("解析代理URL失败: %w", err)
	}

	// 简化：使用 SOCKS5 方式
	// 大多数 HTTP 代理也支持 SOCKS5 协议
	// proxyParsed.Host 已经包含主机和端口（如果有）
	proxyAddr := proxyParsed.Host

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP代理拨号器失败: %w", err)
	}

	return dialer.Dial, nil
}

// createSocks5Dialer 创建SOCKS5拨号器
func createSocks5Dialer(proxyAddr string) (func(network, addr string) (net.Conn, error), error) {
	// 使用 golang.org/x/net/proxy 的 socks5 拨号器
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("创建SOCKS5拨号器失败: %w", err)
	}

	return dialer.Dial, nil
}
