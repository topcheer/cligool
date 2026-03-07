// +build windows

package main

import (
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
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
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
	Rows        int    `json:"rows,omitempty"`    // 终端行数
	Cols        int    `json:"cols,omitempty"`    // 终端列数
}

// 全局编码转换器
var consoleEncoding encoding.Encoding

// Windows ConPTY 常量
const (
	PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE = 0x00020016
	S_OK                                = 0

	// 控制台输入模式常量
	ENABLE_PROCESSED_INPUT        = 0x0001
	ENABLE_LINE_INPUT             = 0x0002
	ENABLE_ECHO_INPUT             = 0x0004
	ENABLE_WINDOW_INPUT           = 0x0008
	ENABLE_MOUSE_INPUT            = 0x0010
	ENABLE_INSERT_MODE            = 0x0020
	ENABLE_QUICK_EDIT_MODE        = 0x0040
	ENABLE_EXTENDED_FLAGS         = 0x0080
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200

	// 输入记录类型
	KEY_EVENT                 = 0x0001
	MOUSE_EVENT               = 0x0002
	WINDOW_BUFFER_SIZE_EVENT  = 0x0004
	MENU_EVENT                = 0x0008
	FOCUS_EVENT               = 0x0010

	// 标准句柄常量
	STD_INPUT_HANDLE  = uintptr(0xFFFFFFF6) - 10
	STD_OUTPUT_HANDLE = uintptr(0xFFFFFFF6) - 11
)

// Windows API 函数
var (
	modkernel32                 = windows.NewLazySystemDLL("kernel32.dll")
	procCreatePseudoConsole     = modkernel32.NewProc("CreatePseudoConsole")
	procResizePseudoConsole     = modkernel32.NewProc("ResizePseudoConsole")
	procClosePseudoConsole      = modkernel32.NewProc("ClosePseudoConsole")
	procInitializeProcThreadAttributeList = modkernel32.NewProc("InitializeProcThreadAttributeList")
	procUpdateProcThreadAttribute          = modkernel32.NewProc("UpdateProcThreadAttribute")
	procDeleteProcThreadAttributeList     = modkernel32.NewProc("DeleteProcThreadAttributeList")
	procGetStdHandle            = modkernel32.NewProc("GetStdHandle")
	procGetConsoleMode          = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode          = modkernel32.NewProc("SetConsoleMode")
	procReadConsoleInput        = modkernel32.NewProc("ReadConsoleInputW")
	procGetConsoleScreenBufferInfo = modkernel32.NewProc("GetConsoleScreenBufferInfo")
)

// _CONSOLE_SCREEN_BUFFER_INFO Windows控制台屏幕缓冲区信息
type _CONSOLE_SCREEN_BUFFER_INFO struct {
	dwSize              _COORD
	dwCursorPosition    _COORD
	wAttributes         uint16
	srWindow            windows.SmallRect
	dwMaximumWindowSize _COORD
}

// _INPUT_RECORD Windows输入记录
type _INPUT_RECORD struct {
	EventType uint16
	Event     [16]byte
}

// _KEY_EVENT_RECORD Windows键盘事件记录
type _KEY_EVENT_RECORD struct {
	bKeyDown          int32
	wRepeatCount      uint16
	wVirtualKeyCode   uint16
	wVirtualScanCode  uint16
	UnicodeChar       uint16
	dwControlKeyState uint32
}

// _COORD Windows坐标结构
type _COORD struct {
	X, Y int16
}

// Pack 将COORD打包成单个uintptr（高16位是Y，低16位是X）
func (c *_COORD) Pack() uintptr {
	return uintptr((int32(c.Y) << 16) | int32(c.X))
}

// _HPCON 伪控制台句柄
type _HPCON windows.Handle

// pseudoConsole Windows伪终端
type pseudoConsole struct {
	handle     _HPCON
	cmdIn      windows.Handle
	cmdOut     windows.Handle
	ptyIn      windows.Handle
	ptyOut     windows.Handle
}

