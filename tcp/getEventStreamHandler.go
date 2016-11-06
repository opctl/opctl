package tcp

import (
  "github.com/gorilla/websocket"
  "net/http"
  "github.com/opspec-io/engine/core"
  "encoding/json"
  "github.com/opspec-io/sdk-golang/pkg/models"
)

func newGetEventStreamHandler(
coreApi core.Core,
) http.Handler {

  return &getEventStreamHandler{
    coreApi:coreApi,
    upgrader:websocket.Upgrader{
      ReadBufferSize:4096,
      WriteBufferSize:4096,
      CheckOrigin: func(r *http.Request) bool {
        return true
      },
    },
  }

}

type getEventStreamHandler struct {
  coreApi  core.Core
  upgrader websocket.Upgrader
}

func (this getEventStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  conn, err := this.upgrader.Upgrade(w, r, nil)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  defer conn.Close()

  eventChannel := make(chan models.Event)

  err = this.coreApi.GetEventStream(
    eventChannel,
  )
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }

  isWebsocketClosedChannel := make(chan bool, 1)
  defer func() {
    isWebsocketClosedChannel <- true
  }()

  go func() {
    for {
      select {
      case <-isWebsocketClosedChannel:
        return
      case event, isEventChannelOpen := <-eventChannel:
        if (!isEventChannelOpen) {
          // guard event channel closed
          return
        }

        eventBytes, err := json.Marshal(event)
        if (nil != err) {
          http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        err = conn.WriteMessage(websocket.TextMessage, eventBytes);
        if (nil != err) {
          http.Error(w, err.Error(), http.StatusInternalServerError)
        }
      }
    }
  }()

  for {
    _, _, err := conn.ReadMessage()
    if (err != nil) {
      // return on read error
      return
    }
  }

}
