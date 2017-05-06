package node

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/node/core"
	"github.com/opctl/opctl/node/tcp"
	"github.com/opctl/opctl/util/containerprovider/docker"
	"github.com/opctl/opctl/util/pubsub"
	"os"
	"path"
)

func New() {
	containerProvider, err := docker.New()
	if nil != err {
		panic(err)
	}

	dataDirPath, err := dataDirPath()
	if nil != err {
		panic(err)
	}

	lockFile := lockfile.New()
	// ensure we're the only node around these parts
	err = lockFile.Lock(lockFilePath(dataDirPath))
	if nil != err {
		pIdOExistingNode := lockFile.PIdOfOwner(lockFilePath(dataDirPath))
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIdOExistingNode))
	}

	// cleanup [legacy] opspec.engine container if exists; ignore errors
	containerProvider.DeleteContainerIfExists("opspec.engine")

	// cleanup existing DCG (dynamic call graph) data
	dcgDataDirPath := dcgDataDirPath(dataDirPath)
	err = os.RemoveAll(dcgDataDirPath)
	if nil != err {
		panic(fmt.Errorf("unable to cleanup DCG (dynamic call graph) data at path: %v\n", dcgDataDirPath))
	}

	err = tcp.New(
		core.New(
			pubsub.New(pubsub.NewEventRepo(eventDbPath(dcgDataDirPath))),
			containerProvider,
		),
	).Start()
	if nil != err {
		panic(err)
	}

}

func dataDirPath() (string, error) {
	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return "", err
	}

	return path.Join(
		perUserAppDataPath,
		"opctl",
	), nil
}

func dcgDataDirPath(dataDirPath string) string {
	return path.Join(
		dataDirPath,
		"dcg",
	)
}

func eventDbPath(dcgDataDirPath string) string {
	return path.Join(
		dcgDataDirPath,
		"event.db",
	)
}

func lockFilePath(dataDirPath string) string {
	return path.Join(
		dataDirPath,
		"pid.lock",
	)
}
