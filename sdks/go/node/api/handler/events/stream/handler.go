package stream

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	iwebsocket "github.com/golang-interfaces/github.com-gorilla-websocket"
	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

//counterfeiter:generate -o fakes/handler.go . Handler
type Handler interface {
	Handle(
		res http.ResponseWriter,
		req *http.Request,
	)
}

// NewHandler returns an initialized Handler instance
func NewHandler(
	node node.Node,
) Handler {
	return _handler{
		node: node,
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
	node     node.Node
	upgrader iwebsocket.Upgrader
}

func (hdlr _handler) Handle(
	httpResp http.ResponseWriter,
	httpReq *http.Request,
) {
	query := httpReq.URL.Query()
	req := &model.GetEventStreamReq{}
	if query.Has("ordered") {
		req.Filter.Ordered = true
		if value := query.Get("ordered"); value != "" {
			ordered, err := strconv.ParseBool(value)
			if err != nil {
				http.Error(httpResp, "ordered must be a boolean", http.StatusBadRequest)
				return
			}
			req.Filter.Ordered = ordered
		}
	}

	conn, err := hdlr.upgrader.Upgrade(httpResp, httpReq, nil)
	if err != nil {
		http.Error(httpResp, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	if sinceString := query.Get("since"); sinceString != "" {
		sinceTime, err := time.Parse(time.RFC3339, sinceString)
		if err != nil {
			// return on any error
			return
		}
		req.Filter.Since = &sinceTime
	}

	if rootsString := query.Get("roots"); rootsString != "" {
		rootsArray := strings.Split(rootsString, ",")
		req.Filter.Roots = rootsArray
	}

	// ack is opt in; enables client to apply back pressure to server so it doesn't get flooded
	isAckRequested := query.Has("ack")

	ctx, cancel := context.WithCancel(httpReq.Context())
	defer cancel()
	var orderedAcks chan struct{}
	if req.Filter.Ordered {
		conn.SetReadLimit(512)
		if isAckRequested {
			orderedAcks = make(chan struct{}, 1)
		}
		// Keep reading while the stream is idle so disconnects cancel the subscription.
		go func() {
			defer conn.Close()
			defer cancel()
			if !isAckRequested {
				if _, _, err := conn.NextReader(); err == nil {
					conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(
						websocket.CloseUnsupportedData, "ordered event streams do not accept messages",
					), time.Now().Add(time.Second))
				}
				return
			}
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
				select {
				case orderedAcks <- struct{}{}:
				default:
					// A client can have at most one ACK pending.
					return
				}
			}
		}()
	}

	// @TODO: handle err channel
	eventChannel, _ := hdlr.node.GetEventStream(
		ctx,
		req,
	)

	for {
		var event model.Event
		var isEventChannelOpen bool
		select {
		case <-ctx.Done():
			return
		case event, isEventChannelOpen = <-eventChannel:
		}
		if !isEventChannelOpen {
			// guard event channel closed
			return
		}

		err := conn.WriteJSON(event)
		if err != nil {
			// return on any error
			return
		}

		if orderedAcks != nil {
			select {
			case <-ctx.Done():
				return
			case <-orderedAcks:
			}
		} else if isAckRequested {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// return on any error
				return
			}
		}
	}

}
