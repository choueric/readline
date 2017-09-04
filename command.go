package main

import (
	"fmt"
	"io"
)

type Cmd struct {
	name string
	subs []*Cmd
}

func Item(name string, subs ...*Cmd) *Cmd {
	return &Cmd{
		name: name,
		subs: subs,
	}
}

func (cmd *Cmd) doPrintTree(w io.Writer, depth int, hasSibling []bool) {
	for i := 0; i < depth; i++ {
		if i != depth-1 {
			if hasSibling[i] {
				fmt.Fprintf(w, "│   ")
			} else {
				fmt.Fprintf(w, "    ")
			}
		} else {
			if hasSibling[i] {
				fmt.Fprintln(w, "├── "+cmd.name)
			} else {
				fmt.Fprintln(w, "└── "+cmd.name)
			}
		}
	}

	if depth == 0 {
		fmt.Fprintln(w, cmd.name)
	}

	length := len(cmd.subs)
	for i, sub := range cmd.subs {
		if i == length-1 {
			hasSibling = append(hasSibling, false)
		} else {
			hasSibling = append(hasSibling, true)
		}

		sub.doPrintTree(w, depth+1, hasSibling)
		hasSibling = append(hasSibling[:len(hasSibling)-1])
	}
}

// PrintTree prints the tree graphic started from cmd.
func (cmd *Cmd) printTree(w io.Writer) {
	var hasSibling []bool
	cmd.doPrintTree(w, 0, hasSibling)
}
