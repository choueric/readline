package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	PROMPT = ">> "
)

func printPrompt(w *bufio.Writer) {
	fmt.Fprintf(w, PROMPT)
	w.Flush()
}

func executeCmdline(w *bufio.Writer, line []byte) {
	fmt.Fprintf(w, "\n[%s]\n", string(line))
}

func handleTab(w *bufio.Writer, line []byte) {
	fmt.Fprintf(w, "\n autocomplete \n")
	fmt.Fprintf(w, "%s%s", PROMPT, string(line))
	w.Flush()
}

func main() {
	stdinFd := int(os.Stdin.Fd())
	oldState, err := MakeRaw(stdinFd)
	if err != nil {
		panic(err)
	}
	defer restoreTerm(stdinFd, oldState)

	end := false
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	line := make([]byte, 1024)

	printPrompt(w)
	for !end {
		c, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				end = true
				fmt.Fprintf(w, "got EOF\n")
				w.Flush()
				continue
			}
			fmt.Fprintf(w, "error: %v", err)
			break
		}

		//fmt.Fprintf(w, "[%d]", c)
		switch c {
		case CharInterrupt:
			fmt.Fprintf(w, "\ngot Interrupt(Ctrl+C)\n")
			line = line[0:0]
			fmt.Fprintf(w, "\n")
			printPrompt(w)
		case CharEOF:
			if len(line) == 0 {
				end = true
				fmt.Fprintf(w, "\ngot EOF(Ctrl+D)\n")
				w.Flush()
			} else {
				line = line[0:0]
				fmt.Fprintf(w, "\n")
				printPrompt(w)
			}
		case CharEnter:
			executeCmdline(w, line)
			line = line[0:0]
			printPrompt(w)
		case CharTab:
			handleTab(w, line)
		default:
			line = append(line, c)
			fmt.Fprintf(w, "%c", c)
			w.Flush()
		}
	}
}
