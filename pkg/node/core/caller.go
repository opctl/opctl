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
		inboundScope map[string]*model.Data,
		scg *model.Scg,
		opPkgRef string,
		rootOpId string,
	) (
		outboundScope map[string]*model.Data,
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
	inboundScope map[string]*model.Data,
	scg *model.Scg,
	opPkgRef string,
	rootOpId string,
) (
	outboundScope map[string]*model.Data,
	err error,
) {

	switch {
	case nil != scg.Container:
		outboundScope, err = this.containerCaller.Call(
			inboundScope,
			nodeId,
			scg.Container,
			opPkgRef,
			rootOpId,
		)
	case nil != scg.Op:
		outboundScope, err = this.opCaller.Call(
			inboundScope,
			nodeId,
			path.Join(filepath.Dir(opPkgRef), scg.Op.Ref),
			rootOpId,
		)
	case len(scg.Parallel) > 0:
		err = this.parallelCaller.Call(
			inboundScope,
			rootOpId,
			opPkgRef,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		outboundScope, err = this.serialCaller.Call(
			inboundScope,
			rootOpId,
			opPkgRef,
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
