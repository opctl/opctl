package model

import "time"

// Event represents a distributed state change
type Event struct {
	CallEnded                *CallEndedEvent                `json:"callEnded,omitempty"`
	ContainerExited          *ContainerExitedEvent          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty"`
	OpEnded                  *OpEndedEvent                  `json:"opEnded,omitempty"`
	OpStarted                *OpStartedEvent                `json:"opStarted,omitempty"`
	OpErred                  *OpErredEvent                  `json:"opErred,omitempty"`
	Timestamp                time.Time                      `json:"timestamp"`
	ParallelCallEnded        *ParallelCallEndedEvent        `json:"parallelCallEnded,omitempty"`
	SerialCallEnded          *SerialCallEndedEvent          `json:"serialCallEnded,omitempty"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

// CallEndedEvent represents a call ended; no further events will occur for the call
type CallEndedEvent struct {
	CallID     string `json:"callId`
	RootCallID string `json:"rootCallId"`
}

// ContainerExitedEvent represents the exit of a containerized process; no further events will occur for the container
type ContainerExitedEvent struct {
	ImageRef    string            `json:"imageRef"`
	ExitCode    int64             `json:"exitCode"`
	RootOpID    string            `json:"rootOpId"`
	ContainerID string            `json:"containerId"`
	OpRef       string            `json:"opRef"`
	Outputs     map[string]*Value `json:"outputs"`
}

type ContainerStartedEvent struct {
	ImageRef    string `json:"imageRef"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdErrWrittenToEvent represents a single write to a containers std err.
type ContainerStdErrWrittenToEvent struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// OpErredEvent represents an op encountering an error condition
type OpErredEvent struct {
	RootOpID string `json:"rootOpId"`
	Msg      string `json:"msg"`
	OpID     string `json:"opId"`
	OpRef    string `json:"opRef"`
}

// OpEndedEvent represents the end of an op; no further events will occur for the op.
type OpEndedEvent struct {
	RootOpID string            `json:"rootOpId"`
	OpID     string            `json:"opId"`
	OpRef    string            `json:"opRef"`
	Outcome  string            `json:"outcome"`
	Outputs  map[string]*Value `json:"outputs"`
}

// OpStartedEvent represents the start of an op
type OpStartedEvent struct {
	RootOpID string `json:"rootOpId"`
	OpID     string `json:"opId"`
	OpRef    string `json:"opRef"`
}

// ParallelCallEndedEvent represents the exit of a parallel call; no further events will occur for the call.
type ParallelCallEndedEvent struct {
	RootOpID string `json:"rootOpId"`
	CallID   string `json:"callId"`
}

// SerialCallEndedEvent represents the exit of a serial call; no further events will occur for the call.
type SerialCallEndedEvent struct {
	RootOpID string            `json:"rootOpId"`
	CallID   string            `json:"callId"`
	Outputs  map[string]*Value `json:"outputs"`
}
