package core

//go:generate counterfeiter -o ./fakeCaller.go --fake-name fakeCaller ./ caller

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
)

type caller interface {
	Call(
		nodeId string,
		inputs chan *variable,
		outputs chan *variable,
		scg *model.Scg,
		pkgRef string,
		rootOpId string,
	) (
		err error,
	)
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
	nodeId string,
	inputs chan *variable,
	outputs chan *variable,
	scg *model.Scg,
	pkgRef string,
	rootOpId string,
) (
	err error,
) {

	if nil == scg {
		// No Op; equivalent to an empty fn body in a programming language
		return
	}

	switch {
	case nil != scg.Container:
		err = this.containerCaller.Call(
			inputs,
			outputs,
			nodeId,
			scg.Container,
			pkgRef,
			rootOpId,
		)
	case nil != scg.Op:
		err = this.opCaller.Call(
			inputs,
			outputs,
			nodeId,
			path.Join(filepath.Dir(pkgRef), scg.Op.Ref),
			rootOpId,
		)
	case len(scg.Parallel) > 0:
		close(outputs)
		err = this.parallelCaller.Call(
			inputs,
			rootOpId,
			pkgRef,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		err = this.serialCaller.Call(
			inputs,
			outputs,
			rootOpId,
			pkgRef,
			scg.Serial,
		)
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", scg)
	}

	return

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
