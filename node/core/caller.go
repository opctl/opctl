package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

type caller interface {
	// Call executes a call
	Call(
		id string,
		scope map[string]*model.Data,
		scg *model.SCG,
		pkgRef string,
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
	scope map[string]*model.Data,
	scg *model.SCG,
	pkgRef string,
	rootOpId string,
) error {

	if nil == scg {
		// No Op; equivalent to an empty fn body in a programming language
		return nil
	}

	switch {
	case nil != scg.Container:
		return this.containerCaller.Call(
			scope,
			id,
			scg.Container,
			pkgRef,
			rootOpId,
		)
	case nil != scg.Op:
		return this.opCaller.Call(
			scope,
			id,
			filepath.Dir(pkgRef),
			rootOpId,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return this.parallelCaller.Call(
			id,
			scope,
			rootOpId,
			pkgRef,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return this.serialCaller.Call(
			id,
			scope,
			rootOpId,
			pkgRef,
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
