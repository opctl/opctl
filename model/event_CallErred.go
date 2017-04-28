package model

type CallErredEventBase struct {
	*CallEventBase
	Msg string
}

type CallErredEvent struct {
	*CallErredEventBase
	Container *ContainerCallErredEvent `json:"container,omitempty"`
	Op        *OpCallErredEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallErredEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallErredEvent    `json:"serial,omitempty"`
}

type ContainerCallErredEvent struct{}

type OpCallErredEvent struct{}

type ParallelCallErredEvent struct{}

type SerialCallErredEvent struct{}
