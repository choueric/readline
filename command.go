package main

import "strings"

type CmdHandler func([]string, interface{}) error

type Cmd struct {
	synopsis string
	handler  CmdHandler
	data     interface{}
}

// TODO: sort the commands
func helpHandler(args []string, data interface{}) error {
	inst := data.(*Instance)
	inst.Println("Help:")
	for n, c := range inst.cmds {
		inst.Printf("  %s: %s\n", n, c.synopsis)
	}

	return nil
}

func executeCmdline(inst *Instance) bool {
	if len(inst.line) == 0 {
		return false
	}
	cmdline := strings.Fields(string(inst.line))
	if len(cmdline) == 0 {
		inst.Log("parse input line [%s] failed\n", string(inst.line))
		return false
	}
	inst.Log("[%v]\n", cmdline)

	switch cmdline[0] {
	case "help":
		helpHandler(cmdline, inst)
	case "exit", "quit":
		return true
	default:
		c, ok := inst.cmds[cmdline[0]]
		if ok {
			c.handler(cmdline, c.data)
		} else {
			inst.Printf("Invalid command [%s]\n", cmdline[0])
			helpHandler(cmdline, inst)
		}
	}

	return false
}
