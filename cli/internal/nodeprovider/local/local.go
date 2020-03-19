package local

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// New returns an initialized "local" node provider
func New() nodeprovider.NodeProvider {
	dataDir, newDataDirErr := datadir.New(nil)
	if nil != newDataDirErr {
		panic(newDataDirErr)
	}

	return nodeProvider{
		dataDir:  dataDir,
		lockfile: lockfile.New(),
		os:       ios.New(),
	}
}

type nodeProvider struct {
	dataDir  datadir.DataDir
	lockfile lockfile.LockFile
	os       ios.IOS
}
