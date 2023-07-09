package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang-text-editor/common"
	"golang-text-editor/utils"
)

const LINE_BUF_MAX = 255 // 1行のバッファサイズ

// エディタ構造体
type Editor struct {
	FilePath    string          // ファイルのパス
	Cursor      *Cursor         // 現在のカーソル位置
	Lines       []*common.GapBuffer // 行リスト
	TabSize     uint8           // タブサイズ (0~255)
	NL          utils.NLCode          // 改行文字識別番号
	IsSaved     bool            // セーブ済みフラグ
	TopLineNum  uint          // 現在表示中の最上行
}

// カーソル構造体
type Cursor struct {
	Row uint16
	Col uint16
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
	editor.Lines = make([]*common.GapBuffer, 0)
	editor.TabSize = tabSize
	editor.NL = -1
	editor.IsSaved = false
	editor.TopLineNum = 1
	return
}

// 行ノードのポインタを1つ進める
func (e *Editor) MoveNextRow() {
	if e.Cursor.Row < uint16(len(e.Lines)) {
		e.Cursor.Row++
	}
}

// 行ノードのポインタを1つ戻す
func (e *Editor) MovePrevRow() {
	if e.Cursor.Row > 1 {
		e.Cursor.Row--
	}
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
		// TODO: 改行文字が混在している時の対応
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

		e.Lines = append(e.Lines, common.NewGapBuffer(replacedRune, LINE_BUF_MAX))
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		cnt++
	}
	e.MoveNextRow()
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
