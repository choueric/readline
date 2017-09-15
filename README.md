# readline

Golang imitation of library readline for command line operations, including
line-editing shortcut keys and auto-completion.

Support multi-byte encoding.

# Usage

# Implementation

Completable elements may include commands, arguments, file names and other
entities, depending on the specific interpreter and its configuration.

# TODO
 
## feature

- [X] add self-defined auto-complete interface.
- [X] add auto-complete interface for list files and directories.
- [X] use completer
- [ ] multi-line input strings, single-line input strings that are long enough
	  to wrap,

- [o] separate it to a library 
- [ ] can not auto-complete multi sub-commands as options, like
	  `git add --a --b a.c`, here, '--a', '--b' and 'a.c(fs-completer)' are
	  sub-commands.
- [ ] Make pretty output formate of candidates
- [ ] Add test case

## bug

- [ ] can not auto-complete 'ls //usr/l'
- [ ] fs-completer: when tab enter and complete to 'ls /usr/ ', there is an 
      extra space after '/', which make user to delete this space and then
	  enter space to auto-complete for the next path node.
