package model

type EventCallStarted struct {
	Container *ContainerCallStartedEvent `json:"container,omitempty"`
	Op        *OpCallStartedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallStartedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallStartedEvent    `json:"serial,omitempty"`
}

type ContainerCallStartedEvent struct {
	*EventCallBase
}

type OpCallStartedEvent struct {
	*EventCallBase
}

type ParallelCallStartedEvent struct {
	*EventCallBase
}

type SerialCallStartedEvent struct {
	*EventCallBase
}
