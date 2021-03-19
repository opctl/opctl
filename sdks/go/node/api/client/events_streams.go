package client

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

func (c apiClient) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (<-chan model.Event, error) {

	reqURL := c.baseURL
	reqURL.Scheme = "ws"
	reqURL.Path = path.Join(reqURL.Path, api.URLEvents_Stream)

	queryValues := reqURL.Query()
	if req.Filter.Since != nil {
		queryValues.Add("since", req.Filter.Since.Format(time.RFC3339))
	}
	if req.Filter.Roots != nil {
		queryValues.Add("roots", strings.Join(req.Filter.Roots, ","))
	}
	reqURL.RawQuery = queryValues.Encode()

	wsConn, _, err := c.wsDialer.DialContext(
		ctx,
		reqURL.String(),
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
			var event model.Event
			err := wsConn.ReadJSON(&event)
			if err != nil {
				return
			}
			eventStream <- event
		}
	}()

	return eventStream, err

}
