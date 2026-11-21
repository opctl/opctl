package attach

import "golang.org/x/sys/unix"

const getTermios = unix.TIOCGETA
const setTermios = unix.TIOCSETA

func flushInput(fd int) error {
	const fread = 1
	return unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, fread)
}
