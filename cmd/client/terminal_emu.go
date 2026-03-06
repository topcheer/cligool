// +build !windows

package main

import (
	"bytes"
)

// TerminalEmulator 终端仿真器，用于处理终端能力查询
type TerminalEmulator struct {
	buffer bytes.Buffer
}

// NewTerminalEmulator 创建终端仿真器
func NewTerminalEmulator() *TerminalEmulator {
	return &TerminalEmulator{}
}

// Process 处理输入数据，拦截终端查询并返回响应
// 返回: (输出数据, 查询响应数据)
func (te *TerminalEmulator) Process(data []byte) ([]byte, []byte) {
	te.buffer.Write(data)

	var output bytes.Buffer
	var responses bytes.Buffer

	// 将缓冲区数据转换为字节数组以便处理
	buf := te.buffer.Bytes()

	for len(buf) > 0 {
		// 检查是否是转义序列开始
		if buf[0] == 0x1b && len(buf) >= 2 {
			// 处理各种查询序列
			if handled, resp := te.handleQuery(buf); handled {
				if len(resp) > 0 {
					responses.Write(resp)
				}
				// 跳过已处理的序列
				// 需要找到序列的结束位置
				consumed := te.getSequenceLength(buf)
				buf = buf[consumed:]
				continue
			}
		}

		// 不是查询序列，输出第一个字节
		output.WriteByte(buf[0])
		buf = buf[1:]
	}

	// 更新缓冲区
	te.buffer.Reset()
	te.buffer.Write(buf)

	return output.Bytes(), responses.Bytes()
}

// handleQuery 处理终端查询序列
// 返回: (是否处理, 响应数据)
func (te *TerminalEmulator) handleQuery(data []byte) (bool, []byte) {
	if len(data) < 2 || data[0] != 0x1b {
		return false, nil
	}

	// ESC c - DA1 查询 (设备属性)
	if data[1] == 'c' {
		return true, []byte("\x1b[?1;2c")
	}

	// ESC [c - DA1 查询 (CSI 形式)
	if data[1] == '[' && len(data) >= 3 {
		if data[2] == 'c' || (data[2] == '0' && len(data) >= 4 && data[3] == 'c') {
			return true, []byte("\x1b[?1;2c")
		}
	}

	// ESC [>c - DA2 查询 (终端版本)
	if data[1] == '[' && len(data) >= 4 && data[2] == '>' {
		if data[3] == 'c' || (data[3] == '0' && len(data) >= 5 && data[4] == 'c') {
			return true, []byte("\x1b[>0;272;0c")
		}
	}

	// ESC [6n - CPR 查询 (光标位置)
	if len(data) >= 4 && data[1] == '[' && data[2] == '6' && data[3] == 'n' {
		return true, []byte("\x1b[1;1R")
	}

	// ESC [>q - XTVERSION 查询
	if len(data) >= 4 && data[1] == '[' && data[2] == '>' && data[3] == 'q' {
		return true, []byte("\x1b[>0;272;0c")
	}

	// ESC P+q - XTGETTCAP 查询 (终端能力)
	if data[1] == 'P' && len(data) >= 3 && data[2] == '+' {
		// 找到序列结束 (ST 字符: ESC \)
		end := bytes.Index(data, []byte{0x1b, '\\'})
		if end > 0 {
			return true, nil // 返回空响应
		}
		// 序列不完整，暂不处理
		return false, nil
	}

	return false, nil
}

// getSequenceLength 获取转义序列的长度
func (te *TerminalEmulator) getSequenceLength(data []byte) int {
	if len(data) < 2 || data[0] != 0x1b {
		return 1
	}

	// CSI 序列: ESC [ ... 终止字符
	if data[1] == '[' {
		for i := 2; i < len(data); i++ {
			c := data[i]
			// CSI 终止字符范围: 0x40-0x7E
			if c >= 0x40 && c <= 0x7E {
				return i + 1
			}
		}
		// 序列不完整，暂时不处理
		return 0
	}

	// ESC ] ... ST 或 BEL (OSC 序列)
	if data[1] == ']' {
		st := bytes.Index(data, []byte{0x1b, '\\'})
		if st > 0 {
			return st + 2
		}
		// 查找 BEL
		for i := 2; i < len(data); i++ {
			if data[i] == 0x07 {
				return i + 1
			}
		}
		// 序列不完整
		return 0
	}

	// ESC P ... ST (DCS 序列)
	if data[1] == 'P' {
		st := bytes.Index(data, []byte{0x1b, '\\'})
		if st > 0 {
			return st + 2
		}
		return 0
	}

	// 简单的转义序列 (2-3 字节)
	if len(data) >= 2 {
		// 检查是否是已知的简单序列
		switch data[1] {
		case 'c', 'D', 'E', 'H', 'M', '7', '8', '>', '=':
			return 2
		case '(':
			if len(data) >= 3 {
				return 3
			}
		}
	}

	// 默认处理为单字节
	return 1
}

// Flush 刷新缓冲区，返回所有剩余数据
func (te *TerminalEmulator) Flush() []byte {
	data := te.buffer.Bytes()
	te.buffer.Reset()
	return data
}

// Reset 重置仿真器状态
func (te *TerminalEmulator) Reset() {
	te.buffer.Reset()
}
