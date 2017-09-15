package readline

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type viewTerm struct {
	prompt string
	w      *bufio.Writer // for output
}

// make stdin raw mode and use stdout as output
func (vt *viewTerm) init(prompt string) error {
	vt.w = bufio.NewWriter(os.Stdout)
	vt.prompt = prompt
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

func (vt *viewTerm) forwardCursor() {
	fmt.Fprintf(vt.w, "\033[2C")
	time.Sleep(time.Second * 1)
}

func (vt *viewTerm) backwardCursor() {
	fmt.Fprintf(vt.w, "\033[2D")
	time.Sleep(time.Second * 1)
}

func (vt *viewTerm) insertRune(r rune) {
	fmt.Fprintf(vt.w, string(r))
	vt.w.Flush()
	time.Sleep(time.Second * 1)
}
