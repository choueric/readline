# readline

library for command line auto-completion 

# TODO
 
- [X] add self-defined auto-complete interface.
- [X] add auto-complete interface for list files and directories.

- [X] use completer
- [X] fix FsAcFunc
    1. add space at the end. like '/usr '.
	2. list the candidates with whole path.

- [ ] can not auto-complete multi sub-commands as options, like
      `git add --a --b a.c`, here, '--a', '--b' and 'a.c(fs-completer)' are sub-commands.
- [ ] Make pretty output formate of candidates
- [ ] Add test case

- [o] separate it to a library
