package readline

import "regexp"

func stripEscapeCode(s string) string {
	reg, _ := regexp.Compile("\033[[A-Za-z0-9]*m")
	return reg.ReplaceAllString(s, "")
}
