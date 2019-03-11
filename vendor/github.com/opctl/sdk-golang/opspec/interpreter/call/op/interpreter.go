package op

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/inputs"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
	"github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type Interpreter interface {
	// Interpret interprets an SCGOpCall into a DCGOpCall
	Interpret(
		scope map[string]*model.Value,
		scgOpCall *model.SCGOpCall,
		opID string,
		parentOpHandle model.DataHandle,
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
		inputsInterpreter:   inputs.NewInterpreter(),
		opOpDotYmlGetter:    dotyml.NewGetter(),
		stringInterpreter:   stringPkg.NewInterpreter(),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}
}

type _interpreter struct {
	dcgScratchDir       string
	data                data.Data
	dataCachePath       string
	inputsInterpreter   inputs.Interpreter
	opOpDotYmlGetter    dotyml.Getter
	stringInterpreter   stringPkg.Interpreter
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
	opID string,
	parentOpHandle model.DataHandle,
	rootOpID string,
) (*model.DCGOpCall, error) {

	var pkgPullCreds *model.PullCreds
	if scgPullCreds := scgOpCall.PullCreds; nil != scgPullCreds {
		pkgPullCreds = &model.PullCreds{}
		var err error
		evaluatedUsername, err := itp.stringInterpreter.Interpret(scope, scgPullCreds.Username, parentOpHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *evaluatedUsername.String

		evaluatedPassword, err := itp.stringInterpreter.Interpret(scope, scgPullCreds.Password, parentOpHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Password = *evaluatedPassword.String
	}

	parentOpDirPath := parentOpHandle.Path()
	opHandle, err := itp.data.Resolve(
		context.TODO(),
		scgOpCall.Ref,
		itp.data.NewFSProvider(filepath.Dir(*parentOpDirPath)),
		itp.data.NewGitProvider(itp.dataCachePath, pkgPullCreds),
	)
	if nil != err {
		return nil, err
	}

	opDotYml, err := itp.opOpDotYmlGetter.Get(
		context.TODO(),
		opHandle,
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
			OpHandle: opHandle,
		},
		OpID:         opID,
		ChildCallID:  childCallID,
		ChildCallSCG: opDotYml.Run,
	}

	dcgOpCall.Inputs, err = itp.inputsInterpreter.Interpret(
		scgOpCall.Inputs,
		opDotYml.Inputs,
		parentOpHandle,
		*opHandle.Path(),
		scope,
		filepath.Join(itp.dcgScratchDir, opID),
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret call to %v; error was: %v", dcgOpCall.OpHandle.Ref(), err)
	}

	return dcgOpCall, nil
}
