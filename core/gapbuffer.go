package core

// https://tech.pjin.jp/blog/2020/11/18/buffer-4

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
func (gBuf GapBuffer) Length() int {
	return gBuf.size - gBuf.gapSize
}

// バッファが空かどうかの判定
func (gBuf GapBuffer) IsEmpty() bool {
	if gBuf.Length() == 0 {
		return true
	}
	return false
}

// バッファのruneが一致するかの判定
func (gBuf GapBuffer) Check(idx int, data []rune) bool {
	for i, elm := range data {
		if gBuf.Get(idx+i) != elm {
			return false
		}
	}
	return true
}

// バッファのruneを取得
func (gBuf GapBuffer) Get(idx int) rune {
	internalIdx := idx
	if idx >= gBuf.gapIdx {
		internalIdx += gBuf.gapSize
	}
	if idx >= gBuf.Length() {
		return 0
	}
	return gBuf.buf[internalIdx]
}

func (gBuf GapBuffer) GetFrom(startIdx int, endIdx int) (out []rune) {
	if endIdx - 1 < gBuf.Length() {
		for i := startIdx; i < endIdx; i++ {
			out = append(out, gBuf.Get(i))
		}
	}
	return
}

// バッファのruneをすべて取得
func (gBuf GapBuffer) GetAll() (out []rune) {
	out = gBuf.GetFrom(0, gBuf.Length())
	return
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
	for i, elm := range data {
		gBuf.Insert(idx+i, elm)
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

func (gBuf *GapBuffer) EraseFrom(startIdx int, endIdx int) {
	if endIdx - 1 < gBuf.Length() {
		for i := startIdx; i < endIdx; i++ {
			gBuf.Erase(startIdx)
		}
	}
}

// バッファから複数のruneを削除
func (gBuf *GapBuffer) EraseAll(idx int, num int) {
	for i := 0; i < num; i++ {
		gBuf.Erase(idx)
	}
}

// バッファの末尾にruneを追加
func (gBuf *GapBuffer) Append(ch rune) bool {
	return gBuf.Insert(gBuf.Length(), ch)
}

// バッファの末尾に複数のruneを追加
func (gBuf *GapBuffer) AppendAll(data []rune) {
	gBuf.InsertAll(gBuf.Length(), data)
}
