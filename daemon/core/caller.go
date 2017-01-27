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
		args map[string]*model.Data,
		scg *model.Scg,
		opRef string,
		opGraphId string,
	) (
		outputs map[string]*model.Data,
		err error,
	)
}

func newCaller(
	containerCaller containerCaller,
) _caller {
	return _caller{
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
	args map[string]*model.Data,
	scg *model.Scg,
	opRef string,
	opGraphId string,
) (
	outputs map[string]*model.Data,
	err error,
) {

	switch {
	case nil != scg.Container:
		outputs, err = this.containerCaller.Call(
			args,
			nodeId,
			scg.Container,
			opRef,
			opGraphId,
		)
	case nil != scg.Op:
		outputs, err = this.opCaller.Call(
			args,
			nodeId,
			path.Join(filepath.Dir(opRef), scg.Op.Ref),
			opGraphId,
		)
	case len(scg.Parallel) > 0:
		err = this.parallelCaller.Call(
			args,
			opGraphId,
			opRef,
			scg.Parallel,
		)
	case len(scg.Serial) > 0:
		err = this.serialCaller.Call(
			args,
			opGraphId,
			opRef,
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
