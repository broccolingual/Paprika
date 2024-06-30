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

func (e *Event) Close() {
	close(e.Key)
	close(e.Signal)
}

// 入力キーの読み取り
// TODO: 入力が早すぎる場合にチャネルが閉じる問題を解決する
func (e *Event) ScanInput(exit <-chan interface{}) error {
	buf := make([]byte, 8)
	for {
		select {
			case <-exit:
				return nil
			default:
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
				} else {
					return err
				}
		}
	}
}

// OSシグナルの通知
func (e *Event) NotifySignal() {
	signal.Notify(e.Signal, os.Interrupt, syscall.SIGWINCH)
}
