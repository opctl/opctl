package pkg

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this pkg) CreateCollection(
	req model.CreateCollectionReq,
) (err error) {

	err = this.fileSystem.AddDir(
		req.Path,
	)
	if nil != err {
		return
	}

	var opCollection = model.CollectionManifest{
		Manifest: model.Manifest{
			Description: req.Description,
			Name:        req.Name,
		},
	}

	opManifestBytes, err := this.yaml.From(&opCollection)
	if nil != err {
		return
	}

	err = this.fileSystem.SaveFile(
		path.Join(req.Path, NameOfCollectionManifestFile),
		opManifestBytes,
	)

	return

}
