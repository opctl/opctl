package call

import (
	"context"
	"errors"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret a spec into a call
func Interpret(
	ctx context.Context,
	scope map[string]*model.Value,
	callSpec *model.CallSpec,
	id string,
	opPath string,
	parentID *string,
	rootCallID string,
	dataDirPath string,
) (*model.Call, error) {
	var name *string
	if callSpec.Name != nil {
		value, err := str.Interpret(scope, *callSpec.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to interpret call name: %w", err)
		}
		if value.String == nil {
			return nil, errors.New("call name not interpretable to string")
		}
		name = value.String
	}
	call := &model.Call{
		ID:       id,
		Name:     name,
		Needs:    callSpec.Needs,
		ParentID: parentID,
		RootID:   rootCallID,
	}
	var err error

	if callSpec.If != nil {
		callIf, err := predicates.Interpret(
			*callSpec.If,
			scope,
		)
		if err != nil {
			return nil, err
		}

		call.If = &callIf

		if !callIf {
			// end interpretation early since call will be skipped
			return call, err
		}
	}

	switch {
	case callSpec.Container != nil:
		call.Container, err = container.Interpret(
			scope,
			callSpec.Container,
			id,
			opPath,
			dataDirPath,
		)
		return call, err
	case callSpec.Op != nil:
		call.Op, err = op.Interpret(
			ctx,
			scope,
			callSpec.Op,
			id,
			opPath,
			dataDirPath,
		)
		return call, err
	case callSpec.Parallel != nil:
		call.Parallel = *callSpec.Parallel
		return call, nil
	case callSpec.ParallelLoop != nil:
		call.ParallelLoop, err = parallelloop.Interpret(
			*callSpec.ParallelLoop,
			scope,
		)
		return call, err
	case callSpec.Serial != nil:
		call.Serial = *callSpec.Serial
		return call, nil
	case callSpec.SerialLoop != nil:
		call.SerialLoop, err = serialloop.Interpret(
			*callSpec.SerialLoop,
			scope,
		)
		return call, err
	default:
		return nil, fmt.Errorf("invalid call graph '%+v'", callSpec)
	}
}
