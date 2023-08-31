package main

import (
	"os"
	"os/signal"
	"syscall"
)

type Event struct {
	Key chan rune
	Signal chan os.Signal
}

func NewEvent() *Event {
	e := new(Event)
	e.Key = make(chan rune, 1)
	e.Signal = make(chan os.Signal, 1)
	return e
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

// OSシグナルの通知
func (e *Event) NotifySignal() {
	signal.Notify(e.Signal, os.Interrupt, syscall.SIGWINCH)
}
