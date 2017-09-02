// +build darwin

package main

import (
	"syscall"
	"unsafe"
)

func getTermios(fd int) (*syscall.Termios, error) {
	termios := new(syscall.Termios)
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCGETA, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func setTermios(fd int, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
