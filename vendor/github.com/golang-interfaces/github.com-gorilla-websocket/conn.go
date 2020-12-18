package iwebsocket

//go:generate counterfeiter -o fakeConn.go --fake-name FakeConn ./ Conn

// Conn is implemented by websocket.Conn
type Conn interface {
// Close closes the underlying network connection without sending or waiting for a close frame.
Close() error
}
