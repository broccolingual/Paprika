package main

import "fmt"

func main() {
	gBuf := NewGapBuffer([]rune("ありがとうございました"))
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

	gBuf.moveGap(8)
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

	gBuf.moveGap(2)
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

	gBuf.moveGap(8)
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

	gBuf.Insert(0, rune('大'))
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

	gBuf.Insert(1, rune('変'))
	fmt.Println(gBuf)
	fmt.Println(string(gBuf.Data), gBuf.GetSize())

}

type GapBuffer struct {
	Size    uint
	GapIdx  uint
	GapSize uint
	Data    []rune
}

func NewGapBuffer(data []rune) *GapBuffer {
	gBuf := new(GapBuffer)
	gBuf.Size = 320
	gBuf.GapIdx = 256
	gBuf.GapSize = 64
	gBuf.Data = data
	gBuf.initGap()
	return gBuf
}

func (gBuf *GapBuffer) GetSize() uint {
	return uint(len(string(gBuf.Data)))
}

func (gBuf *GapBuffer) initGap() {
	gBuf.Data = append(gBuf.Data, make([]rune, gBuf.GapIdx-uint(len(gBuf.Data))+gBuf.GapSize)...)
}

func (gBuf *GapBuffer) moveGap(idx uint) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	oldGapIdx := gBuf.GapIdx
	gBuf.GapIdx = idx
	if oldGapIdx < idx {
		buf := gBuf.Data[(oldGapIdx + gBuf.GapSize):(idx + gBuf.GapSize)]
		for i, el := range buf {
			gBuf.Data[oldGapIdx+uint(i)] = el
		}
		for i := oldGapIdx + gBuf.GapSize; i < idx+gBuf.GapSize; i++ {
			gBuf.Data[uint(i)] = 0
		}
	} else if oldGapIdx > idx {
		buf := gBuf.Data[idx:oldGapIdx]
		for i, el := range buf {
			gBuf.Data[idx+gBuf.GapSize+uint(i)] = el
		}
		for i := idx; i < oldGapIdx; i++ {
			gBuf.Data[uint(i)] = 0
		}
	} else {
		return
	}
}

// 未使用
func (gBuf *GapBuffer) Get(idx uint) rune {
	if idx >= gBuf.GapIdx {
		idx += gBuf.GapSize
	}
	return gBuf.Data[idx]
}

func (gBuf *GapBuffer) Insert(idx uint, ch rune) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.Data[gBuf.GapIdx] = ch
	gBuf.GapIdx++
	gBuf.GapSize--
}

func (gBuf *GapBuffer) Erase(idx uint) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.GapSize++
}
