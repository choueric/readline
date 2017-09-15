package readline

import (
	"fmt"
	"io"
)

type ExecuteFunc func(line string, data interface{}) bool

type Instance struct {
	line     cmdLine
	view     viewTerm
	input    inputTerm
	root     Completer
	lastChar rune
	execute  ExecuteFunc
	data     interface{}
	Debug    bool
}

func New(prompt string) (*Instance, error) {
	inst := &Instance{}
	err := inst.view.init(prompt)
	if err != nil {
		return nil, err
	}

	err = inst.input.init()
	if err != nil {
		return nil, err
	}

	return inst, nil
}

func (inst *Instance) Destroy() {
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
		inst.Printf("\n++ %s\n", fmt.Sprintf(format, v...))
	}
}

////////////////////////////////////////////////////////////////////////////////

func (inst *Instance) cmdReset() {
	inst.line.reset()
	inst.view.reset()
}

func (inst *Instance) cmdInsert(c rune) {
	inst.line.insert(c)
	inst.view.insert(c)
}

func (inst *Instance) cmdDel() {
	c := inst.line.del()
	inst.view.del(c)
}

func (inst *Instance) cmdBackspace() {
	c := inst.line.backspace()
	inst.view.backspace(c)
}

func (inst *Instance) cmdForwardCursor() {
	c := inst.line.forwardCursor()
	inst.view.forwardCursor(c, inst.line.columnWidth())
}

func (inst *Instance) cmdBackwardCursor() {
	c := inst.line.backwardCursor()
	inst.view.backwardCursor(c)
}

////////////////////////////////////////////////////////////////////////////////

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
