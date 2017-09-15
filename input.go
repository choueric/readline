package readline

import (
	"bufio"
	"io"
	"os"
)

/*
   Clear Screen: \u001b[{n}J clears the screen
       n=0 clears from cursor until end of screen,
       n=1 clears from cursor to beginning of screen
       n=2 clears entire screen
   Clear Line: \u001b[{n}K clears the current line
       n=0 clears from cursor to end of line
       n=1 clears from cursor to start of line
       n=2 clears entire line
*/

const (
	ESC = 0x1b // \033
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
}

func interruptHandler(inst *Instance) (byte, bool) {
	inst.Printf("\ngot Interrupt(Ctrl+C)\n")
	inst.clearLine()
	inst.view.printPrompt()
	return CharInterrupt, false
}

func backspaceHandler(inst *Instance) (byte, bool) {
	if inst.line.Len() == 0 {
		return CharBackspace, false
	}
	inst.lineDel()
	return CharBackspace, false
}

func enterHandler(inst *Instance) (byte, bool) {
	inst.Print("\n")
	end := inst.execute(inst.line.String(), inst.data)
	if !end {
		inst.clearLine()
		inst.view.printPrompt()
	}
	return CharEnter, end
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
				inst.lineAdd(' ')
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

func eofHandler(inst *Instance) (byte, bool) {
	if inst.line.Len() == 0 {
		inst.Printf("\ngot EOF(Ctrl+D)\n")
		return CharEOF, true
	}
	inst.clearLine()
	inst.view.printPrompt()
	return CharEOF, false
}

func InputLoop(inst *Instance) {
	inst.view.printPrompt()
	end := false
	for !end {
		c, err := inst.input.readByte()
		if err != nil {
			if err == io.EOF {
				inst.Printf("got EOF\n")
			} else {
				inst.Printf("error: %v\n", err)
			}
			end = true
		}

		//inst.Log("[%d]", c)
		handler, ok := inputMap[c]
		key := c
		if ok {
			key, end = handler(inst)
		} else {
			inst.line.insert(c)
		}
		inst.lastKey = key
	}
}
