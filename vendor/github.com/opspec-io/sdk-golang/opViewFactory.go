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
filesystem Filesystem,
yamlCodec yamlCodec,
) opViewFactory {

  return &_opViewFactory{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _opViewFactory struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
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
  err = this.yamlCodec.FromYaml(
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

