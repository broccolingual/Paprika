package main

import "os"

func main() {
	window := NewWindow()
	window.Term.EnableASB()
	window.Term.EnableRawMode()
	defer window.Term.DisableRawMode()
	defer window.Term.DisableASB()
	defer window.Term.EnableCursor()

	for i, path := range os.Args {
		if i != 0 {
			window.AddTab(path)
		}
	}

	for _, tab := range window.Tabs {
		tab.LoadFile()
	}

	go window.readKeys()
	window.detectKeys()
}
