package tcp

import (
  "github.com/gorilla/websocket"
  "net/http"
  "github.com/opctl/engine/core"
  "encoding/json"
  coreModels "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/tcp/models"
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

  eventChannel := make(chan coreModels.Event)

  go func() {
    for {

      event := <-eventChannel

      eventBytes, err := json.Marshal(models.NewEventMsg(event))
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
