package local

import (
	"bufio"
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vos"
	"path/filepath"
	"strconv"
)

type nodeRepo interface {
	// stores a process
	Add(processId int) (err error)
	DeleteIfExists()
	GetIfExists() (processId int)
}

func newNodeRepo(
	appDataPath appdatapath.AppDataPath,
	fs vfs.Vfs,
) nodeRepo {
	return &_nodeRepo{
		appDataPath: appDataPath,
		fs:          fs,
	}
}

type _nodeRepo struct {
	appDataPath appdatapath.AppDataPath
	fs          vfs.Vfs
	os          vos.Vos
}

func (this *_nodeRepo) Add(processId int) (err error) {
	err = this.fs.MkdirAll(this.nodeDataDirPath(), 0700)
	if nil != err {
		panic(err)
	}

	file, err := this.fs.Create(this.nodePidFilePath())
	if nil != err {
		// expected under race conditions
		return
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(processId))
	if nil != err {
		return err
	}

	return
}

func (this *_nodeRepo) DeleteIfExists() {
	err := this.fs.RemoveAll(this.nodeDataDirPath())
	if nil != err {
		fmt.Printf("Unable to delete local node info; Error was: %v\n", err.Error())
	}
}

func (this *_nodeRepo) GetIfExists() (processId int) {
	// open pid file
	file, err := this.fs.Open(this.nodePidFilePath())
	if nil != err {
		fmt.Printf("Unable to obtain local node info; Error was: %v\n", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		processId, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Unable to obtain local node info; Error was: %v\n", err.Error())
		}
	}
	return
}

func (this *_nodeRepo) nodePidFilePath() string {
	return filepath.Join(this.nodeDataDirPath(), ".pid")
}

func (this *_nodeRepo) nodeDataDirPath() string {
	return filepath.Join(
		appdatapath.New().PerUser(),
		"opctl",
	)
}
