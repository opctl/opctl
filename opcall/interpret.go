package opcall

import (
	"bytes"
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
		pkgPullCreds.Username, err = oc.string.Interpret(scope, scgPullCreds.Username, parentPkgHandle)
		if nil != err {
			return nil, err
		}

		pkgPullCreds.Password, err = oc.string.Interpret(scope, scgPullCreds.Password, parentPkgHandle)
		if nil != err {
			return nil, err
		}
	}

	pkgHandle, err := oc.pkg.Resolve(
		scgOpCall.Pkg.Ref,
		oc.pkg.NewFSProvider(filepath.Dir(parentPkgHandle.Ref())),
		oc.pkg.NewGitProvider(oc.pkgCachePath, pkgPullCreds),
	)
	if nil != err {
		return nil, err
	}

	opManifest, err := oc.pkg.GetManifest(pkgHandle)
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId:  rootOpId,
			PkgHandle: pkgHandle,
		},
		OpId:         opId,
		ChildCallId:  oc.uuid.NewV4().String(),
		ChildCallSCG: opManifest.Run,
	}

	var argErrors []error
	dcgOpCall.Inputs, argErrors = oc.inputs.Interpret(
		scgOpCall.Inputs,
		opManifest.Inputs,
		parentPkgHandle,
		pkgHandle.Ref(),
		scope,
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
