package node

import (
	"context"
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
	"github.com/opspec-io/sdk-golang/node/core"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime/docker"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"os"
	"path"
	"path/filepath"
)

func New() {
	rootFSPath, err := rootFSPath()
	if nil != err {
		panic(err)
	}

	lockFile := lockfile.New()
	// ensure we're the only node around these parts
	err = lockFile.Lock(lockFilePath(rootFSPath))
	if nil != err {
		pIdOExistingNode := lockFile.PIdOfOwner(lockFilePath(rootFSPath))
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIdOExistingNode))
	}

	containerRuntime, err := docker.New()
	if nil != err {
		panic(err)
	}

	// cleanup [legacy] opspec.engine container if exists; ignore errors
	containerRuntime.DeleteContainerIfExists("opspec.engine")

	// cleanup existing pkg cache
	pkgCachePath := filepath.Join(rootFSPath, "pkgs")
	if err := os.RemoveAll(pkgCachePath); nil != err {
		panic(fmt.Errorf("unable to cleanup pkg cache at path: %v\n", pkgCachePath))
	}

	// cleanup existing DCG (dynamic call graph) data
	dcgDataDirPath := dcgDataDirPath(rootFSPath)
	if err := os.RemoveAll(dcgDataDirPath); nil != err {
		panic(fmt.Errorf("unable to cleanup DCG (dynamic call graph) data at path: %v\n", dcgDataDirPath))
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpErrChannel :=
		newHttpListener(core.New(
			pubsub.New(pubsub.NewBadgerDBEventStore(eventDbPath(dcgDataDirPath))),
			containerRuntime,
			rootFSPath,
		)).
			Listen(ctx)

	libP2PErrChannel :=
		newLibP2PListener().
			Listen(ctx)

	select {
	case httpErr := <-httpErrChannel:
		panic(httpErr)
	case libP2PErr := <-libP2PErrChannel:
		panic(libP2PErr)
	}

}

// fsRootPath returns the root fs path for the node
func rootFSPath() (string, error) {
	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return "", err
	}

	return path.Join(
		perUserAppDataPath,
		"opctl",
	), nil
}

func dcgDataDirPath(rootFSPath string) string {
	return path.Join(
		rootFSPath,
		"dcg",
	)
}

func eventDbPath(dcgDataDirPath string) string {
	return path.Join(
		dcgDataDirPath,
		"event.db",
	)
}

func lockFilePath(rootFSPath string) string {
	return path.Join(
		rootFSPath,
		"pid.lock",
	)
}
