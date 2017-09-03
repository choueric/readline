package main

func acAllCmds(inst *Instance) {
	inst.Print("\n")
	for n, _ := range inst.cmds {
		inst.Printf("%s\t", n)
	}
	inst.Print("\n")
}
