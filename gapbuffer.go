package main

type GapBuffer struct {
	Size    int
	GapIdx  int
	GapSize int
	Data    []rune
}

func NewGapBuffer(data []rune, bufSize int) *GapBuffer {
	gBuf := new(GapBuffer)
	gBuf.Size = bufSize
	gBuf.GapIdx = len(data)
	gBuf.GapSize = bufSize - gBuf.GapIdx
	gBuf.Data = data
	gBuf.initGap()
	return gBuf
}

func (gBuf *GapBuffer) GetSize() int {
	return gBuf.Size - gBuf.GapSize
}

func (gBuf *GapBuffer) initGap() {
	gBuf.Data = append(gBuf.Data, make([]rune, gBuf.GapIdx-len(gBuf.Data)+gBuf.GapSize)...)
}

func (gBuf *GapBuffer) moveGap(idx int) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	oldGapIdx := gBuf.GapIdx
	gBuf.GapIdx = idx
	if oldGapIdx < idx {
		buf := make([]rune, idx-oldGapIdx)
		_ = copy(buf, gBuf.Data[(oldGapIdx+gBuf.GapSize):(idx+gBuf.GapSize)])
		for i := 0; i < len(buf); i++ {
			gBuf.Data[oldGapIdx+i] = buf[i]
		}
	} else if oldGapIdx > idx {
		buf := make([]rune, oldGapIdx-idx)
		_ = copy(buf, gBuf.Data[idx:oldGapIdx])
		for i := 0; i < len(buf); i++ {
			gBuf.Data[idx+gBuf.GapSize+i] = buf[i]
		}
	} else {
		return
	}
}

func (gBuf *GapBuffer) Get(idx int) rune {
	if idx >= gBuf.GapIdx {
		idx += gBuf.GapSize
	}
	return gBuf.Data[idx]
}

func (gBuf *GapBuffer) GetAll() []rune {
	var tmp []rune
	for i := 0; i < int(gBuf.Size-gBuf.GapSize); i++ {
		tmp = append(tmp, gBuf.Get(i))
	}
	return tmp
}

func (gBuf *GapBuffer) Insert(idx int, ch rune) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.Data[gBuf.GapIdx] = ch
	gBuf.GapIdx++
	gBuf.GapSize--
}

func (gBuf *GapBuffer) Erase(idx int) {
	if idx < 0 || idx > gBuf.Size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.GapSize++
}
