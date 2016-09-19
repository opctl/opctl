package opspec

//go:generate counterfeiter -o ./fakeGetEventStreamUseCase.go --fake-name fakeGetEventStreamUseCase ./ getEventStreamUseCase

import (
  "fmt"
  "log"
  "github.com/gorilla/websocket"
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/sdk-golang/adapters"
)

type getEventStreamUseCase interface {
  Execute(
  ) (stream chan models.Event, err error)
}

func newGetEventStreamUseCase(
engineHost adapters.EngineHost,
jsonFormat format,
) getEventStreamUseCase {

  return &_getEventStreamUseCase{
    engineHost:engineHost,
    jsonFormat:jsonFormat,
  }

}

type _getEventStreamUseCase struct {
  engineHost adapters.EngineHost
  jsonFormat format
}

func (this _getEventStreamUseCase) Execute(
) (eventStream chan models.Event, err error) {

  eventStream = make(chan models.Event, 1000)

  protocolRelativeBaseUrl, err := this.engineHost.GetEngineProtocolRelativeBaseUrl()
  if (nil != err) {
    return
  }

  c, _, err := websocket.DefaultDialer.Dial(
    fmt.Sprintf("ws:%v/event-stream", protocolRelativeBaseUrl),
    nil,
  )
  if (err != nil) {
    fmt.Println(err)
  }

  go func() {
    defer c.Close()
    for {

      _, bytes, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        return
      }

      var event models.Event
      err = this.jsonFormat.To(bytes, &event)
      if (nil != err) {
        fmt.Printf("json.Unmarshal err: %v \n", err)
      }
      eventStream <- event

    }
  }()

  return

}
