package model

type CallStartedEvent struct {
	*CallEventBase
	Container *ContainerCallStartedEvent `json:"container,omitempty"`
	Op        *OpCallStartedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallStartedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallStartedEvent    `json:"serial,omitempty"`
}

type ContainerCallStartedEvent struct{}

type OpCallStartedEvent struct{}

type ParallelCallStartedEvent struct{}

type SerialCallStartedEvent struct{}
