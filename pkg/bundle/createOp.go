package bundle

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

func (this _bundle) CreateOp(
req models.CreateOpReq,
) (err error) {

  err = this.fileSystem.AddDir(
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

  err = this.fileSystem.SaveFile(
    path.Join(req.Path, NameOfOpManifestFile),
    opManifestBytes,
  )

  return

}
