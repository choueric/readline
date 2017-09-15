package readline

import "regexp"

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

// TODO: improve this method
func columnWidth(c rune) int {
	w := len(string(c))
	if w == 1 {
		return 1
	} else {
		return 2
	}
}
