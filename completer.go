package main

import "strings"

func acAllCmds(inst *Instance) {
	inst.Print("\n")
	for n, _ := range inst.cmds {
		inst.Printf("%s\t", n)
	}
	inst.Print("\n")
}

func getCandidates(inst *Instance, args []string) []string {
	var candidates []string
	for n, _ := range inst.cmds {
		candidates = append(candidates, n)
	}

	// trace the command path to determine the candidates
	for i, c := range args {
		if i != len(args)-1 {
			// TODO: update candidates for this level
		} else {
			var newCandidates []string
			for _, candidate := range candidates {
				if strings.HasPrefix(candidate, c) {
					newCandidates = append(newCandidates, candidate)
				}
				candidates = newCandidates
			}
		}
	}

	return candidates
}

func matchCandidates(candidates []string, arg string) []string {

	return nil
}

func acOneCmd(inst *Instance, args []string) error {
	count := len(args)
	if count == 0 {
		panic("args is empty")
	}

	candidates := getCandidates(inst, args)
	if len(candidates) == 0 {
		panic("candidates is empty")
	}
	if len(candidates) == 1 {
		cmd := candidates[0]
		index := strings.Index(cmd, args[count-1]) + len(args[count-1])
		inst.Printf("%s ", cmd[index:len(cmd)])
		for i := index; i < len(cmd); i++ {
			inst.line = append(inst.line, byte(cmd[i]))
		}
	}

	candidates = matchCandidates(candidates, args[count-1])

	return nil
}
