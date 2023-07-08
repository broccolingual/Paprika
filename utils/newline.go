package utils

type NLCode int8 // 改行文字タイプの定義

const (
	CR NLCode = iota
	LF 
	CRLF
)

func GetNLCode(runes []rune) NLCode {
	if runes[len(runes)-1] == rune('\r') {
		return CR
	} else if runes[len(runes)-1] == rune('\n') {
		if runes[len(runes)-2] == rune('\r') {
			return CRLF
		}
		return LF
	}
	return -1
}
