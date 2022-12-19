package main

const (
	NL_LF = iota
	NL_CRLF
)

// Define rawnode object
type RowNode struct {
	Prev *RowNode
	Next *RowNode
	Row  GapBuffer
}

func NewRowsList() *RowNode {
	dummy := new(RowNode)
	dummy.Prev = dummy
	dummy.Next = dummy
	dummy.Row = nil
	return dummy
}

func (e *RowNode) IsRoot() bool {
	if e.Row == nil {
		return true
	}
	return false
}

func (e *RowNode) Append(data []rune, bufSize int) {
	_new := new(RowNode)
	_new.Next = e.Next
	e.Next = _new
	_new.Prev = e
	_new.Next.Prev = _new
	_new.Row = NewGapBuffer(data, bufSize)
}

func (e *RowNode) Insert(data []rune, bufSize int) *RowNode {
	_new := new(RowNode)
	_new.Next = e.Next
	_new.Prev = e
	_new.Row = NewGapBuffer(data, bufSize)
	e.Next.Prev = _new
	e.Next = _new
	return _new
}

func (e *RowNode) Delete() *RowNode {
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev
	return e.Prev
}
