package local

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/util/lockfile"
	"github.com/opspec-io/opctl/util/vos"
	"path"
)

func New() nodeprovider.NodeProvider {
	return nodeProvider{
		lockfile: lockfile.New(),
		os:       vos.New(),
	}
}

type nodeProvider struct {
	lockfile lockfile.LockFile
	os       vos.Vos
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
