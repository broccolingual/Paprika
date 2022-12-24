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
	fmt.Printf("\033[1m%4d\033[m  \033[48;5;235m%s\033[m", lineNum, Highlighter(Tokenize(rowData), ".go", true))
}

func (w *Window) DrawUnfocusRow(lineNum int, rowData string) {
	w.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[38;5;239m%4d\033[m  %s", lineNum, Highlighter(Tokenize(rowData), ".go", false))
}

// TODO: 開始行とカーソルの位置が一致しない問題の修正
func (w *Window) DrawAll(num uint16) {
	cTab := w.Tabs[w.TabIdx]

	w.Term.InitCursorPos()
	pNode := cTab.GetRowFromNum(num)
	for i := 0; i < w.MaxRows-1; i++ {
		if pNode.IsRoot() {
			break
		}
		if pNode == cTab.CurrentNode {
			w.DrawFocusRow(i+int(num), string(pNode.Row.GetAll()))
		} else {
			w.DrawUnfocusRow(i+int(num), string(pNode.Row.GetAll()))
		}
		pNode = pNode.Next
	}
}

func (w *Window) UpdateStatusBar() {
	cTab := w.Tabs[w.TabIdx]

	w.Term.MoveCursorPos(1, uint16(w.MaxRows))
	w.Term.ClearRow()
	fmt.Print("\033[48;5;25m")
	for i := 0; i < w.MaxCols; i++ {
		fmt.Print(" ")
	}
	var nl string
	switch cTab.NL {
	case NL_CRLF:
		nl = "CRLF"
	case NL_LF:
		nl = "LF"
	default:
		nl = "Unknown"
	}
	var sf string
	switch cTab.SaveFlag {
	case true:
		sf = "Saved"
	case false:
		sf = "*Not saved"
	}
	fmt.Print("\033[m")
	w.Term.MoveCursorPos(1, uint16(w.MaxRows))
	fmt.Printf("\033[48;5;25m\033[1m %s\033[m\033[48;5;25m [%d/%d]", cTab.FilePath, w.TabIdx+1, len(w.Tabs))
	fmt.Printf(" | Ln %d, Col %d | Tab Size: %d | %s", cTab.Cursor.Row, cTab.Cursor.Col, cTab.TabSize, nl)
	fmt.Printf(" | %s", sf)
	fmt.Print("\033[m")
}

func (w *Window) Reflesh() {
	cTab := w.Tabs[w.TabIdx]

	w.Term.ClearAll()
	w.DrawAll(5)
	w.UpdateStatusBar()
	w.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}

func (w *Window) RefleshCursorOnly() {
	cTab := w.Tabs[w.TabIdx]

	w.UpdateStatusBar()
	w.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}
