package main

// Define gapbuffer object
type GapBuffer struct {
	size    int
	gapIdx  int
	gapSize int
	buf     []rune
}

func NewGapBuffer(data []rune, bufSize int) *GapBuffer {
	gBuf := new(GapBuffer)
	gBuf.size = bufSize
	gBuf.gapIdx = len(data)
	gBuf.gapSize = bufSize - gBuf.gapIdx
	gBuf.buf = data
	gBuf.initGap()
	return gBuf
}

func (gBuf *GapBuffer) GetSize() int {
	return gBuf.size - gBuf.gapSize
}

func (gBuf *GapBuffer) initGap() {
	gBuf.buf = append(gBuf.buf, make([]rune, gBuf.gapIdx-len(gBuf.buf)+gBuf.gapSize)...)
}

func (gBuf *GapBuffer) moveGap(idx int) {
	if idx < 0 || idx > gBuf.size {
		return
	}
	oldGapIdx := gBuf.gapIdx
	gBuf.gapIdx = idx
	if oldGapIdx < idx {
		buf := make([]rune, idx-oldGapIdx)
		_ = copy(buf, gBuf.buf[(oldGapIdx+gBuf.gapSize):(idx+gBuf.gapSize)])
		for i := 0; i < len(buf); i++ {
			gBuf.buf[oldGapIdx+i] = buf[i]
		}
	} else if oldGapIdx > idx {
		buf := make([]rune, oldGapIdx-idx)
		_ = copy(buf, gBuf.buf[idx:oldGapIdx])
		for i := 0; i < len(buf); i++ {
			gBuf.buf[idx+gBuf.gapSize+i] = buf[i]
		}
	}
}

func (gBuf *GapBuffer) Get(idx int) rune {
	if idx >= gBuf.gapIdx {
		idx += gBuf.gapSize
	}
	return gBuf.buf[idx]
}

func (gBuf *GapBuffer) GetAll() []rune {
	var tmp []rune
	for i := 0; i < int(gBuf.size-gBuf.gapSize); i++ {
		tmp = append(tmp, gBuf.Get(i))
	}
	return tmp
}

func (gBuf *GapBuffer) Insert(idx int, ch rune) {
	if idx < 0 || idx > gBuf.size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.buf[gBuf.gapIdx] = ch
	gBuf.gapIdx++
	gBuf.gapSize--
}

func (gBuf *GapBuffer) Erase(idx int) {
	if idx < 0 || idx > gBuf.size {
		return
	}
	gBuf.moveGap(idx)
	gBuf.gapSize++
}
