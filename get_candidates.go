package readline

import (
	"errors"
	"fmt"
	"strings"
)

// TODO: print with a pretty format
func printCandidates(inst *Instance, cp Completer, candidates []string) {
	args := strings.Fields(string(inst.line))
	prefix := args[len(args)-1]
	inst.Print("\n")
	for _, n := range candidates {
		inst.Printf("%s\t", cp.modifyCandidate(prefix, n))
	}
	inst.Print("\n")
}

// @cp is the parent cmd node and must be cmd-completer.
// find the sub-completer whose name is @cmd
func findSubCompleter(inst *Instance, cmd string, cp Completer) Completer {
	if cp == nil {
		panic("Completer is nil")
	}

	inst.Log("findSubCompleter: find [%s] in %v\n", cmd,
		func(subs []Completer) []string {
			names := make([]string, 0)
			for _, c := range subs {
				names = append(names, c.name())
			}
			return names
		}(cp.subs()))

	if cp.isSp() {
		panic("Sp Completer must not have sub completers")
	}

	subs := cp.subs()
	for _, c := range subs {
		if !c.isSp() && c.name() == cmd {
			return c
		}
	}

	return nil
}

// search @cp's sub-completers and collect candidates
func getCandidatesFromSubs(inst *Instance, cp Completer) ([]string, bool, error) {
	if cp == nil {
		return nil, false, errors.New("Completer is nil")
	}

	var sp Completer
	var candidates []string
	subs := cp.subs()
	for _, c := range subs {
		if c.isSp() {
			sp = c
			break
		} else {
			candidates = append(candidates, c.name())
		}
	}

	if sp != nil {
		cans, end := sp.getCandidates(string(inst.line))
		return cans, end, nil
	}

	return candidates, true, nil
}

// find all commands in @cmds which have prefix of @arg
func getCandidatesByPrefix(inst *Instance, arg string, cp Completer) (Completer, []string, bool, error) {
	var candidates []string
	var sp Completer
	end := true
	currentCp := cp
	cps := cp.subs()
	for _, c := range cps {
		if c.isSp() {
			sp = c
		} else {
			if strings.HasPrefix(c.name(), arg) {
				candidates = append(candidates, c.name())
			}
		}
	}

	if len(candidates) == 0 && sp != nil {
		cans, _end := sp.getCandidates(string(inst.line))
		inst.Log("prefix: %s, end: %v, candidates: [%v]\n", arg, _end, cans)
		for _, v := range cans {
			if strings.HasPrefix(v, arg) {
				candidates = append(candidates, v)
			}
		}
		if len(candidates) == 1 {
			_end = true
		}
		end = _end
		currentCp = sp
	}

	return currentCp, candidates, end, nil
}

// check the inst.line and find the candidates
func getCandidates(inst *Instance) (Completer, []string, bool, error) {
	args := strings.Fields(string(inst.line))
	count := len(inst.line)
	inst.Log("args = %v, len(line) = %d\n", args, count)

	cp := inst.root
	candidates, end, err := getCandidatesFromSubs(inst, cp)
	if err != nil {
		return nil, nil, end, err
	}

	for i, arg := range args {

		lastArg := i == len(args)-1
		partialArg := inst.line[count-1] != ' '
		inst.Log("process cmd [%s], lastArg: %v, partialArg: %v\n",
			arg, lastArg, partialArg)

		if lastArg && partialArg {
			cp, candidates, end, err = getCandidatesByPrefix(inst, arg, cp)
			if err != nil {
				return nil, nil, end, err
			}
		} else {
			cp = findSubCompleter(inst, arg, cp)
			if cp == nil {
				inst.Log("can not find %s", arg)
				return nil, nil, false, errors.New(fmt.Sprintf("can not find %s", arg))
			}
			if lastArg {
				candidates, end, err = getCandidatesFromSubs(inst, cp)
				if err != nil {
					return nil, nil, end, err
				}
			}
		}
	}
	inst.Log("candidates: %v\n", candidates)
	return cp, candidates, end, nil
}

func doComplete(inst *Instance, candidate string) {
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
}

// [ls]: ls -> [ls ]
// return: whether is the auto-complete end
func completeWhole(inst *Instance, candidate string) bool {
	doComplete(inst, candidate)
	return true
}

// if all candidates have the same prefix, complete the common part
// if so, return true
// e.g. [clean], [clone]: c -> [cl]
func completePartial(inst *Instance, candidates []string) {
	prefix := lcp(candidates)
	if len(prefix) != 0 {
		doComplete(inst, prefix)
	}
}

// Via: https://rosettacode.org/wiki/Longest_common_prefix
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
