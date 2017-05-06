package local

import (
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-interfaces/vos"
	"github.com/golang-utils/lockfile"
	"github.com/opctl/opctl/nodeprovider"
	"path"
)

func New() nodeprovider.NodeProvider {

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		panic(err)
	}

	return nodeProvider{
		lockfile:     lockfile.New(),
		os:           vos.New(),
		lockFilePath: path.Join(perUserAppDataPath, "opctl", "pid.lock"),
	}
}

type nodeProvider struct {
	lockfile     lockfile.LockFile
	os           vos.VOS
	lockFilePath string
}
