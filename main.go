package main

import (
	"flag"
	"fmt"
	"os"
)

func executeCmdline(cmdline string, data interface{}) bool {
	fmt.Printf("\nexecute [%v]\n", cmdline)
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
		Item("echo"),
		Item("help"),
	)
	inst.PrintTree(os.Stdout)

	inputLoop(inst)
}
