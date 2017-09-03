package main

import (
	"bufio"
	"fmt"
	"os"
)

type Instance struct {
	cmds     map[string]*Cmd
	line     []byte
	w        *bufio.Writer
	r        *bufio.Reader
	fd       int
	oldState *State
	prompt   string
}

func (inst *Instance) Init(in *os.File, out *os.File) error {
	inst.fd = int(in.Fd())
	oldState, err := MakeRaw(inst.fd)
	if err != nil {
		panic(err)
		return err
	}

	inst.r = bufio.NewReader(in)
	inst.w = bufio.NewWriter(out)
	inst.oldState = oldState

	inst.cmds = make(map[string]*Cmd)
	inst.AddCmd("help", &Cmd{"print this message", helpHandler, inst})
	inst.AddCmd("exit", &Cmd{"exit programe", nil, nil})
	inst.AddCmd("quit", &Cmd{"exit programe", nil, nil})
	return nil
}

func (inst *Instance) Deinit() {
	restoreTerm(inst.fd, inst.oldState)
}

func (inst *Instance) AddCmd(name string, cmd *Cmd) error {
	inst.cmds[name] = cmd
	return nil
}

func (inst *Instance) printPrompt() {
	inst.Print("\n" + inst.prompt)
}

func (inst *Instance) Printf(format string, v ...interface{}) {
	fmt.Fprintf(inst.w, format, v...)
	inst.Flush()
}

func (inst *Instance) Log(format string, v ...interface{}) {
	fmt.Fprintf(inst.w, "\n++ %s", fmt.Sprintf(format, v...))
	inst.Flush()
}

func (inst *Instance) Print(v ...interface{})   { fmt.Fprint(inst.w, v...); inst.Flush() }
func (inst *Instance) Println(v ...interface{}) { fmt.Fprintln(inst.w, v...); inst.Flush() }
func (inst *Instance) Flush()                   { inst.w.Flush() }
