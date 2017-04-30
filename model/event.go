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
	ContainerStdErrEOFRead   *ContainerStdErrEOFReadEvent   `json:"containerStdErrEofRead,omitempty" yaml:"containerStdErrEofRead,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty" yaml:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutEOFRead   *ContainerStdOutEOFReadEvent   `json:"containerStdOutEofRead,omitempty" yaml:"containerStdOutEofRead,omitempty"`
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
	case nil != e.ContainerStdErrEOFRead:
		pre = e.ContainerStdErrEOFRead.String()
	case nil != e.ContainerStdErrWrittenTo:
		pre = e.ContainerStdErrWrittenTo.String()
	case nil != e.ContainerStdOutEOFRead:
		pre = e.ContainerStdOutEOFRead.String()
	case nil != e.ContainerStdOutWrittenTo:
		pre = e.ContainerStdOutWrittenTo.String()
	case nil != e.OutputInitialized:
		pre = e.OutputInitialized.String()
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
	Data        []byte
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// implement fmt.Stringer interface
func (csewte ContainerStdErrWrittenToEvent) String() string {
	return string(csewte.Data)
}

// ContainerStdErrEOFReadEvent represents EOF being read from a containers std err.
// This communicates, no further ContainerStdErrWrittenToEvent's will occur.
type ContainerStdErrEOFReadEvent struct {
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// implement fmt.Stringer interface
func (cseere ContainerStdErrEOFReadEvent) String() string {
	return ""
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	ImageRef    string
	Data        []byte
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// implement fmt.Stringer interface
func (csowte ContainerStdOutWrittenToEvent) String() string {
	return string(csowte.Data)
}

// ContainerStdOutEOFReadEvent represents EOF being read from a containers std out.
// This communicates no further ContainerStdOutWrittenToEvent's will occur.
type ContainerStdOutEOFReadEvent struct {
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// implement fmt.Stringer interface
func (csoere ContainerStdOutEOFReadEvent) String() string {
	return ""
}

type OutputInitializedEvent struct {
	Name     string
	CallId   string
	Value    *Data
	RootOpId string
}

// implement fmt.Stringer interface
func (oie OutputInitializedEvent) String() string {
	return ""
}
