package create

//go:generate counterfeiter -o ./fakeInvoker.go --fake-name FakeInvoker ./ Invoker

import (
	"context"
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/model"
	"github.com/opctl/sdk-golang/node/core"
	"github.com/opctl/sdk-golang/node/core/containerruntime/docker"
	"github.com/opctl/sdk-golang/util/pubsub"
	"os"
	"path"
	"path/filepath"
)

type Invoker interface {
	Invoke(
		opts model.NodeCreateOpts,
	)
}

// NewInvoker returns a new invoker for a node create cmd
func NewInvoker() Invoker {
	return _invoker{}
}

type _invoker struct{}

func (ivkr _invoker) Invoke(
	opts model.NodeCreateOpts,
) {
	dataDirPath, err := dataDirPath(opts)
	if nil != err {
		panic(err)
	}

	lockFile := lockfile.New()
	// ensure we're the only node around these parts
	err = lockFile.Lock(lockFilePath(dataDirPath))
	if nil != err {
		pIDOfExistingNode := lockFile.PIdOfOwner(lockFilePath(dataDirPath))
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIDOfExistingNode))
	}

	containerRuntime, err := docker.New()
	if nil != err {
		panic(err)
	}

	// cleanup [legacy] op cache (if it exists)
	legacyOpCachePath := filepath.Join(dataDirPath, "pkgs")
	if err := os.RemoveAll(legacyOpCachePath); nil != err {
		panic(fmt.Errorf("unable to cleanup op cache at path: %v\n", legacyOpCachePath))
	}

	// cleanup op cache
	opCachePath := filepath.Join(dataDirPath, "ops")
	if err := os.RemoveAll(opCachePath); nil != err {
		panic(fmt.Errorf("unable to cleanup op cache at path: %v\n", opCachePath))
	}

	// cleanup existing DCG (dynamic call graph) data
	dcgDataDirPath := dcgDataDirPath(dataDirPath)
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
			dataDirPath,
		)).
			Listen(ctx)

	select {
	case httpErr := <-httpErrChannel:
		panic(httpErr)
	}

}

// dataDirPath returns path of dir used to store node data
func dataDirPath(
	opts model.NodeCreateOpts,
) (string, error) {
	if nil != opts.DataDir {
		return filepath.Abs(*opts.DataDir)
	}

	if dataDirEnvVar, ok := os.LookupEnv("OPCTL_DATA_DIR"); ok {
		return filepath.Abs(dataDirEnvVar)
	}

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
