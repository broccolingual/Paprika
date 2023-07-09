package main

import (
	"fmt"

	"golang-text-editor/common"
	"golang-text-editor/utils"
)

type View struct {
	Term    *common.UnixTerm
	Event	  *Event
	Tabs    []*Editor
	TabIdx  int
	MaxRows uint16
	MaxCols uint16
}

func NewView() *View {
	v := new(View)
	v.Term = common.NewUnixTerm()
	v.Tabs = make([]*Editor, 0)
	v.TabIdx = 0
	v.Event = NewEvent()
	return v
}

func (v *View) MainLoop() {
	e := v.Event
	go e.ScanInput() // キー入力の読み取り用
	go e.GetWinSize() // 画面サイズの更新
	go e.NotifySignal() // シグナルの読み取り

	Loop:
		for {
			select {
			case r := <-e.Key: // キーイベント受け取り
				exitCode := v.processInput(r)
				if exitCode != 0 {
					break Loop
				}
			case ws := <-e.WindowSize: // 画面サイズ変更イベント受け取り
				if v.IsWinSizeChanged(ws) {
					v.ChangeWinSize(ws)
					v.Reflesh('\\')
				}
			case <- e.Signal: // OSシグナルの受け取り
				// termiosでシグナルを無効化しているため動作しない
			}
		}
}

// 画面サイズが変更されたかを検知
func (v *View) IsWinSizeChanged(ws WinSize) (isChanged bool) {
	isChanged = false
	if v.MaxRows != ws.Row || v.MaxCols != ws.Col {
		isChanged = true
	}
	return
}

// 画面サイズの変更
func (v *View) ChangeWinSize(ws WinSize) {
	v.MaxRows = ws.Row
	v.MaxCols = ws.Col
}

// タブの追加
func (v *View) AddTab(filePath string) {
	v.Tabs = append(v.Tabs, NewEditor(filePath, 4))
}

// タブの削除
func (v *View) DeleteTab() bool {
	v.Tabs = append(v.Tabs[:v.TabIdx], v.Tabs[v.TabIdx+1:]...)
	if len(v.Tabs) == 0 {
		return false
	}
	if !v.PrevTab() {
		v.NextTab()
	}
	return true
}

// 指定したインデックスのタブに移動
func (v *View) MoveTab(idx int) bool {
	if idx >= 0 && idx < len(v.Tabs) {
		v.TabIdx = idx
		return true
	}
	return false
}

// 次のタブへ移動
func (v *View) NextTab() bool {
	return v.MoveTab(v.TabIdx + 1)
}

// 前のタブへ移動
func (v *View) PrevTab() bool {
	return v.MoveTab(v.TabIdx - 1)
}

// 現在のタブのオブジェクトの取得
func (v *View) GetCurrentTab() *Editor {
	return v.Tabs[v.TabIdx]
}

// TODO: 描画関数の修正
func (v *View) DrawFocusRow(lineNum int, rowData string) {
	v.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[48;5;235m")
	for i := 0; i < int(v.MaxCols); i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("\033[m")
	v.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[1m%4d\033[m  \033[48;5;235m%s\033[m", lineNum, Highlighter(Tokenize(rowData), ".go", true))
}

func (v *View) DrawUnfocusRow(lineNum int, rowData string) {
	v.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[38;5;239m%4d\033[m  %s", lineNum, Highlighter(Tokenize(rowData), ".go", false))
}

func (v *View) DrawAllRow() {
	cTab := v.Tabs[v.TabIdx]
	v.Term.InitCursorPos()
	for i := 1; i < int(v.MaxRows); i++ {
		v.Term.MoveCursorPos(1, uint16(i))
		fmt.Printf("\033[38;5;239m%4d\033[m  %s", i, string(cTab.Lines[i-1].GetAll()))
	}
}

// TODO: 開始行とカーソルの位置が一致しない問題の修正
func (v *View) DrawAll() {
	v.DrawAllRow()
}

func (v *View) UpdateStatusBar(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.Term.MoveCursorPos(1, uint16(v.MaxRows))
	v.Term.ClearRow()
	fmt.Print("\033[48;5;25m")
	for i := 0; i < int(v.MaxCols); i++ {
		fmt.Print(" ")
	}
	var nl string
	switch cTab.NL {
	case utils.CRLF:
		nl = "CRLF"
	case utils.CR:
		nl = "CR"
	case utils.LF:
		nl = "LF"
	default:
		nl = "Unknown"
	}
	var sf string
	switch cTab.IsSaved {
	case true:
		sf = "Saved"
	case false:
		sf = "*Not saved"
	}
	fmt.Print("\033[m")
	v.Term.MoveCursorPos(1, uint16(v.MaxRows))
	fmt.Printf("\033[48;5;25m\033[1m %s\033[m\033[48;5;25m [%d/%d]", cTab.FilePath, v.TabIdx+1, len(v.Tabs))
	fmt.Printf(" | Ln %d, Col %d | Tab Size: %d | %s", cTab.Cursor.Row, cTab.Cursor.Col, cTab.TabSize, nl)
	fmt.Printf(" | %s | Unicode %U", sf, inputRune)
	fmt.Print("\033[m")
}

func (v *View) Reflesh(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.Term.ClearAll()
	v.DrawAll()
	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}

func (v *View) RefleshCursorOnly(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}