// newPseudoConsole 创建一个新的伪控制台
func newPseudoConsole(width, height int16) (*pseudoConsole, error) {
	pc := &pseudoConsole{}

	// 创建输入管道 (ptyIn -> cmdIn)
	if err := windows.CreatePipe(&pc.ptyIn, &pc.cmdIn, nil, 0); err != nil {
		return nil, fmt.Errorf("创建输入管道失败: %w", err)
	}

	// 创建输出管道 (cmdOut -> ptyOut)
	if err := windows.CreatePipe(&pc.cmdOut, &pc.ptyOut, nil, 0); err != nil {
		windows.CloseHandle(pc.ptyIn)
		windows.CloseHandle(pc.cmdIn)
		return nil, fmt.Errorf("创建输出管道失败: %w", err)
	}

	// 创建COORD并打包
	coord := &_COORD{
		X: width,
		Y: height,
	}

	// 调用CreatePseudoConsole
	// 注意：第一个参数是打包后的COORD值，不是指针
	hr, _, _ := procCreatePseudoConsole.Call(
		coord.Pack(),    // size (packed COORD)
		uintptr(pc.ptyIn),  // hInput
		uintptr(pc.ptyOut), // hOutput
		0,                // dwFlags
		uintptr(unsafe.Pointer(&pc.handle)),
	)

	if hr != S_OK {
		windows.CloseHandle(pc.ptyIn)
		windows.CloseHandle(pc.cmdIn)
		windows.CloseHandle(pc.ptyOut)
		windows.CloseHandle(pc.cmdOut)
		return nil, fmt.Errorf("CreatePseudoConsole失败: HRESULT=0x%X", hr)
	}

	// 关闭不需要的管道端
	windows.CloseHandle(pc.ptyIn)
	windows.CloseHandle(pc.ptyOut)
	pc.ptyIn = 0
	pc.ptyOut = 0

	return pc, nil
}

// close 关闭伪控制台
func (pc *pseudoConsole) close() error {
	if pc.handle != 0 {
		procClosePseudoConsole.Call(uintptr(pc.handle))
		pc.handle = 0
	}
	if pc.cmdIn != 0 {
		windows.CloseHandle(pc.cmdIn)
		pc.cmdIn = 0
	}
	if pc.cmdOut != 0 {
		windows.CloseHandle(pc.cmdOut)
		pc.cmdOut = 0
	}
	return nil
}

// setConsoleRawMode 设置控制台为原始模式
func setConsoleRawMode() (windows.Handle, uint32, error) {
	// 获取标准输入句柄
	ret, _, err := procGetStdHandle.Call(STD_INPUT_HANDLE)
	if ret == 0 {
		return 0, 0, fmt.Errorf("获取标准输入句柄失败: %v", err)
	}
	stdinHandle := windows.Handle(ret)

	// 获取当前控制台模式
	var originalMode uint32
	ret1, _, err := procGetConsoleMode.Call(uintptr(stdinHandle), uintptr(unsafe.Pointer(&originalMode)))
	if ret1 == 0 {
		return 0, 0, fmt.Errorf("获取控制台模式失败: %v", err)
	}

	// 只禁用行输入和回显，保留其他所有标志
	// 这样可以让方向键等特殊按键正常工作
	rawMode := originalMode &^ (ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT)

	ret2, _, err := procSetConsoleMode.Call(uintptr(stdinHandle), uintptr(rawMode))
	if ret2 == 0 {
		return 0, 0, fmt.Errorf("设置控制台原始模式失败: %v", err)
	}

	return stdinHandle, originalMode, nil
}

// restoreConsoleMode 恢复控制台模式
func restoreConsoleMode(stdinHandle windows.Handle, mode uint32) error {
	ret, _, err := procSetConsoleMode.Call(uintptr(stdinHandle), uintptr(mode))
	if ret == 0 {
		return fmt.Errorf("恢复控制台模式失败: %v", err)
	}
	return nil
}

