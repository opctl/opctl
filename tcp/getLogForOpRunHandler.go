package tcp

import (
  "github.com/gorilla/websocket"
  "net/http"
  "github.com/dev-op-spec/engine/core"
  "net/url"
  "github.com/chrisdostert/mux"
  "encoding/json"
  "github.com/dev-op-spec/engine/core/models"
  "time"
)

func newGetLogForOpRunHandler(
coreApi core.Api,
) http.Handler {

  return &getLogForOpRunHandler{
    coreApi:coreApi,
  }

}

type getLogForOpRunHandler struct {
  coreApi core.Api
}

func (this getLogForOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  opRunId, err := url.QueryUnescape(mux.Vars(r)["opRunId"])
  if (nil != err) {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

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

  logChannel := make(chan *models.LogEntry)

  go func() {
    for {

      logEntry, isOpen := <-logChannel
      if (isOpen) {

        logEntryBytes, err := json.Marshal(logEntry)
        if (nil != err) {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          conn.Close()
        }

        err = conn.WriteMessage(websocket.TextMessage, logEntryBytes);
        if (nil != err) {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          conn.Close()
        }

      }else {

        conn.WriteControl(
          websocket.CloseMessage,
          websocket.FormatCloseMessage(websocket.CloseNormalClosure, "success"),
          time.Time{},
        )
        conn.Close()

        return

      }

    }
  }()

  err = this.coreApi.GetLogForOpRun(
    opRunId,
    logChannel,
  )
  if (nil != err) {
    conn.Close()
  }
}