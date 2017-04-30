package model

import (
	"fmt"
	"strings"
)

const (
	CallOutcomeSuccessful   = "SUCCESSFUL"
	CallOutcomeUnsuccessful = "UNSUCCESSFUL"
	CallOutcomeCancelled    = "CANCELLED"
)

type CallEndedEventBase struct {
	*CallEventBase
	Outcome string
}

// implement fmt.Stringer interface
func (ceeb CallEndedEventBase) String() string {
	return strings.Join(
		[]string{
			ceeb.CallEventBase.String(),
			fmt.Sprintf("Outcome='%v'", ceeb.Outcome),
		},
		" ",
	)
}

type CallEndedEvent struct {
	*CallEndedEventBase
	Container *ContainerCallEndedEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallEndedEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallEndedEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallEndedEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cee CallEndedEvent) String() string {
	var pre string
	switch {
	case nil != cee.Container:
		pre = cee.Container.String()
	case nil != cee.Op:
		pre = cee.Op.String()
	case nil != cee.Parallel:
		pre = cee.Parallel.String()
	case nil != cee.Serial:
		pre = cee.Serial.String()
	}

	return strings.Join(
		[]string{
			pre,
			cee.CallEndedEventBase.String(),
		},
		" ",
	)
}

type ContainerCallEndedEvent struct {
	ExitCode int
}

// implement fmt.Stringer interface
func (ccee ContainerCallEndedEvent) String() string {
	return strings.Join(
		[]string{
			"ContainerCallEnded",
			fmt.Sprintf("ExitCode='%v'", ccee.ExitCode),
		},
		" ",
	)
}

type OpCallEndedEvent struct{}

// implement fmt.Stringer interface
func (ocee OpCallEndedEvent) String() string {
	return "OpCallEnded"
}

type ParallelCallEndedEvent struct{}

// implement fmt.Stringer interface
func (pcee ParallelCallEndedEvent) String() string {
	return "ParallelCallEnded"
}

type SerialCallEndedEvent struct{}

// implement fmt.Stringer interface
func (scee SerialCallEndedEvent) String() string {
	return "SerialCallEnded"
}
