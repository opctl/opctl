package stream

//go:generate counterfeiter -o ./fakeHandler.go --fake-name FakeHandler ./ Handler

import (
	"context"
	"github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/gorilla/websocket"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core"
	"net/http"
	"strings"
	"time"
)

type Handler interface {
	Handle(
		res http.ResponseWriter,
		req *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	core core.Core,
) Handler {
	return _handler{
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

type _handler struct {
	core     core.Core
	upgrader iwebsocket.Upgrader
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	conn, err := hdlr.upgrader.Upgrade(httpResp, httpReq, nil)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	req := &model.GetEventStreamReq{Filter: model.EventFilter{}}
	if sinceString := httpReq.URL.Query().Get("since"); "" != sinceString {
		sinceTime, err := time.Parse(time.RFC3339, sinceString)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusBadRequest)
			return
		}
		req.Filter.Since = &sinceTime
	}

	if rootsString := httpReq.URL.Query().Get("roots"); "" != rootsString {
		rootsArray := strings.Split(rootsString, ",")
		req.Filter.Roots = rootsArray
	}

	// ack is opt in; enables client to apply back pressure to server so it doesn't get flooded
	_, isAckRequested := httpReq.URL.Query()["ack"]

	ctx, cancel := context.WithCancel(httpReq.Context())
	defer cancel()

	// @TODO: handle err channel
	eventChannel, _ := hdlr.core.GetEventStream(
		ctx,
		req,
	)

	for {
		event, isEventChannelOpen := <-eventChannel
		if !isEventChannelOpen {
			// guard event channel closed
			return
		}

		err := conn.WriteJSON(event)
		if nil != err {
			http.Error(httpResp, err.Error(), http.StatusInternalServerError)
			return
		}

		if isAckRequested {
			_, _, err := conn.ReadMessage()
			if nil != err {
				http.Error(httpResp, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

}