// readConsoleInput 读取控制台输入（包括方向键等特殊按键）
func readConsoleInput(stdinHandle windows.Handle) ([]byte, error) {
	var record _INPUT_RECORD
	var count uint32

	ret, _, err := procReadConsoleInput.Call(
		uintptr(stdinHandle),
		uintptr(unsafe.Pointer(&record)),
		uintptr(1),
		uintptr(unsafe.Pointer(&count)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("ReadConsoleInput 失败: %v", err)
	}
	if count == 0 {
		return nil, nil
	}

	// 处理键盘事件
	if record.EventType == KEY_EVENT {
		// 解析键盘事件
		keyEvent := (*_KEY_EVENT_RECORD)(unsafe.Pointer(&record.Event[0]))

		// 只处理按键按下事件
		if keyEvent.bKeyDown == 0 {
			return nil, nil
		}

		// 将 Unicode 字符转换为字节
		if keyEvent.UnicodeChar != 0 {
			// 普通字符输入
			buf := make([]byte, 4)
			n := utf8.EncodeRune(buf, rune(keyEvent.UnicodeChar))
			return buf[:n], nil
		} else {
			// 特殊按键（方向键、功能键等），转换为 ANSI 转义序列
			var seq string
			switch keyEvent.wVirtualKeyCode {
			case 0x26: // VK_UP
				seq = "\x1b[A"
			case 0x28: // VK_DOWN
				seq = "\x1b[B"
			case 0x25: // VK_LEFT
				seq = "\x1b[D"
			case 0x27: // VK_RIGHT
				seq = "\x1b[C"
			case 0x0D: // VK_RETURN
				seq = "\r"
			case 0x08: // VK_BACK
				seq = "\x7f"
			case 0x09: // VK_TAB
				seq = "\t"
			case 0x1B: // VK_ESCAPE
				seq = "\x1b"
			case 0x2E: // VK_DELETE
				seq = "\x1b[3~"
			case 0x24: // VK_HOME
				seq = "\x1bOH"
			case 0x23: // VK_END
				seq = "\x1bOF"
			case 0x21: // VK_PRIOR (Page Up)
				seq = "\x1b[5~"
			case 0x22: // VK_NEXT (Page Down)
				seq = "\x1b[6~"
			default:
				// 其他按键不处理
				return nil, nil
			}
			return []byte(seq), nil
		}
	}

	return nil, nil
}

// getConsoleSize 获取当前控制台大小
func getConsoleSize() (cols, rows int16, err error) {
	// 尝试方法1: 使用 golang.org/x/term 的跨平台方法
	if width, height, e := term.GetSize(int(os.Stdout.Fd())); e == nil {
		return int16(width), int16(height), nil
	}

	// 尝试方法2: 使用 Windows API
	// 获取标准输出句柄
	ret, _, _ := procGetStdHandle.Call(STD_OUTPUT_HANDLE)
	if ret == 0 {
		return 0, 0, fmt.Errorf("获取标准输出句柄失败")
	}
	stdoutHandle := windows.Handle(ret)

	// 获取控制台屏幕缓冲区信息
	var info _CONSOLE_SCREEN_BUFFER_INFO
	ret1, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(stdoutHandle),
		uintptr(unsafe.Pointer(&info)),
	)
	if ret1 == 0 {
		return 0, 0, fmt.Errorf("获取控制台屏幕缓冲区信息失败")
	}

	// 窗口大小
	cols = int16(info.srWindow.Right - info.srWindow.Left + 1)
	rows = int16(info.srWindow.Bottom - info.srWindow.Top + 1)

	// 确保返回的值是合理的
	if cols <= 0 || rows <= 0 {
		return 0, 0, fmt.Errorf("检测到无效的终端大小: %dx%d", cols, rows)
	}

	return cols, rows, nil
}

// write 向伪控制台写入数据
func (pc *pseudoConsole) write(data []byte) (int, error) {
	if pc.cmdIn == 0 {
		return 0, fmt.Errorf("输入管道未打开")
	}

	var written uint32
	err := windows.WriteFile(pc.cmdIn, data, &written, nil)
	return int(written), err
}

// read 从伪控制台读取数据
func (pc *pseudoConsole) read(buf []byte) (int, error) {
	if pc.cmdOut == 0 {
		return 0, fmt.Errorf("输出管道未打开")
	}

	var read uint32
	err := windows.ReadFile(pc.cmdOut, buf, &read, nil)
	return int(read), err
}

// resize 调整伪控制台大小
func (pc *pseudoConsole) resize(cols, rows int16) error {
	if pc.handle == 0 {
		return fmt.Errorf("伪控制台未初始化")
	}

	// 创建COORD并打包
	coord := &_COORD{
		X: cols,
		Y: rows,
	}

	// 调用ResizePseudoConsole
	hr, _, _ := procResizePseudoConsole.Call(
		uintptr(pc.handle),
		coord.Pack(),
	)

	if hr != S_OK {
		return fmt.Errorf("ResizePseudoConsole失败: HRESULT=0x%X", hr)
	}

	return nil
}

