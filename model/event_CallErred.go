package model

import (
	"fmt"
	"strings"
)

type CallErredEventBase struct {
	*CallEventBase
	Msg string
}

// implement fmt.Stringer interface
func (ceeb CallErredEventBase) String() string {
	return strings.Join(
		[]string{
			ceeb.CallEventBase.String(),
			fmt.Sprintf("Msg='%v'", ceeb.Msg),
		},
		" ",
	)
}

type CallErredEvent struct {
	*CallErredEventBase
	Container *ContainerCallErredEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallErredEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallErredEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallErredEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cee CallErredEvent) String() string {
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
			cee.CallErredEventBase.String(),
		},
		" ",
	)
}

type ContainerCallErredEvent struct{}

// implement fmt.Stringer interface
func (ccee ContainerCallErredEvent) String() string {
	return "ContainerCallErred"
}

type OpCallErredEvent struct{}

// implement fmt.Stringer interface
func (ocee OpCallErredEvent) String() string {
	return "OpCallErred"
}

type ParallelCallErredEvent struct{}

// implement fmt.Stringer interface
func (pcee ParallelCallErredEvent) String() string {
	return "ParallelCallErred"
}

type SerialCallErredEvent struct{}

// implement fmt.Stringer interface
func (scee SerialCallErredEvent) String() string {
	return "SerialCallErred"
}
