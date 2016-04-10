package core

import (
  "github.com/dev-op-spec/engine/core/ports"
)

type compositionRoot interface {
  RunOpUseCase() runOpUseCase
  AddOpUseCase() addOpUseCase
  AddSubOpUseCase() addSubOpUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  GetLogForOpRunUseCase() getLogForOpRunUseCase
  ListOpsUseCase() listOpsUseCase
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  eventStream := newEventStream()

  // factories
  pathToOpsDirFactory := newPathToOpsDirFactory()
  pathToOpDirFactory := newPathToOpDirFactory(pathToOpsDirFactory)
  pathToOpFileFactory := newPathToOpFileFactory(pathToOpDirFactory)
  uniqueStringFactory := newUniqueStringFactory()

  yamlCodec := newYamlCodec()

  opRunLogFeed := newOpRunLogFeed()

  // use cases
  runOpUseCase := newRunOpUseCase(
    eventStream,
    filesys,
    containerEngine,
    opRunLogFeed,
    uniqueStringFactory,
    yamlCodec,
  )

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

  getEventStreamUseCase := newGetEventStreamUseCase(
    eventStream,
  )

  getLogForOpRunUseCase := newGetLogForOpRunUseCase(
    opRunLogFeed,
  )

  listOpsUseCase := newListOpsUseCase(
    filesys,
    pathToOpFileFactory,
    pathToOpsDirFactory,
    yamlCodec,
  )

  setDescriptionOfOpUseCase := newSetDescriptionOfOpUseCase(
    filesys,
    pathToOpFileFactory,
    yamlCodec,
  )

  compositionRoot = &_compositionRoot{
    runOpUseCase: runOpUseCase,
    addOpUseCase: addOpUseCase,
    addSubOpUseCase: addSubOpUseCase,
    getEventStreamUseCase:getEventStreamUseCase,
    getLogForOpRunUseCase: getLogForOpRunUseCase,
    listOpsUseCase: listOpsUseCase,
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  runOpUseCase           runOpUseCase
  addOpUseCase              addOpUseCase
  addSubOpUseCase           addSubOpUseCase
  getEventStreamUseCase     getEventStreamUseCase
  getLogForOpRunUseCase     getLogForOpRunUseCase
  listOpsUseCase            listOpsUseCase
  setDescriptionOfOpUseCase setDescriptionOfOpUseCase
}

func (this _compositionRoot) RunOpUseCase() runOpUseCase {
  return this.runOpUseCase
}

func (this _compositionRoot) AddOpUseCase() addOpUseCase {
  return this.addOpUseCase
}

func (this _compositionRoot) AddSubOpUseCase() addSubOpUseCase {
  return this.addSubOpUseCase
}

func (this _compositionRoot) GetEventStreamUseCase() getEventStreamUseCase {
  return this.getEventStreamUseCase
}

func (this _compositionRoot) GetLogForOpRunUseCase() getLogForOpRunUseCase {
  return this.getLogForOpRunUseCase
}

func (this _compositionRoot) ListOpsUseCase() listOpsUseCase {
  return this.listOpsUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
