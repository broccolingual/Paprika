package main

import (
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
	w.Draw()
	w.UpdateStatusBar()
	for {
		r := <-w.KeyChan
		switch r {
		case CTRL_A:
		case CTRL_B:
		case CTRL_C:
		case CTRL_D:
		case CTRL_E:
		case CTRL_F:
		case CTRL_G:
		case CTRL_H:
		case CTRL_I:
		case CTRL_J:
		case CTRL_K:
		case CTRL_L:
		case CTRL_M: // Enter
			w.Editor.Cursor.Col += 1
		case CTRL_N:
		case CTRL_O:
		case CTRL_P:
		case CTRL_Q:
		case CTRL_R:
		case CTRL_S: // Save
			_ = w.Editor.SaveNew("./bin/Makefile.bak", NL_CRLF)
		case CTRL_T:
		case CTRL_U:
		case CTRL_V:
		case CTRL_W:
		case CTRL_X: // Exit
			return
		case CTRL_Y:
		case CTRL_Z:
		case ESC:
			return
		case 32: // Space
			w.Editor.Cursor.Row += 1
			w.Editor.CurrentNode.Row.Insert(int(w.Editor.Cursor.Row-2), r)
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		case 127: // Backspace
			if w.Editor.Cursor.Row > 1 {
				w.Editor.Cursor.Row -= 1
				w.Editor.CurrentNode.Row.Erase(int(w.Editor.Cursor.Row - 1))
			}
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		case KEY_UP:
			if w.Editor.Cursor.Col > 1 {
				w.Editor.Cursor.Col -= 1
				w.Editor.Cursor.Row = 1
				w.Editor.CurrentNode = w.Editor.CurrentNode.Prev
			}
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		case KEY_DOWN:
			w.Editor.Cursor.Col += 1
			w.Editor.Cursor.Row = 1
			w.Editor.CurrentNode = w.Editor.CurrentNode.Next
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		case KEY_RIGHT:
			if w.Editor.Cursor.Row <= uint16(w.Editor.CurrentNode.Row.GetSize()) {
				w.Editor.Cursor.Row += 1
			}
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		case KEY_LEFT:
			if w.Editor.Cursor.Row > 1 {
				w.Editor.Cursor.Row -= 1
			}
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		default:
			w.Editor.Cursor.Row += 1
			w.Editor.CurrentNode.Row.Insert(int(w.Editor.Cursor.Row-2), r)
			w.Clear()
			w.Draw()
			w.UpdateStatusBar()
			w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
		}
	}
}
