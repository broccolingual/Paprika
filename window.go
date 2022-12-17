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

func (w *Window) ClearLine(col uint16) {
	w.MoveCursorPos(0, col)
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
	pNode := w.Editor.Root
	cnt := 1
	if pNode.Prev == pNode.Next {
		return
	}
	for {
		pNode = pNode.Next
		if pNode.Row == nil {
			break
		}
		if pNode == w.Editor.CurrentNode {
			fmt.Printf("> %s\n", string(pNode.Row.GetAll()))
		} else {
			fmt.Printf("  %s\n", string(pNode.Row.GetAll()))
		}
		cnt++
	}
}
