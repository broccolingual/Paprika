package common

type RowNode _RowNode

// 行ノードの定義 (双方向連結リスト)
type _RowNode struct {
	fRoot bool // ルートノードフラグ
	Prev  *RowNode
	Next  *RowNode
	Buf   *GapBuffer
}

// ルートノードの生成
func NewRootNode() *RowNode {
	root := new(RowNode)
	root.fRoot = true
	root.Prev = root
	root.Next = root
	root.Buf = NewGapBuffer(make([]rune, 0), 0)
	return root
}

// 子ノードの生成
func NewChildNode(data []rune, bufSize int) *RowNode {
	child := new(RowNode)
	child.fRoot = false
	child.Prev = child
	child.Next = child
	child.Buf = NewGapBuffer(data, bufSize)
	return child
}

// ノードがルートノードかどうか判定
func (e *RowNode) IsRoot() bool {
	if e.fRoot {
		return true
	}
	return false
}

// ノードを挿入
func (e *RowNode) Insert(data []rune, bufSize int) *RowNode {
	_new := NewChildNode(data, bufSize)
	_new.Next = e.Next
	_new.Prev = e
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
