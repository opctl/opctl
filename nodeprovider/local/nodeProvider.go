package local

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/lockfile"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vos"
	"path"
)

func New() nodeprovider.NodeProvider {
	_fs := osfs.New()

	return nodeProvider{
		lockfile: lockfile.New(),
		os:       vos.New(_fs),
	}
}

type nodeProvider struct {
	lockfile lockfile.LockFile
	os       vos.VOS
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
