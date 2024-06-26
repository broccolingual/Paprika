package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/broccolingual/Xanadu/core"
	"github.com/broccolingual/Xanadu/utils"
)

const LINE_BUF_MAX = 256 // 1行のバッファサイズ

// エディタ構造体
type Editor struct {
	FilePath    string          // ファイルのパス
	Cursor      *Cursor         // 現在のカーソル位置
	Lines       []*core.GapBuffer // 行リスト TODO: ファイル全体をGapBufferで管理する
	TabSize     uint8           // タブサイズ (0~255)
	NL          utils.NLCode          // 改行文字識別番号
	IsSaved     bool            // セーブ済みフラグ
	ScrollRow   uint            // 現在表示中の最上行
}

// カーソル構造体
type Cursor struct {
	Row uint
	Col uint
}

// 新しいカーソルの取得
func NewCursor() (cursor *Cursor) {
	cursor = new(Cursor)
	cursor.Row = 1
	cursor.Col = 1
	return
}

// 新しいエディタの取得
func NewEditor(filePath string, tabSize uint8) (editor *Editor) {
	editor = new(Editor)
	editor.FilePath = filePath
	editor.Cursor = NewCursor()
	editor.Lines = make([]*core.GapBuffer, 0)
	editor.TabSize = tabSize
	editor.NL = -1
	editor.IsSaved = true
	editor.ScrollRow = 1
	return
}

func (e *Editor) InsertLine(idx uint) {
	e.Lines = append(e.Lines[:idx], append([]*core.GapBuffer{core.NewGapBuffer([]rune{}, LINE_BUF_MAX)}, e.Lines[idx:]...)...)
}

func (e *Editor) DeleteLine(idx uint) {
	if idx < uint(len(e.Lines)-1) {
		copy(e.Lines[idx:], e.Lines[idx+1:])
	}
	e.Lines[len(e.Lines)-1] = nil
	e.Lines = e.Lines[:len(e.Lines)-1]
}

func (e *Editor) IsFirstRow() bool {
	return e.Cursor.Row <= 1
}

func (e *Editor) IsLastRow() bool {
	return e.Cursor.Row >= uint(len(e.Lines)) - 1
}

func (e *Editor) MoveNextRow() {
	if !e.IsLastRow() {
		e.Cursor.Row++
	}
}

func (e *Editor) MovePrevRow() {
	if !e.IsFirstRow() {
		e.Cursor.Row--
	}
}

func (e *Editor) MoveTargetRow(row uint) {
	e.Cursor.Row = row
}

func (e *Editor) MoveHeadRow() {
	e.MoveTargetRow(1)
}

func (e *Editor) MoveTailRow() {
	e.MoveTargetRow(uint(len(e.Lines))-1)
}

func (e *Editor) IsTargetRow(rowNum uint) bool {
	return rowNum == e.Cursor.Row
}

func (e *Editor) ScrollDown() {
	if !e.IsLastRow() {
		e.ScrollRow++
	}
}

func (e *Editor) ScrollUp() {
	if !e.IsFirstRow() {
		e.ScrollRow--
	}
}

func (e *Editor) ScrollTargetRow(row uint) {
	e.ScrollRow = row
}

func (e *Editor) ScrollHead() {
	e.ScrollTargetRow(1)
}

func (e *Editor) ScrollTail() {
	e.ScrollTargetRow(uint(len(e.Lines))-1)
}

func (e *Editor) IsFirstCol() bool {
	return e.Cursor.Col <= 1
}

func (e *Editor) IsLastCol() bool {
	return e.Cursor.Col > uint(e.Lines[e.Cursor.Row-1].Length())
}

func (e *Editor) MoveNextCol() {
	if !e.IsLastCol() {
		e.Cursor.Col++
	}
}

func (e *Editor) MovePrevCol() {
	if !e.IsFirstCol() {
		e.Cursor.Col--
	}
}

func (e *Editor) MoveTargetCol(col uint) {
	e.Cursor.Col = col
}

func (e *Editor) MoveHeadCol() {
	e.MoveTargetCol(1)
}

func (e *Editor) MoveTailCol() {
	e.MoveTargetCol(uint(e.Lines[e.Cursor.Row-1].Length())+1)
}

func (e *Editor) GetCurrentMaxCol() uint {
	return uint(e.Lines[e.Cursor.Row-1].Length())
}

// エディタに指定されたパスのファイルをロードして、行ノードを構成
func (e *Editor) LoadFile() {
	fp, err := os.Open(e.FilePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// conv tab to string
	var tabStr string
	for i := 0; i < int(e.TabSize); i++ {
		tabStr += " "
	}

	reader := bufio.NewReaderSize(fp, LINE_BUF_MAX)
	cnt := 0
	for {
		line, err := reader.ReadString(byte('\n'))                    // '\n'で分割
		replacedStr := strings.ReplaceAll(string(line), "\t", tabStr) // タブをスペースに変換
		replacedRune := []rune(replacedStr)
		if cnt == 0 { // 改行文字の判定
			e.NL = utils.GetNLCode(replacedRune)
		}
		switch e.NL { // 改行文字の削除
		case utils.CRLF:
			if len(replacedRune) >= 2 {
				replacedRune = replacedRune[:len(replacedRune)-2]
			}
		case utils.CR:
		case utils.LF:
			if len(replacedRune) >= 1 {
				replacedRune = replacedRune[:len(replacedRune)-1]
			}
		}

		e.Lines = append(e.Lines, core.NewGapBuffer(replacedRune, LINE_BUF_MAX))
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		cnt++
	}
}

// エディタに指定されたパスで上書き保存
func (e *Editor) SaveOverwrite(nl utils.NLCode) (saveBytes int) {
	saveBytes = e.saveFile(e.FilePath, nl)
	return
}

// 新しくファイルを保存
func (e *Editor) SaveNew(filePath string, nl utils.NLCode) (saveBytes int) {
	saveBytes = e.saveFile(filePath, nl)
	return
}

// ファイルを保存
func (e *Editor) saveFile(filePath string, nl utils.NLCode) (saveBytes int) {
	fp, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var buf string

	cnt := 0
	for _, row := range e.Lines {
		cnt++
		buf += fmt.Sprintf("%s", string(row.GetAll()))
		switch nl {
		case utils.CRLF:
			buf += "\r\n"
		case utils.CR:
			buf += "\r"
		case utils.LF:
			buf += "\n"
		default:
			buf += "\n"
		}
	}

	saveBytes, err = fp.Write([]byte(buf))
	if err != nil {
		panic(err)
	}
	return
}
