package node

import "context"

/**
driver is a generic interface for things which expose opctl via some protocol
*/
type driver interface {
	Drive(
		ctx context.Context,
	) <-chan error
}
