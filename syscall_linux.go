// +build linux

package blugo

import (
	"syscall"
)

func ioctl(fd int, req ...uintptr) error {
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		req[0],
		req[1],
	); errno != 0 {
		return errno
	}
	return nil
}
