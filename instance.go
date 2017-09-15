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

func (inst *Instance) Log(format string, v ...interface{}) {
	if inst.Debug {
		inst.Printf("\n++ %s", fmt.Sprintf(format, v...))
	}
}

func (inst *Instance) resetCmdline() {
	inst.line.reset()
	inst.Print("\n")
}

func (inst *Instance) PrintTree(w io.Writer) { printTree(inst.root, w) }

func (inst *Instance) Printf(format string, v ...interface{}) {
	fmt.Fprintf(inst.view.w, format, v...)
	inst.view.w.Flush()
}

func (inst *Instance) Print(v ...interface{}) {
	fmt.Fprint(inst.view.w, v...)
	inst.view.w.Flush()
}

func (inst *Instance) Println(v ...interface{}) {
	fmt.Fprintln(inst.view.w, v...)
	inst.view.w.Flush()
}

func (inst *Instance) Error(v ...interface{}) {
	fmt.Fprint(inst.view.w, v...)
	inst.view.w.Flush()
}

func (inst *Instance) Errorf(format string, v ...interface{}) {
	fmt.Fprintf(inst.view.w, format, v...)
	inst.view.w.Flush()
}
