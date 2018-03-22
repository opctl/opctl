package opcall

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/inputs"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"path/filepath"
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
	rootFSPath string,
) Interpreter {
	return _interpreter{
		dcgScratchDir:       filepath.Join(rootFSPath, "dcg"),
		expression:          expression.New(),
		data:                data.New(),
		opOpDotYmlGetter:    dotyml.NewGetter(),
		dataCachePath:       filepath.Join(rootFSPath, "pkgs"),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		inputsInterpreter:   inputs.NewInterpreter(),
	}
}

type _interpreter struct {
	dcgScratchDir       string
	expression          expression.Expression
	data                data.Data
	opOpDotYmlGetter    dotyml.Getter
	dataCachePath       string
	uniqueStringFactory uniquestring.UniqueStringFactory
	inputsInterpreter   inputs.Interpreter
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
		evaluatedUsername, err := itp.expression.EvalToString(scope, scgPullCreds.Username, parentOpHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *evaluatedUsername.String

		evaluatedPassword, err := itp.expression.EvalToString(scope, scgPullCreds.Password, parentOpHandle)
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
