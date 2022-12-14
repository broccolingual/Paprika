package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type Window struct {
	Name    string
	Cursor  *Cursor
	KeyChan chan rune
	Rows    *RowNode
}

type Cursor struct {
	Row uint16
	Col uint16
}

func NewWindow(name string) *Window {
	w := new(Window)
	w.Name = name
	w.Cursor = NewCursor()
	w.KeyChan = make(chan rune)
	return w
}

func NewCursor() *Cursor {
	c := new(Cursor)
	c.Row = 1
	c.Col = 1
	return c
}

func (w *Window) Clear() {
	syscall.Write(0, []byte("\033[2J"))
}

func (w *Window) ClearLine(col uint16) {
	w.MoveCursorPos(0, col)
	syscall.Write(0, []byte("\033[2K"))
}

func (w *Window) InitCursorPos() {
	syscall.Write(0, []byte("\033[1;1H"))
}

// row: 1~, col: 1~
func (w *Window) MoveCursorPos(row uint16, col uint16) {
	syscall.Write(0, []byte(fmt.Sprintf("\033[%d;%dH", col, row)))
}

func (w *Window) Display() {
	cnt := 0
	tmp := w.Rows.Next
	for {
		if tmp == w.Rows {
			break
		}
		cnt++
		fmt.Printf("\033[3m%4d\033[0m | %s\n", cnt, string(tmp.Row.GetAll()))
		tmp = tmp.Next
	}
}

// Define termios(unix) object
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

// Define original tty object
type Tty struct {
	original *termios
}

// Define winsize object
type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Set termios attribute
func tcSetAttr(fd uintptr, termios *termios) error {
	// TCSETS+1 == TCSETSW, because TCSAFLUSH doesn't exist
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS+1), uintptr(unsafe.Pointer(termios))); err != 0 {
		return err
	}
	return nil
}

// Get termios attribute
func tcGetAttr(fd uintptr) *termios {
	var termios = &termios{}
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios))); err != 0 {
		log.Fatalf("Problem getting terminal attributes: %s\n", err)
	}
	return termios
}

// Enable Raw / Non Canonical Mode
// termios - Iflag
// ^BRKINT: BrakeをNullバイトとして読み込む
// INPCK: 入力のパリティチェック有効化
// ICRNL: ^IGNCRの場合、CRをNLで置換
// ISTRIP: 8bit目を落とす
// IXON: 出力のXON/XOFFフロー制御の有効化

// termios - Cflag
// CS8: 文字サイズを8bitに指定

// termios - Lflag
// ECHO: 入力された文字をエコー
// ICANON: カノニカルモードの有効化
// IEXTEN: 実装依存の入力処理の有効化
// ISIG: シグナルを発生させる(Ctrl+C/Z etc...)
func (t *Tty) EnableRawMode() {
	t.original = tcGetAttr(os.Stdin.Fd())
	var raw termios
	raw = *t.original
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	raw.Cflag |= syscall.CS8
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	raw.Cc[syscall.VMIN+1] = 0
	raw.Cc[syscall.VTIME+1] = 1
	if e := tcSetAttr(os.Stdin.Fd(), &raw); e != nil {
		log.Fatalf("Problem enabling raw mode: %s\n", e)
	}
}

// Disable Raw / Non Canonical mode
func (t *Tty) DisableRawMode() {
	if e := tcSetAttr(os.Stdin.Fd(), t.original); e != nil {
		log.Fatalf("Problem disabling raw mode: %s\n", e)
	}
}

// Get console window size
func GetWinSize() (uint16, uint16) {
	var w WinSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&w)))
	return uint16(w.Row), uint16(w.Col)
}

// Enable Alternative Screen Buffer
func EnableASB() {
	syscall.Write(0, []byte("\033[?1049h"))
}

// Disable Alternative Screen Buffer
func DisableASB() {
	syscall.Write(0, []byte("\033[?1049l"))
}
