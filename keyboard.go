package main

import (
	"fmt"
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

func (v *View) processInput(r rune) uint8 {
	cTab := v.GetCurrentTab() // Current Tab
	v.Term.DisableCursor()
	defer v.Term.EnableCursor()
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
		if !cTab.IsLastRow() {
			v.ScrollDown()
		}
		cTab.InsertLine(uint(cTab.Cursor.Row))
		tmp := cTab.Lines[cTab.Cursor.Row-1].GetFrom(int(cTab.Cursor.Col-1), cTab.Lines[cTab.Cursor.Row-1].GetSize())
		cTab.Lines[cTab.Cursor.Row-1].EraseFrom(int(cTab.Cursor.Col-1), cTab.Lines[cTab.Cursor.Row-1].GetSize())
		cTab.MoveNextRow()
		cTab.MoveHeadCol()
		cTab.Lines[cTab.Cursor.Row-1].AppendAll(tmp)
		v.Reflesh()
	case CTRL_N:
	case CTRL_O: // Move Top
		cTab.MoveHeadRow()
		cTab.ScrollHead()
		v.RefleshTextField()
	case CTRL_P: // Move Bottom
		cTab.MoveTailRow()
		cTab.ScrollTail()
		v.RefleshTextField()
	case CTRL_Q:
	case CTRL_R: // Prev Tab
		v.PrevTab()
		v.Reflesh()
	case CTRL_S: // Save
		_ = cTab.SaveNew(fmt.Sprintf("./bin/%s.bak", filepath.Base(cTab.FilePath)), cTab.NL)
		cTab.IsSaved = true
		v.UpdateTabBar()
	case CTRL_T: // Next Tab
		v.NextTab()
		v.Reflesh()
	case CTRL_U:
	case CTRL_V: // Paste
	case CTRL_W:
	case CTRL_X: // Exit
		return 1
	case CTRL_Y: // Delete Tab
		if !v.DeleteTab() {
			return 1
		}
		v.Reflesh()
	case CTRL_Z: // Comment Out
	case ESC:
		return 1
	case SPACE: // Space
		cTab.IsSaved = false
		cTab.Lines[cTab.Cursor.Row-1].Insert(int(cTab.Cursor.Col-1), rune(' '))
		cTab.MoveNextCol()
		v.RefleshTargetRow(cTab.Cursor.Row)
		v.RefleshCursor()
		v.UpdateTabBar()
	case BACKSPACE: // Backspace
		if !cTab.IsFirstRow() {
			v.ScrollUp()
		}
		if !cTab.Lines[cTab.Cursor.Row-1].IsEmpty() && cTab.Cursor.Col != 1 { // 行に何か入力されている場合
			cTab.IsSaved = false
			cTab.Lines[cTab.Cursor.Row-1].Erase(int(cTab.Cursor.Col-2))
			cTab.MovePrevCol()
			v.RefleshTargetRow(cTab.Cursor.Row)
			v.RefleshCursor()
			v.UpdateTabBar()
		} else { // 行に何も入力されていない場合
			tmp := cTab.Lines[cTab.Cursor.Row-1].GetAll()
			cTab.DeleteLine(uint(cTab.Cursor.Row-1))
			cTab.MovePrevRow()
			cTab.MoveTailCol()
			cTab.Lines[cTab.Cursor.Row-1].AppendAll(tmp)
			v.Reflesh()
		}
	case KEY_UP: // Scroll Up
		if !cTab.IsFirstRow() {
			cTab.MovePrevRow()
			v.ScrollUp()
		}
	case KEY_DOWN: // Scroll Down
		if !cTab.IsLastRow() {
			cTab.MoveNextRow()
			v.ScrollDown()
		}
	case KEY_RIGHT:
		if !cTab.IsLastCol() {
			cTab.MoveNextCol()
			v.RefleshCursor()
			v.UpdateStatusBar()
		}
	case KEY_LEFT:
		if !cTab.IsFirstCol() {
			cTab.MovePrevCol()
			v.RefleshCursor()
			v.UpdateStatusBar()
		}
	default:
		cTab.IsSaved = false
		cTab.Lines[cTab.Cursor.Row-1].Insert(int(cTab.Cursor.Col-1), r)
		cTab.MoveNextCol()
		v.RefleshTargetRow(cTab.Cursor.Row)
		v.RefleshCursor()
		v.UpdateTabBar()
	}
	return 0
}
