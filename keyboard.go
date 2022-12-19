package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func (w *Window) detectKeys() {
	w.Reflesh()
	for {
		r := <-w.KeyChan
		w.Term.DisableCursor()
		switch r {
		case CTRL_A: // For test
			w.Tabs[w.TabIdx].CurrentNode.Delete()
			if w.Tabs[w.TabIdx].Cursor.Row > 1 {
				w.Tabs[w.TabIdx].Cursor.Row -= 1
				w.Tabs[w.TabIdx].MovePrevRow()
			} else {
				w.Tabs[w.TabIdx].MoveNextRow()
			}
			w.Reflesh()
		case CTRL_B:
		case CTRL_C: // Copy
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
			w.Tabs[w.TabIdx].CurrentNode = w.Tabs[w.TabIdx].CurrentNode.Insert(make([]rune, 0), LINE_BUF_MAX)
			w.Tabs[w.TabIdx].Cursor.Row += 1
			w.Reflesh()
		case CTRL_N:
		case CTRL_O:
		case CTRL_P:
		case CTRL_Q:
		case CTRL_R: // Prev Tab
			w.PrevTab()
			w.Reflesh()
		case CTRL_S: // Save
			_ = w.Tabs[w.TabIdx].SaveNew(fmt.Sprintf("./bin/%s.bak", filepath.Base(w.Tabs[w.TabIdx].FilePath)), w.Tabs[w.TabIdx].NL)
		case CTRL_T: // Next Tab
			w.NextTab()
			w.Reflesh()
		case CTRL_U:
		case CTRL_V: // Paste
		case CTRL_W:
		case CTRL_X: // Exit
			return
		case CTRL_Y: // Delete Tab
			if !w.DeleteTab() {
				return
			}
			w.Reflesh()
		case CTRL_Z:
		case ESC:
			return
		case 32: // Space
			w.Tabs[w.TabIdx].Cursor.Col += 1
			w.Tabs[w.TabIdx].CurrentNode.Row.Insert(int(w.Tabs[w.TabIdx].Cursor.Col-2), r)
			w.Reflesh()
		case 127: // Backspace
			if w.Tabs[w.TabIdx].Cursor.Col > 1 {
				w.Tabs[w.TabIdx].Cursor.Col -= 1
				w.Tabs[w.TabIdx].CurrentNode.Row.Erase(int(w.Tabs[w.TabIdx].Cursor.Col - 1))
			}
			w.Term.LineClear()
			w.DrawFocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		case KEY_UP:
			w.Term.LineClear()
			w.DrawUnfocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			if w.Tabs[w.TabIdx].Cursor.Row > 1 {
				w.Tabs[w.TabIdx].Cursor.Row -= 1
				w.Tabs[w.TabIdx].Cursor.Col = 1
				w.Tabs[w.TabIdx].MovePrevRow()
			}
			w.Term.MoveCursorPos(1, w.Tabs[w.TabIdx].Cursor.Row)
			w.Term.LineClear()
			w.DrawFocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		case KEY_DOWN:
			w.Term.LineClear()
			w.DrawUnfocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			if w.Tabs[w.TabIdx].Cursor.Row <= w.Tabs[w.TabIdx].Rows {
				w.Tabs[w.TabIdx].Cursor.Row += 1
				w.Tabs[w.TabIdx].Cursor.Col = 1
				w.Tabs[w.TabIdx].MoveNextRow()
			}
			w.Term.MoveCursorPos(1, w.Tabs[w.TabIdx].Cursor.Row)
			w.Term.LineClear()
			w.DrawFocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		case KEY_RIGHT:
			if w.Tabs[w.TabIdx].Cursor.Col <= uint16(w.Tabs[w.TabIdx].CurrentNode.Row.GetSize()) {
				w.Tabs[w.TabIdx].Cursor.Col += 1
			}
			w.RefleshCursorOnly()
		case KEY_LEFT:
			if w.Tabs[w.TabIdx].Cursor.Col > 1 {
				w.Tabs[w.TabIdx].Cursor.Col -= 1
			}
			w.RefleshCursorOnly()
		default:
			w.Tabs[w.TabIdx].Cursor.Col += 1
			w.Tabs[w.TabIdx].CurrentNode.Row.Insert(int(w.Tabs[w.TabIdx].Cursor.Col-2), r)
			w.Term.LineClear()
			w.DrawFocusRow(int(w.Tabs[w.TabIdx].Cursor.Row), string(w.Tabs[w.TabIdx].CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		}
		w.Term.EnableCursor()
	}
}
