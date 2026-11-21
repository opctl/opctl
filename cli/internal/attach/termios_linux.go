package attach

import "golang.org/x/sys/unix"

const getTermios = unix.TCGETS
const setTermios = unix.TCSETS

func flushInput(fd int) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIFLUSH)
}
