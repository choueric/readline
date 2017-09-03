package main

type CmdHandler func([]string, interface{}) error

type Cmd struct {
	synopsis string
	handler  CmdHandler
	data     interface{}
}

// TODO: sort the commands
func helpHandler(args []string, data interface{}) error {
	inst := data.(*Instance)
	inst.Println("Help:")
	for n, c := range inst.cmds {
		inst.Printf("  %s: %s\n", n, c.synopsis)
	}

	return nil
}
