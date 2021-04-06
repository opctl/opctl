package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/parallelLoopCaller.go . parallelLoopCaller
type parallelLoopCaller interface {
	// Executes a parallel loop call
	Call(
		parentCtx context.Context,
		id string,
		inboundScope map[string]*model.Value,
		callSpecParallelLoop model.ParallelLoopCallSpec,
		opPath string,
		parentCallID *string,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newParallelLoopCaller(caller caller) parallelLoopCaller {
	return _parallelLoopCaller{
		caller: caller,
	}
}

type _parallelLoopCaller struct {
	caller caller
}

func (plpr _parallelLoopCaller) Call(
	parentCtx context.Context,
	id string,
	inboundScope map[string]*model.Value,
	callSpecParallelLoop model.ParallelLoopCallSpec,
	opPath string,
	parentCallID *string,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	// setup cancellation
	parallelLoopCtx, cancelParallelLoop := context.WithCancel(parentCtx)
	defer cancelParallelLoop()

	childCallIndex := 0
	childCallIndexByID := map[string]int{}

	type childResult struct {
		CallID  string
		Err     error
		Outputs map[string]*model.Value
	}
	childResults := make(chan childResult)

	// This waitgroup ensures all child goroutines are allowed to clean up
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		childCallID, err := uniquestring.Construct()
		if err != nil {
			// end run immediately on any error
			return nil, err
		}

		childCallScope, err := iteration.Scope(
			childCallIndex,
			inboundScope,
			callSpecParallelLoop.Range,
			callSpecParallelLoop.Vars,
		)
		if err != nil {
			return nil, err
		}

		// interpret iteration of the loop
		callParallelLoop, err := parallelloop.Interpret(
			callSpecParallelLoop,
			childCallScope,
		)
		if err != nil {
			return nil, err
		}

		if parallelloop.IsIterationComplete(childCallIndex, *callParallelLoop) {
			break
		}

		childCallIndexByID[childCallID] = childCallIndex

		wg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					childResults <- childResult{
						CallID:  childCallID,
						Err:     fmt.Errorf("panic: %v", r),
						Outputs: nil,
					}
				}
			}()

			defer wg.Done()

			outputs, err := plpr.caller.Call(
				parallelLoopCtx,
				childCallID,
				childCallScope,
				&callSpecParallelLoop.Run,
				opPath,
				parentCallID,
				rootCallID,
			)
			if parallelLoopCtx.Err() != nil {
				// context has been cancelled, so skip reporting results
				return
			}
			childResults <- childResult{
				CallID:  childCallID,
				Err:     err,
				Outputs: outputs,
			}
		}()

		childCallIndex++
	}

	if len(childCallIndexByID) == 0 {
		return nil, nil
	}

	childCallOutputsByIndex := map[int]map[string]*model.Value{}
	outboundScope := inboundScope

	for {
		select {
		case <-parallelLoopCtx.Done():
			return nil, parallelLoopCtx.Err()

		case result := <-childResults:
			if result.Err != nil {
				cancelParallelLoop()
				close(childResults)
				return nil, result.Err
			}

			if childCallIndex, isChildCallEnded := childCallIndexByID[result.CallID]; isChildCallEnded {
				childCallOutputsByIndex[childCallIndex] = result.Outputs
			}

			if len(childCallOutputsByIndex) == len(childCallIndexByID) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(childCallIndexByID); i++ {
					callOutputs := childCallOutputsByIndex[i]
					for varName, varData := range callOutputs {
						outboundScope[varName] = varData
					}
				}

				return loop.DeScope(
					inboundScope,
					callSpecParallelLoop.Range,
					callSpecParallelLoop.Vars,
					outboundScope,
				), nil
			}
		}
	}
}
