package common

import (
	"fmt"
	"os"
	"syscall"
	
	"golang.org/x/sys/unix"
	"github.com/pkg/term/termios"
)

type UnixTerm _UnixTerm

type _UnixTerm struct {
	origTtyState *unix.Termios
}

func NewUnixTerm() *UnixTerm {
	term := new(UnixTerm)
	return term
}

// Set termios attribute
func (term *UnixTerm) tcSetAttr(attr *unix.Termios) error {
	return termios.Tcsetattr(uintptr(os.Stdin.Fd()), termios.TCSANOW, attr)
}

// Get termios attribute
func (term *UnixTerm) tcGetAttr() (*unix.Termios, error) {
	var attr unix.Termios
	err := termios.Tcgetattr(uintptr(os.Stdin.Fd()), &attr)
	return &attr, err
}

// Rawモード/非カノニカルモードの有効化
// https://linuxjm.osdn.jp/html/LDP_man-pages/man3/termios.3.html
func (term *UnixTerm) EnableRawMode() {
	term.origTtyState, _ = term.tcGetAttr()
	var attr unix.Termios
	termios.Cfmakeraw(&attr)
	term.tcSetAttr(&attr)
}

// Rawモード/非カノニカルモードの無効化
func (term *UnixTerm) DisableRawMode() {
	term.tcSetAttr(term.origTtyState)
}

// エスケープシーケンスの送信
func (term *UnixTerm) setAttr(code string) {
	syscall.Write(0, []byte(code))
}

// Alternative Screen Bufferの有効化
func (term *UnixTerm) EnableAlternativeScreenBuffer() {
	term.setAttr("\033[?1049h")
}

// Alternative Screen Bufferの無効化
func (term *UnixTerm) DisableAlternativeScreenBuffer() {
	term.setAttr("\033[?1049l")
}

// カーソルの有効化
func (term *UnixTerm) EnableCursor() {
	term.setAttr("\033[?25h")
}

// カーソルの無効化
func (term *UnixTerm) DisableCursor() {
	term.setAttr("\033[?25l")
}

// カーソル以降をすべて消去
func (term *UnixTerm) ClearAfterCursor() {
	term.setAttr("\033[0J")
}

// カーソル以前をすべて消去
func (term *UnixTerm) ClearBeforeCursor() {
	term.setAttr("\033[1J")
}

// スクリーンをすべて消去
func (term *UnixTerm) ClearAll() {
	term.setAttr("\033[2J")
}

// その行のカーソルの右端を消去
func (term *UnixTerm) ClearRowRight() {
	term.setAttr("\033[0K")
}

// その行カーソルの左端を消去
func (term *UnixTerm) ClearRowLeft() {
	term.setAttr("\033[1K")
}

// その行を消去
func (term *UnixTerm) ClearRow() {
	term.setAttr("\033[2K")
}

// 1行1列にカーソルを移動
func (term *UnixTerm) InitCursorPos() {
	term.setAttr("\033[1;1H")
}

// 対象行・列にカーソルを移動
// row: 1~, col: 1~
func (term *UnixTerm) MoveCursorPos(col uint, row uint) {
	term.setAttr(fmt.Sprintf("\033[%d;%dH", row, col))
}

func (term *UnixTerm) ScrollDown(n uint8) {
	term.setAttr(fmt.Sprintf("\033[%dS", n))
}

func (term *UnixTerm) ScrollUp(n uint8) {
	term.setAttr(fmt.Sprintf("\033[%dT", n))
}

func (term *UnixTerm) SetColor(c uint8) {
	term.setAttr(fmt.Sprintf("\033[38;5;%dm", c))
}

func (term *UnixTerm) SetBGColor(c uint8) {
	term.setAttr(fmt.Sprintf("\033[48;5;%dm", c))
}

func (term *UnixTerm) SetBold() {
	term.setAttr("\033[1m")
}

func (term *UnixTerm) SetItalic() {
	term.setAttr("\033[3m")
}

func (term *UnixTerm) SetUnderbar() {
	term.setAttr("\033[4m")
}

func (term *UnixTerm) SetBlink() {
	term.setAttr("\033[5m")
}

func (term *UnixTerm) SetFastBlink() {
	term.setAttr("\033[6m")
}

func (term *UnixTerm) SetInversion() {
	term.setAttr("\033[7m")
}

func (term *UnixTerm) SetHide() {
	term.setAttr("\033[8m")
}

func (term *UnixTerm) ResetStyle() {
	term.setAttr("\033[m")
}
