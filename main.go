package main

import "os"

func main() {
	tty := new(Tty)
	EnableASB()
	tty.EnableRawMode()
	defer tty.DisableRawMode()
	defer DisableASB()

	window := NewWindow(os.Args[1])
	window.Editor.LoadFile()
	window.Editor.CurrentNode = window.Editor.CurrentNode.Next

	go window.readKeys()
	window.switchKeys()
}
