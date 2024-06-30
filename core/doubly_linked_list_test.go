package core

import "testing"

func Test_DLL_IsRoot(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"Test #1", true},
		{"Test #2", false},
		{"Test #3", true},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewDllRoot(nil)
			if i == 0 {
			} else if i == 1 {
				node = node.Insert(nil)
			} else if i == 2 {
				node = node.Insert(nil)
				node = node.next
			}
			if got := node.IsRoot(); got != tt.want {
				t.Errorf("IsRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DLL_Insert(t *testing.T) {
	type fields struct {
		root *DoublyLinkedList
		prev *DoublyLinkedList
		next *DoublyLinkedList
		buf *GapBuffer
	}
	type args struct {
		gBuf *GapBuffer
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{nil, nil, nil, NewGapBuffer([]rune("A"), 64)}, args{NewGapBuffer([]rune("B"), 64)}, "B"},
		{"Test #2", fields{nil, nil, nil, NewGapBuffer([]rune("Hello"), 64)}, args{NewGapBuffer([]rune("World"), 64)}, "World"},
		{"Test #3", fields{nil, nil, nil, NewGapBuffer([]rune("あいう"), 64)}, args{NewGapBuffer([]rune("えお"), 64)}, "えお"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewDllRoot(tt.fields.buf)
			newDll := root.Insert(tt.args.gBuf)
			if got := string(newDll.buf.GetAll()); got != tt.want {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DLL_Remove(t *testing.T) {
	type fields struct {
		root *DoublyLinkedList
		prev *DoublyLinkedList
		next *DoublyLinkedList
		buf *GapBuffer
	}
	tests := []struct {
		name string
		fields fields
		want string
	}{
		{"Test #1", fields{nil, nil, nil, NewGapBuffer([]rune("A"), 64)}, "A"},
		{"Test #2", fields{nil, nil, nil, NewGapBuffer([]rune("Hello"), 64)}, "Hello"},
		{"Test #3", fields{nil, nil, nil, NewGapBuffer([]rune("あいう"), 64)}, "あいう"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewDllRoot(tt.fields.buf)
			newDll := root.Insert(NewGapBuffer([]rune("B"), 64))
			newDll = newDll.Remove()
			if got := string(newDll.buf.GetAll()); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DLL_Length(t *testing.T) {
	tests := []struct {
		name string
		numOfInsert int
		want int
	}{
		{"Test #1", 0, 1},
		{"Test #2", 1, 2},
		{"Test #3", 50, 51},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewDllRoot(nil)
			for i := 0; i < tt.numOfInsert; i++ {
				root = root.Insert(nil)
			}
			if got := root.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DLL_GetIdx(t *testing.T) {
	tests := []struct {
		name string
		numOfInsert int
		want int
	}{
		{"Test #1", 0, 0},
		{"Test #2", 1, 1},
		{"Test #3", 50, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := NewDllRoot(nil)
			for i := 0; i < tt.numOfInsert; i++ {
				root = root.Insert(nil)
			}
			if got := root.GetIdx(); got != tt.want {
				t.Errorf("GetIdx() = %v, want %v", got, tt.want)
			}
		})
	}
}
