package model

const (
	CallOutcomeSuccessful   = "SUCCESSFUL"
	CallOutcomeUnsuccessful = "UNSUCCESSFUL"
	CallOutcomeCancelled    = "CANCELLED"
)

type CallEndedEventBase struct {
	*CallEventBase
	Outcome string
}

type CallEndedEvent struct {
	*CallEndedEventBase
	Container *ContainerCallEndedEvent `json:"container,omitempty"`
	Op        *OpCallEndedEvent        `json:"op,omitempty"`
	Parallel  *ParallelCallEndedEvent  `json:"parallel,omitempty"`
	Serial    *SerialCallEndedEvent    `json:"serial,omitempty"`
}

type ContainerCallEndedEvent struct{}

type OpCallEndedEvent struct{}

type ParallelCallEndedEvent struct{}

type SerialCallEndedEvent struct{}
