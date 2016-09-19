package opspec

//go:generate counterfeiter -o ./fakeOpViewFactory.go --fake-name fakeOpViewFactory ./ opViewFactory

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type opViewFactory interface {
  Construct(
  opBundlePath string,
  ) (
  opView models.OpView,
  err error,
  )
}

func newOpViewFactory(
filesystem filesystem,
yaml format,
) opViewFactory {

  return &_opViewFactory{
    filesystem:filesystem,
    yaml:yaml,
  }

}

type _opViewFactory struct {
  filesystem filesystem
  yaml       format
}

func (this _opViewFactory) Construct(
opBundlePath string,
) (
opView models.OpView,
err error,
) {

  opManifestPath := path.Join(opBundlePath, NameOfOpManifestFile)

  opManifestBytes, err := this.filesystem.GetBytesOfFile(
    opManifestPath,
  )
  if (nil != err) {
    return
  }

  opManifest := models.OpManifest{}
  err = this.yaml.To(
    opManifestBytes,
    &opManifest,
  )
  if (nil != err) {
    return
  }

  opView = *models.NewOpView(
    opManifest.Description,
    opManifest.Inputs,
    opManifest.Name,
    opManifest.Run,
    opManifest.Version,
  )

  return

}

