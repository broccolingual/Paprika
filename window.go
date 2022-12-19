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
	Editor  *Editor
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

func NewWindow(filePath string) *Window {
	w := new(Window)
	w.Term = NewUnixTerm()
	w.Editor = NewEditor(filePath, 4)
	w.KeyChan = make(chan rune)
	w.UpdateWinSize()
	return w
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
	fmt.Printf("%4d  %s", lineNum, rowData)
}

func (w *Window) DrawAll() {
	w.Term.InitCursorPos()
	pNode := w.Editor.Root
	if pNode.Prev == pNode.Next {
		return
	}
	for i := 1; i < w.MaxRows; i++ {
		pNode = pNode.Next
		if pNode.Row == nil {
			break
		}
		if pNode == w.Editor.CurrentNode {
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
	switch w.Editor.NL {
	case NL_CRLF:
		nl = "CRLF"
	case NL_LF:
		nl = "LF"
	default:
		nl = "Unknown"
	}
	fmt.Print("\033[m")
	w.Term.MoveCursorPos(1, uint16(w.MaxRows))
	fmt.Printf("\033[48;5;25m\033[1m %s\033[m\033[48;5;25m | Ln %d, Col %d | Tab Size: %d | %s", w.Editor.FilePath, w.Editor.Cursor.Row, w.Editor.Cursor.Col, w.Editor.TabSize, nl)
	fmt.Print("\033[m")
}

func (w *Window) Reflesh() {
	w.Term.ScreenClear()
	w.DrawAll()
	w.UpdateStatusBar()
	w.Term.MoveCursorPos(w.Editor.Cursor.Col+6, w.Editor.Cursor.Row)
}

func (w *Window) RefleshCursorOnly() {
	w.UpdateStatusBar()
	w.Term.MoveCursorPos(w.Editor.Cursor.Col+6, w.Editor.Cursor.Row)
}
