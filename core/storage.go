package core

import (
  "github.com/opspec-io/sdk-golang/pkg/models"
  "sync"
)

type storage interface {
  addOpRunStartedEvent(opRunStartedEvent models.Event)
  listOpRunStartedEventsWithRootId(rootOpId string) []models.Event
  isRootOpRunKilled(rootOpRunId string) bool
  deleteOpRunsWithRootId(rootOpRunId string)
}

func newStorage() storage {

  return &_storage{
    opRunStartedEventsByRootIdMutex:sync.RWMutex{},
    opRunStartedEventsByRootId:make(map[string][]models.Event),
  }

}

type _storage struct {
  opRunStartedEventsByRootIdMutex sync.RWMutex
  opRunStartedEventsByRootId      map[string][]models.Event
}

func (this *_storage) addOpRunStartedEvent(opRunStartedEvent models.Event) {
  this.opRunStartedEventsByRootIdMutex.Lock()
  defer this.opRunStartedEventsByRootIdMutex.Unlock()

  this.opRunStartedEventsByRootId[opRunStartedEvent.OpRunStarted.RootOpRunId] = append(
    this.opRunStartedEventsByRootId[opRunStartedEvent.OpRunStarted.RootOpRunId],
    opRunStartedEvent,
  )
}

func (this *_storage) listOpRunStartedEventsWithRootId(rootId string) []models.Event {
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
