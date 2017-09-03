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
	inst.line = inst.line[0:0]
	inst.Printf("\n")
	inst.printPrompt()
	return false
}

func backspaceHandler(inst *Instance) bool {
	inst.Print("\b \b")
	inst.line = inst.line[0 : len(inst.line)-1]
	return false
}

func enterHandler(inst *Instance) bool {
	if executeCmdline(inst, inst.line) != 0 {
		return true
	}
	inst.line = inst.line[0:0]
	inst.printPrompt()
	return false
}

func tabHandler(inst *Instance) bool {
	inst.Log(" autocomplete \n")
	inst.Printf("%s%s", inst.prompt, string(inst.line))
	return false
}

func eofHandler(inst *Instance) bool {
	if len(inst.line) == 0 {
		inst.Printf("\ngot EOF(Ctrl+D)\n")
		return true
	}
	inst.line = inst.line[0:0]
	inst.Printf("\n")
	inst.printPrompt()
	return false
}

func helpHandler(inst *Instance) {
	inst.Println("Help:")
	for _, c := range inst.cmds {
		inst.Printf("  %s: %s\n", c.name, c.synopsis)
	}
	inst.Print("  exit, help: exit program\n")
}

func executeCmdline(inst *Instance, line []byte) int {
	ret := 0
	cmdline := strings.Fields(string(line))
	if len(cmdline) == 0 {
		inst.Log("parse input line failed: %s", string(line))
		return ret
	}
	inst.Log("[%v]\n", cmdline)

	switch cmdline[0] {
	case "help":
		helpHandler(inst)
	case "exit", "quit":
		return 1
	default:
		var cmd *Cmd
		for _, c := range inst.cmds {
			if cmdline[0] == c.name {
				cmd = c
			}
		}
		if cmd != nil {
			cmd.handler(cmdline)
		} else {
			helpHandler(inst)
		}
	}

	return ret
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
				inst.Printf("error: %v", err)
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
