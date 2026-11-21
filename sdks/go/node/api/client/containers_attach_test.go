package client

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

type failingTerminalConn struct {
	net.Conn
	writesLeft int
	err        error
}

func (c *failingTerminalConn) Write(data []byte) (int, error) {
	if c.writesLeft == 0 {
		return 0, c.err
	}
	if c.writesLeft > 0 {
		c.writesLeft--
	}
	return c.Conn.Write(data)
}

func TestContainerInputReportsBytesSentBeforeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	writeErr := errors.New("connection write failed")
	var transport *failingTerminalConn
	dialer := websocket.Dialer{
		// Fit each WebSocket message in one transport write.
		WriteBufferSize: 128 * 1024,
		NetDial: func(network, addr string) (net.Conn, error) {
			conn, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			transport = &failingTerminalConn{Conn: conn, writesLeft: -1, err: writeErr}
			return transport, nil
		},
	}
	conn, _, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	transport.writesLeft = 1
	input := &containerInput{conn: conn}
	n, err := input.Write(make([]byte, 64*1024+1))
	if n != 64*1024 || !errors.Is(err, writeErr) {
		t.Fatalf("Write returned (%d, %v), want (65536, %v)", n, err, writeErr)
	}
}
