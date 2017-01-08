package core

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
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
	bundle bundle.Bundle,
	containerEngine containerengine.ContainerEngine,
	eventBus eventbus.EventBus,
	nodeRepo nodeRepo,
	opCaller opCaller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) caller {

	objectUnderConstruction := &_caller{
		containerCaller: newContainerCaller(bundle, containerEngine, eventBus, nodeRepo),
		opCaller:        opCaller,
	}

	objectUnderConstruction.parallelCaller = newParallelCaller(
		objectUnderConstruction,
		uniqueStringFactory,
	)

	objectUnderConstruction.serialCaller = newSerialCaller(
		objectUnderConstruction,
		uniqueStringFactory,
	)

	return objectUnderConstruction

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
	case nil != scg.Op:
		outputs, err = this.opCaller.Call(
			args,
			nodeId,
			path.Join(filepath.Dir(opRef), scg.Op.Ref),
			opGraphId,
		)
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", scg)
	}

	return

}
