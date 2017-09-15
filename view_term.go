package readline

import (
	"bufio"
	"fmt"
	"os"
)

type viewTerm struct {
	prompt    string
	promptLen int           // cursor length of prompt
	w         *bufio.Writer // for output
	curPos    int           // cursor position on terminal
}

// make stdin raw mode and use stdout as output
func (vt *viewTerm) init(prompt string) error {
	vt.w = bufio.NewWriter(os.Stdout)
	vt.prompt = prompt
	vt.promptLen = len(stripEscapeCode(prompt))
	return nil
}

func (vt *viewTerm) deinit() {
}

////////////////////////////////////////////////////////////////////////////////

func (vt *viewTerm) reset() {
	vt.curPos = 0
	fmt.Fprint(vt.w, "\n")
}

func (vt *viewTerm) insert(c rune) {
	vt.curPos += columnWidth(c)
}

func (vt *viewTerm) del(c rune) {
}

func (vt *viewTerm) backspace(c rune) {
	vt.curPos = max(0, vt.curPos-columnWidth(c))
}

func (vt *viewTerm) forwardCursor(c rune, curLen int) {
	vt.curPos = min(curLen, vt.curPos+columnWidth(c))
}

func (vt *viewTerm) backwardCursor(c rune) {
	vt.curPos = max(0, vt.curPos-columnWidth(c))
}

////////////////////////////////////////////////////////////////////////////////

func (vt *viewTerm) printPrompt() {
	fmt.Fprint(vt.w, "\n"+vt.prompt)
	vt.w.Flush()
}

func (vt *viewTerm) flush() {
	vt.w.Flush()
}

////////////////////////////////////////////////////////////////////////////////

func (vt *viewTerm) clearLine() {
	fmt.Fprint(vt.w, "\033[1000D")
	fmt.Fprint(vt.w, "\033[0K")
}

func (vt *viewTerm) setCursor(pos int) {
	fmt.Fprint(vt.w, "\033[1000D")
	pos += vt.promptLen
	if pos != 0 {
		fmt.Fprintf(vt.w, "\033[%dC", pos)
	}
}
