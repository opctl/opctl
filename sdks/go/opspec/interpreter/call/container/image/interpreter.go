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
		scgContainerCallImage *model.SCGContainerCallImage,
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
	scgContainerCallImage *model.SCGContainerCallImage,
	scratchDir string,
) (*model.DCGContainerCallImage, error) {

	if nil == scgContainerCallImage {
		return nil, fmt.Errorf("image required")
	}

	// try to interpret as dir
	src, err := itp.dirInterpreter.Interpret(
		scope,
		scgContainerCallImage.Ref,
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
		scgContainerCallImage.Ref,
	)
	if nil != err {
		return nil, err
	}

	dcgContainerCallImage.Ref = ref.String

	if nil != scgContainerCallImage.PullCreds {
		username, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.PullCreds.Username)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds username; error was: %v", err)
		}

		password, err := itp.stringInterpreter.Interpret(scope, scgContainerCallImage.PullCreds.Password)
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
