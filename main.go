package main

import "os"

func main() {
	window := NewWindow()
	tty := new(Tty)
	window.Term.EnableASB()
	tty.EnableRawMode()
	defer tty.DisableRawMode()
	defer window.Term.DisableASB()
	defer window.Term.EnableCursor()

	window.AddTab(os.Args[1])
	window.AddTab(os.Args[2])

	for _, tab := range window.Tabs {
		tab.LoadFile()
		tab.CurrentNode = tab.CurrentNode.Next
	}

	go window.readKeys()
	window.detectKeys()
}
