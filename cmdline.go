package readline

import "fmt"

type cmdLine struct {
	line   []rune // command line edit buffer
	curPos int    // cursor position, i.e the index of the rune slice
	// But it can be len(line), which means at the end of line
}

// insert a rune right before the current cursor position
// and the cursor position plus one
func (cl *cmdLine) insertRune(r rune) {
	cl.line = append(cl.line, 0)
	copy(cl.line[cl.curPos+1:], cl.line[cl.curPos:])
	cl.line[cl.curPos] = r
	cl.curPos++
}

// delete the rune under the cursor, then cursor position does not change
func (cl *cmdLine) deleteRune() {
	if len(cl.line) == 0 || cl.curPos == len(cl.line) {
		return
	}
	cl.line = append(cl.line[:cl.curPos], cl.line[cl.curPos+1:]...)
}

// backspace the rune before the cursor, then cursor position reduces one
func (cl *cmdLine) backspaceRune() {
	if len(cl.line) == 0 || cl.curPos == 0 {
		return
	}
	cl.curPos--
	cl.line = append(cl.line[:cl.curPos], cl.line[cl.curPos+1:]...)
}

// forward the cursor (rightward actually)
func (cl *cmdLine) forwardCursor() int {
	if cl.curPos == len(cl.line) {
		return cl.curPos
	}
	cl.curPos++
	return cl.curPos
}

// backward the cursor (leftward)
func (cl *cmdLine) backwardCursor() int {
	if cl.curPos == 0 {
		return cl.curPos
	}
	cl.curPos--
	return cl.curPos
}

func (cl *cmdLine) Len() int {
	return len(cl.line)
}

func (cl cmdLine) String() string {
	return string(cl.line)
}

func (cl cmdLine) prettyStr() string {
	return fmt.Sprintf("(%s, %d)", string(cl.line), cl.curPos)
}
