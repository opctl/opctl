package image

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
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
		dirInterpreter:    dir.NewInterpreter(),
		stringInterpreter: stringPkg.NewInterpreter(),
	}
}

type _interpreter struct {
	dirInterpreter    dir.Interpreter
	stringInterpreter stringPkg.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCallImage *model.SCGContainerCallImage,
	opHandle model.DataHandle,
) (*model.DCGContainerCallImage, error) {

	if nil == scgContainerCallImage {
		return nil, fmt.Errorf("image required")
	}

	if nil != scgContainerCallImage.Src {
		src, err := itp.dirInterpreter.Interpret(
			scope,
			*scgContainerCallImage.Src,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image src; error was: %v", err)
		}

		return &model.DCGContainerCallImage{
			Src: src,
		}, nil
	}

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
			return nil, fmt.Errorf("error encountered interpreting image pullcreds username; error was: %v", err)
		}

		password, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.PullCreds.Password, opHandle)
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
