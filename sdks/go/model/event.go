package model

import (
	"time"

	"github.com/ipld/go-ipld-prime"
)

// Event represents a distributed state change
type Event struct {
	AuthAdded                *AuthAdded                `json:"authAdded,omitempty"`
	CallEnded                *CallEnded                `json:"callEnded,omitempty"`
	CallStarted              *CallStarted              `json:"callStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenTo `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenTo `json:"containerStdOutWrittenTo,omitempty"`
	CallKillRequested        *CallKillRequested        `json:"callKillRequested,omitempty"`
	Timestamp                time.Time                 `json:"timestamp"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

// AuthAdded represents auth was added for external resources
type AuthAdded struct {
	Auth Auth `json:"auth"`
}

// CallKillRequested represents a request was made to kill an op; a CallEnded event may follow
type CallKillRequested struct {
	Request KillOpReq `json:"request"`
}

// CallEnded represents a call ended; no further events will occur for the call
type CallEnded struct {
	Call    Call              `json:"call"`
	Ref     string            `json:"ref"`
	Error   *CallEndedError   `json:"error,omitempty"`
	Outputs map[string]*ipld.Node `json:"outputs"`
	Outcome string            `json:"outcome"`
}

// CallStarted represents the start of an op
type CallStarted struct {
	Call Call   `json:"call"`
	Ref  string `json:"ref"`
}

// CallEndedError represents an error associated w/ an ended call
type CallEndedError struct {
	Message string `json:"message"`
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
