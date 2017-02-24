package node

import (
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/pkg/node/core"
	"github.com/opspec-io/opctl/pkg/node/tcp"
	"github.com/opspec-io/opctl/util/containerprovider/docker"
	"github.com/opspec-io/opctl/util/lockfile"
	"github.com/opspec-io/opctl/util/vfs/os"
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

	// ensure we've got a clean scratch dir
	dcgDirPath := path.Join(dataDirPath(), "dcgs")
	err = os.New().RemoveAll(dcgDirPath)
	if nil != err {
		panic(fmt.Errorf("unable to cleanup path: %v\n", dcgDirPath))
	}

	err = tcp.New(
		core.New(containerProvider),
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

func lockFilePath() string {
	return path.Join(
		dataDirPath(),
		"pid.lock",
	)
}
