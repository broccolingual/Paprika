package core

type DoublyLinkedList _DoublyLinkedList

type _DoublyLinkedList struct {
	prev *DoublyLinkedList
	next *DoublyLinkedList
	buf *GapBuffer
}

func NewDoublyLinkedList(gBuf *GapBuffer) *DoublyLinkedList {
	dll := new(DoublyLinkedList)
	dll.prev = dll
	dll.next = dll
	dll.buf = gBuf
	return dll
}

// 1行後に追加
func (dll *DoublyLinkedList) Append(gBuf *GapBuffer) *DoublyLinkedList {
	newDll := NewDoublyLinkedList(gBuf)
	newDll.next = dll.next
	dll.next = newDll
	newDll.prev = dll
	return dll
}

// 1行削除
func (dll *DoublyLinkedList) Remove() *DoublyLinkedList {
	dll.prev.next = dll.next
	dll.next.prev = dll.prev
	return dll.next
}

// 要素数のカウント
func (dll *DoublyLinkedList) Length() int {
	curDll := dll.next
	cnt := 1
	for curDll != dll {
		curDll = curDll.next
		cnt++
	}
	return cnt
}