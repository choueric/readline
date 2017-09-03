package main

import (
	"io"
	"strings"
)

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

func handleTab(inst *Instance, line []byte) {
	inst.Log(" autocomplete \n")
	inst.Printf("%s%s", inst.prompt, string(line))
}

func inputLoop(inst *Instance) {
	line := make([]byte, 0)
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
		switch c {
		case CharInterrupt:
			inst.Printf("\ngot Interrupt(Ctrl+C)\n")
			line = line[0:0]
			inst.Printf("\n")
			inst.printPrompt()
		case CharEOF:
			if len(line) == 0 {
				inst.Printf("\ngot EOF(Ctrl+D)\n")
				end = true
			} else {
				line = line[0:0]
				inst.Printf("\n")
				inst.printPrompt()
			}
		case CharEnter:
			ret := executeCmdline(inst, line)
			if ret != 0 {
				end = true
				break
			}
			line = line[0:0]
			inst.printPrompt()
		case CharTab:
			handleTab(inst, line)
		case CharBackspace:
			inst.Print("\b \b")
			line = line[0 : len(line)-1]
		default:
			line = append(line, c)
			inst.Printf("%c", c)
		}
	}
}
