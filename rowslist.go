package main

const (
	NL_LF = iota
	NL_CRLF
)

// 行ノードの定義 (双方向連結リスト)
type RowNode struct {
	Prev *RowNode
	Next *RowNode
	Row  GapBuffer
}

// ルートノードの生成
func NewRowsList() *RowNode {
	dummy := new(RowNode)
	dummy.Prev = dummy
	dummy.Next = dummy
	dummy.Row = nil
	return dummy
}

// ノードがルートノードかどうか判定
func (e *RowNode) IsRoot() bool {
	if e.Row == nil {
		return true
	}
	return false
}

// ノードを追加
func (e *RowNode) Append(data []rune, bufSize int) {
	_new := new(RowNode)
	_new.Next = e.Next
	e.Next = _new
	_new.Prev = e
	_new.Next.Prev = _new
	_new.Row = NewGapBuffer(data, bufSize)
}

// ノードを挿入
func (e *RowNode) Insert(data []rune, bufSize int) *RowNode {
	_new := new(RowNode)
	_new.Next = e.Next
	_new.Prev = e
	_new.Row = NewGapBuffer(data, bufSize)
	e.Next.Prev = _new
	e.Next = _new
	return _new
}

// ノードを削除
func (e *RowNode) Delete() *RowNode {
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev
	return e.Prev
}
