package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const LINE_BUF_MAX = 255

type Editor struct {
	FilePath    string
	Cursor      *Cursor
	Root        *RowNode
	CurrentNode *RowNode
	TabSize     uint8
	NL          int
	Rows        uint16
}

type Cursor struct {
	Row uint16
	Col uint16
}

func NewCursor() (cursor *Cursor) {
	cursor = new(Cursor)
	cursor.Row = 1
	cursor.Col = 1
	return
}

func NewEditor(filePath string, tabSize uint8) (editor *Editor) {
	editor = new(Editor)
	editor.FilePath = filePath
	editor.Cursor = NewCursor()
	editor.Root = NewRowsList()
	editor.CurrentNode = editor.Root
	editor.TabSize = tabSize
	editor.NL = -1
	editor.Rows = 0
	return
}

func GetNL(row []rune) int {
	if row[len(row)-1] == rune('\n') {
		if row[len(row)-2] == rune('\r') {
			return NL_CRLF
		}
		return NL_LF
	}
	return -1
}

func (e *Editor) MoveNextRow() *RowNode {
	e.CurrentNode = e.CurrentNode.Next
	return e.CurrentNode
}

func (e *Editor) MovePrevRow() *RowNode {
	e.CurrentNode = e.CurrentNode.Prev
	return e.CurrentNode
}

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

		e.Root.Prev.Append(replacedRune, LINE_BUF_MAX)
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

func (e *Editor) SaveOverwrite(nl int) (saveBytes int) {
	saveBytes = e.saveFile(e.FilePath, nl)
	return
}

func (e *Editor) SaveNew(filePath string, nl int) (saveBytes int) {
	saveBytes = e.saveFile(filePath, nl)
	return
}

func (e *Editor) saveFile(filePath string, nl int) (saveBytes int) {
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
		buf += fmt.Sprintf("%s", string(tmp.Row.GetAll()))
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
