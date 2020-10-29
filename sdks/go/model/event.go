package model

import "time"

// Event represents a distributed state change
type Event struct {
	CallEnded                *CallEnded                `json:"callEnded,omitempty"`
	CallStarted              *CallStarted              `json:"callStarted,omitempty"`
	ContainerExited          *ContainerExited          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStarted         `json:"containerStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenTo `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenTo `json:"containerStdOutWrittenTo,omitempty"`
	OpKillRequested          *OpKillRequested          `json:"opKillRequested,omitempty"`
	Timestamp                time.Time                 `json:"timestamp"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
	CallTypeOp         = "Op"
)

// OpKillRequested represents a request was made to kill an op; a CallEnded event may follow
type OpKillRequested struct {
	Request KillOpReq `json:"request"`
}

// CallEnded represents a call ended; no further events will occur for the call
type CallEnded struct {
	CallID     string            `json:"callId"`
	CallType   string            `json:"callType"`
	Ref        string            `json:"ref"`
	Error      *CallEndedError   `json:"error,omitempty"`
	Outputs    map[string]*Value `json:"outputs"`
	Outcome    string            `json:"outcome"`
	RootCallID string            `json:"rootCallId"`
}

// CallStarted represents the start of an op
type CallStarted struct {
	CallID     string `json:"callId"`
	CallType   string `json:"callType"`
	RootCallID string `json:"rootCallId"`
	OpRef      string `json:"opRef"`
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
	RootCallID  string            `json:"rootCallId"`
	ContainerID string            `json:"containerId"`
	OpRef       string            `json:"opRef"`
	Outputs     map[string]*Value `json:"outputs"`
}

// ContainerStarted represents the start of a container
type ContainerStarted struct {
	ImageRef    string `json:"imageRef"`
	RootCallID  string `json:"rootCallId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdErrWrittenTo represents a single write to a containers std err.
type ContainerStdErrWrittenTo struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootCallID  string `json:"rootCallId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}

// ContainerStdOutWrittenTo represents a single write to a containers std out.
type ContainerStdOutWrittenTo struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootCallID  string `json:"rootCallId"`
	ContainerID string `json:"containerId"`
	OpRef       string `json:"opRef"`
}
