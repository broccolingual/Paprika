package main

import (
	"os"
	"syscall"
	"unsafe"
	"time"
)

type Event struct {
	Key chan rune
	WindowSize chan WinSize
}

func NewEvent() *Event {
	e := new(Event)
	e.Key = make(chan rune)
	e.WindowSize = make(chan WinSize)
	return e
}

// Define winsize object
type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Get console window size
func (e *Event) UpdateWinSize() {
	for {
		var ws WinSize
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
		e.WindowSize <- ws
		time.Sleep(time.Millisecond * 100)
	}
}
