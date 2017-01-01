package engineclient

import (
  "fmt"
  "log"
  "github.com/gorilla/websocket"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "strings"
)

func (this _engineClient) GetEventStream(
req *model.GetEventStreamReq,
) (eventStream chan model.Event, err error) {

  eventStream = make(chan model.Event, 1000)

  protocolRelativeBaseUrl, err := this.engineProvider.GetEngineProtocolRelativeBaseUrl()
  if (nil != err) {
    return
  }

  // construct query params
  queryParams := []string{}
  if filters := req.Filters; nil != req.Filters {
    var filtersBytes []byte
    filtersBytes, err = this.jsonFormat.From(filters)
    if (nil != err) {
      return
    }
    queryParams = append(queryParams, fmt.Sprintf("filters=%v", string(filtersBytes)))
  }

  c, _, err := websocket.DefaultDialer.Dial(
    fmt.Sprintf("ws:%v/event-stream?%v", protocolRelativeBaseUrl, strings.Join(queryParams, "&")),
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
