package main

import (
	"fmt"
	"os"
)

func lsHandler(args []string) error {
	fmt.Println("= lsHandler =")

	return nil
}

func main() {
	inst := &Instance{
		prompt: "\033[32m>>\033[0m ",
	}
	inst.Init(os.Stdin, os.Stdout)
	defer inst.Deinit()

	inst.AddCmd(&Cmd{"ls", "list files and directory", lsHandler})

	inputLoop(inst)
}
