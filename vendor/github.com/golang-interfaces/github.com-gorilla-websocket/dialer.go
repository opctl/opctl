package iwebsocket

//go:generate counterfeiter -o fakeDialer.go --fake-name FakeDialer ./ Dialer

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// Dialer is implemented by websocket.Dialer
type Dialer interface {
	// Dial creates a new client connection. Use requestHeader to specify the
	// origin (Origin), subprotocols (Sec-WebSocket-Protocol) and cookies (Cookie).
	// Use the response.Header to get the selected subprotocol
	// (Sec-WebSocket-Protocol) and cookies (Set-Cookie).
	//
	// If the WebSocket handshake fails, ErrBadHandshake is returned along with a
	// non-nil *http.Response so that callers can handle redirects, authentication,
	// etcetera. The response body may not contain the entire response and does not
	// need to be closed by the application.
	Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}
