package main

import "os"

func main() {
	tty := new(Tty)
	EnableASB()
	tty.EnableRawMode()
	defer tty.DisableRawMode()
	defer DisableASB()

	w := NewWindow(os.Args[1])
	w.InitCursorPos()
	w.Rows = LoadFile(os.Args[1])

	// get key event
	go w.readKeys()

	w.switchKeys()
}
