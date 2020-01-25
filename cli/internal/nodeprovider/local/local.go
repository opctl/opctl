package local

import (
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"path"
)

// New returns an initialized "local" node provider
func New() nodeprovider.NodeProvider {
	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		panic(err)
	}

	return nodeProvider{
		lockfile:     lockfile.New(),
		os:           ios.New(),
		lockFilePath: path.Join(perUserAppDataPath, "opctl", "pid.lock"),
	}
}

type nodeProvider struct {
	lockfile     lockfile.LockFile
	os           ios.IOS
	lockFilePath string
}
