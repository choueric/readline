package main

import (
	"io"
	"strings"
)

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
	inst.Print("\b \b")
	inst.line = inst.line[0 : len(inst.line)-1]
	return false
}

func enterHandler(inst *Instance) bool {
	line := strings.Trim(string(inst.line), " ")
	end := executeCmdline(inst, []byte(line))
	if !end {
		inst.clearLine()
		inst.printPrompt()
	}
	return end
}

func tabHandler(inst *Instance) bool {
	line := strings.Trim(string(inst.line), " ")
	if len(line) == 0 {
		acAllCmds(inst)
		inst.Printf("%s%s", inst.prompt, string(inst.line))
		return false
	}

	args := strings.Fields(string(line))
	if len(args) != 0 {
		acOneCmd(inst, args)
		return false
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
			inst.line = append(inst.line, c)
			inst.Printf("%c", c)
		}
	}
}
