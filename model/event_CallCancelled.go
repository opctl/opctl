package model

import "strings"

// CallCancelledEvent represents cancellation of a call request.
type CallCancelledEvent struct {
	*CallEventBase
	Container *ContainerCallCancelledEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallCancelledEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallCancelledEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallCancelledEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cce CallCancelledEvent) String() string {
	var pre string
	switch {
	case nil != cce.Container:
		pre = cce.Container.String()
	case nil != cce.Op:
		pre = cce.Op.String()
	case nil != cce.Parallel:
		pre = cce.Parallel.String()
	case nil != cce.Serial:
		pre = cce.Serial.String()
	}

	return strings.Join(
		[]string{
			pre,
			cce.CallEventBase.String(),
		},
		" ",
	)
}

type ContainerCallCancelledEvent struct{}

// implement fmt.Stringer interface
func (ccce ContainerCallCancelledEvent) String() string {
	return "ContainerCallCancelled"
}

type OpCallCancelledEvent struct{}

// implement fmt.Stringer interface
func (occe OpCallCancelledEvent) String() string {
	return "OpCallCancelled"
}

type ParallelCallCancelledEvent struct{}

// implement fmt.Stringer interface
func (pcce ParallelCallCancelledEvent) String() string {
	return "ParallelCallCancelled"
}

type SerialCallCancelledEvent struct{}

// implement fmt.Stringer interface
func (scce SerialCallCancelledEvent) String() string {
	return "SerialCallCancelled"
}
