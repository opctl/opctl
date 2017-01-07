package tcp

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	pongWait       = 60 * time.Second
	maxMessageSize = 1024 * 1024
)

type websocketClient struct {
	ws *websocket.Conn
}

func (c *websocketClient) readPump() {
	defer func() {
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
