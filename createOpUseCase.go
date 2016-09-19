package opspec

//go:generate counterfeiter -o ./fakeCreateOpUseCase.go --fake-name fakeCreateOpUseCase ./ createOpUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type createOpUseCase interface {
  Execute(
  req models.CreateOpReq,
  ) (err error)
}

func newCreateOpUseCase(
filesystem filesystem,
yaml format,
) createOpUseCase {

  return &_createOpUseCase{
    filesystem:filesystem,
    yaml:yaml,
  }

}

type _createOpUseCase struct {
  filesystem filesystem
  yaml       format
}

func (this _createOpUseCase) Execute(
req models.CreateOpReq,
) (err error) {

  err = this.filesystem.AddDir(
    req.Path,
  )
  if (nil != err) {
    return
  }

  var opManifest = models.OpManifest{
    Manifest:models.Manifest{
      Description:req.Description,
      Name:req.Name,
    },
  }

  opManifestBytes, err := this.yaml.From(&opManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    path.Join(req.Path, "op.yml"),
    opManifestBytes,
  )

  return

}
