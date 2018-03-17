package opcall

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/inputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"path/filepath"
)

type Interpreter interface {
	// Interpret interprets an SCGOpCall into a DCGOpCall
	Interpret(
		scope map[string]*model.Value,
		scgOpCall *model.SCGOpCall,
		opId string,
		parentOpDirHandle model.DataHandle,
		rootOpId string,
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
		pkg:                 pkg.New(),
		pkgCachePath:        filepath.Join(rootFSPath, "pkgs"),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		inputsInterpreter:   inputs.NewInterpreter(),
	}
}

type _interpreter struct {
	dcgScratchDir       string
	expression          expression.Expression
	data                data.Data
	pkg                 pkg.Pkg
	pkgCachePath        string
	uniqueStringFactory uniquestring.UniqueStringFactory
	inputsInterpreter   inputs.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
	opId string,
	parentOpDirHandle model.DataHandle,
	rootOpId string,
) (*model.DCGOpCall, error) {

	var pkgPullCreds *model.PullCreds
	if scgPullCreds := scgOpCall.Pkg.PullCreds; nil != scgPullCreds {
		pkgPullCreds = &model.PullCreds{}
		var err error
		evaluatedUsername, err := itp.expression.EvalToString(scope, scgPullCreds.Username, parentOpDirHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *evaluatedUsername.String

		evaluatedPassword, err := itp.expression.EvalToString(scope, scgPullCreds.Password, parentOpDirHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Password = *evaluatedPassword.String
	}

	parentOpDirPath := parentOpDirHandle.Path()
	opDirHandle, err := itp.data.Resolve(
		context.TODO(),
		scgOpCall.Pkg.Ref,
		itp.data.NewFSProvider(filepath.Dir(*parentOpDirPath)),
		itp.data.NewGitProvider(itp.pkgCachePath, pkgPullCreds),
	)
	if nil != err {
		return nil, err
	}

	opManifest, err := itp.pkg.GetManifest(opDirHandle)
	if nil != err {
		return nil, err
	}

	childCallId, err := itp.uniqueStringFactory.Construct()
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: model.DCGBaseCall{
			RootOpId:   rootOpId,
			DataHandle: opDirHandle,
		},
		OpId:         opId,
		ChildCallId:  childCallId,
		ChildCallSCG: opManifest.Run,
	}

	dcgOpCall.Inputs, err = itp.inputsInterpreter.Interpret(
		scgOpCall.Inputs,
		opManifest.Inputs,
		parentOpDirHandle,
		*opDirHandle.Path(),
		scope,
		filepath.Join(itp.dcgScratchDir, opId),
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret call to %v; error was: %v", dcgOpCall.DataHandle.Ref(), err)
	}

	return dcgOpCall, nil
}
