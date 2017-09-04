package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func executeCmdline(line string, data interface{}) bool {
	cmdline := strings.Fields(line)
	if len(cmdline) == 0 {
		fmt.Printf("invalid command [%s]\n", line)
		return false
	}
	switch cmdline[0] {
	case "exit":
		return true
	default:
		fmt.Printf("\nexecute [%v]\n", cmdline)
	}
	return false
}

func main() {
	inst := &Instance{
		prompt: "\033[32m>>\033[0m ",
	}
	inst.Init(os.Stdin, os.Stdout)
	defer inst.Deinit()

	flag.BoolVar(&inst.debug, "d", false, "eanble debug")
	flag.Parse()

	inst.SetExecute(executeCmdline, nil)
	inst.SetCmds(
		Item("ls"),
		Item("lsblk"),
		Item("git",
			Item("clone"),
			Item("clean"),
			Item("log",
				Item("all"),
				Item("verbose"),
			),
		),
		Item("exit"),
		Item("help"),
	)
	inst.PrintTree(os.Stdout)

	inputLoop(inst)
}
