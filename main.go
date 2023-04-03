package main

import (
	"fmt"
	"os"
)

func main() {
	// 引数にファイルが指定されていない場合
	if len(os.Args) == 1 {
		fmt.Println("Error: ファイル名が指定されていません")
		os.Exit(0)
	}

	window := NewWindow()
	window.Term.EnableASB()
	window.Term.EnableRawMode()
	defer window.Term.DisableRawMode()
	defer window.Term.DisableASB()
	defer window.Term.EnableCursor()

	// 引数のパスをタブに追加
	for i, path := range os.Args {
		if i == 0 {
			continue
		}
		pathInfo, _ := os.Stat(path)
		if pathInfo.IsDir() == false {
			window.AddTab(path)
		} else {
			// ディレクトリの場合の処理
		}
	}

	// 全てのタブのファイルをロード
	for _, tab := range window.Tabs {
		tab.LoadFile()
	}

	go window.readKeys() // キー入力の読み取り用goroutine
	window.detectKeys()  // キー入力の識別
}
