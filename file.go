package main

import (
	"bufio"
	"io"
	"os"
)

type File struct {
	Name  string
	Lines []string
}

func NewFile(name string) *File {
	f := new(File)
	f.Name = name
	return f
}

func LoadFile(name string) *File {
	fp, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	f := NewFile(name)

	reader := bufio.NewReaderSize(fp, 512)
	for {
		line, _, err := reader.ReadLine()
		f.Lines = append(f.Lines, string(line))
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return f
}

func (f *File) InsertNewLine(col uint, str string) {
	f.Lines = append(f.Lines[:col+2], f.Lines[col+1:]...)
	f.Lines[col+1] = str
}
