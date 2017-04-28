package model

type CallRequestedEvent struct {
	*CallEventBase
	Container *ContainerCallRequestedEvent `json:"container,omitempty"`
	Op        *OpCallRequestedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallRequestedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallRequestedEvent    `json:"serial,omitempty"`
}

type ContainerCallRequestedEvent struct {
	ImageRef string
}

type OpCallRequestedEvent struct {
	PkgRef string
}

type ParallelCallRequestedEvent struct{}

type SerialCallRequestedEvent struct{}
