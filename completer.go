package main

import (
	"errors"
	"fmt"
	"strings"
)

func printCandidates(inst *Instance, candidates []string) {
	inst.Print("\n")
	for _, n := range candidates {
		inst.Printf("%s\t", n)
	}
	inst.Print("\n")
}

// @cmd is the parent cmd node, recurse its sub commands and find out if
// there is a sub command matches to @arg.
func findCmd(arg string, cmd *Cmd) *Cmd {
	for _, c := range cmd.subs {
		if c.name == arg {
			return c
		}
	}

	return nil
}

// get all @cmd's sub commands and treat as candidates
// TODO: add error return
func getCmdSubs(cmd *Cmd) []string {
	if cmd == nil {
		return nil
	}
	var candidates []string
	for _, c := range cmd.subs {
		candidates = append(candidates, c.name)
	}

	return candidates
}

// find all sub commands of @cmd which have prefix of @arg
func matchCandidates(arg string, cmd *Cmd) []string {
	if cmd == nil {
		return nil
	}

	var candidates []string
	for _, c := range cmd.subs {
		if strings.HasPrefix(c.name, arg) {
			candidates = append(candidates, c.name)
		}
	}

	return candidates
}

// check the inst.line and find the candidates
func getCandidates(inst *Instance) ([]string, error) {
	args := strings.Fields(string(inst.line))
	count := len(inst.line)
	inst.Log("args = %v, len(line) = %d\n", args, count)

	cmdNode := inst.cmdRoot
	candidates := getCmdSubs(cmdNode)

	for i, arg := range args {
		inst.Log("process cmd [%s]\n", arg)

		lastArg := i == len(args)-1
		partialArg := inst.line[count-1] != ' '

		if lastArg && partialArg {
			candidates = matchCandidates(arg, cmdNode)
		} else {
			cmdNode = findCmd(arg, cmdNode)
			if cmdNode == nil {
				return candidates, errors.New(fmt.Sprintf("can not find %s", arg))
			}
			if lastArg {
				candidates = getCmdSubs(cmdNode)
			}
		}
	}
	inst.Log("candidates: %v\n", candidates)
	return candidates, nil
}

func doComplete(inst *Instance, candidate string, space bool) {
	args := strings.Fields(string(inst.line))
	count := len(args)
	if count == 0 {
		panic("when do complete, the input can not be empty fields")
	}
	todo := args[len(args)-1]
	index := strings.Index(candidate, todo) + len(todo)
	for i := index; i < len(candidate); i++ {
		inst.lineAdd(byte(candidate[i]))
	}
	if space {
		inst.lineAdd(' ')
	}
}

// [ls]: ls -> [ls ]
func completeWhole(inst *Instance, candidate string) {
	doComplete(inst, candidate, true)
}

// if all candidates have the same prefix, complete the common part
// if so, return true
func completePartial(inst *Instance, candidates []string) bool {
	prefix := lcp(candidates)
	if len(prefix) != 0 {
		doComplete(inst, prefix, false)
	}
	return false
}

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
func lcp(l []string) string {
	// Special cases first
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}
