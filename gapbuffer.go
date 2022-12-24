package main

type GapBuffer interface {
	GetSize() int
	Get(idx int) rune
	GetAll() []rune
	Insert(idx int, ch rune) bool
	Erase(idx int) bool
}

// ギャップバッファ構造体
type gapBuffer struct {
	size    int    // バッファサイズ
	gapIdx  int    // ギャップの開始インデックス
	gapSize int    // ギャップサイズ
	buf     []rune // バッファ
}

// 新しいギャップバッファの取得
func NewGapBuffer(data []rune, bufSize int) GapBuffer {
	gBuf := new(gapBuffer)
	gBuf.size = bufSize
	gBuf.gapIdx = len(data)
	gBuf.gapSize = bufSize - gBuf.gapIdx
	gBuf.buf = data
	gBuf.initGap()
	return gBuf
}

// ギャップバッファの初期化 (バッファのギャップ部分を0埋め)
func (gBuf *gapBuffer) initGap() {
	gBuf.buf = append(gBuf.buf, make([]rune, gBuf.gapIdx-len(gBuf.buf)+gBuf.gapSize)...)
}

// 指定したインデックスにギャップを移動
func (gBuf *gapBuffer) moveGap(idx int) {
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

// バッファ内のruneのサイズを取得
func (gBuf *gapBuffer) GetSize() int {
	return gBuf.size - gBuf.gapSize
}

// バッファのruneを取得
func (gBuf *gapBuffer) Get(idx int) rune {
	if idx >= gBuf.gapIdx {
		idx += gBuf.gapSize
	}
	return gBuf.buf[idx]
}

// バッファのruneをすべて取得
func (gBuf *gapBuffer) GetAll() []rune {
	var tmp []rune
	for i := 0; i < int(gBuf.size-gBuf.gapSize); i++ {
		tmp = append(tmp, gBuf.Get(i))
	}
	return tmp
}

// バッファにruneを挿入
func (gBuf *gapBuffer) Insert(idx int, ch rune) bool {
	if idx < 0 || idx > gBuf.size {
		return false
	}
	gBuf.moveGap(idx)
	gBuf.buf[gBuf.gapIdx] = ch
	gBuf.gapIdx++
	gBuf.gapSize--
	return true
}

// バッファのruneを削除
func (gBuf *gapBuffer) Erase(idx int) bool {
	if idx < 0 || idx > gBuf.size {
		return false
	}
	gBuf.moveGap(idx)
	gBuf.gapSize++
	return true
}
