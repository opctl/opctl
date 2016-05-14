package core

//go:generate counterfeiter -o ./fakeCompositionRoot.go --fake-name fakeCompositionRoot ./ compositionRoot

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

type compositionRoot interface {
  RunOpUseCase() runOpUseCase
  AddOpUseCase() addOpUseCase
  AddSubOpUseCase() addSubOpUseCase
  GetEventStreamUseCase() getEventStreamUseCase
  KillOpRunUseCase() killOpRunUseCase
  ListOpsUseCase() listOpsUseCase
  SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase
}

func newCompositionRoot(
containerEngine ports.ContainerEngine,
filesys ports.Filesys,
) (compositionRoot compositionRoot, err error) {

  /* factories */
  pathToOpsDirFactory := newPathToOpsDirFactory()

  pathToOpDirFactory := newPathToOpDirFactory(pathToOpsDirFactory)

  pathToOpFileFactory := newPathToOpFileFactory(pathToOpDirFactory)

  uniqueStringFactory := newUniqueStringFactory()

  /* components */
  eventStream := newEventStream()

  yamlCodec := newYamlCodec()

  logger := func(logEntryEmittedEvent models.LogEntryEmittedEvent) {
    eventStream.Publish(logEntryEmittedEvent)
  }

  opRunner := newOpRunner(
    containerEngine,
    eventStream,
    filesys,
    logger,
    uniqueStringFactory,
    yamlCodec,
  )

  /* use cases */
  runOpUseCase := newRunOpUseCase(
    opRunner,
    uniqueStringFactory,
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

  killOpRunUseCase := newKillOpRunUseCase(
    opRunner,
    uniqueStringFactory,
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
    killOpRunUseCase:killOpRunUseCase,
    listOpsUseCase: listOpsUseCase,
    setDescriptionOfOpUseCase: setDescriptionOfOpUseCase,
  }

  return

}

type _compositionRoot struct {
  runOpUseCase              runOpUseCase
  addOpUseCase              addOpUseCase
  addSubOpUseCase           addSubOpUseCase
  getEventStreamUseCase     getEventStreamUseCase
  killOpRunUseCase          killOpRunUseCase
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

func (this _compositionRoot) KillOpRunUseCase() killOpRunUseCase {
  return this.killOpRunUseCase
}

func (this _compositionRoot) ListOpsUseCase() listOpsUseCase {
  return this.listOpsUseCase
}

func (this _compositionRoot) SetDescriptionOfOpUseCase() setDescriptionOfOpUseCase {
  return this.setDescriptionOfOpUseCase
}
