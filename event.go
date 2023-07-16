package main

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
	"time"
)

const (
	INTERVAL_UPDATE_WINSIZE = time.Millisecond * 500
)

type Event struct {
	Key chan rune
	WindowSize chan WinSize
	Signal chan os.Signal
}

func NewEvent() *Event {
	e := new(Event)
	e.Key = make(chan rune, 1)
	e.WindowSize = make(chan WinSize, 1)
	e.Signal = make(chan os.Signal, 1)
	return e
}

// ウィンドウサイズオブジェクトの定義
type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// 入力キーの読み取り
func (e *Event) ScanInput() {
	buf := make([]byte, 64)
	for {
		if n, err := os.Stdin.Read(buf); err == nil {
			b := buf[:n]
			for {
				r, n := parseKey(b)
				if n == 0 {
					break
				}
				e.Key <- r
				b = b[n:]
			}
		}
	}
}

// コンソールのウィンドウサイズの取得
func (e *Event) GetWinSize() {
	for {
		var ws WinSize
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
		e.WindowSize <- ws
		time.Sleep(INTERVAL_UPDATE_WINSIZE)
	}
}

// OSシグナルの通知
func (e *Event) NotifySignal() {
	signal.Notify(e.Signal, os.Interrupt)
}
