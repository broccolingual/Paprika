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

func (v *View) processInput(r rune) uint8 {
	cTab := v.Tabs[v.TabIdx] // Current Tab
	v.Reflesh('\\')
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
	case CTRL_N:
	case CTRL_O:
	case CTRL_P:
	case CTRL_Q:
	case CTRL_R: // Prev Tab
		v.PrevTab()
		v.Reflesh(r)
	case CTRL_S: // Save
		_ = cTab.SaveNew(fmt.Sprintf("./bin/%s.bak", filepath.Base(cTab.FilePath)), cTab.NL)
		cTab.IsSaved = true
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
		cTab.IsSaved = false
		v.Reflesh(r)
	case CTRL_Z: // Comment Out
	case ESC:
		return 1
	case SPACE: // Space
	case BACKSPACE: // Backspace
	case KEY_UP:
		cTab.MovePrevRow()
		v.RefleshCursorOnly(r)
	case KEY_DOWN:
		cTab.MoveNextRow()
		v.RefleshCursorOnly(r)
	case KEY_RIGHT:
	case KEY_LEFT:
	default:
		switch r { // Completion
		case rune('('):
		case rune('{'):
		case rune('['):
		case rune('\''):
		case rune('"'):
		case rune('\t'):
		default:
		}
	}
	return 0
}
