package core

import (
  "github.com/opctl/engine/core/models"
)

type storage interface {
  addOpRunStartedEvent(opRunStartedEvent models.OpRunStartedEvent)
  listOpRunStartedEventsWithRootId(rootOpId string) []models.OpRunStartedEvent
  isRootOpRunKilled(rootOpRunId string) bool
  deleteOpRunsWithRootId(rootOpRunId string)
}

func newStorage() storage {

  return &_storage{
    opRunStartedEventsByRootId:make(map[string][]models.OpRunStartedEvent),
  }

}

type _storage struct {
  opRunStartedEventsByRootId map[string][]models.OpRunStartedEvent
}

func (this *_storage) addOpRunStartedEvent(opRunStartedEvent models.OpRunStartedEvent) {
  this.opRunStartedEventsByRootId[opRunStartedEvent.RootOpRunId()] = append(
    this.opRunStartedEventsByRootId[opRunStartedEvent.RootOpRunId()],
    opRunStartedEvent,
  )
}

func (this *_storage) listOpRunStartedEventsWithRootId(rootId string) []models.OpRunStartedEvent {
  return this.opRunStartedEventsByRootId[rootId]
}

func (this *_storage) isRootOpRunKilled(rootOpRunId string) bool {
  _, isRunning := this.opRunStartedEventsByRootId[rootOpRunId]
  return !isRunning
}

func (this *_storage) deleteOpRunsWithRootId(rootId string) {
  delete(this.opRunStartedEventsByRootId, rootId)
}
