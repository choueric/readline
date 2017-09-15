package readline

type cmdLine struct {
	line  []rune // command line edit buffer
	index int    // current index of line
}

////////////////////////////////////////////////////////////////////////////////

func (cl *cmdLine) reset() {
	cl.line = make([]rune, 0)
	cl.index = 0
}

// insert a character right before the current cursor position
// and the cursor position plus one
func (cl *cmdLine) insert(c rune) {
	cl.line = append(cl.line, 0)
	copy(cl.line[cl.index+1:], cl.line[cl.index:])
	cl.line[cl.index] = c
	cl.index++
}

// delete the character under the cursor, then cursor position does not change
func (cl *cmdLine) del() {
	if len(cl.line) == 0 || cl.index == len(cl.line) {
		return
	}
	cl.line = append(cl.line[:cl.index], cl.line[cl.index+1:]...)
}

// backspace the character before the cursor, then cursor position reduces one
// return the deleted character
func (cl *cmdLine) backspace() rune {
	if len(cl.line) == 0 || cl.index == 0 {
		return 0
	}
	cl.index--
	c := cl.line[cl.index]
	cl.line = append(cl.line[:cl.index], cl.line[cl.index+1:]...)
	return c
}

// forward the cursor (rightward actually)
func (cl *cmdLine) forwardCursor() rune {
	if cl.index == cl.charNum() {
		return 0
	}
	c := cl.line[cl.index]
	cl.index = min(len(cl.line), cl.index+1)
	return c
}

// backward the cursor (leftward)
func (cl *cmdLine) backwardCursor() rune {
	if cl.index == 0 {
		return 0
	}
	c := cl.line[cl.index-1]
	cl.index = max(0, cl.index-1)
	return c
}

////////////////////////////////////////////////////////////////////////////////

// return the number of characters
func (cl *cmdLine) charNum() int {
	return len(cl.line)
}

// return the width in column of current input
func (cl *cmdLine) columnWidth() int {
	num := 0
	for _, c := range cl.line {
		num += columnWidth(c)
	}
	return num
}

func (cl cmdLine) String() string {
	return string(cl.line)
}
