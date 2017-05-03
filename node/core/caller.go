package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

type caller interface {
	Call(
		callId string,
		inboundScope map[string]*model.Data,
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

// Executes/runs an op
func (this _caller) Call(
	callId string,
	inboundScope map[string]*model.Data,
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
			inboundScope,
			callId,
			scg.Container,
			pkgRef,
			rootOpId,
		)
	case nil != scg.Op:
		return this.opCaller.Call(
			inboundScope,
			callId,
			path.Dir(pkgRef),
			rootOpId,
			scg.Op,
		)
	case len(scg.Parallel) > 0:
		return this.parallelCaller.Call(
			callId,
			inboundScope,
			rootOpId,
			pkgRef,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		return this.serialCaller.Call(
			callId,
			inboundScope,
			rootOpId,
			pkgRef,
			scg.Serial,
		)
	default:
		fmt.Printf("Invalid call graph %+v\n", scg)
	}

	return nil

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
