package main

import (
	"fmt"
	"os"
)

func lsHandler(args []string, data interface{}) error {
	fmt.Println("= ls =")

	return nil
}

func echoHandler(args []string, data interface{}) error {
	fmt.Println("= echo =")

	return nil
}

func main() {
	inst := &Instance{
		prompt: "\033[32m>>\033[0m ",
	}
	inst.Init(os.Stdin, os.Stdout)
	defer inst.Deinit()

	inst.AddCmd("ls", &Cmd{"list files and directory", lsHandler, nil})
	inst.AddCmd("echo", &Cmd{"echo what you input", echoHandler, nil})

	inputLoop(inst)
}
