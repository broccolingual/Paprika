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
	Lines       *core.DoublyLinkedList // 行リスト
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
	editor.CurrentLine = *core.NewDllRoot()
	editor.TabSize = tabSize
	editor.NL = -1
	editor.IsSaved = true
	editor.ScrollRow = 1
	return
}

func (e *Editor) InsertLine(line *core.DoublyLinkedList) {
	line.Insert(*core.NewGapBuffer([]rune{}, LINE_BUF_MAX))
}

func (e *Editor) DeleteLine(line *core.DoublyLinkedList) {
	line.Remove()
}

func (e *Editor) IsFirstRow() bool {
	return e.CurrentLine.IsRoot()
}

func (e *Editor) IsLastRow() bool {
	return e.CurrentLine.Next().IsRoot()
}

func (e *Editor) MoveNextRow() {
	if !e.IsLastRow() {
		e.CurrentLine = e.CurrentLine.Next()
	}
}

func (e *Editor) MovePrevRow() {
	if !e.IsFirstRow() {
		e.CurrentLine = e.CurrentLine.Prev()
	}
}

func (e *Editor) MoveTargetRow(row uint) {
	e.Cursor.Row = row
}

func (e *Editor) MoveHeadRow() {
	e.CurrentLine = e.CurrentLine.Root()
}

func (e *Editor) MoveTailRow() {
	e.CurrentLine.Prev()
}

func (e *Editor) IsTargetRow(rowNum uint) bool {
	return rowNum == e.CurrentLine.GetIdx() + 1
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
	e.MoveTargetCol(GetCurrentMaxCol()+1)
}

func (e *Editor) GetCurrentMaxCol() uint {
	return uint(e.CurrentLine.GetBuf().Length())
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

		e.CurrentLine.Insert(*core.NewGapBuffer(replacedRune, LINE_BUF_MAX))
		e.CurrentLine = e.CurrentLine.Next()
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

	for e.CurrentLine.IsRoot() == false {
		buf += fmt.Sprintf("%s", string(e.CurrentLine.GetBuf().GetAll()))
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
		e.CurrentLine = e.CurrentLine.Next()
	}

	saveBytes, err = fp.Write([]byte(buf))
	if err != nil {
		panic(err)
	}
	return
}
