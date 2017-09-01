package main

import (
	"io"
	"strings"
)

const (
	PROMPT = ">> "
)

func printPrompt(inst *Instance) {
	inst.Printf(PROMPT)
	inst.Flush()
}

func helpHandler(inst *Instance) {
	inst.Println("Help:")
	for _, c := range inst.cmds {
		inst.Printf("  %s: %s\n", c.name, c.synopsis)
	}
	inst.Print("  exit, help: exit program\n")
}

func executeCmdline(inst *Instance, line []byte) int {
	cmdline := strings.Trim(string(line), " ")
	inst.Log("[%s]\n", cmdline)

	switch cmdline {
	case "help":
		helpHandler(inst)
	case "exit", "quit":
		return 1
	default:
	}

	return 0
}

func handleTab(inst *Instance, line []byte) {
	inst.Log(" autocomplete \n")
	inst.Printf("%s%s", PROMPT, string(line))
	inst.Flush()
}

func inputLoop(inst *Instance) {
	line := make([]byte, 0)
	printPrompt(inst)
	end := false
	for !end {
		c, err := inst.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				inst.Printf("got EOF\n")
				inst.Flush()
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
			printPrompt(inst)
		case CharEOF:
			if len(line) == 0 {
				inst.Printf("\ngot EOF(Ctrl+D)\n")
				inst.Flush()
				end = true
			} else {
				line = line[0:0]
				inst.Printf("\n")
				printPrompt(inst)
			}
		case CharEnter:
			ret := executeCmdline(inst, line)
			if ret != 0 {
				end = true
				break
			}
			line = line[0:0]
			printPrompt(inst)
		case CharTab:
			handleTab(inst, line)
		default:
			line = append(line, c)
			inst.Printf("%c", c)
			inst.Flush()
		}
	}
}
