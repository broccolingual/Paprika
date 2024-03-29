package main

import (
	"fmt"
	"syscall"

	"github.com/broccolingual/Xanadu/core"
	"github.com/broccolingual/Xanadu/utils"
)

type View struct {
	Term    *core.UnixTerm
	Event	  *Event
	Tabs    []*Editor
	TabIdx  int
	WinRow	uint16
	WinCol	uint16
}

func NewView() *View {
	v := new(View)
	v.Term = core.NewUnixTerm()
	v.Event = NewEvent()
	v.Tabs = make([]*Editor, 0)
	v.TabIdx = 0
	v.WinCol = 0
	v.WinRow = 0
	return v
}

func (v *View) MainLoop() {
	e := v.Event
	e.NotifySignal() // シグナルの読み取り

	exit := make(chan interface{})
	defer close(exit)
	go e.ScanInput(exit) // キー入力の読み取り用

	Loop:
		for {
			select {
			case r := <-e.Key: // キーイベント受け取り
				exitCode := v.processInput(r)
				if exitCode != 0 {
					break Loop
				}
			case sig := <-e.Signal: // OSシグナルの受け取り
				switch sig {
					case syscall.SIGWINCH:
						v.UpdateWinSize()
						v.Reflesh()
				}
			default:
			}
		}
}

// 画面サイズの変更
func (v *View) UpdateWinSize() {
	row, col := v.Term.GetWinSize()
	v.WinRow = row
	v.WinCol = col
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

func (v *View) DrawRow(vPos uint, lineNum uint) {
	cTab := v.GetCurrentTab()
	defer v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, vPos)
	v.Term.ClearRow()
	v.Term.SetColor(240)
	fmt.Printf("%4d  ", lineNum)
	v.Term.ResetStyle()
	fmt.Printf("%s", string(cTab.Lines[lineNum-1].GetAll()))
}

func (v *View) DrawFocusRow(vPos uint, lineNum uint) {
	cTab := v.GetCurrentTab()
	defer v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, vPos)
	v.Term.ClearRow()
	v.Term.SetBGColor(235)
	for i := 0; i < int(v.WinCol); i++ {
		fmt.Print(" ")
	}
	v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, vPos)
	v.Term.SetBold()
	fmt.Printf("%4d  ", lineNum)
	v.Term.ResetStyle()
	v.Term.SetBGColor(235)
	fmt.Printf("%s", string(cTab.Lines[lineNum-1].GetAll()))
}

func (v *View) DrawAllRow() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	defer v.Term.ResetStyle()
	v.Term.InitCursorPos()
	for i := 1; i < int(v.WinRow - 1); i++ {
		cLineNum := int(cTab.ScrollRow) + i - 1
		if cLineNum >= len(cTab.Lines) {
			break
		}
		if cTab.IsTargetRow(uint(cLineNum)) {
			v.DrawFocusRow(uint(i+1), uint(cLineNum))
		} else {
			v.DrawRow(uint(i+1), uint(cLineNum))
		}
	}
}

func (v *View) UpdateTabBar() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	defer v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, 1)
	v.Term.ClearRow()
	v.Term.SetBGColor(235)
	for i := 0; i < int(v.WinCol); i++ {
		fmt.Print(" ")
	}
	v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, 1)
	for i, tab := range v.Tabs {
		if i == v.TabIdx {
			v.Term.ResetStyle()
			v.Term.SetBold()
			v.Term.SetColor(25)
			fmt.Printf(" %s ", tab.FilePath)
			v.Term.ResetStyle()
			if !tab.IsSaved {
				fmt.Print("* ")
			}
			fmt.Print("|")
		} else {
			v.Term.ResetStyle()
			v.Term.SetBGColor(235)
			fmt.Printf(" %s ", tab.FilePath)
			if !tab.IsSaved {
				fmt.Print("* ")
			}
			fmt.Print("|")
		}
	}
}

func (v *View) UpdateStatusBar() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	defer v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, uint(v.WinRow))
	v.Term.ClearRow()
	v.Term.SetBGColor(25)
	for i := 0; i < int(v.WinCol); i++ {
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
	v.Term.ResetStyle()
	v.Term.MoveCursorPos(1, uint(v.WinRow))
	v.Term.SetBGColor(25)
	fmt.Printf(" Ln %d, Col %d | Tab Size: %d | %s", cTab.Cursor.Row, cTab.Cursor.Col, cTab.TabSize, nl)
}

func (v *View) Reflesh() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	v.Term.ClearAll()
	v.UpdateTabBar()
	v.DrawAllRow()
	v.UpdateStatusBar()
}

func (v *View) RefleshTextField() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	v.Term.MoveCursorPos(1, 2)
	v.Term.ClearAfterCursor()
	v.DrawAllRow()
	v.UpdateStatusBar()
}

func (v *View) RefleshTargetRow(rowNum uint) {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
	v.Term.MoveCursorPos(1, uint(rowNum-cTab.ScrollRow+2))
	v.Term.ClearRow()
	if cTab.IsTargetRow(rowNum) {
		v.DrawFocusRow(uint(rowNum-cTab.ScrollRow+2), rowNum)
	} else {
		v.DrawRow(uint(rowNum-cTab.ScrollRow+2), rowNum)
	}
}

func (v *View) RefleshCursor() {
	cTab := v.GetCurrentTab()
	defer v.Term.MoveCursorPos(cTab.Cursor.Col+6, cTab.Cursor.Row-cTab.ScrollRow+2)
}

func (v *View) ScrollUp() {
	cTab := v.GetCurrentTab()
	prevCol := cTab.Cursor.Col
	if cTab.ScrollRow >= cTab.Cursor.Row {
		cTab.ScrollUp()
		v.Reflesh()
	} else {
		v.RefleshTargetRow(cTab.Cursor.Row + 1)
		v.RefleshTargetRow(cTab.Cursor.Row)
		v.RefleshCursor()
	}
	if prevCol > cTab.GetCurrentMaxCol() { // Cursor is on the last column
		cTab.MoveTailCol()
		v.RefleshCursor()
	}
	v.UpdateStatusBar()
}

func (v *View) ScrollDown() {
	cTab := v.GetCurrentTab()
	prevCol := cTab.Cursor.Col
	if cTab.ScrollRow + uint(v.WinRow) - 3 <= cTab.Cursor.Row {
		cTab.ScrollDown()
		v.Reflesh()
	} else {
		v.RefleshTargetRow(cTab.Cursor.Row - 1)
		v.RefleshTargetRow(cTab.Cursor.Row)
		v.RefleshCursor()
	}
	if prevCol > cTab.GetCurrentMaxCol() { // Cursor is on the last column
		cTab.MoveTailCol()
		v.RefleshCursor()
	}
	v.UpdateStatusBar()
}
