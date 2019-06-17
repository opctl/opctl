package image

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"errors"

	"github.com/opctl/sdk-golang/model"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallImage *model.SCGContainerCallImage,
		opHandle model.DataHandle,
	) (*model.DCGContainerCallImage, error)
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
	scope map[string]*model.Value,
	scgContainerCallImage *model.SCGContainerCallImage,
	opHandle model.DataHandle,
) (*model.DCGContainerCallImage, error) {

	if nil == scgContainerCallImage {
		return nil, errors.New("image required")
	}

	// construct image
	ref, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.Ref, opHandle)
	if nil != err {
		return nil, err
	}

	dcgContainerCallImage := &model.DCGContainerCallImage{
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

		dcgContainerCallImage.PullCreds = &model.PullCreds{
			Username: *username.String,
			Password: *password.String,
		}
	}
	return dcgContainerCallImage, nil
}
