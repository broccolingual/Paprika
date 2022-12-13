package main

import "fmt"

func main() {
	f := LoadFile("Makefile")
	dummy := NewRowsList()

	// append
	for _, line := range f.Lines {
		dummy.Prev.Append([]rune(line), 256)
	}

	// display
	cnt := 0
	tmp := dummy.Next
	for {
		if tmp == dummy {
			break
		}
		cnt++
		fmt.Printf("%4d [ %s\n", cnt, string(tmp.Row.GetAll()))
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
	_new.Next = e.Next
	e.Next = _new
	_new.Prev = e
	_new.Row = NewGapBuffer(data, bufSize)
}
