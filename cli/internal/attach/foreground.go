package attach

import (
	"os"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// IsInteractive reports whether opctl can use the terminal for container input.
func IsInteractive(stdin, stdout *os.File) bool {
	fd := int(stdin.Fd())
	foreground, err := unix.IoctlGetInt(fd, unix.TIOCGPGRP)
	return canForwardInput(term.IsTerminal(fd), term.IsTerminal(int(stdout.Fd())), foreground, unix.Getpgrp(), err)
}

func canForwardInput(stdinTTY, stdoutTTY bool, foregroundGroup, processGroup int, err error) bool {
	return stdinTTY && stdoutTTY && err == nil && foregroundGroup == processGroup
}
