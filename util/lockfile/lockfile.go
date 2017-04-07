package lockfile

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ LockFile

import (
	"bufio"
	"fmt"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/pscanary"
	"github.com/virtual-go/vos"
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
	_fs := osfs.New()
	_os := vos.New(_fs)

	return lockfile{
		fs:       _fs,
		os:       _os,
		psCanary: pscanary.New(_os),
	}
}

type lockfile struct {
	fs       fs.FS
	os       vos.VOS
	psCanary pscanary.PsCanary
}

func (this lockfile) Lock(filepath string) error {
	err := this.fs.MkdirAll(path.Dir(filepath), 0700)
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

func (this lockfile) readLockFile(
	filepath string,
) int {
	// open lockfile
	file, err := this.fs.Open(filepath)
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

func (this lockfile) writeLockFile(filepath string) error {
	// create lockfile
	file, err := this.fs.Create(filepath)
	if nil != err {
		return err
	}
	defer file.Close()

	// write PID
	_, err = file.WriteString(strconv.Itoa(this.os.Getpid()))
	return err
}
