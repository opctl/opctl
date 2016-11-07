package engineclient

import (
  "fmt"
  "log"
  "github.com/gorilla/websocket"
  "github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _engineClient) GetEventStream(
) (eventStream chan model.Event, err error) {

  eventStream = make(chan model.Event, 1000)

  protocolRelativeBaseUrl, err := this.engineProvider.GetEngineProtocolRelativeBaseUrl()
  if (nil != err) {
    return
  }

  c, _, err := websocket.DefaultDialer.Dial(
    fmt.Sprintf("ws:%v/event-stream", protocolRelativeBaseUrl),
    nil,
  )
  if (err != nil) {
    return
  }

  go func() {
    defer c.Close()
    for {

      _, bytes, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        return
      }

      var event model.Event
      err = this.jsonFormat.To(bytes, &event)
      if (nil != err) {
        fmt.Printf("json.Unmarshal err: %v \n", err)
      }
      eventStream <- event

    }
  }()

  return

}
