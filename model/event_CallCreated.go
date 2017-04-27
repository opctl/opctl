package model

type EventCallCreated struct {
	Container *ContainerCallCreatedEvent `json:"container,omitempty"`
	Op        *OpCallCreatedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallCreatedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallCreatedEvent    `json:"serial,omitempty"`
}

type ContainerCallCreatedEvent struct {
	*EventCallBase
	ImageRef string
}

type OpCallCreatedEvent struct {
	*EventCallBase
	PkgRef string
}

type ParallelCallCreatedEvent struct {
	*EventCallBase
}

type SerialCallCreatedEvent struct {
	*EventCallBase
}
