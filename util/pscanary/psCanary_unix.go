// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package pscanary

import (
	"os"
	"syscall"
)

// implementation based on github.com/nightlyone/lockfile
func (this psCanary) IsAlive(processId int) bool {
	proc, err := os.FindProcess(processId)
	if nil != err {
		return false
	}

	if err := proc.Signal(syscall.Signal(0)); nil != err {
		return false
	}
	return true
}