// monitorWindowSize 监控窗口大小变化
func monitorWindowSize(pc *pseudoConsole, stopChan chan struct{}) {
	// 保存上一次的大小
	var lastCols, lastRows int16 = -1, -1

	for {
		select {
		case <-stopChan:
			return
		default:
			// 定时检查控制台大小（每500ms）
			time.Sleep(500 * time.Millisecond)

			// 检查当前控制台大小
			cols, rows, err := getConsoleSize()
			if err != nil {
				continue
			}

			// 如果大小变化，调整 ConPTY
			if cols != lastCols || rows != lastRows {
				pc.resize(cols, rows)
				lastCols = cols
				lastRows = rows
			}
		}
	}
}

func main() {
	// 初始化控制台编码
	consoleEncoding = getConsoleEncoding()

	log.Println("CliGool Windows客户端启动...")

	serverURL := flag.String("server", "https://cligool.zty8.cn", "中继服务器URL")
	sessionID := flag.String("session", "", "会话ID")
	cols := flag.Int("cols", 0, "终端列数（0=自动检测）")
	rows := flag.Int("rows", 0, "终端行数（0=自动检测）")
	execCmd := flag.String("cmd", "", "直接执行的命令（如 claude, gemini 等）")
	execArgs := flag.String("args", "", "传递给命令的参数（可选，用空格分隔）")
	flag.Parse()

	// 构建完整的命令行
	var commandPath string
	var cmdArgs []string
	if execCmd != "" {
		log.Printf("准备启动命令: %s", execCmd)
		commandPath = execCmd
		// 解析参数
		if execArgs != "" {
			cmdArgs = strings.Fields(execArgs)
		}
	} else {
		log.Println("准备启动cmd.exe...")
		commandPath = "cmd.exe"
	}

	// 自动检测控制台大小
	actualCols, actualRows, err := getConsoleSize()
	if err != nil {
		if *cols == 0 {
			*cols = 120
		}
		if *rows == 0 {
			*rows = 80
		}
	} else {
		if *cols == 0 {
			*cols = int(actualCols)
		}
		if *rows == 0 {
			*rows = int(actualRows)
		}
	}

	sid := *sessionID
	if sid == "" {
		sid = uuid.New().String()
	}

	log.Println("会话ID:", sid)

	// 显示连接信息
	printHeader(sid, *serverURL)

	log.Println("开始连接WebSocket...")

	// 启动WebSocket并运行终端会话
	if err := runTerminalSession(*serverURL, sid, *cols, *rows, commandPath, cmdArgs); err != nil {
		log.Printf("终端会话失败: %v", err)
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}
}

