package common

type DLNode _DLNode

// 行ノードの定義 (双方向連結リスト)
type _DLNode struct {
	fRoot bool // ルートノードフラグ
	prev  *DLNode
	next  *DLNode
	buf   *GapBuffer
}

// ルートノードの生成
func NewRootNode() *DLNode {
	root := new(DLNode)
	root.fRoot = true
	root.prev = root
	root.next = root
	root.buf = NewGapBuffer(make([]rune, 0), 0)
	return root
}

// 子ノードの生成
func NewChildNode(data []rune, bufSize int) *DLNode {
	child := new(DLNode)
	child.fRoot = false
	child.prev = child
	child.next = child
	child.buf = NewGapBuffer(data, bufSize)
	return child
}

// ノードがルートノードかどうか判定
func (e *DLNode) IsRoot() bool {
	if e.fRoot {
		return true
	}
	return false
}

// ノードを挿入
func (e *DLNode) Insert(data []rune, bufSize int) *DLNode {
	_new := NewChildNode(data, bufSize)
	_new.next = e.next
	_new.prev = e
	e.next.prev = _new
	e.next = _new
	return _new
}

// ノードを削除
func (e *DLNode) Delete() *DLNode {
	e.prev.next = e.next
	e.next.prev = e.prev
	return e.prev
}

func (e *DLNode) Next() *DLNode {
	return e.next
}

func (e *DLNode) Prev() *DLNode {
	return e.prev
}

func (e *DLNode) GetBuf() *GapBuffer {
	return e.buf
}
