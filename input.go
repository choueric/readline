package readline

import (
	"bufio"
	"io"
	"os"
)

type inputTerm struct {
	r        *bufio.Reader // for input
	termFd   int           // for setting the terminal input
	oldState *State        // old state of terminal
}

func (it *inputTerm) init() error {
	it.termFd = int(os.Stdin.Fd())
	oldState, err := MakeRaw(it.termFd)
	if err != nil {
		panic(err)
		return err
	}
	it.oldState = oldState

	it.r = bufio.NewReader(os.Stdin)

	return nil
}

func (it *inputTerm) deinit() {
	restoreTerm(it.termFd, it.oldState)
}

func (it *inputTerm) readByte() (byte, error) {
	return it.r.ReadByte()
}

func (it *inputTerm) readChar() (rune, int, error) {
	return it.r.ReadRune()
}

type inputHandler func(*Instance) (rune, bool)

var inputMap = map[rune]inputHandler{
	CharEOF:       eofHandler,
	CharTab:       tabHandler,
	CharEnter:     enterHandler,
	CharBackspace: backspaceHandler,
	CharInterrupt: interruptHandler,
	CharESC:       escapeHandler,
}

func eofHandler(inst *Instance) (rune, bool) {
	if inst.line.charNum() == 0 {
		inst.Printf("\n^D\n")
		return CharEOF, true
	}
	inst.cmdReset()
	return CharEOF, false
}

func tabHandler(inst *Instance) (rune, bool) {
	key := CharTab
	if inst.lastChar != CharTab { // First tab
		_, candidates, end, err := getCandidates(inst)
		if err != nil {
			inst.Log("1st tab error: %v\n", err)
			return key, false
		}
		switch len(candidates) {
		case 0:
			inst.Log("TODO: can not happen\n")
		case 1:
			completeWhole(inst, candidates[0])
			if end {
				inst.cmdInsert(' ')
				key = ' '
			}
		default:
			completePartial(inst, candidates)
		}
	} else { // Second Tab
		cp, candidates, _, err := getCandidates(inst)
		if err != nil {
			inst.Log("2nd tab error: %v\n", err)
			return key, false
		}
		switch len(candidates) {
		case 0, 1:
			inst.Log("TODO: can not happen\n")
		default:
			inst.Log("multi candidates\n")
			printCandidates(inst, cp, candidates)
			inst.Printf("%s%s", inst.view.prompt, inst.line.String())
		}
	}

	return key, false
}

func enterHandler(inst *Instance) (rune, bool) {
	inst.Print("\n")
	end := inst.execute(inst.line.String(), inst.data)
	if !end {
		inst.cmdReset()
	}
	return CharEnter, end
}

func backspaceHandler(inst *Instance) (rune, bool) {
	inst.cmdBackspace()
	return CharBackspace, false
}

func interruptHandler(inst *Instance) (rune, bool) {
	inst.cmdReset()
	inst.Printf("^C\n")
	return CharInterrupt, false
}

func escapeHandler(inst *Instance) (rune, bool) {
	c1, _ := inst.input.readByte()
	c2, _ := inst.input.readByte()
	inst.Log("  %q%q", c1, c2)
	if c1 != '[' {
		panic("wrong patter for escape code")
	}

	switch c2 {
	case 'C': // arrow right
		inst.cmdForwardCursor()
	case 'D': // arrow left
		inst.cmdBackwardCursor()
	case '3': // delete
		c, _ := inst.input.readByte()
		if c != '~' {
			panic("delete need ~ at the end.")
		}
		inst.cmdDel()
	default:
	}

	return CharESC, false
}

func InputLoop(inst *Instance) {
	inst.line.reset()
	inst.view.printPrompt()
	for {
		c, w, err := inst.input.readChar()
		if err != nil {
			if err == io.EOF {
				inst.Printf("got EOF error\n")
			} else {
				inst.Printf("error: %v\n", err)
			}
			break
		}

		inst.Log("[%q(%d):%d]", c, c, w)
		handler, ok := inputMap[c]
		if ok {
			key, end := handler(inst)
			if end {
				break
			}
			inst.lastChar = key
		} else {
			inst.cmdInsert(c)
			inst.lastChar = c
		}

		inst.view.clearLine()
		inst.Print(inst.view.prompt + inst.line.String())
		inst.view.setCursor(inst.view.curPos)
		inst.view.flush()
	}
}
