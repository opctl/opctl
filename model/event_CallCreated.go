package model

type EventCallCreated struct {
	*EventCallBase
	Container *ContainerCallCreatedEvent `json:"container,omitempty"`
	Op        *OpCallCreatedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallCreatedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallCreatedEvent    `json:"serial,omitempty"`
}

type ContainerCallCreatedEvent struct {
	ImageRef string
}

type OpCallCreatedEvent struct {
	PkgRef string
}

type ParallelCallCreatedEvent struct{}

type SerialCallCreatedEvent struct{}
