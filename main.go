package main

import "os"

func main() {
	window := NewWindow(os.Args[1])
	tty := new(Tty)
	window.Term.EnableASB()
	tty.EnableRawMode()
	defer tty.DisableRawMode()
	defer window.Term.DisableASB()
	defer window.Term.EnableCursor()

	window.Editor.LoadFile()
	window.Editor.CurrentNode = window.Editor.CurrentNode.Next

	go window.readKeys()
	window.detectKeys()
}
