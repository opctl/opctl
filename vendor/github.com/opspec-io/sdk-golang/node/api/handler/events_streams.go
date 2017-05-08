package handler

import (
	"encoding/json"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
)

func newGetEventStreamHandler(
	core core.Core,
) http.Handler {

	return getEventStreamHandler{
		core: core,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

}

type getEventStreamHandler struct {
	core     core.Core
	upgrader iwebsocket.Upgrader
}

func (gesh getEventStreamHandler) ServeHTTP(httpResp http.ResponseWriter, httpReq *http.Request) {
	conn, err := gesh.upgrader.Upgrade(httpResp, httpReq, nil)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	// inspired by https://docs.docker.com/engine/reference/api/docker_remote_handler_v1.24/#/monitor-dockers-events
	req := &model.GetEventStreamReq{}
	if filterJson := httpReq.URL.Query().Get("filter"); "" != filterJson {
		req.Filter = &model.EventFilter{}
		err = json.Unmarshal([]byte(filterJson), req.Filter)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusBadRequest)
			return
		}
	}

	eventChannel := make(chan *model.Event)

	err = gesh.core.GetEventStream(
		req,
		eventChannel,
	)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	for {
		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			// guard event channel closed
			return
		}

		eventBytes, err := json.Marshal(event)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, eventBytes)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
