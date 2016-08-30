package core

import (
  "github.com/opspec-io/engine/core/models"
  "sync"
)

type storage interface {
  addOpRunStartedEvent(opRunStartedEvent models.OpRunStartedEvent)
  listOpRunStartedEventsWithRootId(rootOpId string) []models.OpRunStartedEvent
  isRootOpRunKilled(rootOpRunId string) bool
  deleteOpRunsWithRootId(rootOpRunId string)
}

func newStorage() storage {

  return &_storage{
    opRunStartedEventsByRootIdMutex:sync.RWMutex{},
    opRunStartedEventsByRootId:make(map[string][]models.OpRunStartedEvent),
  }

}

type _storage struct {
  opRunStartedEventsByRootIdMutex sync.RWMutex
  opRunStartedEventsByRootId      map[string][]models.OpRunStartedEvent
}

func (this *_storage) addOpRunStartedEvent(opRunStartedEvent models.OpRunStartedEvent) {
  this.opRunStartedEventsByRootIdMutex.Lock()
  defer this.opRunStartedEventsByRootIdMutex.Unlock()

  this.opRunStartedEventsByRootId[opRunStartedEvent.RootOpRunId()] = append(
    this.opRunStartedEventsByRootId[opRunStartedEvent.RootOpRunId()],
    opRunStartedEvent,
  )
}

func (this *_storage) listOpRunStartedEventsWithRootId(rootId string) []models.OpRunStartedEvent {
  this.opRunStartedEventsByRootIdMutex.RLock()
  defer this.opRunStartedEventsByRootIdMutex.RUnlock()

  return this.opRunStartedEventsByRootId[rootId]
}

func (this *_storage) isRootOpRunKilled(rootOpRunId string) bool {
  this.opRunStartedEventsByRootIdMutex.RLock()
  defer this.opRunStartedEventsByRootIdMutex.RUnlock()

  _, isRunning := this.opRunStartedEventsByRootId[rootOpRunId]
  return !isRunning
}

func (this *_storage) deleteOpRunsWithRootId(rootId string) {
  this.opRunStartedEventsByRootIdMutex.Lock()
  defer this.opRunStartedEventsByRootIdMutex.Unlock()

  delete(this.opRunStartedEventsByRootId, rootId)
}
