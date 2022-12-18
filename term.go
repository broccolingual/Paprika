package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

type UnixTerm interface {
	EnableCursor()
	DisableCursor()
	ScreenClear()
	LineClear()
	InitCursorPos()
	MoveCursorPos()
	EnableASB()
	DisableASB()
}

type _UnixTerm struct{}

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

// Enable Alternative Screen Buffer
func EnableASB() {
	syscall.Write(0, []byte("\033[?1049h"))
}

// Disable Alternative Screen Buffer
func DisableASB() {
	syscall.Write(0, []byte("\033[?1049l"))
}
