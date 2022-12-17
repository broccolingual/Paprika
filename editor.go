package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Editor struct {
	FilePath    string
	Cursor      *Cursor
	Root        *RowNode
	CurrentNode *RowNode
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

func NewEditor(filePath string) (editor *Editor) {
	editor = new(Editor)
	editor.FilePath = filePath
	editor.Cursor = NewCursor()
	editor.Root = NewRowsList()
	editor.CurrentNode = editor.Root
	return
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

	reader := bufio.NewReaderSize(fp, 512)
	for {
		line, _, err := reader.ReadLine()
		e.Root.Prev.Append([]rune(string(line)), 512)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
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
	var fp *os.File
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fp, err = os.Create(filePath)
	} else {
		fp, err = os.Open(filePath)
	}
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
		if nl == NL_LF {
			buf += "\n"
		} else if nl == NL_CRLF {
			buf += "\r\n"
		} else {
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
