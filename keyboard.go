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
	SPACE     = 32
	BACKSPACE = 127
	KEY_UP    = 1001
	KEY_DOWN  = 1002
	KEY_RIGHT = 1003
	KEY_LEFT  = 1004
)

var COMPLETION_LIST map[rune]rune = map[rune]rune{
	'{':  '}',
	'[':  ']',
	'(':  ')',
	'"':  '"',
	'\'': '\'',
}

func parseKey(b []byte) (rune, int) {
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

func (v *View) scanInput() {
	buf := make([]byte, 64)
	for {
		if n, err := os.Stdin.Read(buf); err == nil {
			b := buf[:n]
			for {
				r, n := parseKey(b)
				if n == 0 {
					break
				}
				v.Event.Key <- r
				b = b[n:]
			}
		}
	}
}

func (v *View) processInput(r rune) uint8 {
	v.Reflesh('\\')
	cTab := v.Tabs[v.TabIdx] // Current Tab
	v.Term.DisableCursor()
	switch r {
	case CTRL_A:
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
			cTab.CurrentNode = cTab.CurrentNode.Insert(cTab.CurrentNode.Buf.GetAll(), LINE_BUF_MAX)
			cTab.MovePrevRow()
			tmp := cTab.CurrentNode.Buf.GetSize() + 1
			for i := int(cTab.Cursor.Col); i < tmp; i++ {
				cTab.CurrentNode.Buf.Erase(int(cTab.Cursor.Col) - 1)
			}
			cTab.MoveNextRow()
			for i := 1; i < int(cTab.Cursor.Col); i++ {
				cTab.CurrentNode.Buf.Erase(0)
			}
			cTab.Cursor.Col = 1
			cTab.Cursor.Row++
		} else {
			if cTab.CurrentNode.Buf.GetSize() != 0 {
				cTab.MovePrevRow()
				cTab.CurrentNode = cTab.CurrentNode.Insert(make([]rune, 0), LINE_BUF_MAX)
			} else {
				cTab.CurrentNode = cTab.CurrentNode.Insert(make([]rune, 0), LINE_BUF_MAX)
				cTab.Cursor.Row++
			}
		}
		cTab.Rows++
		cTab.SaveFlag = false
		v.Reflesh(r)
	case CTRL_N:
	case CTRL_O:
	case CTRL_P:
	case CTRL_Q:
	case CTRL_R: // Prev Tab
		v.PrevTab()
		v.Reflesh(r)
	case CTRL_S: // Save
		_ = cTab.SaveNew(fmt.Sprintf("./bin/%s.bak", filepath.Base(cTab.FilePath)), cTab.NL)
		cTab.SaveFlag = true
		v.RefleshCursorOnly(r)
	case CTRL_T: // Next Tab
		v.NextTab()
		v.Reflesh(r)
	case CTRL_U:
	case CTRL_V: // Paste
	case CTRL_W:
	case CTRL_X: // Exit
		return 1
	case CTRL_Y: // Delete Tab
		if !v.DeleteTab() {
			return 1
		}
		cTab.SaveFlag = false
		v.Reflesh(r)
	case CTRL_Z: // Comment Out
		if cTab.CurrentNode.Buf.Check(0, []rune("// ")) {
			cTab.CurrentNode.Buf.EraseAll(0, 3)
			cTab.Cursor.Col -= 3
		} else {
			cTab.CurrentNode.Buf.InsertAll(0, []rune("// "))
			cTab.Cursor.Col += 3
		}
		v.Term.ClearRow()
		v.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
		cTab.SaveFlag = false
		v.RefleshCursorOnly(r)
	case ESC:
		return 1
	case SPACE: // Space
		cTab.Cursor.Col++
		cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), r)
		cTab.SaveFlag = false
		v.Reflesh(r)
	case BACKSPACE: // Backspace
		if cTab.Cursor.Col != 1 {
			cTab.Cursor.Col--
			cTab.CurrentNode.Buf.Erase(int(cTab.Cursor.Col - 1))
			v.Term.ClearRow()
			v.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
			cTab.SaveFlag = false
			v.RefleshCursorOnly(r)
		} else {
			if cTab.Cursor.Row != 1 {
				var tmp []rune
				if cTab.CurrentNode.Buf.GetSize() != 0 {
					tmp = cTab.CurrentNode.Buf.GetAll()
				}
				cTab.Cursor.Row--
				cTab.CurrentNode = cTab.CurrentNode.Delete()
				origCursorPos := uint16(cTab.CurrentNode.Buf.GetSize() + 1)
				cTab.Cursor.Col = origCursorPos
				for _, ch := range tmp {
					cTab.Cursor.Col++
					cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), ch)
				}
				cTab.Cursor.Col = origCursorPos
				cTab.SaveFlag = false
				v.Reflesh(r)
			}
		}
	case KEY_UP:
		v.Term.ClearRow()
		v.DrawUnfocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
		if cTab.Cursor.Row > 1 {
			cTab.Cursor.Row--
			cTab.Cursor.Col = 1
			cTab.MovePrevRow()
		}
		v.Term.MoveCursorPos(1, cTab.Cursor.Row)
		v.Term.ClearRow()
		v.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
		v.RefleshCursorOnly(r)
	case KEY_DOWN:
		v.Term.ClearRow()
		if cTab.GetLineNumFromNode(cTab.CurrentNode) >= cTab.TopRowNum+uint16(v.MaxRows)-2 {
			cTab.TopRowNum++
			v.Reflesh(r)
		} else {
			v.DrawUnfocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
			if cTab.Cursor.Row <= cTab.Rows {
				cTab.Cursor.Row++
				cTab.Cursor.Col = 1
				cTab.MoveNextRow()
			}
			v.Term.MoveCursorPos(1, cTab.Cursor.Row)
			v.Term.ClearRow()
			v.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
			v.RefleshCursorOnly(r)
		}
	case KEY_RIGHT:
		if cTab.Cursor.Col <= uint16(cTab.CurrentNode.Buf.GetSize()) {
			cTab.Cursor.Col++
		}
		v.RefleshCursorOnly(r)
	case KEY_LEFT:
		if cTab.Cursor.Col > 1 {
			cTab.Cursor.Col--
		}
		v.RefleshCursorOnly(r)
	default:
		cTab.Cursor.Col++
		cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), r)
		switch r { // Completion
		case rune('('):
			cTab.Cursor.Col++
			cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), COMPLETION_LIST[r])
			cTab.Cursor.Col--
		case rune('{'):
			cTab.Cursor.Col++
			cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), COMPLETION_LIST[r])
			cTab.Cursor.Col--
		case rune('['):
			cTab.Cursor.Col++
			cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), COMPLETION_LIST[r])
			cTab.Cursor.Col--
		case rune('\''):
			cTab.Cursor.Col++
			cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), COMPLETION_LIST[r])
			cTab.Cursor.Col--
		case rune('"'):
			cTab.Cursor.Col++
			cTab.CurrentNode.Buf.Insert(int(cTab.Cursor.Col-2), COMPLETION_LIST[r])
			cTab.Cursor.Col--
		case rune('\t'):
			// TODO: Tab補完の実装 (※反応せず)
			cTab.Cursor.Col += uint16(cTab.TabSize)
		default:
		}
		v.Term.ClearRow()
		v.DrawFocusRow(int(cTab.Cursor.Row), string(cTab.CurrentNode.Buf.GetAll()))
		cTab.SaveFlag = false
		v.RefleshCursorOnly(r)
	}
	v.Term.EnableCursor()
	return 0
}
