package lockfile

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ LockFile

import (
	"bufio"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/pscanary"
	"path"
	"strconv"
)

type LockFile interface {
	// obtains a lock
	Lock(filepath string) error

	// obtains the PId of the current lock owner
	PIdOfOwner(filepath string) int
}

func New() LockFile {

	return lockfile{
		os:       ios.New(),
		psCanary: pscanary.New(),
	}
}

type lockfile struct {
	os       ios.IOS
	psCanary pscanary.PsCanary
}

func (this lockfile) Lock(filepath string) error {
	err := this.os.MkdirAll(path.Dir(filepath), 0700)
	if nil != err {
		return err
	}

	pIdOfOwner := this.PIdOfOwner(filepath)
	if 0 != pIdOfOwner {
		// if an owner exists we've been preempted
		return fmt.Errorf("Unable to obtain lock; currently owned by PId: %v\n", pIdOfOwner)
	}

	return this.writeLockFile(filepath)
}

// 0 means no owner
func (this lockfile) PIdOfOwner(filepath string) int {
	if pIdFromFile := this.readLockFile(filepath); this.psCanary.IsAlive(pIdFromFile) {
		return pIdFromFile
	}
	return 0
}

func (lf lockfile) readLockFile(
	filepath string,
) int {
	// open lockfile
	file, err := lf.os.Open(filepath)
	if nil != err {
		return 0
	}
	defer file.Close()

	// read PID
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if pIdFromFile, err := strconv.Atoi(scanner.Text()); nil == err {
			return pIdFromFile
		}
		break
	}
	return 0
}

func (lf lockfile) writeLockFile(filepath string) error {
	// create lockfile
	file, err := lf.os.Create(filepath)
	if nil != err {
		return err
	}
	defer file.Close()

	// write PID
	_, err = file.WriteString(strconv.Itoa(lf.os.Getpid()))
	return err
}
