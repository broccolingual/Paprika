package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	NL_LF = iota
	NL_CRLF
)

func main() {
	// Alternative Screen Buffer Settings
	EnableASB()
	// defer DisableASB()

	// loading file
	dummy := LoadFile("Makefile")

	// display
	cnt := 0
	tmp := dummy.Next
	for {
		if tmp == dummy {
			break
		}
		cnt++
		fmt.Printf("%4d | %s\n", cnt, string(tmp.Row.GetAll()))
		tmp = tmp.Next
	}

	DisableASB()

	bytesCount := dummy.SaveFile("./bin/Makefile.bak", NL_LF)
	fmt.Printf("Write %d bytes.\n", bytesCount)

	// w := NewWindow("__MAIN__")
	// go w.readKeys()

	// w.switchKeys()
}

// Define rawnode object
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

// Make rawslist from loading file
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

// Save file from rowslist
func (dummy *RowNode) SaveFile(fileName string, nl int) int {
	var fp *os.File
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		fp, err = os.Create(fileName)
	} else {
		fp, err = os.Open(fileName)
	}
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var buf string

	cnt := 0
	tmp := dummy.Next
	for {
		if tmp == dummy {
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

	bytesCount, err := fp.Write([]byte(buf))
	if err != nil {
		panic(err)
	}
	return bytesCount
}
