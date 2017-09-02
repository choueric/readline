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
	inst := &Instance{}
	inst.Init(os.Stdin, os.Stdout)
	defer inst.Deinit()

	inst.AddCmd(&Cmd{"ls", "list files and directory", lsHandler})

	inputLoop(inst)
}
