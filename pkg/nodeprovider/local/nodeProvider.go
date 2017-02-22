package local

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/util/vfs/os"
	"github.com/opspec-io/opctl/util/vos"
)

func New() nodeprovider.NodeProvider {
	return nodeProvider{
		nodeRepo: newNodeRepo(
			appdatapath.New(),
			os.New(),
		),
		os: vos.New(),
	}
}

type nodeProvider struct {
	nodeRepo nodeRepo
	os       vos.Vos
}
