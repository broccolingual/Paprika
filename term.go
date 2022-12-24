package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type UnixTerm interface {
	SetAttr(code string)
	EnableASB()
	DisableASB()
	EnableCursor()
	DisableCursor()
	ClearAfterCursor()
	ClearBeforeCursor()
	ClearAll()
	ClearRowRight()
	ClearRowLeft()
	ClearRow()
	InitCursorPos()
	MoveCursorPos(col uint16, row uint16)
}

type unixTerm struct{}

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

// Rawモード/非カノニカルモードの有効化
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

// Rawモード/非カノニカルモードの無効化
func (t *Tty) DisableRawMode() {
	if e := tcSetAttr(os.Stdin.Fd(), t.original); e != nil {
		log.Fatalf("Problem disabling raw mode: %s\n", e)
	}
}

func NewUnixTerm() UnixTerm {
	ut := new(unixTerm)
	return ut
}

// エスケープシーケンスの送信
func (ut *unixTerm) SetAttr(code string) {
	syscall.Write(0, []byte(code))
}

// Alternative Screen Bufferの有効化
func (ut *unixTerm) EnableASB() {
	ut.SetAttr("\033[?1049h")
}

// Alternative Screen Bufferの無効化
func (ut *unixTerm) DisableASB() {
	ut.SetAttr("\033[?1049l")
}

// カーソルの有効化
func (ut *unixTerm) EnableCursor() {
	ut.SetAttr("\033[?25h")
}

// カーソルの無効化
func (ut *unixTerm) DisableCursor() {
	ut.SetAttr("\033[?25l")
}

// カーソル以降をすべて消去
func (ut *unixTerm) ClearAfterCursor() {
	ut.SetAttr("\033[0J")
}

// カーソル以前をすべて消去
func (ut *unixTerm) ClearBeforeCursor() {
	ut.SetAttr("\033[1J")
}

// スクリーンをすべて消去
func (ut *unixTerm) ClearAll() {
	ut.SetAttr("\033[2J")
}

// その行のカーソルの右端を消去
func (ut *unixTerm) ClearRowRight() {
	ut.SetAttr("\033[0K")
}

// その行カーソルの左端を消去
func (ut *unixTerm) ClearRowLeft() {
	ut.SetAttr("\033[1K")
}

// その行を消去
func (ut *unixTerm) ClearRow() {
	ut.SetAttr("\033[2K")
}

// 1行1列にカーソルを移動
func (ut *unixTerm) InitCursorPos() {
	ut.SetAttr("\033[1;1H")
}

// 対象行・列にカーソルを移動
// row: 1~, col: 1~
func (ut *unixTerm) MoveCursorPos(col uint16, row uint16) {
	ut.SetAttr(fmt.Sprintf("\033[%d;%dH", row, col))
}
