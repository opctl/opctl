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
		opDirHandle model.DataHandle,
		rootOpId string,
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
	opDirHandle model.DataHandle,
	rootOpId string,
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
			opDirHandle,
			rootOpId,
		)
	case nil != scg.Op:
		return this.opCaller.Call(
			scope,
			id,
			opDirHandle,
			rootOpId,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return this.parallelCaller.Call(
			id,
			scope,
			rootOpId,
			opDirHandle,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return this.serialCaller.Call(
			id,
			scope,
			rootOpId,
			opDirHandle,
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
