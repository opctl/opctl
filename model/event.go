package model

import (
	"fmt"
	"strings"
	"time"
)

// Event represents a distributed state change
// Events implement the fmt.Stringer interface in addition to maintaining json and yaml struct tags.
type Event struct {
	Timestamp                time.Time
	CallAccepted             *CallAcceptedEvent             `json:"callAccepted,omitempty" yaml:"callAccepted,omitempty"`
	CallCancelled            *CallCancelledEvent            `json:"callCancelled,omitempty" yaml:"callCancelled,omitempty"`
	CallEnded                *CallEndedEvent                `json:"callEnded,omitempty" yaml:"callEnded,omitempty"`
	CallErred                *CallErredEvent                `json:"callErred,omitempty" yaml:"callErred,omitempty"`
	CallRequested            *CallRequestedEvent            `json:"callRequested,omitempty" yaml:"callRequested,omitempty"`
	CallStarted              *CallStartedEvent              `json:"callStartedEvent,omitempty" yaml:"callStartedEvent,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty" yaml:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty" yaml:"containerStdOutWrittenTo,omitempty"`
	OutputInitialized        *OutputInitializedEvent        `json:"outputInitialized,omitempty" yaml:"outputInitialized,omitempty"`
}

// implement fmt.Stringer interface
func (e Event) String() string {
	var pre string
	switch {
	case nil != e.CallAccepted:
		pre = e.CallAccepted.String()
	case nil != e.CallCancelled:
		pre = e.CallCancelled.String()
	case nil != e.CallEnded:
		pre = e.CallEnded.String()
	case nil != e.CallErred:
		pre = e.CallErred.String()
	case nil != e.CallRequested:
		pre = e.CallRequested.String()
	case nil != e.CallStarted:
		pre = e.CallStarted.String()
	case nil != e.ContainerStdErrWrittenTo:
		return e.ContainerStdErrWrittenTo.String()
	case nil != e.ContainerStdOutWrittenTo:
		return e.ContainerStdOutWrittenTo.String()
	case nil != e.OutputInitialized:
		return ""
	}

	return strings.Join(
		[]string{
			pre,
			fmt.Sprintf("Timestamp='%v'", e.Timestamp.Format(time.RFC3339)),
		},
		" ",
	)
}

type CallEventBase struct {
	CallID     string `json:"callId"`
	RootCallID string `json:"rootCallId"`
}

// implement fmt.Stringer interface
func (ceb CallEventBase) String() string {
	return strings.Join(
		[]string{
			fmt.Sprintf("RootCallId='%v'", ceb.RootCallID),
			fmt.Sprintf("CallId='%v'", ceb.CallID),
		},
		" ",
	)
}

// ContainerStdErrWrittenToEvent represents a single write to a containers std err.
type ContainerStdErrWrittenToEvent struct {
	*CallEventBase
	Data []byte
}

// implement fmt.Stringer interface
func (csewte ContainerStdErrWrittenToEvent) String() string {
	return string(csewte.Data)
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	*CallEventBase
	Data []byte
}

// implement fmt.Stringer interface
func (csowte ContainerStdOutWrittenToEvent) String() string {
	return string(csowte.Data)
}

type OutputInitializedEvent struct {
	*CallEventBase
	Name  string
	Value *Data
}
