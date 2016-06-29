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

  opFilePath := path.Join(opBundlePath, NameOfOpFile)

  opFileBytes, err := this.filesystem.GetBytesOfFile(
    opFilePath,
  )
  if (nil != err) {
    return
  }

  opFile := models.OpFile{}
  err = this.yamlCodec.FromYaml(
    opFileBytes,
    &opFile,
  )
  if (nil != err) {
    return
  }

  params := []models.OpParamView{}
  for paramName, opFileParam := range opFile.Params {
    params = append(
      params,
      *models.NewOpParamView(
        paramName,
        opFileParam.Description,
        opFileParam.IsSecret,
      ),
    )
  }

  subOps := []models.SubOpView{}
  for _, subOp := range opFile.SubOps {
    subOps = append(
      subOps,
      *models.NewSubOpView(
        subOp.IsParallel,
        subOp.Url),
    )
  }

  opView = *models.NewOpView(
    opFile.Description,
    opFile.Name,
    params,
    subOps,
  )

  return

}

