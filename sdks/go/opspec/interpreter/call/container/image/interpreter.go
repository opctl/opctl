package image

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		callContainerImageSpec *model.CallContainerImageSpec,
		scratchDir string,
	) (*model.DCGContainerCallImage, error)
}

// NewInterpreter returns an initialized Interpreter instance.
func NewInterpreter() Interpreter {
	return _interpreter{
		dirInterpreter:    dir.NewInterpreter(),
		stringInterpreter: str.NewInterpreter(),
	}
}

type _interpreter struct {
	dirInterpreter    dir.Interpreter
	stringInterpreter str.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	callContainerImageSpec *model.CallContainerImageSpec,
	scratchDir string,
) (*model.DCGContainerCallImage, error) {

	if nil == callContainerImageSpec {
		return nil, fmt.Errorf("image required")
	}

	// try to interpret as dir
	src, err := itp.dirInterpreter.Interpret(
		scope,
		callContainerImageSpec.Ref,
		scratchDir,
		false,
	)
	if nil == err {
		return &model.DCGContainerCallImage{
			Src: src,
		}, nil
	}

	// fallback to string
	dcgContainerCallImage := &model.DCGContainerCallImage{}
	ref, err := itp.stringInterpreter.Interpret(
		scope,
		callContainerImageSpec.Ref,
	)
	if nil != err {
		return nil, err
	}

	dcgContainerCallImage.Ref = ref.String

	if nil != callContainerImageSpec.PullCreds {
		username, err := itp.stringInterpreter.Interpret(scope, callContainerImageSpec.PullCreds.Username)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds username; error was: %v", err)
		}

		password, err := itp.stringInterpreter.Interpret(scope, callContainerImageSpec.PullCreds.Password)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds password; error was: %v", err)
		}

		dcgContainerCallImage.PullCreds = &model.PullCreds{
			Username: *username.String,
			Password: *password.String,
		}
	}
	return dcgContainerCallImage, nil
}
