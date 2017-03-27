package node

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opctl/opctl/node/core"
	"github.com/opctl/opctl/node/tcp"
	"github.com/opctl/opctl/util/containerprovider/docker"
	"github.com/opctl/opctl/util/lockfile"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/vfs/os"
	"path"
)

func New() {
	containerProvider, err := docker.New()
	if nil != err {
		panic(err)
	}

	lockFile := lockfile.New()

	// ensure we're the only node around these parts
	err = lockFile.Lock(lockFilePath())
	if nil != err {
		pIdOExistingNode := lockFile.PIdOfOwner(lockFilePath())
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIdOExistingNode))
	}

	// cleanup [legacy] opspec.engine container if exists; ignore errors
	containerProvider.DeleteContainerIfExists("opspec.engine")

	// cleanup existing DCG (dynamic call graph) data
	err = os.New().RemoveAll(dcgDataDirPath())
	if nil != err {
		panic(fmt.Errorf("unable to cleanup DCG (dynamic call graph) data at path: %v\n", dcgDataDirPath()))
	}

	err = tcp.New(
		core.New(
			pubsub.New(pubsub.NewEventRepo(eventDbPath())),
			containerProvider,
		),
	).Start()
	if nil != err {
		panic(err)
	}

}

func dataDirPath() string {
	return path.Join(
		appdatapath.New().PerUser(),
		"opctl",
	)
}

func dcgDataDirPath() string {
	return path.Join(
		dataDirPath(),
		"dcg",
	)
}

func eventDbPath() string {
	return path.Join(
		dcgDataDirPath(),
		"event.db",
	)
}

func lockFilePath() string {
	return path.Join(
		dataDirPath(),
		"pid.lock",
	)
}
