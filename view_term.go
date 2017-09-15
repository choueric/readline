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

func (vt *viewTerm) Printf(format string, v ...interface{}) {
	fmt.Fprintf(vt.w, format, v...)
	vt.w.Flush()
}

func (vt *viewTerm) Print(v ...interface{}) {
	fmt.Fprint(vt.w, v...)
	vt.w.Flush()
}

func (vt *viewTerm) Println(v ...interface{}) {
	fmt.Fprintln(vt.w, v...)
	vt.w.Flush()
}

func (vt *viewTerm) Error(v ...interface{}) {
	fmt.Fprint(vt.w, v...)
	vt.w.Flush()
}

func (vt *viewTerm) Errorf(format string, v ...interface{}) {
	fmt.Fprintf(vt.w, format, v...)
	vt.w.Flush()
}

func (vt *viewTerm) printPrompt() {
	fmt.Fprint(vt.w, "\n"+vt.prompt)
	vt.w.Flush()
}

func (vt *viewTerm) flush() {
	vt.w.Flush()
}

////////////////////////////////////////////////////////////////////////////////

func (vt *viewTerm) resetCursor() {
	fmt.Fprint(vt.w, "\033[1000D")
}

func (vt *viewTerm) clearLine() {
	fmt.Fprint(vt.w, "\033[1000D")
	fmt.Fprint(vt.w, "\033[0K")
}

func (vt *viewTerm) setCursor(pos int) {
	pos += vt.promptLen
	if pos != 0 {
		fmt.Fprintf(vt.w, "\033[%dC", pos)
	}
}

func (vt *viewTerm) forwardCursor() {
	fmt.Fprint(vt.w, "\033[2C")
}

func (vt *viewTerm) backwardCursor() {
	fmt.Fprint(vt.w, "\033[2D")
}

func (vt *viewTerm) insertRune(r rune) {
	fmt.Fprint(vt.w, string(r))
	vt.w.Flush()
}
