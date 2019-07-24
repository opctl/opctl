package client

import (
	"context"
	"github.com/opctl/opctl/sdks/go/node/api"
	"github.com/opctl/opctl/sdks/go/types"
	"path"
	"strings"
	"time"
)

func (c client) GetEventStream(
	ctx context.Context,
	req *types.GetEventStreamReq,
) (chan types.Event, error) {

	reqURL := c.baseUrl
	reqURL.Scheme = "ws"
	reqURL.Path = path.Join(reqURL.Path, api.URLEvents_Stream)

	queryValues := reqURL.Query()
	if nil != req.Filter.Since {
		queryValues.Add("since", req.Filter.Since.Format(time.RFC3339))
	}
	if nil != req.Filter.Roots {
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

	eventStream := make(chan types.Event, 1000)
	go func() {
		// ensure web socket closed on exit
		defer wsConn.Close()

		// ensure channel closed on exit
		defer close(eventStream)

		for {
			var event types.Event
			err := wsConn.ReadJSON(&event)
			if nil != err {
				return
			}
			eventStream <- event
		}
	}()

	return eventStream, err

}
