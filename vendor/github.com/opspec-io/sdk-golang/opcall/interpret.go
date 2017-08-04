package opcall

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (this _OpCall) Interpret(
	scope map[string]*model.Value,
	scgOpCall *model.SCGOpCall,
	opId string,
	parentPkgHandle model.PkgHandle,
	rootOpId string,
) (*model.DCGOpCall, error) {

	var username, password string
	if scgPullCreds := scgOpCall.Pkg.PullCreds; nil != scgPullCreds {
		username = this.interpolater.Interpolate(scgPullCreds.Username, scope)
		password = this.interpolater.Interpolate(scgPullCreds.Password, scope)
	}

	pkgHandle, err := this.pkg.Resolve(
		scgOpCall.Pkg.Ref,
		this.pkg.NewFSProvider(filepath.Dir(parentPkgHandle.Ref())),
		this.pkg.NewGitProvider(this.pkgCachePath, &model.PullCreds{Username: username, Password: password}),
	)
	if nil != err {
		return nil, err
	}

	opManifest, err := this.pkg.GetManifest(pkgHandle)
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId:  rootOpId,
			PkgHandle: pkgHandle,
		},
		OpId:         opId,
		ChildCallId:  this.uuid.NewV4().String(),
		ChildCallSCG: opManifest.Run,
	}

	var argErrors []error
	dcgOpCall.Inputs, argErrors = this.inputs.Interpret(
		scgOpCall.Inputs,
		opManifest.Inputs,
		parentPkgHandle.Ref(),
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
