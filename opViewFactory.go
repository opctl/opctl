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

  opBundleManifestPath := path.Join(opBundlePath, NameOfOpBundleManifest)

  opBundleManifestBytes, err := this.filesystem.GetBytesOfFile(
    opBundleManifestPath,
  )
  if (nil != err) {
    return
  }

  opBundleManifest := models.OpBundleManifest{}
  err = this.yamlCodec.FromYaml(
    opBundleManifestBytes,
    &opBundleManifest,
  )
  if (nil != err) {
    return
  }

  opView = *models.NewOpView(
    opBundleManifest.Description,
    opBundleManifest.Inputs,
    opBundleManifest.Name,
    opBundleManifest.Run,
    opBundleManifest.Version,
  )

  return

}

