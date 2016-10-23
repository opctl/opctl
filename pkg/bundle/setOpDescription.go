package bundle

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

func (this _bundle) SetOpDescription(
req models.SetOpDescriptionReq,
) (err error) {

  pathToOpManifest := path.Join(req.PathToOp, NameOfOpManifestFile)

  opBytes, err := this.fileSystem.GetBytesOfFile(
    pathToOpManifest,
  )
  if (nil != err) {
    return
  }

  opManifest := models.OpManifest{}
  err = this.yaml.To(
    opBytes,
    &opManifest,
  )
  if (nil != err) {
    return
  }

  opManifest.Description = req.Description

  opBytes, err = this.yaml.From(&opManifest)
  if (nil != err) {
    return
  }

  err = this.fileSystem.SaveFile(
    pathToOpManifest,
    opBytes,
  )

  return

}
