package node

import "context"

/**
listener is a generic interface for things which expose opctl via some protocol
*/
type listener interface {
	Listen(
		ctx context.Context,
	) <-chan error
}
