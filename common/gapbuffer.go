package common

type GapBuffer _GapBuffer

// ギャップバッファ構造体
type _GapBuffer struct {
	size    int    // バッファサイズ
	gapIdx  int    // ギャップの開始インデックス
	gapSize int    // ギャップサイズ
	buf     []rune // バッファ
}

// 新しいギャップバッファの取得
func NewGapBuffer(data []rune, bufSize int) *GapBuffer {
	gBuf := new(GapBuffer)
	gBuf.size = bufSize
	gBuf.gapIdx = len(data)
	gBuf.gapSize = bufSize - gBuf.gapIdx
	gBuf.buf = data
	gBuf.initGap()
	return gBuf
}

// ギャップバッファの初期化 (バッファのギャップ部分を0埋め)
func (gBuf *GapBuffer) initGap() {
	gBuf.buf = append(gBuf.buf, make([]rune, gBuf.gapIdx-len(gBuf.buf)+gBuf.gapSize)...)
}

// 指定したインデックスにギャップを移動
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

// バッファ内のruneのサイズを取得
func (gBuf *GapBuffer) GetSize() int {
	return gBuf.size - gBuf.gapSize
}

// バッファのruneが一致するかの判定
func (gBuf *GapBuffer) Check(idx int, data []rune) bool {
	for i := 0; i < len(data); i++ {
		if gBuf.Get(idx+i) != data[i] {
			return false
		}
	}
	return true
}

// バッファのruneを取得
func (gBuf *GapBuffer) Get(idx int) rune {
	if idx >= gBuf.gapIdx {
		idx += gBuf.gapSize
	}
	if idx >= gBuf.GetSize() {
		// TODO: Overflowの対応
	}
	return gBuf.buf[idx]
}

// バッファのruneをすべて取得
func (gBuf *GapBuffer) GetAll() []rune {
	var tmp []rune
	for i := 0; i < int(gBuf.size-gBuf.gapSize); i++ {
		tmp = append(tmp, gBuf.Get(i))
	}
	return tmp
}

// バッファにruneを挿入
func (gBuf *GapBuffer) Insert(idx int, ch rune) bool {
	if idx < 0 || idx > gBuf.size {
		return false
	}
	gBuf.moveGap(idx)
	gBuf.buf[gBuf.gapIdx] = ch
	gBuf.gapIdx++
	gBuf.gapSize--
	return true
}

// バッファに複数のruneを挿入
func (gBuf *GapBuffer) InsertAll(idx int, data []rune) {
	for i := 0; i < len(data); i++ {
		gBuf.Insert(i+idx, data[i])
	}
}

// バッファのruneを削除
func (gBuf *GapBuffer) Erase(idx int) bool {
	if idx < 0 || idx > gBuf.size {
		return false
	}
	gBuf.moveGap(idx)
	gBuf.gapSize++
	return true
}

// バッファから複数のruneを削除
func (gBuf *GapBuffer) EraseAll(idx int, num int) {
	for i := 0; i < num; i++ {
		gBuf.Erase(idx)
	}
}

// バッファが空かどうかの判定
func (gBuf *GapBuffer) IsBlank() bool {
	if gBuf.GetSize() == 0 {
		return true
	}
	return false
}