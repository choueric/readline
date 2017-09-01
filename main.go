package main

import (
	"bufio"
	"fmt"
	"os"
)

func lsHandler(args []string) error {
	fmt.Println("= lsHandler =")

	return nil
}

func main() {
	stdinFd := int(os.Stdin.Fd())
	oldState, err := MakeRaw(stdinFd)
	if err != nil {
		panic(err)
	}
	defer restoreTerm(stdinFd, oldState)

	inst := &Instance{}
	inst.r = bufio.NewReader(os.Stdin)
	inst.w = bufio.NewWriter(os.Stdout)
	inst.AddCmd(&Cmd{"ls", "list files and directory", lsHandler})

	inputLoop(inst)
}
