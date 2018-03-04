package opcall

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (oc _OpCall) Interpret(
	scope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
	opId string,
	parentPkgHandle model.PkgHandle,
	rootOpId string,
) (*model.DCGOpCall, error) {

	var pkgPullCreds *model.PullCreds
	if scgPullCreds := scgOpCall.Pkg.PullCreds; nil != scgPullCreds {
		pkgPullCreds = &model.PullCreds{}
		var err error
		evaluatedUsername, err := oc.expression.EvalToString(scope, scgPullCreds.Username, parentPkgHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *evaluatedUsername.String

		evaluatedPassword, err := oc.expression.EvalToString(scope, scgPullCreds.Password, parentPkgHandle)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Password = *evaluatedPassword.String
	}

	parentPkgPath := parentPkgHandle.Path()
	pkgHandle, err := oc.pkg.Resolve(
		context.TODO(),
		scgOpCall.Pkg.Ref,
		oc.pkg.NewFSProvider(filepath.Dir(*parentPkgPath)),
		oc.pkg.NewGitProvider(oc.pkgCachePath, pkgPullCreds),
	)
	if nil != err {
		return nil, err
	}

	opManifest, err := oc.pkg.GetManifest(pkgHandle)
	if nil != err {
		return nil, err
	}

	childCallId, err := oc.uniqueStringFactory.Construct()
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: model.DCGBaseCall{
			RootOpId:  rootOpId,
			PkgHandle: pkgHandle,
		},
		OpId:         opId,
		ChildCallId:  childCallId,
		ChildCallSCG: opManifest.Run,
	}

	var argErrors []error
	dcgOpCall.Inputs, argErrors = oc.inputs.Interpret(
		scgOpCall.Inputs,
		opManifest.Inputs,
		parentPkgHandle,
		*pkgHandle.Path(),
		scope,
		filepath.Join(oc.dcgScratchDir, opId),
	)
	if len(argErrors) > 0 {
		messageBuffer := bytes.NewBufferString("")
		for _, validationError := range argErrors {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		messageBuffer.WriteString(`
`)
		return nil, fmt.Errorf(`
-
  error(s) occurred interpreting call to %v:
%v
-`, dcgOpCall.PkgHandle.Ref(), messageBuffer.String())
	}

	return dcgOpCall, nil
}
