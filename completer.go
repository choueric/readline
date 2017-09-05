package main

import (
	"fmt"
	"io"
)

type Completer interface {
	name() string
	isSp() bool
	subs() []Completer
	getCandidates(string) ([]string, bool)
	modifyCandidate(string, string) string // for output
}

func doPrintTree(c Completer, w io.Writer, depth int, hasSibling []bool) {
	for i := 0; i < depth; i++ {
		if i != depth-1 {
			if hasSibling[i] {
				fmt.Fprintf(w, "│   ")
			} else {
				fmt.Fprintf(w, "    ")
			}
		} else {
			if hasSibling[i] {
				fmt.Fprintln(w, "├── "+c.name())
			} else {
				fmt.Fprintln(w, "└── "+c.name())
			}
		}
	}

	if depth == 0 {
		fmt.Fprintln(w, c.name())
	}

	subs := c.subs()
	length := len(subs)
	for i, sub := range subs {
		if i == length-1 {
			hasSibling = append(hasSibling, false)
		} else {
			hasSibling = append(hasSibling, true)
		}

		doPrintTree(sub, w, depth+1, hasSibling)
		hasSibling = append(hasSibling[:len(hasSibling)-1])
	}
}

// PrintTree prints the tree graphic started from cmd.
func printTree(c Completer, w io.Writer) {
	var hasSibling []bool
	doPrintTree(c, w, 0, hasSibling)
}
