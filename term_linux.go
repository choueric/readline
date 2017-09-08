package readline

import (
	"syscall"
	"unsafe"
)

// +build linux

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
