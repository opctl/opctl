package model

import "time"

type Event struct {
	CallCancelled            *CallCancelledEvent            `json:"callCancelled,omitempty"`
	CallCreated              *EventCallCreated              `json:"callCreated,omitempty"`
	CallEnded                *CallEndedEvent                `json:"callEnded,omitempty"`
	CallErred                *CallErredEvent                `json:"callErred,omitempty"`
	CallStartedEvent         *EventCallStarted              `json:"callStartedEvent,omitempty"`
	ContainerExited          *ContainerExitedEvent          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStarted,omitempty"`
	ContainerStdErrEOFRead   *ContainerStdErrEOFReadEvent   `json:"containerStdErrEofRead,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutEOFRead   *ContainerStdOutEOFReadEvent   `json:"containerStdOutEofRead,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty"`
	OutputInitialized        *OutputInitializedEvent        `json:"outputInitialized,omitempty"`
}

type EventBase struct {
	Timestamp time.Time
}

type EventCallBase struct {
	*EventBase
	CallID     string `json:"callId"`
	RootCallID string `json:"rootCallId"`
}

// represents a containerized process exiting
type ContainerExitedEvent struct {
	*EventBase
	ImageRef    string
	ExitCode    int
	RootOpId    string
	ContainerId string
	PkgRef      string
}

type ContainerStartedEvent struct {
	*EventBase
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdErrWrittenToEvent represents a single write to a containers std err.
type ContainerStdErrWrittenToEvent struct {
	*EventBase
	ImageRef    string
	Data        []byte
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdErrEOFReadEvent represents EOF being read from a containers std err.
// This communicates, no further ContainerStdErrWrittenToEvent's will occur.
type ContainerStdErrEOFReadEvent struct {
	*EventBase
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	*EventBase
	ImageRef    string
	Data        []byte
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdOutEOFReadEvent represents EOF being read from a containers std out.
// This communicates no further ContainerStdOutWrittenToEvent's will occur.
type ContainerStdOutEOFReadEvent struct {
	*EventBase
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

type OutputInitializedEvent struct {
	*EventBase
	Name     string
	CallId   string
	Value    *Data
	RootOpId string
}
