package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/opspec-io/sdk-golang/model"
	"net/http"
	"strings"
	"time"
)

func (hdlr _handler) events_streams(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	conn, err := hdlr.upgrader.Upgrade(httpResp, httpReq, nil)
	if nil != err {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	req := &model.GetEventStreamReq{}
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

	eventChannel := make(chan *model.Event)

	err = hdlr.core.GetEventStream(
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
