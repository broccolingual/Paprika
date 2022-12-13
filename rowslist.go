package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	dummy := LoadFile("Makefile")
	// insert
	cnt := 0
	tmp := dummy.Next
	for {
		if tmp == dummy {
			break
		}
		if cnt == 3 {
			tmp.Insert([]rune("THIS IS APPEND LINE !"), 512)
		}
		cnt++
		tmp = tmp.Next
	}

	// display
	cnt = 0
	tmp = dummy.Next
	for {
		if tmp == dummy {
			break
		}
		cnt++
		fmt.Printf("%4d | %s\n", cnt, string(tmp.Row.GetAll()))
		tmp = tmp.Next
	}
}

type RowNode struct {
	Prev *RowNode
	Next *RowNode
	Row  *GapBuffer
}

func NewRowsList() *RowNode {
	dummy := new(RowNode)
	dummy.Prev = dummy
	dummy.Next = dummy
	dummy.Row = nil
	return dummy
}

func (e *RowNode) Append(data []rune, bufSize int) {
	_new := new(RowNode)
	_new.Next = e.Next
	e.Next = _new
	_new.Prev = e
	_new.Next.Prev = _new
	_new.Row = NewGapBuffer(data, bufSize)
}

func (e *RowNode) Insert(data []rune, bufSize int) {
	_new := new(RowNode)
	_new.Prev = e.Prev
	_new.Next = e
	e.Prev = _new
	_new.Row = NewGapBuffer(data, bufSize)
}

func LoadFile(fileName string) *RowNode {
	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	dummy := NewRowsList()

	reader := bufio.NewReaderSize(fp, 512)
	for {
		line, _, err := reader.ReadLine()
		dummy.Prev.Append([]rune(string(line)), 512)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return dummy
}
