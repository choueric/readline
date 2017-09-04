package main

import "io"

type inputHandler func(*Instance) bool

var inputMap = map[byte]inputHandler{
	CharEOF:       eofHandler,
	CharTab:       tabHandler,
	CharEnter:     enterHandler,
	CharBackspace: backspaceHandler,
	CharInterrupt: interruptHandler,
}

func interruptHandler(inst *Instance) bool {
	inst.Printf("\ngot Interrupt(Ctrl+C)\n")
	inst.clearLine()
	inst.printPrompt()
	return false
}

func backspaceHandler(inst *Instance) bool {
	if len(inst.line) == 0 {
		return false
	}
	inst.lineDel()
	return false
}

func enterHandler(inst *Instance) bool {
	end := inst.execute(string(inst.line), inst.data)
	if !end {
		inst.clearLine()
		inst.printPrompt()
	}
	return end
}

func tabHandler(inst *Instance) bool {
	if inst.lastKey != CharTab { // First tab
		candidates, err := getCandidates(inst)
		if err != nil {
			inst.Log("1st tab error: %v\n", err)
			return false
		}
		switch len(candidates) {
		case 0:
			inst.Log("TODO: use auto-completer interface\n")
		case 1:
			completeWhole(inst, candidates[0])
		default:
			completePartial(inst, candidates)
		}
	} else { // Second Tab
		candidates, err := getCandidates(inst)
		if err != nil {
			inst.Log("2nd tab error: %v\n", err)
			return false
		}
		switch len(candidates) {
		case 0:
			inst.Log("TODO: use filelist completer\n")
		case 1:
			panic("This can not be happen")
		default:
			printCandidates(inst, candidates)
			inst.Printf("%s%s", inst.prompt, string(inst.line))
		}
	}

	return false
}

func eofHandler(inst *Instance) bool {
	if len(inst.line) == 0 {
		inst.Printf("\ngot EOF(Ctrl+D)\n")
		return true
	}
	inst.clearLine()
	inst.printPrompt()
	return false
}

func inputLoop(inst *Instance) {
	inst.printPrompt()
	end := false
	for !end {
		c, err := inst.r.ReadByte()
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
		if ok {
			end = handler(inst)
		} else {
			inst.lineAdd(c)
		}
		inst.lastKey = c
	}
}
