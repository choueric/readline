package readline

import (
	"io/ioutil"
	"path"
	"strings"
)

const currentDir = "."

type fsComplete struct {
}

func ListFs() Completer {
	return &fsComplete{}
}

func (n *fsComplete) name() string      { return "*fsComplete*" }
func (n *fsComplete) isSp() bool        { return true }
func (n *fsComplete) subs() []Completer { return make([]Completer, 0) }

func (n *fsComplete) getCandidates(line string) ([]string, bool) {
	var candidates []string
	end := true
	dir, dirPrefix := getDir(line)
	//fmt.Printf("[%s, %s]\n", dir, dirPrefix)
	files, _ := ioutil.ReadDir(dir)

	join := joinAbsolute
	switch dirPrefix {
	case "":
		join = joinDirect
	case ".":
		join = joinCurrent
	}
	for _, f := range files {
		name := f.Name()
		if f.IsDir() {
			end = false
		}
		candidates = append(candidates, join(dirPrefix, name, f.IsDir()))
	}

	return candidates, end
}

func (n *fsComplete) modifyCandidate(prefix string, input string) string {
	index := strings.LastIndex(prefix, "/")
	return input[index+1 : len(input)]
}

type joinFunc func(string, string, bool) string

func joinDirect(p string, n string, isDir bool) string {
	if isDir {
		return n + "/"
	}
	return n
}

func joinCurrent(p string, n string, isDir bool) string {
	if isDir {
		return p + "/" + n + "/"
	}
	return p + "/" + n
}

func joinAbsolute(p string, n string, isDir bool) string {
	if isDir {
		return path.Join(p, n) + "/"
	}
	return path.Join(p, n)
}

func getDir(line string) (string, string) {
	if line[len(line)-1] == ' ' {
		return currentDir, ""
	}

	args := strings.Fields(line)
	if len(args) == 0 {
		return currentDir, ""
	}
	arg := args[len(args)-1]
	index := strings.LastIndex(arg, "/")
	switch index {
	case -1:
		return currentDir, ""
	case 0:
		return "/", "/"
	default:
		return arg[0:index], arg[0:index]
	}
}
