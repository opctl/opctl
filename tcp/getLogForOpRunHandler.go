package tcp

import (
  "github.com/gorilla/websocket"
  "net/http"
  "fmt"
)

func newGetLogForOpRunHandler() http.Handler {

  return &getLogForOpRunHandler{}

}

func print_binary(s []byte) {
  fmt.Printf("Received b:");
  for n := 0; n < len(s); n++ {
    fmt.Printf("%d,", s[n]);
  }
  fmt.Printf("\n");
}

type getLogForOpRunHandler struct{}

func (this getLogForOpRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
      return true
    },
  }

  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    //log.Println(err)
    return
  }

  for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      return
    }

    print_binary(p)

    err = conn.WriteMessage(messageType, p);
    if err != nil {
      return
    }
  }
}