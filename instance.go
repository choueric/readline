package readline

import (
	"fmt"
	"io"
)

type ExecuteFunc func(line string, data interface{}) bool

type Instance struct {
	line    cmdLine
	view    viewTerm
	input   inputTerm
	root    Completer
	lastKey byte
	execute ExecuteFunc
	data    interface{}
	Debug   bool
}

func (inst *Instance) Init(prompt string) error {
	err := inst.view.init(prompt)
	if err != nil {
		return err
	}

	err = inst.input.init()
	if err != nil {
		return err
	}

	return nil
}

func (inst *Instance) Deinit() {
	inst.view.deinit()
	inst.input.deinit()
}

func (inst *Instance) SetCompleter(subs ...Completer) {
	inst.root = Cmd("", subs...)
}

func (inst *Instance) SetExecute(f ExecuteFunc, data interface{}) {
	inst.execute = f
	inst.data = data
}

////////////////////////////////////////////////////////////////////////////////

func (inst *Instance) Printf(format string, v ...interface{}) {
	inst.view.Printf(format, v...)
}

func (inst *Instance) Print(v ...interface{}) {
	inst.view.Print(v...)
}

func (inst *Instance) Println(v ...interface{}) {
	inst.view.Println(v...)
}

func (inst *Instance) Log(format string, v ...interface{}) {
	if inst.Debug {
		inst.view.Printf("\n++ %s", fmt.Sprintf(format, v...))
	}
}

func (inst *Instance) Error(v ...interface{}) {
	inst.view.Print(v...)
}

func (inst *Instance) Errorf(format string, v ...interface{}) {
	inst.view.Printf(format, v...)
}

////////////////////////////////////////////////////////////////////////////////

func (inst *Instance) resetCmdline() {
	inst.line.reset()
	inst.Print("\n")
}

func (inst *Instance) lineAdd(c byte) {
	inst.view.Printf("%c", c)
	inst.line.insert(c)
	inst.lastKey = c
}

func (inst *Instance) PrintTree(w io.Writer) { printTree(inst.root, w) }
