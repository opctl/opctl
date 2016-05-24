package core

//go:generate counterfeiter -o ./fakeAddOpUseCase.go --fake-name fakeAddOpUseCase ./ addOpUseCase

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/ports"
)

type addOpUseCase interface {
  Execute(
  req models.AddOpReq,
  ) (err error)
}

func newAddOpUseCase(
filesys ports.Filesys,
pathToOpDirFactory pathToOpDirFactory,
pathToOpFileFactory pathToOpFileFactory,
yamlCodec yamlCodec,
) addOpUseCase {

  return &_addOpUseCase{
    filesys:filesys,
    pathToOpDirFactory:pathToOpDirFactory,
    pathToOpFileFactory:pathToOpFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addOpUseCase struct {
  filesys             ports.Filesys
  pathToOpDirFactory  pathToOpDirFactory
  pathToOpFileFactory pathToOpFileFactory
  yamlCodec           yamlCodec
}

func (this _addOpUseCase) Execute(
req models.AddOpReq,
) (err error) {

  pathToOpDir := this.pathToOpDirFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.CreateDir(pathToOpDir)
  if (nil != err) {
    return
  }

  var opFile = models.OpFile{
    Description:req.Description,
    Name:req.Name,
  }

  opFileBytes, err := this.yamlCodec.toYaml(&opFile)
  if (nil != err) {
    return
  }

  pathToOpFile := this.pathToOpFileFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.SaveFile(
    pathToOpFile,
    opFileBytes,
  )

  return

}
