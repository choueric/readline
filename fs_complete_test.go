package readline

import "testing"

type TestPair struct {
	input     string
	expect    string
	dirPrefix string
}

func Test_getDir(t *testing.T) {
	testPair := []TestPair{
		{"a", currentDir, ""},
		{"add ", currentDir, ""},
		{"./", currentDir, "."},
		{".", currentDir, ""},
		{"./a", currentDir, "."},
		{"/u", "/", "/"},
		{"/usr", "/", "/"},
		{"/usr/", "/usr", "/usr"},
		{"/usr/local/bin", "/usr/local", "/usr/local"},
	}

	for i, p := range testPair {
		output, dirPrefix := getDir(p.input)
		if output != p.expect {
			t.Errorf("[%d] expect wrong. ([%s] -> [%s]), output: [%s]",
				i, p.input, p.expect, output)
		}
		if dirPrefix != p.dirPrefix {
			t.Errorf("[%d] prefix wrong. ([%s] -> [%s]), output: [%s]",
				i, p.input, p.dirPrefix, dirPrefix)
		}
	}
}
