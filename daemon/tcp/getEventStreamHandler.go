package tcp

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/opctl/daemon/core"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"net/http"
)

func newGetEventBusHandler(
	core core.Core,
) http.Handler {

	return &getEventBusHandler{
		core: core,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

}

type getEventBusHandler struct {
	core     core.Core
	upgrader websocket.Upgrader
}

func (this getEventBusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := this.upgrader.Upgrade(w, r, nil)
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	eventChannel := make(chan model.Event)

	// inspired by https://docs.docker.com/engine/reference/api/docker_remote_api_v1.24/#/monitor-dockers-events
	req := &model.GetEventStreamReq{}
	if filterJson := r.URL.Query().Get("filter"); "" != filterJson {
		req.Filter = &model.EventFilter{}
		err = json.Unmarshal([]byte(filterJson), req.Filter)
		if nil != err {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = this.core.GetEventStream(
		req,
		eventChannel,
	)
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
				if !isEventChannelOpen {
					// guard event channel closed
					return
				}

				eventBytes, err := json.Marshal(event)
				if nil != err {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = conn.WriteMessage(websocket.TextMessage, eventBytes)
				if nil != err {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// return on read error
			return
		}
	}

}
