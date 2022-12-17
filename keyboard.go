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
	w.InitCursorPos()
	fmt.Printf("\033[7mEditing \"%s\"\033[0m", w.Editor.FilePath)
	w.MoveCursorPos(1, 2)
	w.Draw()
	for {
		r := <-w.KeyChan
		w.InitCursorPos()
		w.ClearLine(1)
		fmt.Printf("\033[7mEditing \"%s\" | INPUT: ", w.Editor.FilePath)
		// rows, cols := GetWinSize()
		// w.MoveCursorPos(uint16(cols/2)-10, uint16(rows/2))
		// fmt.Printf("%d(UTF-8) -> ", r)
		switch r {
		case CTRL_A:
			fmt.Printf("CTRL_A")
		case CTRL_B:
			fmt.Printf("CTRL_B")
		case CTRL_C:
			fmt.Printf("CTRL_C")
		case CTRL_D:
			fmt.Printf("CTRL_D")
		case CTRL_E:
			fmt.Printf("CTRL_E")
		case CTRL_F:
			fmt.Printf("CTRL_F")
		case CTRL_G:
			fmt.Printf("CTRL_G")
		case CTRL_H:
			fmt.Printf("CTRL_H")
		case CTRL_I:
			fmt.Printf("CTRL_I")
		case CTRL_J:
			fmt.Printf("CTRL_J")
		case CTRL_K:
			fmt.Printf("CTRL_K")
		case CTRL_L:
			fmt.Printf("CTRL_L")
		case CTRL_M:
			fmt.Printf("CTRL_M or ENTER")
			w.Editor.Cursor.Col += 1
		case CTRL_N:
			fmt.Printf("CTRL_N")
		case CTRL_O:
			fmt.Printf("CTRL_O")
		case CTRL_P:
			fmt.Printf("CTRL_P")
		case CTRL_Q:
			fmt.Printf("CTRL_Q")
		case CTRL_R:
			fmt.Printf("CTRL_R")
		case CTRL_S: // SAVE
			fmt.Printf("CTRL_S")
			_ = w.Editor.SaveNew("./bin/Makefile.bak", NL_CRLF)
		case CTRL_T:
			fmt.Printf("CTRL_T")
		case CTRL_U:
			fmt.Printf("CTRL_U")
		case CTRL_V:
			fmt.Printf("CTRL_V")
		case CTRL_W:
			fmt.Printf("CTRL_W")
		case CTRL_X: // EXIT
			fmt.Printf("CTRL_X")
			return
		case CTRL_Y:
			fmt.Printf("CTRL_Y")
		case CTRL_Z:
			fmt.Printf("CTRL_Z")
		case ESC:
			fmt.Printf("ESC")
			return
		case 32:
			fmt.Printf("SPACE")
		case 127:
			fmt.Printf("BACKSPACE")
		case KEY_UP:
			fmt.Printf("KEY_UP")
			if w.Editor.Cursor.Col > 1 {
				w.Editor.Cursor.Col -= 1
			}
			fmt.Printf(" | Row: %d, Col: %d\n", w.Editor.Cursor.Col, w.Editor.Cursor.Row)
			w.MoveCursorPos(w.Editor.Cursor.Row+7, w.Editor.Cursor.Col+1)
		case KEY_DOWN:
			fmt.Printf("KEY_DOWN")
			w.Editor.Cursor.Col += 1
			fmt.Printf(" | Row: %d, Col: %d\n", w.Editor.Cursor.Col, w.Editor.Cursor.Row)
			w.MoveCursorPos(w.Editor.Cursor.Row+7, w.Editor.Cursor.Col+1)
		case KEY_RIGHT:
			fmt.Printf("KEY_RIGHT")
			w.Editor.Cursor.Row += 1
			fmt.Printf(" | Row: %d, Col: %d\n", w.Editor.Cursor.Col, w.Editor.Cursor.Row)
			w.MoveCursorPos(w.Editor.Cursor.Row+7, w.Editor.Cursor.Col+1)
		case KEY_LEFT:
			fmt.Printf("KEY_LEFT")
			if w.Editor.Cursor.Row > 1 {
				w.Editor.Cursor.Row -= 1
			}
			fmt.Printf(" | Row: %d, Col: %d\033[0m\n", w.Editor.Cursor.Col, w.Editor.Cursor.Row)
			w.MoveCursorPos(w.Editor.Cursor.Row+7, w.Editor.Cursor.Col+1)
		default:
			fmt.Printf(string(r))
		}
	}
}
