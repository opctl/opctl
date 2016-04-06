package core

import (
  "github.com/dev-op-spec/engine/core/ports"
)

type compositionRoot interface {
  AddOpUseCase() addOpUseCase
  AddSubOpUseCase() addSubOpUseCase
  ListOpsUseCase() listOpsUseCase
  RunOpUseCase() runOpUseCase
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  // factories
  pathToOpsDirFactory := newPathToOpsDirFactory()
  pathToOpDirFactory := newPathToOpDirFactory(pathToOpsDirFactory)
  pathToOpFileFactory := newPathToOpFileFactory(pathToOpDirFactory)
  uniqueStringFactory := newUniqueStringFactory()

  yamlCodec := newYamlCodec()

  // use cases
  addOpUseCase := newAddOpUseCase(
    filesys,
    pathToOpDirFactory,
    pathToOpFileFactory,
    yamlCodec,
  )

  addSubOpUseCase := newAddSubOpUseCase(
    filesys,
    pathToOpFileFactory,
    yamlCodec,
  )

  listOpsUseCase := newListOpsUseCase(
    filesys,
    pathToOpFileFactory,
    pathToOpsDirFactory,
    yamlCodec,
  )

  runOpUseCase := newRunOpUseCase(
    filesys,
    containerEngine,
    uniqueStringFactory,
    yamlCodec,
  )

  setDescriptionOfOpUseCase := newSetDescriptionOfOpUseCase(
    filesys,
    pathToOpFileFactory,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    addOpUseCase: addOpUseCase,
    addSubOpUseCase: addSubOpUseCase,
    listOpsUseCase: listOpsUseCase,
    runOpUseCase: runOpUseCase,
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  addOpUseCase              addOpUseCase
  addSubOpUseCase           addSubOpUseCase
  listOpsUseCase            listOpsUseCase
  runOpUseCase              runOpUseCase
  setDescriptionOfOpUseCase setDescriptionOfOpUseCase
}

func (this _compositionRoot) AddOpUseCase() addOpUseCase {
  return this.addOpUseCase
}

func (this _compositionRoot) AddSubOpUseCase() addSubOpUseCase {
  return this.addSubOpUseCase
}

func (this _compositionRoot) ListOpsUseCase() listOpsUseCase {
  return this.listOpsUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
