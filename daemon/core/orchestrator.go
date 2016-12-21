package core

import (
	"fmt"
	"github.com/opspec-io/engine/pkg/containerengine"
	"github.com/opspec-io/engine/util/eventbus"
	"github.com/opspec-io/engine/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
	"path/filepath"
)

type orchestrator interface {
	Execute(
		nodeId string,
		args map[string]*model.Arg,
		callGraph *model.CallGraph,
		opRef string,
		opGraphId string,
	) (
		err error,
	)
}

func newOrchestrator(
	bundle bundle.Bundle,
	containerEngine containerengine.ContainerEngine,
	eventBus eventbus.EventBus,
	nodeRepo nodeRepo,
	opOrchestrator opOrchestrator,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) orchestrator {

	objectUnderConstruction := &_orchestrator{
		containerOrchestrator: newContainerOrchestrator(bundle, containerEngine, eventBus, nodeRepo),
		opOrchestrator:        opOrchestrator,
	}

	objectUnderConstruction.parallelOrchestrator = newParallelOrchestrator(
		objectUnderConstruction,
		uniqueStringFactory,
	)

	objectUnderConstruction.serialOrchestrator = newSerialOrchestrator(
		objectUnderConstruction,
		uniqueStringFactory,
	)

	return objectUnderConstruction

}

type _orchestrator struct {
	containerOrchestrator containerOrchestrator
	opOrchestrator        opOrchestrator
	parallelOrchestrator  parallelOrchestrator
	serialOrchestrator    serialOrchestrator
}

// Executes/runs an op
func (this _orchestrator) Execute(
	nodeId string,
	args map[string]*model.Arg,
	callGraph *model.CallGraph,
	opRef string,
	opGraphId string,
) (
	err error,
) {

	switch {
	case nil != callGraph.Container:
		{
			err = this.containerOrchestrator.Execute(
				args,
				nodeId,
				callGraph.Container,
				opRef,
				opGraphId,
			)
		}
	case len(callGraph.Parallel) > 0:
		{
			err = this.parallelOrchestrator.Execute(
				args,
				opGraphId,
				opRef,
				callGraph.Parallel,
			)
		}
	case len(callGraph.Serial) > 0:
		{
			err = this.serialOrchestrator.Execute(
				args,
				opGraphId,
				opRef,
				callGraph.Serial,
			)
		}
	case nil != callGraph.Op:
		{
			err = this.opOrchestrator.Execute(
				args,
				nodeId,
				path.Join(filepath.Dir(opRef), callGraph.Op.Ref),
				opGraphId,
			)
		}
	default:
		err = fmt.Errorf("Invalid call graph %+v\n", callGraph)
	}

	return

}
