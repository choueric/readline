package main

import (
	"io"
	"syscall"
	"unsafe"
)

const (
	CharInterrupt = 3    // Ctrl+C
	CharEOF       = 4    // Ctrl+D
	CharTab       = '\t' // 9
	CharEnter     = '\r' // 13
)

type State struct {
	termios syscall.Termios
}

// These constants are declared here, rather than importing
// them from the syscall package as some syscall packages, even
// on linux, for example gccgo, do not declare them.
const ioctlReadTermios = 0x5401  // syscall.TCGETS
const ioctlWriteTermios = 0x5402 // syscall.TCSETS

func getTermios(fd int) (*syscall.Termios, error) {
	termios := new(syscall.Termios)
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func setTermios(fd int, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	_, err := getTermios(fd)
	return err == nil
}

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) {
	var oldState State

	if termios, err := getTermios(fd); err != nil {
		return nil, err
	} else {
		oldState.termios = *termios
	}

	newState := oldState.termios
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	// newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8

	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	return &oldState, setTermios(fd, &newState)
}

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) {
	termios, err := getTermios(fd)
	if err != nil {
		return nil, err
	}

	return &State{termios: *termios}, nil
}

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func restoreTerm(fd int, state *State) error {
	return setTermios(fd, &state.termios)
}

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) {
	oldState, err := getTermios(fd)
	if err != nil {
		return nil, err
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO
	newState.Lflag |= syscall.ICANON | syscall.ISIG
	newState.Iflag |= syscall.ICRNL
	if err := setTermios(fd, newState); err != nil {
		return nil, err
	}

	defer func() {
		setTermios(fd, oldState)
	}()

	var buf [16]byte
	var ret []byte
	for {
		n, err := syscall.Read(fd, buf[:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			if len(ret) == 0 {
				return nil, io.EOF
			}
			break
		}
		if buf[n-1] == '\n' {
			n--
		}
		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return ret, nil
}
