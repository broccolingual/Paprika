package core

type DoublyLinkedList _DoublyLinkedList

type _DoublyLinkedList struct {
	root *DoublyLinkedList
	prev *DoublyLinkedList
	next *DoublyLinkedList
	buf *GapBuffer
}

func NewDllRoot(gBuf *GapBuffer) *DoublyLinkedList {
	dll := new(DoublyLinkedList)
	dll.root = dll
	dll.prev = dll
	dll.next = dll
	dll.buf = gBuf
	return dll
}

func NewDllNode(gBuf *GapBuffer, rootNode *DoublyLinkedList) *DoublyLinkedList {
	dll := new(DoublyLinkedList)
	dll.root = rootNode
	dll.prev = dll
	dll.next = dll
	dll.buf = gBuf
	return dll
}

// ルートノードかどうかの判定
func (dll *DoublyLinkedList) IsRoot() bool {
	return dll.root == dll
}

// 1行後に追加
func (dll *DoublyLinkedList) Insert(gBuf *GapBuffer) *DoublyLinkedList {
	newNode := NewDllNode(gBuf, dll.root)
	dll.next.prev = newNode
	newNode.next = dll.next
	newNode.prev = dll
	dll.next = newNode
	return newNode
}

// 1行削除
func (dll *DoublyLinkedList) Remove() *DoublyLinkedList {
	if dll.root == dll {
		return dll
	}
	dll.prev.next = dll.next
	dll.next.prev = dll.prev
	return dll.next
}

// 要素数のカウント
func (dll *DoublyLinkedList) Length() int {
	currentNode := dll.root.next
	cnt := 1
	for currentNode.IsRoot() == false {
		currentNode = currentNode.next
		cnt++
	}
	return cnt
}

func (dll *DoublyLinkedList) GetIdx() int {
	if dll.IsRoot() {
		return 0
	}
	currentNode := dll.root.next
	idx := 1
	for currentNode != dll {
		currentNode = currentNode.next
		idx++
	}
	return idx
}

func (dll *DoublyLinkedList) GetBuf() *GapBuffer {
	return dll.buf
}

func (dll *DoublyLinkedList) Root() *DoublyLinkedList {
	return dll.root
}

func (dll *DoublyLinkedList) Next() *DoublyLinkedList {
	return dll.next
}

func (dll *DoublyLinkedList) Prev() *DoublyLinkedList {
	return dll.prev
}
