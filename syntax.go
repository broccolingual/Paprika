package main

import (
	"fmt"
	"strings"
)

var goReserved = []string{
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"if",
	"import",
	"interface",
	"map",
	"goto",
	"package",
	"range",
	"return",
	"type",
	"select",
	"struct",
	"switch",
	"var",
}

func containsToken(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Tokenize(rowData string) []string {
	return strings.Split(rowData, " ")
}

func Highlighter(token []string, lang string, focus bool) string {
	var concat string
	switch lang {
	case ".go":
		for _, ch := range token {
			if containsToken(goReserved, ch) {
				if focus == false {
					concat += fmt.Sprintf("\033[38;5;3m%s\033[m ", ch)
				} else {
					concat += fmt.Sprintf("\033[38;5;3m%s\033[m\033[48;5;235m ", ch)
				}
			} else {
				concat += ch + " "
			}
		}
	default:
	}
	return concat
}
