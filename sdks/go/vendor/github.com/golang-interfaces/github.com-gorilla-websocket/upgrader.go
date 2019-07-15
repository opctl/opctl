package iwebsocket

//go:generate counterfeiter -o fakeUpgrader.go --fake-name FakeUpgrader ./ Upgrader

import (
  "net/http"
  "github.com/gorilla/websocket"
)

// Upgrader is implemented by websocket.Upgrader
type Upgrader interface {
// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// application negotiated subprotocol (Sec-Websocket-Protocol).
//
// If the upgrade fails, then Upgrade replies to the client with an HTTP error
// response.
Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}
