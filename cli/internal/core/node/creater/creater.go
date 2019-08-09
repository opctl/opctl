package creater

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeCreater.go --fake-name FakeCreater ./ Creater

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/node/core"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

// Creater exposes the "node create" sub command
type Creater interface {
	Create(
		opts model.NodeCreateOpts,
	)
}

// New returns an initialized "node create" command
func New() Creater {
	return _creater{}
}

type _creater struct{}

func (ivkr _creater) Create(
	opts model.NodeCreateOpts,
) {
	dataDirPath, err := dataDirPath(opts)
	if nil != err {
		panic(err)
	}

	dcgDataDirPath := dcgDataDirPath(dataDirPath)

	ivkr.initDataDir(
		dataDirPath,
		dcgDataDirPath,
	)

	containerRuntime, err := docker.New()
	if nil != err {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpErrChannel :=
		newHTTPListener(core.New(
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

func (ivkr _creater) initDataDir(
	dataDirPath,
	dcgDataDirPath string,
) {
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
	if err := os.RemoveAll(dcgDataDirPath); nil != err {
		panic(fmt.Errorf("unable to cleanup DCG (dynamic call graph) data at path: %v\n", dcgDataDirPath))
	}

	if err := os.MkdirAll(dataDirPath, 0775|os.ModeSetgid); nil != err {
		panic(fmt.Errorf("unable to create OPCTL_DATA_DIR at path: %v; error was %v\n", dcgDataDirPath, err))
	}

	if err := os.Chmod(dataDirPath, 0775|os.ModeSetgid); nil != err {
		panic(fmt.Errorf("unable to setgid of OPCTL_DATA_DIR at path: %v; error was %v\n", dcgDataDirPath, err))
	}

	lockFile := lockfile.New()
	// ensure we're the only node around these parts
	if err := lockFile.Lock(lockFilePath(dataDirPath)); nil != err {
		pIDOfExistingNode := lockFile.PIdOfOwner(lockFilePath(dataDirPath))
		panic(fmt.Errorf("node already running w/ PId: %v\n", pIDOfExistingNode))
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
