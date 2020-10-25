package model

import "time"

// Event represents a distributed state change
type Event struct {
	CallEnded                *CallEndedEvent                `json:"callEnded,omitempty"`
	ContainerExited          *ContainerExitedEvent          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty"`
	OpKillRequested          *OpKillRequested               `json:"opKillRequested,omitempty"`
	OpEnded                  *OpEndedEvent                  `json:"opEnded,omitempty"`
	OpStarted                *OpStartedEvent                `json:"opStarted,omitempty"`
	Timestamp                time.Time                      `json:"timestamp"`
	ParallelCallEnded        *ParallelCallEndedEvent        `json:"parallelCallEnded,omitempty"`
	ParallelLoopCallEnded    *ParallelLoopCallEndedEvent    `json:"parallelLoopCallEnded,omitempty"`
	SerialCallEnded          *SerialCallEndedEvent          `json:"serialCallEnded,omitempty"`
	SerialLoopCallEnded      *SerialLoopCallEndedEvent      `json:"serialLoopCallEnded,omitempty"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

// OpKillRequested represents a request was made to kill an op; a CallKilledEvent may follow
type OpKillRequested struct {
	Request KillOpReq `json:"request"`
}

// CallEndedEvent represents a call ended; no further events will occur for the call
type CallEndedEvent struct {
	CallID     string               `json:"callId"`
	Ref        string               `json:"ref"`
	Error      *CallEndedEventError `json:"error,omitempty"`
	Outputs    map[string]*Value    `json:"outputs"`
	RootCallID string               `json:"rootCallId"`
}

// CallEndedEventError represents an error of associated w/ an ended call
type CallEndedEventError struct {
	Message string `json:"message"`
}

// ContainerExitedEvent represents the exit of a containerized process; no further events will occur for the container
type ContainerExitedEvent struct {
	ImageRef    string               `json:"imageRef"`
	Error       *CallEndedEventError `json:"error,omitempty"`
	ExitCode    int64                `json:"exitCode"`
	RootOpID    string               `json:"rootOpId"`
	ContainerID string               `json:"containerId"`
	OpRef       string               `json:"opRef"`
	Outputs     map[string]*Value    `json:"outputs"`
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
	Error    *CallEndedEventError `json:"error,omitempty"`
	RootOpID string               `json:"rootOpId"`
	OpID     string               `json:"opId"`
	OpRef    string               `json:"opRef"`
	Outcome  string               `json:"outcome"`
	Outputs  map[string]*Value    `json:"outputs"`
}

// OpStartedEvent represents the start of an op
type OpStartedEvent struct {
	RootOpID string `json:"rootOpId"`
	OpID     string `json:"opId"`
	OpRef    string `json:"opRef"`
}

// ParallelCallEndedEvent represents the exit of a parallel call; no further events will occur for the call.
type ParallelCallEndedEvent struct {
	CallID   string               `json:"callId"`
	Error    *CallEndedEventError `json:"error,omitempty"`
	Outputs  map[string]*Value    `json:"outputs"`
	RootOpID string               `json:"rootOpId"`
}

// ParallelLoopCallEndedEvent represents the exit of a parallel loop call; no further events will occur for the call.
type ParallelLoopCallEndedEvent struct {
	CallID   string               `json:"callId"`
	Error    *CallEndedEventError `json:"error,omitempty"`
	Outputs  map[string]*Value    `json:"outputs"`
	RootOpID string               `json:"rootOpId"`
}

// SerialCallEndedEvent represents the exit of a serial call; no further events will occur for the call.
type SerialCallEndedEvent struct {
	CallID   string               `json:"callId"`
	Error    *CallEndedEventError `json:"error,omitempty"`
	Outputs  map[string]*Value    `json:"outputs"`
	RootOpID string               `json:"rootOpId"`
}

// SerialLoopCallEndedEvent represents the exit of a serial loop call; no further events will occur for the call.
type SerialLoopCallEndedEvent struct {
	CallID   string               `json:"callId"`
	Error    *CallEndedEventError `json:"error,omitempty"`
	Outputs  map[string]*Value    `json:"outputs"`
	RootOpID string               `json:"rootOpId"`
}
