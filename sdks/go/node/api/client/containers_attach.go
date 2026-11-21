package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/urltemplates"
)

func (c apiClient) AttachContainer(ctx context.Context, containerID string) (node.ContainerInput, error) {
	reqURL := c.baseURL
	if reqURL.Scheme == "https" {
		reqURL.Scheme = "wss"
	} else {
		reqURL.Scheme = "ws"
	}
	endpoint := strings.Replace(urltemplates.Containers_Attach, "{containerID}", url.PathEscape(containerID), 1)
	reqURL.RawPath = strings.TrimRight(reqURL.EscapedPath(), "/") + endpoint
	reqURL.Path, _ = url.PathUnescape(reqURL.RawPath)

	conn, response, err := c.wsDialer.DialContext(ctx, reqURL.String(), nil)
	if err != nil {
		if response != nil {
			defer response.Body.Close()
			switch response.StatusCode {
			case http.StatusNotFound:
				return nil, node.ErrContainerInputNotFound
			case http.StatusConflict:
				return nil, node.ErrContainerInputAttached
			}
		}
		return nil, fmt.Errorf("attach container: %w", err)
	}

	input := &containerInput{conn: conn, done: make(chan struct{})}
	go func() {
		defer input.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	return input, nil
}

type containerInput struct {
	conn      *websocket.Conn
	done      chan struct{}
	writeMu   sync.Mutex
	closeOnce sync.Once
	closeErr  error
}

const inputWriteTimeout = 2 * time.Second

func (input *containerInput) Write(data []byte) (int, error) {
	input.writeMu.Lock()
	defer input.writeMu.Unlock()
	var sent int
	for len(data) > 0 {
		chunk := data[:min(len(data), node.MaxContainerInputMessageSize)]
		if err := input.conn.SetWriteDeadline(time.Now().Add(inputWriteTimeout)); err != nil {
			return sent, err
		}
		if err := input.conn.WriteMessage(websocket.BinaryMessage, chunk); err != nil {
			return sent, err
		}
		sent += len(chunk)
		data = data[len(chunk):]
	}
	return sent, nil
}

func (input *containerInput) Resize(rows, cols uint) error {
	input.writeMu.Lock()
	defer input.writeMu.Unlock()
	if err := input.conn.SetWriteDeadline(time.Now().Add(inputWriteTimeout)); err != nil {
		return err
	}
	return input.conn.WriteJSON(model.TermSize{Rows: rows, Cols: cols})
}

func (input *containerInput) Done() <-chan struct{} { return input.done }

func (input *containerInput) Close() error {
	input.closeOnce.Do(func() {
		input.closeErr = input.conn.Close()
		close(input.done)
	})
	return input.closeErr
}
