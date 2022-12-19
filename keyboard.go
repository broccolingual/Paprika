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
		cTab := w.Tabs[w.TabIdx] // Current Tab
		w.Term.DisableCursor()
		switch r {
		case CTRL_A: // For test
			// TODO: Backspaceの挙動
			cTab.CurrentNode.Delete()
			if cTab.Cursor.Row > 1 {
				cTab.Cursor.Row--
				cTab.MovePrevRow()
			} else {
				cTab.MoveNextRow()
			}
			cTab.Rows--
			cTab.SaveFlag = false
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
			if cTab.Cursor.Col != 1 {
				cTab.CurrentNode = cTab.CurrentNode.Insert(cTab.CurrentNode.Row.GetAll(), LINE_BUF_MAX)
				cTab.MovePrevRow()
				tmp := cTab.CurrentNode.Row.GetSize() + 1
				for i := int(cTab.Cursor.Col); i < tmp; i++ {
					cTab.CurrentNode.Row.Erase(int(cTab.Cursor.Col) - 1)
				}
				cTab.MoveNextRow()
				for i := 1; i < int(cTab.Cursor.Col); i++ {
					cTab.CurrentNode.Row.Erase(0)
				}
				cTab.Cursor.Col = 1
				cTab.Cursor.Row++
			} else {
				if cTab.CurrentNode.Row.GetSize() != 0 {
					cTab.MovePrevRow()
					cTab.CurrentNode = cTab.CurrentNode.Insert(make([]rune, 0), LINE_BUF_MAX)
				} else {
					cTab.CurrentNode = cTab.CurrentNode.Insert(make([]rune, 0), LINE_BUF_MAX)
					cTab.Cursor.Row++
				}
			}
			cTab.Rows++
			cTab.SaveFlag = false
			w.Reflesh()
		case CTRL_N:
		case CTRL_O:
		case CTRL_P:
		case CTRL_Q:
		case CTRL_R: // Prev Tab
			w.PrevTab()
			w.Reflesh()
		case CTRL_S: // Save
			_ = cTab.SaveNew(fmt.Sprintf("./bin/%s.bak", filepath.Base(cTab.FilePath)), cTab.NL)
			cTab.SaveFlag = true
			w.RefleshCursorOnly()
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
			cTab.SaveFlag = false
			w.Reflesh()
		case CTRL_Z:
		case ESC:
			return
		case 32: // Space
			cTab.Cursor.Col++
			cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), r)
			cTab.SaveFlag = false
			w.Reflesh()
		case 127: // Backspace
			// TODO: Backspaceの挙動実装, 1行目の時問題アリ
			if cTab.Cursor.Col != 1 {
				cTab.Cursor.Col--
				cTab.CurrentNode.Row.Erase(int(cTab.Cursor.Col - 1))
				w.Term.LineClear()
				w.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
				cTab.SaveFlag = false
				w.RefleshCursorOnly()
			} else {
				var tmp []rune
				if cTab.CurrentNode.Row.GetSize() != 0 {
					tmp = cTab.CurrentNode.Row.GetAll()
				}
				cTab.Cursor.Row--
				cTab.CurrentNode = cTab.CurrentNode.Delete()
				origCursorPos := uint16(cTab.CurrentNode.Row.GetSize() + 1)
				cTab.Cursor.Col = origCursorPos
				for _, ch := range tmp {
					cTab.Cursor.Col++
					cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), ch)
				}
				cTab.Cursor.Col = origCursorPos
				cTab.SaveFlag = false
				w.Reflesh()
			}
		case KEY_UP:
			w.Term.LineClear()
			w.DrawUnfocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
			if cTab.Cursor.Row > 1 {
				cTab.Cursor.Row--
				cTab.Cursor.Col = 1
				cTab.MovePrevRow()
			}
			w.Term.MoveCursorPos(1, cTab.Cursor.Row)
			w.Term.LineClear()
			w.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		case KEY_DOWN:
			w.Term.LineClear()
			w.DrawUnfocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
			if cTab.Cursor.Row <= cTab.Rows {
				cTab.Cursor.Row++
				cTab.Cursor.Col = 1
				cTab.MoveNextRow()
			}
			w.Term.MoveCursorPos(1, cTab.Cursor.Row)
			w.Term.LineClear()
			w.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
			w.RefleshCursorOnly()
		case KEY_RIGHT:
			if cTab.Cursor.Col <= uint16(cTab.CurrentNode.Row.GetSize()) {
				cTab.Cursor.Col++
			}
			w.RefleshCursorOnly()
		case KEY_LEFT:
			if cTab.Cursor.Col > 1 {
				cTab.Cursor.Col--
			}
			w.RefleshCursorOnly()
		default:
			cTab.Cursor.Col++
			cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), r)
			switch r { // Complemention
			case rune('('):
				cTab.Cursor.Col++
				cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), rune(')'))
				cTab.Cursor.Col--
			case rune('{'):
				cTab.Cursor.Col++
				cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), rune('}'))
				cTab.Cursor.Col--
			case rune('['):
				cTab.Cursor.Col++
				cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), rune(']'))
				cTab.Cursor.Col--
			case rune('\''):
				cTab.Cursor.Col++
				cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), rune('\''))
				cTab.Cursor.Col--
			case rune('"'):
				cTab.Cursor.Col++
				cTab.CurrentNode.Row.Insert(int(cTab.Cursor.Col-2), rune('"'))
				cTab.Cursor.Col--
			case rune('\t'):
				// Tab補完の実装
			default:
			}
			w.Term.LineClear()
			w.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Row.GetAll()))
			cTab.SaveFlag = false
			w.RefleshCursorOnly()
		}
		w.Term.EnableCursor()
	}
}
