package core

//go:generate counterfeiter -o ./fakeAddSubOpUseCase.go --fake-name fakeAddSubOpUseCase ./ addSubOpUseCase

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/ports"
)

type addSubOpUseCase interface {
  Execute(
  req models.AddSubOpReq,
  ) (err error)
}

func newAddSubOpUseCase(
filesys ports.Filesys,
pathToOpFileFactory pathToOpFileFactory,
yamlCodec yamlCodec,
) addSubOpUseCase {

  return &_addSubOpUseCase{
    filesys:filesys,
    pathToOpFileFactory:pathToOpFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addSubOpUseCase struct {
  filesys             ports.Filesys
  pathToOpFileFactory pathToOpFileFactory
  yamlCodec           yamlCodec
}

func (this _addSubOpUseCase) Execute(
req models.AddSubOpReq,
) (err error) {

  pathToOpFile := this.pathToOpFileFactory.Construct(
    req.ProjectUrl,
    req.OpName,
  )

  opFileBytes, err := this.filesys.GetBytesOfFile(pathToOpFile)
  if (nil != err) {
    return
  }

  opFile := models.OpFile{}

  err = this.yamlCodec.fromYaml(
    opFileBytes,
    &opFile,
  )
  if (nil != err) {
    return
  }

  newOpFileSubOp := models.OpFileSubOp{
    Url:req.SubOpUrl,
  }

  if (len(req.PrecedingSubOpUrl) > 0) {

    subOps := []models.OpFileSubOp{}

    for _, subOp := range opFile.SubOps {

      subOps = append(subOps, subOp)
      if (subOp.Url == req.PrecedingSubOpUrl) {
        subOps = append(subOps, newOpFileSubOp)
      }

    }

    opFile.SubOps = subOps

  } else {

    opFile.SubOps = append(opFile.SubOps, newOpFileSubOp)

  }

  opFileBytes, err = this.yamlCodec.toYaml(&opFile)
  if (nil != err) {
    return
  }

  err = this.filesys.SaveFile(
    pathToOpFile,
    opFileBytes,
  )

  return

}
