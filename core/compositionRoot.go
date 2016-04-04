package core

import (
  "github.com/dev-op-spec/engine/core/ports"
)

type compositionRoot interface {
  AddOperationUseCase() addOperationUseCase
  AddSubOperationUseCase() addSubOperationUseCase
  ListOperationsUseCase() listOperationsUseCase
  RunOperationUseCase() runOperationUseCase
  SetDescriptionOfOperationUseCase() setDescriptionOfOperationUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  // factories
  pathToOperationsDirFactory := newPathToOperationsDirFactory()
  pathToOperationDirFactory := newPathToOperationDirFactory(pathToOperationsDirFactory)
  pathToOperationFileFactory := newPathToOperationFileFactory(pathToOperationDirFactory)
  uniqueStringFactory := newUniqueStringFactory()

  yamlCodec := newYamlCodec()

  // use cases
  addOperationUseCase := newAddOperationUseCase(
    filesys,
    pathToOperationDirFactory,
    pathToOperationFileFactory,
    yamlCodec,
  )

  addSubOperationUseCase := newAddSubOperationUseCase(
    filesys,
    pathToOperationFileFactory,
    yamlCodec,
  )

  listOperationsUseCase := newListOperationsUseCase(
    filesys,
    pathToOperationFileFactory,
    pathToOperationsDirFactory,
    yamlCodec,
  )

  runOperationUseCase := newRunOperationUseCase(
    filesys,
    containerEngine,
    uniqueStringFactory,
    yamlCodec,
  )

  setDescriptionOfOperationUseCase := newSetDescriptionOfOperationUseCase(
    filesys,
    pathToOperationFileFactory,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    addOperationUseCase: addOperationUseCase,
    addSubOperationUseCase: addSubOperationUseCase,
    listOperationsUseCase: listOperationsUseCase,
    runOperationUseCase: runOperationUseCase,
    setDescriptionOfOperationUseCase: setDescriptionOfOperationUseCase,
  }

  return

}

type _compositionRoot struct {
  addOperationUseCase              addOperationUseCase
  addSubOperationUseCase           addSubOperationUseCase
  listOperationsUseCase            listOperationsUseCase
  runOperationUseCase              runOperationUseCase
  setDescriptionOfOperationUseCase setDescriptionOfOperationUseCase
}

func (this _compositionRoot) AddOperationUseCase() addOperationUseCase {
  return this.addOperationUseCase
}

func (this _compositionRoot) AddSubOperationUseCase() addSubOperationUseCase {
  return this.addSubOperationUseCase
}

func (this _compositionRoot) ListOperationsUseCase() listOperationsUseCase {
  return this.listOperationsUseCase
}

func (this _compositionRoot) RunOperationUseCase() runOperationUseCase {
  return this.runOperationUseCase
}

func (this _compositionRoot) SetDescriptionOfOperationUseCase() setDescriptionOfOperationUseCase {
  return this.setDescriptionOfOperationUseCase
}
