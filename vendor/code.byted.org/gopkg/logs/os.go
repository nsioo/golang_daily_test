// +build darwin dragonfly freebsd netbsd openbsd linux !windows,!js

package logs

import (
	"syscall"
	"unsafe"
)

func IsTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
