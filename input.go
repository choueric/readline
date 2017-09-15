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

type inputHandler func(*Instance) (byte, bool)

var inputMap = map[byte]inputHandler{
	CharEOF:       eofHandler,
	CharTab:       tabHandler,
	CharEnter:     enterHandler,
	CharBackspace: backspaceHandler,
	CharInterrupt: interruptHandler,
	CharESC:       escapeHandler,
}

func eofHandler(inst *Instance) (byte, bool) {
	if inst.line.Len() == 0 {
		inst.Printf("\n^D\n")
		return CharEOF, true
	}
	inst.resetCmdline()
	return CharEOF, false
}

func tabHandler(inst *Instance) (byte, bool) {
	key := byte(CharTab)
	if inst.lastKey != CharTab { // First tab
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
				inst.line.insert(' ')
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

func enterHandler(inst *Instance) (byte, bool) {
	inst.Print("\n")
	end := inst.execute(inst.line.String(), inst.data)
	if !end {
		inst.resetCmdline()
	}
	return CharEnter, end
}

func backspaceHandler(inst *Instance) (byte, bool) {
	inst.line.backspace()
	return CharBackspace, false
}

func interruptHandler(inst *Instance) (byte, bool) {
	inst.resetCmdline()
	inst.Printf("^C\n")
	return CharInterrupt, false
}

func escapeHandler(inst *Instance) (byte, bool) {
	c1, _ := inst.input.readByte()
	c2, _ := inst.input.readByte()
	if c1 != '[' {
		panic("wrong patter for escape code")
	}

	switch c2 {
	case 'C': // arrow right
		inst.line.forwardCursor()
	case 'D': // arrow left
		inst.line.backwardCursor()
	default:
	}

	return CharESC, false
}

func InputLoop(inst *Instance) {
	inst.line.reset()
	inst.view.printPrompt()
	for {
		c, err := inst.input.readByte()
		if err != nil {
			if err == io.EOF {
				inst.Printf("got EOF\n")
			} else {
				inst.Printf("error: %v\n", err)
			}
			break
		}

		//inst.Log("[%d]", c)
		handler, ok := inputMap[c]
		if ok {
			key, end := handler(inst)
			if end {
				break
			}
			inst.lastKey = key
		} else {
			inst.line.insert(c)
			inst.lastKey = c
		}

		inst.view.clearLine()
		inst.Print(inst.view.prompt + inst.line.String())
		inst.view.setCursor(inst.line.curPos)
		inst.view.flush()
	}
}
