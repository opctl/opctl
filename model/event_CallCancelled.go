package model

type CallCancelledEvent struct {
	*CallEventBase
	Container *ContainerCallCancelledEvent `json:"container,omitempty"`
	Op        *OpCallCancelledEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallCancelledEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallCancelledEvent    `json:"serial,omitempty"`
}

type ContainerCallCancelledEvent struct{}

type OpCallCancelledEvent struct{}

type ParallelCallCancelledEvent struct{}

type SerialCallCancelledEvent struct{}
