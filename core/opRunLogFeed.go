package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "sync"
  "fmt"
)

func newOpRunLogFeed() opRunLogFeed {

  return &_opRunLogFeed{}

}

type opRunLogFeed interface {
  RegisterSubscriber(
  opRunId string,
  subscriberChannel chan *models.LogEntry,
  )

  RegisterPublisher(
  opRunId string,
  publisherChannel chan *models.LogEntry,
  )
}

// see https://golang.org/src/encoding/json/encode.go#L317 for reference of concurrent map
type _opRunLogFeed struct {
  cachedFeedLookupByOpRunIdRWMutex  sync.RWMutex
  cachedFeedLookupByOpRunId         map[string][]*models.LogEntry

  publisherLookupByOpRunIdRWMutex   sync.RWMutex
  publisherLookupByOpRunId          map[string]chan *models.LogEntry

  subscribersLookupByOpRunIdRWMutex sync.RWMutex
  subscribersLookupByOpRunId        map[string][]chan *models.LogEntry
}

func (this *_opRunLogFeed) RegisterSubscriber(
opRunId string,
subscriberChannel chan *models.LogEntry,
) {

  // return cached
  this.cachedFeedLookupByOpRunIdRWMutex.RLock()
  cachedFeed := this.cachedFeedLookupByOpRunId[opRunId]
  this.cachedFeedLookupByOpRunIdRWMutex.RUnlock()

  for _, logEntry := range cachedFeed {

    subscriberChannel <- logEntry
  }

  this.publisherLookupByOpRunIdRWMutex.RLock()
  existingPublisher := this.publisherLookupByOpRunId[opRunId]
  this.publisherLookupByOpRunIdRWMutex.RUnlock()

  // if no publisher we're done
  if (nil == existingPublisher) {

    close(subscriberChannel)

    return

  }

  this.subscribersLookupByOpRunIdRWMutex.RLock()
  existingSubscribers := this.subscribersLookupByOpRunId[opRunId]
  this.subscribersLookupByOpRunIdRWMutex.RUnlock()

  this.subscribersLookupByOpRunIdRWMutex.Lock()
  if (nil == existingSubscribers) {
    // handle first subscriber
    this.subscribersLookupByOpRunId = map[string][]chan *models.LogEntry{
      opRunId: []chan *models.LogEntry{subscriberChannel},
    }
  }else {
    this.subscribersLookupByOpRunId[opRunId] = append(existingSubscribers, subscriberChannel)
  }
  this.subscribersLookupByOpRunIdRWMutex.Unlock()

}

func (this *_opRunLogFeed) RegisterPublisher(
opRunId string,
publisherChannel chan *models.LogEntry,
) {

  this.cachedFeedLookupByOpRunIdRWMutex.Lock()

  if (nil == this.cachedFeedLookupByOpRunId) {
    // handle cache map init
    this.cachedFeedLookupByOpRunId = map[string][]*models.LogEntry{}
  }

  this.cachedFeedLookupByOpRunIdRWMutex.Unlock()

  this.publisherLookupByOpRunIdRWMutex.Lock()

  if (nil == this.publisherLookupByOpRunId) {

    // handle first publisher
    this.publisherLookupByOpRunId = map[string]chan *models.LogEntry{
      opRunId:publisherChannel,
    }

  }else {

    this.publisherLookupByOpRunId[opRunId] = publisherChannel

  }

  this.publisherLookupByOpRunIdRWMutex.Unlock()

  go func() {

    for {

      logEntry, isOpen := <-publisherChannel
      if (isOpen) {

        // temporary
        fmt.Printf(
          "Timestamp: `%v` | Stream: `%v` | Message: `%v` \n",
          logEntry.Timestamp,
          logEntry.Stream,
          logEntry.Message,
        )

        // cache
        this.cachedFeedLookupByOpRunIdRWMutex.RLock()
        cachedFeed := this.cachedFeedLookupByOpRunId[opRunId]
        this.cachedFeedLookupByOpRunIdRWMutex.RUnlock()

        this.cachedFeedLookupByOpRunIdRWMutex.Lock()
        this.cachedFeedLookupByOpRunId[opRunId] = append(cachedFeed, logEntry)
        this.cachedFeedLookupByOpRunIdRWMutex.Unlock()

        // feed subscribers
        this.subscribersLookupByOpRunIdRWMutex.RLock()
        existingSubscribers := this.subscribersLookupByOpRunId[opRunId]
        this.subscribersLookupByOpRunIdRWMutex.RUnlock()

        for _, subscriber := range existingSubscribers {
          subscriber <- logEntry
        }

      }else {

        // delete publisher map entry
        this.publisherLookupByOpRunIdRWMutex.Lock()
        delete(this.publisherLookupByOpRunId, opRunId)
        this.publisherLookupByOpRunIdRWMutex.Unlock()

        // close subscribers
        this.subscribersLookupByOpRunIdRWMutex.RLock()
        existingSubscribers := this.subscribersLookupByOpRunId[opRunId]
        this.subscribersLookupByOpRunIdRWMutex.RUnlock()

        for _, subscriber := range existingSubscribers {
          close(subscriber)
        }

        // delete subscribers map entry
        this.subscribersLookupByOpRunIdRWMutex.Lock()
        delete(this.subscribersLookupByOpRunId, opRunId)
        this.subscribersLookupByOpRunIdRWMutex.Unlock()

        return

      }

    }

  }()

  return

}
