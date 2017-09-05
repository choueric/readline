package main

import "io"

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
	inst.printPrompt()
	return CharInterrupt, false
}

func backspaceHandler(inst *Instance) (byte, bool) {
	if len(inst.line) == 0 {
		return CharBackspace, false
	}
	inst.lineDel()
	return CharBackspace, false
}

func enterHandler(inst *Instance) (byte, bool) {
	end := inst.execute(string(inst.line), inst.data)
	if !end {
		inst.clearLine()
		inst.printPrompt()
	}
	return CharEnter, end
}

func tabHandler(inst *Instance) (byte, bool) {
	key := byte(CharTab)
	if inst.lastKey != CharTab { // First tab
		candidates, end, err := getCandidates(inst)
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
		candidates, _, err := getCandidates(inst)
		if err != nil {
			inst.Log("2nd tab error: %v\n", err)
			return key, false
		}
		switch len(candidates) {
		case 0, 1:
			inst.Log("TODO: can not happen\n")
		default:
			inst.Log("multi candidates\n")
			printCandidates(inst, candidates)
			inst.Printf("%s%s", inst.prompt, string(inst.line))
		}
	}

	return key, false
}

func eofHandler(inst *Instance) (byte, bool) {
	if len(inst.line) == 0 {
		inst.Printf("\ngot EOF(Ctrl+D)\n")
		return CharEOF, true
	}
	inst.clearLine()
	inst.printPrompt()
	return CharEOF, false
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
		key := c
		if ok {
			key, end = handler(inst)
		} else {
			inst.lineAdd(c)
		}
		inst.lastKey = key
	}
}
