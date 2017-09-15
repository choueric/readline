package readline

import "regexp"

func stripEscapeCode(s string) string {
	reg, _ := regexp.Compile("\033[^m]*m")
	return reg.ReplaceAllString(s, "")
}
