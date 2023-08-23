package common

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type UnixTerm _UnixTerm

type _UnixTerm struct {
	origTermSetting *termios
}

// TODO: デフォルトのtermiosライブラリへの置換

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
// https://linuxjm.osdn.jp/html/LDP_man-pages/man3/termios.3.html
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
func (ut *UnixTerm) EnableRawMode() {
	ut.origTermSetting = tcGetAttr(os.Stdin.Fd())
	var raw termios
	raw = *ut.origTermSetting
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cc[syscall.VMIN+1] = 0
	raw.Cc[syscall.VTIME+1] = 1
	if e := tcSetAttr(os.Stdin.Fd(), &raw); e != nil {
		log.Fatalf("Problem enabling raw mode: %s\n", e)
	}
}

// Rawモード/非カノニカルモードの無効化
func (ut *UnixTerm) DisableRawMode() {
	if e := tcSetAttr(os.Stdin.Fd(), ut.origTermSetting); e != nil {
		log.Fatalf("Problem disabling raw mode: %s\n", e)
	}
}

func NewUnixTerm() *UnixTerm {
	ut := new(UnixTerm)
	return ut
}

// エスケープシーケンスの送信
func (ut *UnixTerm) SetAttr(code string) {
	syscall.Write(0, []byte(code))
}

// Alternative Screen Bufferの有効化
func (ut *UnixTerm) EnableAlternativeScreenBuffer() {
	ut.SetAttr("\033[?1049h")
}

// Alternative Screen Bufferの無効化
func (ut *UnixTerm) DisableAlternativeScreenBuffer() {
	ut.SetAttr("\033[?1049l")
}

// カーソルの有効化
func (ut *UnixTerm) EnableCursor() {
	ut.SetAttr("\033[?25h")
}

// カーソルの無効化
func (ut *UnixTerm) DisableCursor() {
	ut.SetAttr("\033[?25l")
}

// カーソル以降をすべて消去
func (ut *UnixTerm) ClearAfterCursor() {
	ut.SetAttr("\033[0J")
}

// カーソル以前をすべて消去
func (ut *UnixTerm) ClearBeforeCursor() {
	ut.SetAttr("\033[1J")
}

// スクリーンをすべて消去
func (ut *UnixTerm) ClearAll() {
	ut.SetAttr("\033[2J")
}

// その行のカーソルの右端を消去
func (ut *UnixTerm) ClearRowRight() {
	ut.SetAttr("\033[0K")
}

// その行カーソルの左端を消去
func (ut *UnixTerm) ClearRowLeft() {
	ut.SetAttr("\033[1K")
}

// その行を消去
func (ut *UnixTerm) ClearRow() {
	ut.SetAttr("\033[2K")
}

// 1行1列にカーソルを移動
func (ut *UnixTerm) InitCursorPos() {
	ut.SetAttr("\033[1;1H")
}

// 対象行・列にカーソルを移動
// row: 1~, col: 1~
func (ut *UnixTerm) MoveCursorPos(col uint, row uint) {
	ut.SetAttr(fmt.Sprintf("\033[%d;%dH", row, col))
}

func (ut *UnixTerm) ScrollDown(n uint8) {
	ut.SetAttr(fmt.Sprintf("\033[%dS", n))
}

func (ut *UnixTerm) ScrollUp(n uint8) {
	ut.SetAttr(fmt.Sprintf("\033[%dT", n))
}

func (ut *UnixTerm) SetColor(c uint8) {
	ut.SetAttr(fmt.Sprintf("\033[38;5;%dm", c))
}

func (ut *UnixTerm) SetBGColor(c uint8) {
	ut.SetAttr(fmt.Sprintf("\033[48;5;%dm", c))
}

func (ut *UnixTerm) SetBold() {
	ut.SetAttr("\033[1m")
}

func (ut *UnixTerm) SetItalic() {
	ut.SetAttr("\033[3m")
}

func (ut *UnixTerm) SetUnderbar() {
	ut.SetAttr("\033[4m")
}

func (ut *UnixTerm) SetBlink() {
	ut.SetAttr("\033[5m")
}

func (ut *UnixTerm) SetFastBlink() {
	ut.SetAttr("\033[6m")
}

func (ut *UnixTerm) SetInversion() {
	ut.SetAttr("\033[7m")
}

func (ut *UnixTerm) SetHide() {
	ut.SetAttr("\033[8m")
}

func (ut *UnixTerm) ResetStyle() {
	ut.SetAttr("\033[m")
}
