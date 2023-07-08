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
	MaxRows int
	MaxCols int
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
	go e.UpdateWinSize() // 画面サイズの更新
	go e.NotifySignal() // シグナルの読み取り

	v.Reflesh('\\')

	Loop:
		for {
			select {
			case r := <-e.Key: // キーイベント受け取り
				exitCode := v.processInput(r)
				if exitCode != 0 {
					break Loop
				}
			// TODO: キー入力がないと画面サイズが更新されない問題の修正
			case ws := <-e.WindowSize: // 画面サイズ変更イベント受け取り
				v.ChangeWinSize(ws)
			// TODO: 強制終了時の処理
			case <- e.Signal: // OSシグナルの受け取り
			}
		}
}

func (v *View) ChangeWinSize(ws WinSize) {
	v.MaxRows = int(ws.Row)
	v.MaxCols = int(ws.Col)
}

func (v *View) AddTab(filePath string) {
	v.Tabs = append(v.Tabs, NewEditor(filePath, 4))
}

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

func (v *View) MoveTab(idx int) bool {
	if idx >= 0 && idx < len(v.Tabs) {
		v.TabIdx = idx
		return true
	}
	return false
}

func (v *View) NextTab() bool {
	return v.MoveTab(v.TabIdx + 1)
}

func (v *View) PrevTab() bool {
	return v.MoveTab(v.TabIdx - 1)
}

func (v *View) DrawFocusRow(lineNum int, rowData string) {
	v.Term.MoveCursorPos(1, uint16(lineNum))
	fmt.Printf("\033[48;5;235m")
	for i := 0; i < v.MaxCols; i++ {
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
	startRowNum := cTab.TopRowNum
	pNode := cTab.GetNodeFromLineNum(startRowNum)

	v.Term.InitCursorPos()
	for i := 1; i < v.MaxRows; i++ {
		if pNode.IsRoot() {
			return
		}
		v.Term.MoveCursorPos(1, uint16(i))
		fmt.Printf("\033[38;5;239m%4d\033[m  %s", int(startRowNum)+i-1, Highlighter(Tokenize(string(pNode.GetBuf().GetAll())), ".go", false))
		pNode = pNode.Next()
	}
}

// TODO: 開始行とカーソルの位置が一致しない問題の修正
func (v *View) DrawAll() {
	v.DrawAllRow()
}

func (v *View) UpdateStatusBar(inputRune rune) {
	cTab := v.Tabs[v.TabIdx]

	v.Term.MoveCursorPos(1, uint16(v.MaxRows))
	v.Term.ClearRow()
	fmt.Print("\033[48;5;25m")
	for i := 0; i < v.MaxCols; i++ {
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
	cTab := v.Tabs[v.TabIdx]

	v.Term.ClearAll()
	v.DrawAll()
	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}

func (v *View) RefleshCursorOnly(inputRune rune) {
	cTab := v.Tabs[v.TabIdx]

	v.UpdateStatusBar(inputRune)
	v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row)
}
