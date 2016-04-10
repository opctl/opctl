package tcp

import (
  "github.com/gorilla/websocket"
  "net/http"
  "github.com/dev-op-spec/engine/core"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
)

func newGetEventStreamHandler(
coreApi core.Api,
) http.Handler {

  return &getEventStreamHandler{
    coreApi:coreApi,
  }

}

type getEventStreamHandler struct {
  coreApi core.Api
}

func (this getEventStreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
      return true
    },
  }

  conn, err := upgrader.Upgrade(w, r, nil)
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  eventChannel := make(chan models.Event)

  go func() {
    for {

      event := <-eventChannel

      eventBytes, err := json.Marshal(event)
      if (nil != err) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        conn.Close()
      }

      err = conn.WriteMessage(websocket.TextMessage, eventBytes);
      if (nil != err) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        conn.Close()
      }

    }
  }()

  err = this.coreApi.GetEventStream(
    eventChannel,
  )
  if (nil != err) {
    conn.Close()
  }
}