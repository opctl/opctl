package image

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"errors"

	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
	"github.com/opctl/opctl/sdks/go/types"
)

type Interpreter interface {
	Interpret(
		scope map[string]*types.Value,
		scgContainerCallImage *types.SCGContainerCallImage,
		opHandle types.DataHandle,
	) (*types.DCGContainerCallImage, error)
}

// NewInterpreter returns an initialized Interpreter instance.
func NewInterpreter() Interpreter {
	return _interpreter{
		stringInterpreter: stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	stringInterpreter stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*types.Value,
	scgContainerCallImage *types.SCGContainerCallImage,
	opHandle types.DataHandle,
) (*types.DCGContainerCallImage, error) {

	if nil == scgContainerCallImage {
		return nil, errors.New("image required")
	}

	// construct image
	ref, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.Ref, opHandle)
	if nil != err {
		return nil, err
	}

	dcgContainerCallImage := &types.DCGContainerCallImage{
		Ref: *ref.String,
	}

	if nil != scgContainerCallImage.PullCreds {
		username, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.PullCreds.Username, opHandle)
		if nil != err {
			return nil, err
		}

		password, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.PullCreds.Password, opHandle)
		if nil != err {
			return nil, err
		}

		dcgContainerCallImage.PullCreds = &types.PullCreds{
			Username: *username.String,
			Password: *password.String,
		}
	}
	return dcgContainerCallImage, nil
}
