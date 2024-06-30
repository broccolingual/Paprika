package core

import "testing"

func Test_GB_Length(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {}
	tests := []struct {
		name string
		fields fields
		args args
		want int
	}{
		{"Test #1", fields{[]rune("A"), 64}, args{}, 1},
		{"Test #2", fields{[]rune("Hello World !"), 64}, args{}, 13},
		{"Test #3", fields{[]rune("あいう"), 64}, args{}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			if got := gBuf.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_Insert(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		idx int
		ch rune
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("A"), 64}, args{1, rune('B')}, "AB"},
		{"Test #2", fields{[]rune("Hello"), 64}, args{0, rune('!')}, "!Hello"},
		{"Test #3", fields{[]rune("あう"), 64}, args{1, rune('い')}, "あいう"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.Insert(tt.args.idx, tt.args.ch)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_InsertAll(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		idx int
		data []rune
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("A"), 64}, args{1, []rune("BC")}, "ABC"},
		{"Test #2", fields{[]rune("Hello"), 64}, args{5, []rune(" World !")}, "Hello World !"},
		{"Test #3", fields{[]rune("あお"), 64}, args{1, []rune("いうえ")}, "あいうえお"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.InsertAll(tt.args.idx, tt.args.data)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("InsertAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_Erase(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		idx int
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("ABC"), 64}, args{2}, "AB"},
		{"Test #2", fields{[]rune("!Hello"), 64}, args{0}, "Hello"},
		{"Test #3", fields{[]rune("あいえう"), 64}, args{2}, "あいう"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.Erase(tt.args.idx)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("Erase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_EraseAll(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		idx int
		num int
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("ABC"), 64}, args{1, 2}, "A"},
		{"Test #2", fields{[]rune("Hello"), 64}, args{4, 1}, "Hell"},
		{"Test #3", fields{[]rune("あいいいう"), 64}, args{2, 2}, "あいう"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.EraseAll(tt.args.idx, tt.args.num)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("InsertAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_Append(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		ch rune
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("AB"), 64}, args{rune('C')}, "ABC"},
		{"Test #2", fields{[]rune("Hello"), 64}, args{rune('!')}, "Hello!"},
		{"Test #3", fields{[]rune("あい"), 64}, args{rune('う')}, "あいう"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.Append(tt.args.ch)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GB_AppendAll(t *testing.T) {
	type fields struct {
		data []rune
		bufSize int
	}
	type args struct {
		data []rune
	}
	tests := []struct {
		name string
		fields fields
		args args
		want string
	}{
		{"Test #1", fields{[]rune("A"), 64}, args{[]rune("BC")}, "ABC"},
		{"Test #2", fields{[]rune("Hello"), 64}, args{[]rune(" World !")}, "Hello World !"},
		{"Test #3", fields{[]rune("あいう"), 64}, args{[]rune("えお")}, "あいうえお"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gBuf := NewGapBuffer(tt.fields.data, tt.fields.bufSize)
			gBuf.AppendAll(tt.args.data)
			if got := string(gBuf.GetAll()); got != tt.want {
				t.Errorf("AppendAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
