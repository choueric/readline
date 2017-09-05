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

// 1. '', ' '
// 2. 'l'
// 3. 'ls'
// 4. 'git'
// 5. 'git '
// 6. 'git add '
// 7. 'git add co'
// 8. 'git add -'
// 9. 'git add pre': command has priority to fs
func main() {
	inst := &Instance{
		prompt: "\033[32m>>\033[0m ",
	}
	inst.Init(os.Stdin, os.Stdout)
	defer inst.Deinit()

	flag.BoolVar(&inst.debug, "d", false, "eanble debug")
	flag.Parse()

	inst.SetExecute(executeCmdline, nil)
	inst.SetCompleter(
		Cmd("ls",
			ListFs()),
		Cmd("lsblk"),
		Cmd("git",
			Cmd("add",
				Cmd("--intent-to-add"),
				Cmd("--interactive"),
				Cmd("pretty"),
				ListFs()),
			Cmd("clone"),
			Cmd("clean"),
			Cmd("log",
				Cmd("all"),
				Cmd("verbose"))),
		Cmd("exit"),
		Cmd("help"),
	)
	inst.PrintTree(os.Stdout)

	inputLoop(inst)
}
