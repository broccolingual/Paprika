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

	view := NewView()
	view.Term.EnableAlternativeScreenBuffer()
	view.Term.EnableRawMode()
	defer view.Term.DisableRawMode()
	defer view.Term.DisableAlternativeScreenBuffer()
	defer view.Term.EnableCursor()

	// 引数のパスをタブに追加
	for i, path := range os.Args {
		if i == 0 {
			continue
		}
		pathInfo, _ := os.Stat(path)
		if pathInfo.IsDir() == false {
			view.AddTab(path)
		} else {
			// ディレクトリの場合の処理
		}
	}

	// 全てのタブのファイルをロード
	for _, tab := range view.Tabs {
		tab.LoadFile()
	}
	view.MainLoop() //メインループ
}
