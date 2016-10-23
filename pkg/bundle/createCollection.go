package bundle

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

func (this _bundle) CreateCollection(
req models.CreateCollectionReq,
) (err error) {

  err = this.fileSystem.AddDir(
    req.Path,
  )
  if (nil != err) {
    return
  }

  var opCollection = models.CollectionManifest{
    Manifest:models.Manifest{
      Description:req.Description,
      Name:req.Name,
    },
  }

  opManifestBytes, err := this.yaml.From(&opCollection)
  if (nil != err) {
    return
  }

  err = this.fileSystem.SaveFile(
    path.Join(req.Path, NameOfCollectionManifestFile),
    opManifestBytes,
  )

  return

}
