package pkg

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this pkg) CreateOp(
	req model.CreateOpReq,
) (err error) {

	err = this.fileSystem.AddDir(
		req.Path,
	)
	if nil != err {
		return
	}

	var opManifest = model.OpManifest{
		Manifest: model.Manifest{
			Description: req.Description,
			Name:        req.Name,
		},
	}

	opManifestBytes, err := this.yaml.From(&opManifest)
	if nil != err {
		return
	}

	err = this.fileSystem.SaveFile(
		path.Join(req.Path, NameOfOpManifestFile),
		opManifestBytes,
	)

	return

}