func printHeader(sessionID, serverURL string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    CliGool 远程终端                        ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Printf("║ 会话ID: %-48s ║\n", sessionID)
	fmt.Printf("║ Web访问: %-48s ║\n", serverURL+"/session/"+sessionID)
	fmt.Printf("║ 连接状态: %-48s ║\n", "连接中...")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func runTerminalSession(serverURL, sessionID string, cols, rows int, commandPath string, cmdArgs []string) error {
	log.Println("开始建立WebSocket连接...")

	// 建立 WebSocket 连接
	wsURL, _ := buildWebSocketURL(serverURL, sessionID)
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Println("WebSocket URL:", wsURL)

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("WebSocket连接失败: %w", err)
	}
	defer conn.Close()

	log.Println("WebSocket已连接")
	fmt.Println("已连接到中继服务器")
	fmt.Println("现在可以在Web终端中输入命令了")
	fmt.Println("Windows模式：使用ConPTY，支持完整终端特性")
	fmt.Println()

	// 发送系统通知并自动打开浏览器
	webURL := fmt.Sprintf("%s/session/%s", serverURL, sessionID)

	// Windows: 使用简单的方法打开浏览器
	go func() {
		// 方法1: 使用 rundll32 打开 URL
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", webURL)
		if err := cmd.Run(); err != nil {
			// 方法2: 降级到 cmd start
			exec.Command("cmd", "/c", "start", "", webURL).Run()
		}
	}()

	// 可选：显示简单的 Toast 通知（Windows 10+）
	go func() {
		// 使用 PowerShell 的 BurntToast 模块（如果可用）
		psScript := fmt.Sprintf(
			`try { Import-Module BurntToast; New-BurntToastNotification -Title '🌐 CliGool Web 终端' -Message '%s' } catch {}`,
			webURL,
		)
		exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript).Run()
	}()

	log.Printf("✅ 已在浏览器中打开: %s", webURL)
	fmt.Println("📱 已发送系统通知并在浏览器中打开 Web 终端")
	fmt.Println()

	// 创建WebSocket写入channel，确保串行写入
	wsWriteChan := make(chan []byte, 100)

	// 启动WebSocket写入goroutine
	go func() {
		for data := range wsWriteChan {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("WebSocket写入失败: %v", err)
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
		Rows:        rows,
		Cols:        cols,
	}
	jsonData, _ := json.Marshal(initMsg)
	wsWriteChan <- jsonData
	log.Printf("已发送初始化消息: 工作目录=%s, 大小=%dx%d", wd, cols, rows)

	// 创建ConPTY
	log.Println("创建ConPTY...")
	pc, err := newPseudoConsole(int16(cols), int16(rows))
	if err != nil {
		return fmt.Errorf("创建ConPTY失败: %w", err)
	}
	defer pc.close()

	// 初始化ProcThreadAttributeList
	var size uintptr
	_, _, _ = procInitializeProcThreadAttributeList.Call(
		0,
		1,
		0,
		uintptr(unsafe.Pointer(&size)),
	)

	// 分配属性列表内存
	attrListData := make([]byte, size)

	hr, _, _ := procInitializeProcThreadAttributeList.Call(
		uintptr(unsafe.Pointer(&attrListData[0])),
		1,
		0,
		uintptr(unsafe.Pointer(&size)),
	)

	if hr != 1 { // 返回值是BOOL，1表示成功
		return fmt.Errorf("初始化属性列表失败")
	}
	defer procDeleteProcThreadAttributeList.Call(uintptr(unsafe.Pointer(&attrListData[0])))

	// 更新属性列表以包含PseudoConsole
	hr, _, _ = procUpdateProcThreadAttribute.Call(
		uintptr(unsafe.Pointer(&attrListData[0])),
		0,
		PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		uintptr(pc.handle),
		unsafe.Sizeof(pc.handle),
		0,
		0,
	)

	if hr != 1 {
		return fmt.Errorf("更新PseudoConsole属性失败")
	}

	// 创建StartupInfoEx
	si := &windows.StartupInfoEx{
		StartupInfo: windows.StartupInfo{
			Cb: uint32(unsafe.Sizeof(windows.StartupInfoEx{})),
		},
	}
	si.ProcThreadAttributeList = (*windows.ProcThreadAttributeList)(unsafe.Pointer(&attrListData[0]))

	// 查找命令的完整路径
	fullCommandPath := commandPath
	if !strings.Contains(commandPath, "\\") && !strings.Contains(commandPath, "/") {
		// 不包含路径分隔符，尝试在PATH中查找
		if path, err := exec.LookPath(commandPath); err == nil {
			fullCommandPath = path
		}
	}

	// 检查是否是脚本文件（.cmd, .bat, .ps1等）
	var appPath, cmdLineStr string
	ext := strings.ToLower(filepath.Ext(fullCommandPath))

	if ext == ".cmd" || ext == ".bat" {
		appPath = os.Getenv("SystemRoot") + "\\System32\\cmd.exe"
		// 构建完整命令行，包含参数
		if len(cmdArgs) > 0 {
			cmdLineStr = fmt.Sprintf("%s /c \"%s\" %s", appPath, fullCommandPath, strings.Join(cmdArgs, " "))
		} else {
			cmdLineStr = fmt.Sprintf("%s /c \"%s\"", appPath, fullCommandPath)
		}
	} else if ext == ".ps1" {
		appPath = os.Getenv("SystemRoot") + "\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
		// 构建完整命令行，包含参数
		if len(cmdArgs) > 0 {
			cmdLineStr = fmt.Sprintf("%s -ExecutionPolicy Bypass -File \"%s\" %s", appPath, fullCommandPath, strings.Join(cmdArgs, " "))
		} else {
			cmdLineStr = fmt.Sprintf("%s -ExecutionPolicy Bypass -File \"%s\"", appPath, fullCommandPath)
		}
	} else {
		appPath = fullCommandPath
		// 构建完整命令行，包含参数
		if len(cmdArgs) > 0 {
			cmdLineStr = fmt.Sprintf("%s %s", fullCommandPath, strings.Join(cmdArgs, " "))
		} else {
			cmdLineStr = fullCommandPath
		}
	}

	// 转换应用程序路径和命令行为UTF-16
	appName, err := windows.UTF16PtrFromString(appPath)
	if err != nil {
		return fmt.Errorf("转换应用程序路径失败: %w", err)
	}
	cmdLine, err := windows.UTF16PtrFromString(cmdLineStr)
	if err != nil {
		return fmt.Errorf("转换命令行失败: %w", err)
	}

	// 创建进程
	log.Printf("启动命令: %s", appPath)
	var pi windows.ProcessInformation

	err = windows.CreateProcess(
		appName,
		cmdLine,
		nil,
		nil,
		false, // 不继承句柄
		windows.EXTENDED_STARTUPINFO_PRESENT,
		nil,   // 使用父进程环境
		nil,   // 使用当前目录
		&si.StartupInfo,
		&pi,
	)
	if err != nil {
		log.Printf("CreateProcess失败: %v", err)
		return fmt.Errorf("CreateProcess失败: %w", err)
	}
	defer windows.CloseHandle(pi.Thread)
	defer windows.CloseHandle(pi.Process)

	log.Printf("命令已启动: %s, PID: %d", appPath, pi.ProcessId)

	// 启动窗口大小监控
	stopMonitor := make(chan struct{})
	go monitorWindowSize(pc, stopMonitor)
	defer close(stopMonitor)

	// 设置标准输入为原始模式
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("警告: 无法设置终端为原始模式: %v", err)
	} else {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	// 本地终端输入 -> ConPTY
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				continue
			}

			if n == 0 {
				continue
			}

			data := buf[:n]

			// 写入ConPTY
			if _, writeErr := pc.write(data); writeErr != nil {
				continue
			}

			// 发送到WebSocket，让Web端看到本地输入
			msg := TerminalMessage{
				Type:    "input",
				Data:    string(data),
				Session: sessionID,
				UserID:  "client",
				Source:  "local",
			}
			jsonData, _ := json.Marshal(msg)
			wsWriteChan <- jsonData
		}
	}()

	// WebSocket -> ConPTY (网页输入写入ConPTY)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var msg TerminalMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			if msg.Type == "input" && msg.Data != "" {
				data := []byte(msg.Data)
				pc.write(data)
			}
		}
	}()

	// ConPTY -> 本地stdout + WebSocket (ConPTY输出同时显示在两端)
	buf := make([]byte, 4096)
	for {
		n, err := pc.read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("ConPTY读取失败: %w", err)
		}

		data := buf[:n]

		// 转换控制台编码到UTF-8
		converted, err := convertToUTF8(data)
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

