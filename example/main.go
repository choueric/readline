package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/choueric/readline"
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
		fmt.Printf("execute [%v]\n", cmdline)
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
	inst, err := readline.New("\033[32m>>\033[0m ")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inst.Destroy()

	flag.BoolVar(&inst.Debug, "d", false, "eanble debug")
	flag.Parse()

	inst.SetExecute(executeCmdline, nil)
	inst.SetCompleter(
		readline.Cmd("ls",
			readline.ListFs()),
		readline.Cmd("lsblk"),
		readline.Cmd("git",
			readline.Cmd("add",
				readline.Cmd("--intent-to-add"),
				readline.Cmd("--interactive"),
				readline.Cmd("pretty"),
				readline.ListFs()),
			readline.Cmd("clone"),
			readline.Cmd("clean"),
			readline.Cmd("log",
				readline.Cmd("all"),
				readline.Cmd("verbose"))),
		readline.Cmd("exit"),
		readline.Cmd("help"),
	)
	inst.PrintTree(os.Stdout)

	readline.InputLoop(inst)
}
