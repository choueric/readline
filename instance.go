package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type ExecuteFunc func(line string, data interface{}) bool

type Instance struct {
	cmdRoot  *Cmd
	line     []byte
	lastKey  byte
	w        *bufio.Writer
	r        *bufio.Reader
	fd       int
	oldState *State
	prompt   string
	execute  ExecuteFunc
	data     interface{}
	debug    bool
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

	return nil
}

func (inst *Instance) Deinit() {
	restoreTerm(inst.fd, inst.oldState)
}

func (inst *Instance) SetCmds(cmds ...*Cmd) {
	inst.cmdRoot = &Cmd{
		name: "",
		subs: cmds,
	}
}

func (inst *Instance) SetExecute(f ExecuteFunc, data interface{}) {
	inst.execute = f
	inst.data = data
}

func (inst *Instance) printPrompt() {
	inst.Print("\n" + inst.prompt)
}

func (inst *Instance) Printf(format string, v ...interface{}) {
	fmt.Fprintf(inst.w, format, v...)
	inst.Flush()
}

func (inst *Instance) Log(format string, v ...interface{}) {
	if inst.debug {
		fmt.Fprintf(inst.w, "\n++ %s", fmt.Sprintf(format, v...))
		inst.Flush()
	}
}

func (inst *Instance) Print(v ...interface{})   { fmt.Fprint(inst.w, v...); inst.Flush() }
func (inst *Instance) Println(v ...interface{}) { fmt.Fprintln(inst.w, v...); inst.Flush() }
func (inst *Instance) Flush()                   { inst.w.Flush() }

func (inst *Instance) clearLine() { inst.line = inst.line[0:0] }

func (inst *Instance) lineAdd(c byte) {
	inst.Printf("%c", c)
	inst.line = append(inst.line, c)
	inst.lastKey = c
}
func (inst *Instance) lineDel() {
	inst.Print("\b \b")
	inst.line = inst.line[0 : len(inst.line)-1]
}

func (inst *Instance) PrintTree(w io.Writer) { inst.cmdRoot.printTree(w) }
