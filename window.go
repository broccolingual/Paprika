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
	c.Row = 0
	c.Col = 0
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

func (w *Window) MoveCursorPos(row uint16, col uint16) {
	c := fmt.Sprintf("\033[%d;%dH", col+1, row+1)
	syscall.Write(0, []byte(c))
}

func (w *Window) UpdateLine(f *File) {
	w.ClearLine(w.Cursor.Col)
	fmt.Printf("%3d] %s", w.Cursor.Col, f.Lines[w.Cursor.Col])
}

// NEW
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

type Tty struct {
	original *termios
}

type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func tcSetAttr(fd uintptr, termios *termios) error {
	// TCSETS+1 == TCSETSW, because TCSAFLUSH doesn't exist
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS+1), uintptr(unsafe.Pointer(termios))); err != 0 {
		return err
	}
	return nil
}

func tcGetAttr(fd uintptr) *termios {
	var termios = &termios{}
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios))); err != 0 {
		log.Fatalf("Problem getting terminal attributes: %s\n", err)
	}
	return termios
}

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

func (t *Tty) DisableRawMode() {
	if e := tcSetAttr(os.Stdin.Fd(), t.original); e != nil {
		log.Fatalf("Problem disabling raw mode: %s\n", e)
	}
}

func GetWinSize() (uint16, uint16) {
	var w WinSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&w)))
	return uint16(w.Row), uint16(w.Col)
}
