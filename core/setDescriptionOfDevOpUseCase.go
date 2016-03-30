package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type setDescriptionOfDevOpUseCase interface {
  Execute(
  req models.SetDescriptionOfDevOpReq,
  ) (err error)
}

func newSetDescriptionOfDevOpUseCase(
filesys ports.Filesys,
pathToDevOpFileFactory pathToDevOpFileFactory,
yamlCodec yamlCodec,
) setDescriptionOfDevOpUseCase {

  return &_setDescriptionOfDevOpUseCase{
    filesys:filesys,
    pathToDevOpFileFactory:pathToDevOpFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _setDescriptionOfDevOpUseCase struct {
  filesys               ports.Filesys
  pathToDevOpFileFactory pathToDevOpFileFactory
  yamlCodec             yamlCodec
}

func (this _setDescriptionOfDevOpUseCase) Execute(
req models.SetDescriptionOfDevOpReq,
) (err error) {

  pathToDevOpFile := this.pathToDevOpFileFactory.Construct(
    req.ProjectUrl,
    req.DevOpName,
  )

  devOpFileBytes, err := this.filesys.GetBytesOfFile(pathToDevOpFile)
  if (nil != err) {
    return
  }

  devOpFile := devOpFile{}
  err = this.yamlCodec.fromYaml(
    devOpFileBytes,
    &devOpFile,
  )
  if (nil != err) {
    return
  }

  devOpFile.Description = req.Description

  devOpFileBytes, err = this.yamlCodec.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  err = this.filesys.SaveFile(
    pathToDevOpFile,
    devOpFileBytes,
  )

  return

}
