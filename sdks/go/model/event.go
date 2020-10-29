package model

import "time"

// Event represents a distributed state change
type Event struct {
	AuthAdded                *AuthAdded                `json:"authAdded,omitempty"`
	CallEnded                *CallEnded                `json:"callEnded,omitempty"`
	CallStarted              *CallStarted              `json:"callStarted,omitempty"`
	ContainerExited          *ContainerExited          `json:"containerExited,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenTo `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenTo `json:"containerStdOutWrittenTo,omitempty"`
	OpKillRequested          *OpKillRequested          `json:"opKillRequested,omitempty"`
	Timestamp                time.Time                 `json:"timestamp"`
}

const (
	OpOutcomeSucceeded   = "SUCCEEDED"
	OpOutcomeFailed      = "FAILED"
	OpOutcomeKilled      = "KILLED"
	CallTypeOp           = "Op"
	CallTypeSerial       = "Serial"
	CallTypeSerialLoop   = "SerialLoop"
	CallTypeParallel     = "Parallel"
	CallTypeParallelLoop = "ParallelLoop"
)

// AuthAdded represents auth was added for external resources
type AuthAdded struct {
	Auth Auth `json:"auth"`
}

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
	Call       Call   `json:"call"`
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
