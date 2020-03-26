package op

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	// Interpret interprets an SCGOpCall into a DCGOpCall
	Interpret(
		scope map[string]*model.Value,
		scgOpCall *model.SCGOpCall,
		opID string,
		parentOpPath string,
		rootOpID string,
	) (*model.DCGOpCall, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	dataDirPath string,
) Interpreter {
	return _interpreter{
		dcgScratchDir:       filepath.Join(dataDirPath, "dcg"),
		data:                data.New(),
		dataCachePath:       filepath.Join(dataDirPath, "ops"),
		dirInterpreter:      dir.NewInterpreter(),
		inputsInterpreter:   inputs.NewInterpreter(),
		opFileGetter:        opfile.NewGetter(),
		stringInterpreter:   str.NewInterpreter(),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}
}

type _interpreter struct {
	dcgScratchDir       string
	data                data.Data
	dataCachePath       string
	dirInterpreter      dir.Interpreter
	inputsInterpreter   inputs.Interpreter
	opFileGetter        opfile.Getter
	stringInterpreter   str.Interpreter
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
	opID string,
	parentOpPath string,
	rootOpID string,
) (*model.DCGOpCall, error) {

	var pkgPullCreds *model.PullCreds
	if scgPullCreds := scgOpCall.PullCreds; nil != scgPullCreds {
		pkgPullCreds = &model.PullCreds{}
		var err error
		interpretdUsername, err := itp.stringInterpreter.Interpret(scope, scgPullCreds.Username)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *interpretdUsername.String

		interpretdPassword, err := itp.stringInterpreter.Interpret(scope, scgPullCreds.Password)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Password = *interpretdPassword.String
	}

	var opPath string
	if regexp.MustCompile("^\\$\\(.+\\)$").MatchString(scgOpCall.Ref) {
		// attempt to process as a variable reference since its variable reference like.
		dirValue, err := itp.dirInterpreter.Interpret(
			scope,
			scgOpCall.Ref,
		)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image src; error was: %v", err)
		}
		opPath = *dirValue.Dir
	} else {
		opHandle, err := itp.data.Resolve(
			context.TODO(),
			scgOpCall.Ref,
			itp.data.NewFSProvider(filepath.Dir(parentOpPath)),
			itp.data.NewGitProvider(itp.dataCachePath, pkgPullCreds),
		)
		if nil != err {
			return nil, err
		}
		opPath = *opHandle.Path()
	}

	opFile, err := itp.opFileGetter.Get(
		context.TODO(),
		opPath,
	)
	if nil != err {
		return nil, err
	}

	childCallID, err := itp.uniqueStringFactory.Construct()
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: model.DCGBaseCall{
			RootOpID: rootOpID,
			OpPath:   opPath,
		},
		OpID:         opID,
		ChildCallID:  childCallID,
		ChildCallSCG: opFile.Run,
	}

	dcgOpCall.Inputs, err = itp.inputsInterpreter.Interpret(
		scgOpCall.Inputs,
		opFile.Inputs,
		opPath,
		scope,
		filepath.Join(itp.dcgScratchDir, opID),
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret call to %v; error was: %v", scgOpCall.Ref, err)
	}

	return dcgOpCall, nil
}
