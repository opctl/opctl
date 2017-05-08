package client

import (
	"encoding/json"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/url"
	"strings"
)

func (c client) GetEventStream(
	req *model.GetEventStreamReq,
) (chan model.Event, error) {

	// construct query params
	queryParams := []string{}
	if filter := req.Filter; nil != filter {
		var filterBytes []byte
		filterBytes, err := json.Marshal(filter)
		if nil != err {
			return nil, err
		}
		queryParams = append(
			queryParams,
			fmt.Sprintf("filter=%v", url.QueryEscape(string(filterBytes))),
		)
	}

	reqUrl := c.baseUrl
	reqUrl.Scheme = "ws"
	reqUrl.Path = api.Events_StreamsURLTpl
	reqUrl.RawQuery = strings.Join(queryParams, "&")

	wsConn, _, err := c.wsDialer.Dial(
		reqUrl.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	eventStream := make(chan model.Event, 1000)
	go func() {
		// ensure web socket closed on exit
		defer wsConn.Close()

		// ensure channel closed on exit
		defer close(eventStream)

		for {

			_, bytes, err := wsConn.ReadMessage()
			if nil != err {
				return
			}

			var event model.Event
			err = json.Unmarshal(bytes, &event)
			if nil != err {
				return
			}
			eventStream <- event

		}
	}()

	return eventStream, err

}
