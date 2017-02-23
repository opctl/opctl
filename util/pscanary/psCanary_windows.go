package pscanary

import (
	"syscall"
)

// implementation based on github.com/nightlyone/lockfile
const (
	error_invalid_parameter = 87
	code_still_active       = 259
)

func (this psCanary) IsAlive(processId int) bool {
	procHnd, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, true, uint32(processId))
	if err != nil {
		if scerr, ok := err.(syscall.Errno); ok {
			if uintptr(scerr) == error_invalid_parameter {
				return false
			}
		}
	}

	var code uint32
	err = syscall.GetExitCodeProcess(procHnd, &code)
	if err != nil {
		return false
	}

	return code == code_still_active
}
