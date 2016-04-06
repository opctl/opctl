package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type listOpsUseCase interface {
  Execute(
  projectUrl *models.Url,
  ) (ops []models.OpDetailedView, err error)
}

func newListOpsUseCase(
filesys ports.Filesys,
pathToOpFileFactory pathToOpFileFactory,
pathToOpsDirFactory pathToOpsDirFactory,
yamlCodec yamlCodec,
) listOpsUseCase {

  return &_listOpsUseCase{
    filesys:filesys,
    pathToOpFileFactory:pathToOpFileFactory,
    pathToOpsDirFactory:pathToOpsDirFactory,
    yamlCodec:yamlCodec,
  }

}

type _listOpsUseCase struct {
  filesys                    ports.Filesys
  pathToOpFileFactory pathToOpFileFactory
  pathToOpsDirFactory pathToOpsDirFactory
  yamlCodec                  yamlCodec
}

func (this _listOpsUseCase) Execute(
projectUrl *models.Url,
) (ops []models.OpDetailedView, err error) {

  pathToOpsDir := this.pathToOpsDirFactory.Construct(
    projectUrl,
  )

  opDirNames, err := this.filesys.ListNamesOfChildDirs(
    pathToOpsDir,
  )
  if (nil != err) {
    return
  }

  for _, opDirName := range opDirNames {

    pathToOpFile := this.pathToOpFileFactory.Construct(
      projectUrl,
      &opDirName,
    )

    var opFileBytes []byte
    opFileBytes, err = this.filesys.GetBytesOfFile(pathToOpFile)
    if (nil != err) {
      return
    }

    opFile := opFile{}

    err = this.yamlCodec.fromYaml(
      opFileBytes,
      &opFile,
    )
    if (nil != err) {
      return
    }

    opSummaryViews := []models.OpSummaryView{}

    for _, opSubOp := range opFile.SubOps {

      opSummaryView := models.NewOpSummaryView(
        opSubOp.Url,
      )

      opSummaryViews = append(opSummaryViews, *opSummaryView)

    }

    opDetailedView := models.NewOpDetailedView(
      opFile.Description,
      opDirName,
      opSummaryViews,
    )

    ops = append(ops, *opDetailedView)

  }

  return

}
