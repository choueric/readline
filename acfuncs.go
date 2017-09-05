package main

import (
	"io/ioutil"
	"path"
	"strings"
)

const currentDir = "."

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

func FsAcFunc(line string) []string {
	var candidates []string
	dir, dirPrefix := getDir(line)
	//fmt.Printf("\n[%s, %s]\n", dir, dirPrefix)
	files, _ := ioutil.ReadDir(dir)
	switch dirPrefix {
	case "":
		for _, f := range files {
			candidates = append(candidates, f.Name())
		}
	case ".":
		for _, f := range files {
			candidates = append(candidates, dirPrefix+"/"+f.Name())
		}
	default:
		for _, f := range files {
			candidates = append(candidates, path.Join(dirPrefix, f.Name()))
		}
	}

	return candidates
}
