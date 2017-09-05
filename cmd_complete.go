package main

type cmdComplete struct {
	_name string
	_subs []Completer
}

func Cmd(name string, subs ...Completer) Completer {
	return &cmdComplete{
		_name: name,
		_subs: subs,
	}
}

func (n *cmdComplete) name() string      { return n._name }
func (n *cmdComplete) isSp() bool        { return false }
func (n *cmdComplete) subs() []Completer { return n._subs }

func (n *cmdComplete) getCandidates(string) ([]string, bool) {
	return []string{n._name}, true
}

func (n *cmdComplete) modifyCandidate(prefix string, input string) string { return input }