// getConsoleEncoding 获取控制台输出编码
func getConsoleEncoding() encoding.Encoding {
	// 获取控制台输出code page
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleOutputCP := kernel32.NewProc("GetConsoleOutputCP")

	codePage, _, _ := getConsoleOutputCP.Call()
	log.Printf("检测到控制台code page: %d", codePage)

	// 根据code page返回对应的编码
	switch codePage {
	case 932:
		// 日文 Shift-JIS
		return japanese.ShiftJIS
	case 936:
		// 简体中文 GBK
		return simplifiedchinese.GBK
	case 949:
		// 韩文 EUC-KR
		return korean.EUCKR
	case 950:
		// 繁体中文 Big5
		return traditionalchinese.Big5
	case 1252:
		// 西欧 Latin-1
		return charmap.Windows1252
	case 437:
		// 英文 CP437
		return charmap.CodePage437
	default:
		// 默认使用Latin-1，可以处理所有单字节字符
		log.Printf("未知的code page %d，使用Latin-1编码", codePage)
		return charmap.Windows1252
	}
}

// convertToUTF8 将控制台编码的字节转换为UTF-8字符串
func convertToUTF8(data []byte) (string, error) {
	// 首先检查数据是否已经是有效的UTF-8（ConPTY默认输出UTF-8）
	if utf8.Valid(data) {
		return string(data), nil
	}

	// 如果不是UTF-8，尝试使用检测到的控制台编码进行转换
	if consoleEncoding == nil {
		// 如果无法检测编码，作为Latin-1处理
		decoded, _ := charmap.Windows1252.NewDecoder().Bytes(data)
		return string(decoded), nil
	}

	reader := transform.NewReader(bytes.NewReader(data), consoleEncoding.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		// 转换失败时，返回原始字符串
		return string(data), nil
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
					log.Printf("发送ping失败: %v", err)
					return
				}
			}
		}
	}()
}
