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

func (v *View) DrawAllRow() {
	cTab := v.Tabs[v.TabIdx]
	v.Term.InitCursorPos()
	for i := 2; i < int(v.MaxRows); i++ {
		v.Term.MoveCursorPos(1, uint16(i))
		v.Term.SetColor(240)
		fmt.Printf("%4d", i-1)
		v.Term.ResetStyle()
		fmt.Printf("  %s", string(cTab.Lines[i-2].GetAll()))
	}
}

func (v *View) DrawAll() {
	v.DrawAllRow()
}

func (v *View) UpdateTabBar() {
	v.Term.MoveCursorPos(1, 1)
	v.Term.ClearRow()
	v.Term.SetBGColor(240)
	for i := 0; i < int(v.MaxCols); i++ {
		fmt.Print(" ")
	}
	v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, 1)
	for i, tab := range v.Tabs {
		if i == v.TabIdx {
			v.Term.ResetStyle()
			v.Term.SetBold()
			fmt.Printf(" %s |", tab.FilePath)
		} else {
			v.Term.ResetStyle()
			v.Term.SetBGColor(240)
			fmt.Printf(" %s |", tab.FilePath)
		}
	}
	v.Term.ResetStyle()
}

func (v *View) UpdateStatusBar(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.Term.MoveCursorPos(1, uint16(v.MaxRows))
	v.Term.ClearRow()
	v.Term.SetBGColor(25)
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
	v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, uint16(v.MaxRows))
	v.Term.SetBGColor(25)
	fmt.Printf(" Ln %d, Col %d | Tab Size: %d | %s", cTab.Cursor.Row, cTab.Cursor.Col, cTab.TabSize, nl)
	fmt.Printf(" | %s | Unicode %U", sf, inputRune)
	v.Term.ResetStyle()
}

func (v *View) Reflesh(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.Term.ClearAll()
	v.UpdateTabBar()
	v.DrawAll()
	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row+1)
}

func (v *View) RefleshCursorOnly(inputRune rune) {
	cTab := v.GetCurrentTab()
	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row+1)
}
