package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Window struct {
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
	w.Editor = NewEditor(filePath)
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

func (w *Window) Clear() {
	syscall.Write(0, []byte("\033[2J"))
}

func (w *Window) ClearLine() {
	syscall.Write(0, []byte("\033[2K"))
}

func (w *Window) InitCursorPos() {
	syscall.Write(0, []byte("\033[1;1H"))
}

// row: 1~, col: 1~
func (w *Window) MoveCursorPos(row uint16, col uint16) {
	syscall.Write(0, []byte(fmt.Sprintf("\033[%d;%dH", col, row)))
}

func (w *Window) Draw() {
	w.InitCursorPos()
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
			fmt.Printf("> %s\n", string(pNode.Row.GetAll()))
		} else {
			fmt.Printf("  %s\n", string(pNode.Row.GetAll()))
		}
	}
}

func (w *Window) UpdateStatusBar() {
	w.MoveCursorPos(1, uint16(w.MaxRows))
	w.ClearLine()
	tmp := fmt.Sprintf("\033[7m%s | Ln %d, Col %d", w.Editor.FilePath, w.Editor.Cursor.Col, w.Editor.Cursor.Row)
	fmt.Printf(tmp)
	for i := 0; i < w.MaxCols-len(tmp); i++ {
		fmt.Print(" ")
	}
	fmt.Print("\033[m")
}

func (w *Window) Reflesh() {
	w.Clear()
	w.Draw()
	w.UpdateStatusBar()
	w.MoveCursorPos(w.Editor.Cursor.Row+2, w.Editor.Cursor.Col)
}
