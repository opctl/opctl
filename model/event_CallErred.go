package model

type CallErredEventBase struct {
	*EventCallBase
	Msg string
}

type CallErredEvent struct {
	Container *ContainerCallErredEvent `json:"container,omitempty"`
	Op        *OpCallErredEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallErredEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallErredEvent    `json:"serial,omitempty"`
}

type ContainerCallErredEvent struct {
	*CallErredEventBase
}

type OpCallErredEvent struct {
	*CallErredEventBase
}

type ParallelCallErredEvent struct {
	*CallErredEventBase
}

type SerialCallErredEvent struct {
	*CallErredEventBase
}
