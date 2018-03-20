package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

type caller interface {
	// Call executes a call
	Call(
		id string,
		scope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		rootOpID string,
	) error
}

func newCaller(
	containerCaller containerCaller,
) *_caller {
	return &_caller{
		containerCaller: containerCaller,
	}
}

type _caller struct {
	containerCaller containerCaller
	opCaller        opCaller
	parallelCaller  parallelCaller
	serialCaller    serialCaller
}

func (this _caller) Call(
	id string,
	scope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	rootOpID string,
) error {

	if nil == scg {
		// No Op
		return nil
	}

	switch {
	case nil != scg.Container:
		return this.containerCaller.Call(
			scope,
			id,
			scg.Container,
			opHandle,
			rootOpID,
		)
	case nil != scg.Op:
		return this.opCaller.Call(
			scope,
			id,
			opHandle,
			rootOpID,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return this.parallelCaller.Call(
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return this.serialCaller.Call(
			id,
			scope,
			rootOpID,
			opHandle,
			scg.Serial,
		)
	default:
		return fmt.Errorf("Invalid call graph %+v\n", scg)
	}

}

func (this *_caller) setOpCaller(
	opCaller opCaller,
) {
	this.opCaller = opCaller
}

func (this *_caller) setParallelCaller(
	parallelCaller parallelCaller,
) {
	this.parallelCaller = parallelCaller
}

func (this *_caller) setSerialCaller(
	serialCaller serialCaller,
) {
	this.serialCaller = serialCaller
}
