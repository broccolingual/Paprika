package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

const (
	CTRL_A = iota + 1
	CTRL_B
	CTRL_C
	CTRL_D
	CTRL_E
	CTRL_F
	CTRL_G
	CTRL_H
	CTRL_I
	CTRL_J
	CTRL_K
	CTRL_L
	CTRL_M
	CTRL_N
	CTRL_O
	CTRL_P
	CTRL_Q
	CTRL_R
	CTRL_S
	CTRL_T
	CTRL_U
	CTRL_V
	CTRL_W
	CTRL_X
	CTRL_Y
	CTRL_Z
	ESC
)

const (
	KEY_UP    = 1001
	KEY_DOWN  = 1002
	KEY_RIGHT = 1003
	KEY_LEFT  = 1004
)

func (w *Window) parseKey(b []byte) (rune, int) {
	if len(b) == 3 {
		if b[0] == byte(27) && b[1] == '[' {
			switch b[2] {
			case 'A':
				return KEY_UP, 3
			case 'B':
				return KEY_DOWN, 3
			case 'C':
				return KEY_RIGHT, 3
			case 'D':
				return KEY_LEFT, 3
			default:
				return -1, 0
			}
		}
	}
	return utf8.DecodeRune(b)
}

func (w *Window) readKeys() {
	buf := make([]byte, 64)
	for {
		if n, err := os.Stdin.Read(buf); err == nil {
			b := buf[:n]
			for {
				r, n := w.parseKey(b)
				if n == 0 {
					break
				}
				w.KeyChan <- r
				b = b[n:]
			}
		}
	}
}

func (w *Window) switchKeys() {
	for {
		r := <-w.KeyChan
		w.Clear()
		rows, cols := GetWinSize()
		w.MoveCursorPos(uint16(cols/2)-10, uint16(rows/2))
		fmt.Printf("%d(UTF-8) -> ", r)
		switch r {
		case CTRL_A:
			fmt.Println("CTRL_A")
		case CTRL_B:
			fmt.Println("CTRL_B")
		case CTRL_C:
			fmt.Println("CTRL_C")
		case CTRL_D:
			fmt.Println("CTRL_D")
		case CTRL_E:
			fmt.Println("CTRL_E")
		case CTRL_F:
			fmt.Println("CTRL_F")
		case CTRL_G:
			fmt.Println("CTRL_G")
		case CTRL_H:
			fmt.Println("CTRL_H")
		case CTRL_I:
			fmt.Println("CTRL_I")
		case CTRL_J:
			fmt.Println("CTRL_J")
		case CTRL_K:
			fmt.Println("CTRL_K")
		case CTRL_L:
			fmt.Println("CTRL_L")
		case CTRL_M:
			fmt.Println("CTRL_M or ENTER")
			w.Cursor.Col += 1
		case CTRL_N:
			fmt.Println("CTRL_N")
		case CTRL_O:
			fmt.Println("CTRL_O")
		case CTRL_P:
			fmt.Println("CTRL_P")
		case CTRL_Q:
			fmt.Println("CTRL_Q")
		case CTRL_R:
			fmt.Println("CTRL_R")
		case CTRL_S: // SAVE
			fmt.Println("CTRL_S")
		case CTRL_T:
			fmt.Println("CTRL_T")
		case CTRL_U:
			fmt.Println("CTRL_U")
		case CTRL_V:
			fmt.Println("CTRL_V")
		case CTRL_W:
			fmt.Println("CTRL_W")
		case CTRL_X: // EXIT
			fmt.Println("CTRL_X")
			return
		case CTRL_Y:
			fmt.Println("CTRL_Y")
		case CTRL_Z:
			fmt.Println("CTRL_Z")
		case ESC:
			fmt.Println("ESC")
		case 32:
			fmt.Println("SPACE")
		case 127:
			fmt.Println("BACKSPACE")
		case KEY_UP:
			fmt.Println("KEY_UP")
			w.Cursor.Col -= 1
		case KEY_DOWN:
			fmt.Println("KEY_DOWN")
			w.Cursor.Col += 1
		case KEY_RIGHT:
			fmt.Println("KEY_RIGHT")
			w.Cursor.Row += 1
		case KEY_LEFT:
			fmt.Println("KEY_LEFT")
			w.Cursor.Row -= 1
		default:
			fmt.Printf(string(r))
		}
	}
}
