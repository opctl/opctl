package model

const (
	CallOutcomeSuccessful   = "SUCCESSFUL"
	CallOutcomeUnsuccessful = "UNSUCCESSFUL"
	CallOutcomeCancelled    = "CANCELLED"
)

type CallEndedEventBase struct {
	*EventCallBase
	Outcome string
}

type CallEndedEvent struct {
	Container *ContainerCallEndedEvent `json:"container,omitempty"`
	Op        *OpCallEndedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallEndedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallEndedEvent    `json:"serial,omitempty"`
}

type ContainerCallEndedEvent struct {
	*CallEndedEventBase
}

type OpCallEndedEvent struct {
	*CallEndedEventBase
}

type ParallelCallEndedEvent struct {
	*CallEndedEventBase
}

type SerialCallEndedEvent struct {
	*CallEndedEventBase
}
