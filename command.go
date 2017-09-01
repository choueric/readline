package main

import (
	"bufio"
	"fmt"
)

type HandleFunc func([]string) error

type Cmd struct {
	name     string
	synopsis string
	handler  HandleFunc
}

type Instance struct {
	cmds []*Cmd
	w    *bufio.Writer
	r    *bufio.Reader
}

func (inst *Instance) AddCmd(cmd *Cmd) error {
	inst.cmds = append(inst.cmds, cmd)

	return nil
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
