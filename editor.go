package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang-text-editor/common"
)

type nlCode int // 改行文字タイプの定義

const (
	NL_LF nlCode = iota
	NL_CRLF
)

const LINE_BUF_MAX = 255 // 1行のバッファサイズ

// エディタ構造体
type Editor struct {
	FilePath    string          // ファイルのパス
	Cursor      *Cursor         // 現在のカーソル位置
	Root        *common.RowNode // 行ノードのルート(ダミーノード)
	CurrentNode *common.RowNode // 現在の行ノード
	TabSize     uint8           // タブサイズ (0~255)
	NL          nlCode          // 改行文字識別番号
	Rows        uint16          // ファイルの行数 (65534行まで)
	SaveFlag    bool            // セーブ済みフラグ
	TopRowNum   uint16          // 現在表示中の最上行
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
	editor.Root = common.NewRootNode()
	editor.CurrentNode = editor.Root
	editor.TabSize = tabSize
	editor.NL = -1
	editor.Rows = 0
	editor.SaveFlag = false
	editor.TopRowNum = 1
	return
}

// 行から改行文字を判定 (-1は不明の改行文字)
func GetNL(row []rune) nlCode {
	if row[len(row)-1] == rune('\n') {
		if row[len(row)-2] == rune('\r') {
			return NL_CRLF
		}
		return NL_LF
	}
	return -1
}

// 行ノードのポインタを1つ進める
func (e *Editor) MoveNextRow() *common.RowNode {
	e.CurrentNode = e.CurrentNode.Next
	return e.CurrentNode
}

// 行ノードのポインタを1つ戻す
func (e *Editor) MovePrevRow() *common.RowNode {
	e.CurrentNode = e.CurrentNode.Prev
	return e.CurrentNode
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

	rowsCnt := 0
	reader := bufio.NewReaderSize(fp, LINE_BUF_MAX)
	for {
		line, err := reader.ReadString(byte('\n'))                    // '\n'で分割
		replacedStr := strings.ReplaceAll(string(line), "\t", tabStr) // タブをスペースに変換
		replacedRune := []rune(replacedStr)
		if rowsCnt == 0 { // 改行文字の判定
			e.NL = GetNL(replacedRune)
		}
		// TODO: 改行文字が混在している時の対応
		switch e.NL { // 改行文字の削除
		case NL_CRLF:
			if len(replacedRune) >= 2 {
				replacedRune = replacedRune[:len(replacedRune)-2]
			}
		case NL_LF:
			if len(replacedRune) >= 1 {
				replacedRune = replacedRune[:len(replacedRune)-1]
			}
		}

		e.Root.Prev.Insert(replacedRune, LINE_BUF_MAX)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		rowsCnt++
	}
	e.MoveNextRow()
	e.Rows = uint16(rowsCnt)
}

// エディタに指定されたパスで上書き保存
func (e *Editor) SaveOverwrite(nl nlCode) (saveBytes int) {
	saveBytes = e.saveFile(e.FilePath, nl)
	return
}

// 新しくファイルを保存
func (e *Editor) SaveNew(filePath string, nl nlCode) (saveBytes int) {
	saveBytes = e.saveFile(filePath, nl)
	return
}

// ファイルを保存
func (e *Editor) saveFile(filePath string, nl nlCode) (saveBytes int) {
	fp, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var buf string

	cnt := 0
	tmp := e.Root.Next
	for {
		if tmp == e.Root {
			break
		}
		cnt++
		buf += fmt.Sprintf("%s", string(tmp.Buf.GetAll()))
		switch nl {
		case NL_CRLF:
			buf += "\r\n"
		case NL_LF:
			buf += "\n"
		default:
			buf += "\n"
		}
		tmp = tmp.Next
	}

	saveBytes, err = fp.Write([]byte(buf))
	if err != nil {
		panic(err)
	}
	return
}

func (e *Editor) GetNodeFromLineNum(num uint16) *common.RowNode {
	pNode := e.Root
	pNode = pNode.Next
	var cntRow uint16 = 1
	for {
		if pNode.IsRoot() {
			break
		}
		if cntRow == num {
			return pNode
		}
		pNode = pNode.Next
		cntRow++
	}
	return nil
}

func (e *Editor) GetLineNumFromNode(tNode *common.RowNode) uint16 {
	pNode := e.Root
	pNode = pNode.Next
	var cntRow uint16 = 1
	for {
		if pNode.IsRoot() {
			break
		} else if pNode == tNode {
			return cntRow
		}
		pNode = pNode.Next
		cntRow++
	}
	return 0
}
