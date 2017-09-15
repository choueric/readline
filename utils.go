package readline

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

func stripEscapeCode(s string) string {
	reg, _ := regexp.Compile("\033[^m]*m")
	return reg.ReplaceAllString(s, "")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func columnWidth(c rune) int {
	return runewidth.RuneWidth(c)
}
