package image

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallImage *model.SCGContainerCallImage,
		opHandle model.DataHandle,
	) (*model.DCGContainerCallImage, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		expression: expression.New(),
	}
}

type _interpreter struct {
	expression expression.Expression
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallImage *model.SCGContainerCallImage,
	opHandle model.DataHandle,
) (*model.DCGContainerCallImage, error) {
	// construct image
	if scgContainerCallImage := scgContainerCallImage; scgContainerCallImage != nil {
		ref, err := itp.expression.EvalToString(scope, scgContainerCallImage.Ref, opHandle)
		if nil != err {
			return nil, err
		}

		dcgContainerCallImage := &model.DCGContainerCallImage{
			Ref: *ref.String,
		}

		if nil != scgContainerCallImage.PullCreds {
			username, err := itp.expression.EvalToString(scope, scgContainerCallImage.PullCreds.Username, opHandle)
			if nil != err {
				return nil, err
			}

			password, err := itp.expression.EvalToString(scope, scgContainerCallImage.PullCreds.Password, opHandle)
			if nil != err {
				return nil, err
			}

			dcgContainerCallImage.PullCreds = &model.DCGPullCreds{
				Username: *username.String,
				Password: *password.String,
			}
		}
		return dcgContainerCallImage, nil
	}
	return nil, errors.New("image required")
}
