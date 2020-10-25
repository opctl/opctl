package model

import "time"

// Event represents a distributed state change
type Event struct {
	CallEnded                *CallEnded                `json:"callEnded,omitempty"`
	ContainerExited          *ContainerExited          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStarted         `json:"containerStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenTo `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenTo `json:"containerStdOutWrittenTo,omitempty"`
	OpKillRequested          *OpKillRequested          `json:"opKillRequested,omitempty"`
	OpEnded                  *OpEnded                  `json:"opEnded,omitempty"`
	OpStarted                *OpStarted                `json:"opStarted,omitempty"`
	Timestamp                time.Time                 `json:"timestamp"`
	ParallelCallEnded        *ParallelCallEnded        `json:"parallelCallEnded,omitempty"`
	ParallelLoopCallEnded    *ParallelLoopCallEnded    `json:"parallelLoopCallEnded,omitempty"`
	SerialCallEnded          *SerialCallEnded          `json:"serialCallEnded,omitempty"`
	SerialLoopCallEnded      *SerialLoopCallEnded      `json:"serialLoopCallEnded,omitempty"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

// OpKillRequested represents a request was made to kill an op; a CallEnded event may follow
type OpKillRequested struct {
	Request KillOpReq `json:"request"`
}

// CallEnded represents a call ended; no further events will occur for the call
type CallEnded struct {
	CallID     string            `json:"callId"`
	Ref        string            `json:"ref"`
	Error      *CallEndedError   `json:"error,omitempty"`
	Outputs    map[string]*Value `json:"outputs"`
	RootCallID string            `json:"rootCallId"`
}

// CallEndedError represents an error associated w/ an ended call
type CallEndedError struct {
	Message string `json:"message"`
}

// ContainerExited represents the exit of a containerized process; no further events will occur for the container
type ContainerExited struct {
	ImageRef    string            `json:"imageRef"`
	Error       *CallEndedError   `json:"error,omitempty"`
	ExitCode    int64             `json:"exitCode"`
	RootOpID    string            `json:"rootOpId"`
	ContainerID string            `json:"containerId"`
	OpRef       string            `json:"opRef"`
	Outputs     map[string]*Value `json:"outputs"`
}

type ContainerStarted struct {
	ImageRef    string `json:"imageRef"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdErrWrittenTo represents a single write to a containers std err.
type ContainerStdErrWrittenTo struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdOutWrittenTo represents a single write to a containers std out.
type ContainerStdOutWrittenTo struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpID    string `json:"rootOpId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// OpEnded represents the end of an op; no further events will occur for the op.
type OpEnded struct {
	Error    *CallEndedError   `json:"error,omitempty"`
	RootOpID string            `json:"rootOpId"`
	OpID     string            `json:"opId"`
	OpRef    string            `json:"opRef"`
	Outcome  string            `json:"outcome"`
	Outputs  map[string]*Value `json:"outputs"`
}

// OpStarted represents the start of an op
type OpStarted struct {
	RootOpID string `json:"rootOpId"`
	OpID     string `json:"opId"`
	OpRef    string `json:"opRef"`
}

// ParallelCallEnded represents the exit of a parallel call; no further events will occur for the call.
type ParallelCallEnded struct {
	CallID   string            `json:"callId"`
	Error    *CallEndedError   `json:"error,omitempty"`
	Outputs  map[string]*Value `json:"outputs"`
	RootOpID string            `json:"rootOpId"`
}

// ParallelLoopCallEnded represents the exit of a parallel loop call; no further events will occur for the call.
type ParallelLoopCallEnded struct {
	CallID   string            `json:"callId"`
	Error    *CallEndedError   `json:"error,omitempty"`
	Outputs  map[string]*Value `json:"outputs"`
	RootOpID string            `json:"rootOpId"`
}

// SerialCallEnded represents the exit of a serial call; no further events will occur for the call.
type SerialCallEnded struct {
	CallID   string            `json:"callId"`
	Error    *CallEndedError   `json:"error,omitempty"`
	Outputs  map[string]*Value `json:"outputs"`
	RootOpID string            `json:"rootOpId"`
}

// SerialLoopCallEnded represents the exit of a serial loop call; no further events will occur for the call.
type SerialLoopCallEnded struct {
	CallID   string            `json:"callId"`
	Error    *CallEndedError   `json:"error,omitempty"`
	Outputs  map[string]*Value `json:"outputs"`
	RootOpID string            `json:"rootOpId"`
}
