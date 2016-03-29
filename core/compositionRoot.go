package core

import (
  "github.com/dev-op-spec/engine/core/ports"
)

type compositionRoot interface {
  AddDevOpUseCase() addDevOpUseCase
  AddPipelineUseCase() addPipelineUseCase
  AddStageToPipelineUseCase() addStageToPipelineUseCase
  ListDevOpsUseCase() listDevOpsUseCase
  ListPipelinesUseCase() listPipelinesUseCase
  RunDevOpUseCase() runDevOpUseCase
  RunPipelineUseCase() runPipelineUseCase
  SetDescriptionOfDevOpUseCase() setDescriptionOfDevOpUseCase
  SetDescriptionOfPipelineUseCase() setDescriptionOfPipelineUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  // factories
  pathToDevOpsDirFactory := newPathToDevOpsDirFactory()
  pathToDevOpDirFactory := newPathToDevOpDirFactory(pathToDevOpsDirFactory)
  pathToDevOpFileFactory := newPathToDevOpFileFactory(pathToDevOpDirFactory)
  pathToPipelinesDirFactory := newPathToPipelinesDirFactory()
  pathToPipelineDirFactory := newPathToPipelineDirFactory(pathToPipelinesDirFactory)
  pathToPipelineFileFactory := newPathToPipelineFileFactory(pathToPipelineDirFactory)

  runDevOpUseCase := newRunDevOpUseCase(containerEngine, pathToDevOpDirFactory)

  yamlCodec := newYamlCodec()

  // use cases
  addDevOpUseCase := newAddDevOpUseCase(
    filesys,
    pathToDevOpDirFactory,
    pathToDevOpFileFactory,
    yamlCodec,
    containerEngine,
  )

  addPipelineUseCase := newAddPipelineUseCase(
    filesys,
    pathToPipelineDirFactory,
    pathToPipelineFileFactory,
    yamlCodec,
  )

  addStageToPipelineUseCase := newAddStageToPipelineUseCase(
    filesys,
    pathToPipelineFileFactory,
    yamlCodec,
  )

  listDevOpsUseCase := newListDevOpsUseCase(
    filesys,
    pathToDevOpFileFactory,
    pathToDevOpsDirFactory,
    yamlCodec,
  )

  listPipelinesUseCase := newListPipelinesUseCase(
    filesys,
    pathToPipelineFileFactory,
    pathToPipelinesDirFactory,
    yamlCodec,
  )

  runPipelineUseCase := newRunPipelineUseCase(
    filesys,
    pathToPipelineDirFactory,
    pathToPipelineFileFactory,
    yamlCodec,
    runDevOpUseCase,
  )

  setDescriptionOfDevOpUseCase := newSetDescriptionOfDevOpUseCase(
    filesys,
    pathToDevOpFileFactory,
    yamlCodec,
  )

  setDescriptionOfPipelineUseCase :=
  newSetDescriptionOfPipelineUseCase(
    filesys,
    pathToPipelineFileFactory,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    addDevOpUseCase: addDevOpUseCase,
    addPipelineUseCase: addPipelineUseCase,
    addStageToPipelineUseCase: addStageToPipelineUseCase,
    listDevOpsUseCase: listDevOpsUseCase,
    listPipelinesUseCase:listPipelinesUseCase,
    runDevOpUseCase: runDevOpUseCase,
    runPipelineUseCase: runPipelineUseCase,
    setDescriptionOfDevOpUseCase: setDescriptionOfDevOpUseCase,
    setDescriptionOfPipelineUseCase:setDescriptionOfPipelineUseCase,
  }

  return

}

type _compositionRoot struct {
  addDevOpUseCase                 addDevOpUseCase
  addPipelineUseCase              addPipelineUseCase
  addStageToPipelineUseCase       addStageToPipelineUseCase
  listDevOpsUseCase               listDevOpsUseCase
  listPipelinesUseCase            listPipelinesUseCase
  runDevOpUseCase                 runDevOpUseCase
  runPipelineUseCase              runPipelineUseCase
  setDescriptionOfDevOpUseCase    setDescriptionOfDevOpUseCase
  setDescriptionOfPipelineUseCase setDescriptionOfPipelineUseCase
}

func (this _compositionRoot) AddDevOpUseCase() addDevOpUseCase {
  return this.addDevOpUseCase
}

func (this _compositionRoot) AddPipelineUseCase() addPipelineUseCase {
  return this.addPipelineUseCase
}

func (this _compositionRoot) AddStageToPipelineUseCase() addStageToPipelineUseCase {
  return this.addStageToPipelineUseCase
}

func (this _compositionRoot) ListDevOpsUseCase() listDevOpsUseCase {
  return this.listDevOpsUseCase
}

func (this _compositionRoot) ListPipelinesUseCase() listPipelinesUseCase {
  return this.listPipelinesUseCase
}

func (this _compositionRoot) RunDevOpUseCase() runDevOpUseCase {
  return this.runDevOpUseCase
}

func (this _compositionRoot) RunPipelineUseCase() runPipelineUseCase {
  return this.runPipelineUseCase
}

func (this _compositionRoot) SetDescriptionOfDevOpUseCase() setDescriptionOfDevOpUseCase {
  return this.setDescriptionOfDevOpUseCase
}

func (this _compositionRoot) SetDescriptionOfPipelineUseCase() setDescriptionOfPipelineUseCase {
  return this.setDescriptionOfPipelineUseCase
}
