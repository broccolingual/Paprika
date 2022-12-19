package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Window struct {
	Term    UnixTerm
	KeyChan chan rune
	Tabs    []*Editor
	TabIdx  int
	MaxRows int
	MaxCols int
}

// Define winsize object
type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func NewWindow() *Window {
	w := new(Window)
	w.Term = NewUnixTerm()
	w.Tabs = make([]*Editor, 0)
	w.TabIdx = 0
	w.KeyChan = make(chan rune)
	w.UpdateWinSize()
	return w
}

func (w *Window) AddTab(filePath string) {
	w.Tabs = append(w.Tabs, NewEditor(filePath, 4))
}

func (w *Window) DeleteTab() bool {
	w.Tabs = append(w.Tabs[:w.TabIdx], w.Tabs[w.TabIdx+1:]...)
	if len(w.Tabs) == 0 {
		return false
	}
	if !w.PrevTab() {
		w.NextTab()
	}
	return true
}

func (w *Window) MoveTab(idx int) bool {
	if idx >= 0 && idx < len(w.Tabs) {
		w.TabIdx = idx
		return true
	}
	return false
}

func (w *Window) NextTab() bool {
	return w.MoveTab(w.TabIdx + 1)
}

func (w *Window) PrevTab() bool {
	return w.MoveTab(w.TabIdx - 1)
}

// Get console window size
func (w *Window) UpdateWinSize() {
	var ws WinSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	w.MaxRows = int(ws.Row)
	w.MaxCols = int(ws.Col)
}

func (w *Window) DrawFocusRow(lineNum int, rowData string) {
	w.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[48;5;235m")
	for i := 0; i < w.MaxCols; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("\033[m")
	w.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[1m%4d\033[m  \033[48;5;235m%s\033[m", lineNum, rowData)
}

func (w *Window) DrawUnfocusRow(lineNum int, rowData string) {
	w.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[38;5;239m%4d\033[m  %s", lineNum, rowData)
}

func (w *Window) DrawAll() {
	w.Term.InitCursorPos()
	pNode := w.Tabs[w.TabIdx].Root
	if pNode.Prev == pNode.Next {
		return
	}
	for i := 1; i < w.MaxRows; i++ {
		pNode = pNode.Next
		if pNode.Row == nil {
			break
		}
		if pNode == w.Tabs[w.TabIdx].CurrentNode {
			w.DrawFocusRow(i, string(pNode.Row.GetAll()))
		} else {
			w.DrawUnfocusRow(i, string(pNode.Row.GetAll()))
		}
	}
}

func (w *Window) UpdateStatusBar() {
	w.Term.MoveCursorPos(1, uint16(w.MaxRows))
	w.Term.LineClear()
	fmt.Print("\033[48;5;25m")
	for i := 0; i < w.MaxCols; i++ {
		fmt.Print(" ")
	}
	var nl string
	switch w.Tabs[w.TabIdx].NL {
	case NL_CRLF:
		nl = "CRLF"
	case NL_LF:
		nl = "LF"
	default:
		nl = "Unknown"
	}
	fmt.Print("\033[m")
	w.Term.MoveCursorPos(1, uint16(w.MaxRows))
	fmt.Printf("\033[48;5;25m\033[1m %s\033[m\033[48;5;25m [%d/%d] | Ln %d, Col %d | Tab Size: %d | %s", w.Tabs[w.TabIdx].FilePath, w.TabIdx+1, len(w.Tabs), w.Tabs[w.TabIdx].Cursor.Row, w.Tabs[w.TabIdx].Cursor.Col, w.Tabs[w.TabIdx].TabSize, nl)
	fmt.Print("\033[m")
}

func (w *Window) Reflesh() {
	w.Term.ScreenClear()
	w.DrawAll()
	w.UpdateStatusBar()
	w.Term.MoveCursorPos(w.Tabs[w.TabIdx].Cursor.Col+6, w.Tabs[w.TabIdx].Cursor.Row)
}

func (w *Window) RefleshCursorOnly() {
	w.UpdateStatusBar()
	w.Term.MoveCursorPos(w.Tabs[w.TabIdx].Cursor.Col+6, w.Tabs[w.TabIdx].Cursor.Row)
}
