// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package pscanary

import (
	"syscall"
)

// implementation based on github.com/nightlyone/lockfile
func (this _PsCanary) IsAlive(processId int) bool {
	proc, err := this.os.FindProcess(processId)
	if nil != err {
		return false
	}

	if err := proc.Signal(syscall.Signal(0)); nil != err {
		return false
	}
	return true
}
