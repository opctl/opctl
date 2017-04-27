package model

type CallCancelledEvent struct {
	Container *ContainerCallCancelledEvent `json:"container,omitempty"`
	Op        *OpCallCancelledEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallCancelledEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallCancelledEvent    `json:"serial,omitempty"`
}

type ContainerCallCancelledEvent struct {
	*EventCallBase
}

type OpCallCancelledEvent struct {
	*EventCallBase
}

type ParallelCallCancelledEvent struct {
	*EventCallBase
}

type SerialCallCancelledEvent struct {
	*EventCallBase
}
