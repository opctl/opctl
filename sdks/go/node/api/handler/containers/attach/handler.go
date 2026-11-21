// Package attach handles container stdin and terminal resize over WebSocket.
package attach

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/internal/urlpath"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

func NewHandler(node node.Core) http.Handler {
	return _handler{node: node}
}

type _handler struct {
	node node.Core
}

func (hdlr _handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	containerID, err := urlpath.NextSegment(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	action, err := urlpath.NextSegment(req.URL)
	if err != nil || containerID == "" || action != "attach" || req.URL.Path != "" {
		http.NotFound(w, req)
		return
	}

	input, release, err := hdlr.node.ClaimContainerInput(req.Context(), containerID)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, node.ErrContainerInputNotFound):
			status = http.StatusNotFound
		case errors.Is(err, node.ErrContainerInputAttached):
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		release()
		return
	}
	// Gorilla counts payload bytes, so WebSocket frame headers fit outside the limit.
	conn.SetReadLimit(node.MaxContainerInputMessageSize)
	shutdown := sync.OnceFunc(func() {
		conn.Close()
		input.Close()
	})
	defer shutdown()

	stop := make(chan struct{})
	defer close(stop)
	go func() {
		select {
		case <-stop:
		case <-input.Done():
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(
				websocket.CloseNormalClosure, "",
			), time.Now().Add(time.Second))
			// Closing both ends unblocks either a socket read or a pipe write.
			shutdown()
		}
	}()

	for {
		// Write each message before reading the next one. A slow container
		// applies backpressure to the sender, and a close cannot skip input.
		kind, data, err := conn.ReadMessage()
		if err != nil || forward(input, kind, data) != nil {
			return
		}
	}
}

func forward(input node.ContainerInput, kind int, data []byte) error {
	if kind == websocket.BinaryMessage {
		_, err := input.Write(data)
		return err
	}

	var size model.TermSize
	if err := json.Unmarshal(data, &size); err != nil {
		return err
	}
	if size.Rows == 0 || size.Cols == 0 || size.Rows > 65535 || size.Cols > 65535 {
		return errors.New("invalid terminal size")
	}
	return input.Resize(size.Rows, size.Cols)
}
