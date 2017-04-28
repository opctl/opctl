package model

import "time"

type Event struct {
	Timestamp                time.Time
	CallCancelled            *CallCancelledEvent            `json:"callCancelled,omitempty"`
	CallEnded                *CallEndedEvent                `json:"callEnded,omitempty"`
	CallErred                *CallErredEvent                `json:"callErred,omitempty"`
	CallRequested            *CallRequestedEvent            `json:"callRequested,omitempty"`
	CallStarted              *CallStartedEvent              `json:"callStartedEvent,omitempty"`
	ContainerExited          *ContainerExitedEvent          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStarted,omitempty"`
	ContainerStdErrEOFRead   *ContainerStdErrEOFReadEvent   `json:"containerStdErrEofRead,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutEOFRead   *ContainerStdOutEOFReadEvent   `json:"containerStdOutEofRead,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty"`
	OutputInitialized        *OutputInitializedEvent        `json:"outputInitialized,omitempty"`
}

type CallEventBase struct {
	CallID     string `json:"callId"`
	RootCallID string `json:"rootCallId"`
}

// represents a containerized process exiting
type ContainerExitedEvent struct {
	ImageRef    string
	ExitCode    int
	RootOpId    string
	ContainerId string
	PkgRef      string
}

type ContainerStartedEvent struct {
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdErrWrittenToEvent represents a single write to a containers std err.
type ContainerStdErrWrittenToEvent struct {
	ImageRef    string
	Data        []byte
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdErrEOFReadEvent represents EOF being read from a containers std err.
// This communicates, no further ContainerStdErrWrittenToEvent's will occur.
type ContainerStdErrEOFReadEvent struct {
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	ImageRef    string
	Data        []byte
	RootOpId    string
	ContainerId string
	PkgRef      string
}

// ContainerStdOutEOFReadEvent represents EOF being read from a containers std out.
// This communicates no further ContainerStdOutWrittenToEvent's will occur.
type ContainerStdOutEOFReadEvent struct {
	ImageRef    string
	RootOpId    string
	ContainerId string
	PkgRef      string
}

type OutputInitializedEvent struct {
	Name     string
	CallId   string
	Value    *Data
	RootOpId string
}
