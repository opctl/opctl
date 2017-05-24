package opcall

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

func (this _OpCall) Interpret(
	scope map[string]*model.Data,
	scgOpCall *model.SCGOpCall,
	opId string,
	pkgBasePath string,
	rootOpId string,
) (*model.DCGOpCall, error) {

	pkgPath, err := this.getPkgPath(
		pkgBasePath,
		scope,
		scgOpCall.Pkg,
	)
	if nil != err {
		return nil, err
	}

	_pkg, err := this.manifest.Unmarshal(filepath.Join(pkgPath, pkg.OpDotYmlFileName))
	if nil != err {
		return nil, err
	}

	dcgOpCall := &model.DCGOpCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgPath,
		},
		OpId:           opId,
		ChildCallId:    this.uuid.NewV4().String(),
		ChildCallSCG:   _pkg.Run,
		ChildCallScope: map[string]*model.Data{},
	}

	messageBuffer := bytes.NewBufferString("")
	for inputName, scgInputValue := range scgOpCall.Inputs {
		// interpret inputs
		dcgOpCall.ChildCallScope[inputName], err = this.inputInterpreter.Interpret(
			inputName,
			scgInputValue,
			_pkg.Inputs,
			scope,
		)

		if nil != err {
			// aggregate errors
			messageBuffer.WriteString(err.Error())
			err = nil
		}
	}

	if messageBuffer.Len() > 0 {
		return nil, errors.New(messageBuffer.String())
	}

	return dcgOpCall, nil
}

func (this _OpCall) getPkgPath(
	pkgBasePath string,
	scope map[string]*model.Data,
	scgOpCallPkg *model.SCGOpCallPkg,
) (string, error) {
	pkgRef, err := this.pkg.ParseRef(scgOpCallPkg.Ref)
	if nil != err {
		return "", err
	}

	var username, password string
	if scgPullCreds := scgOpCallPkg.PullCreds; nil != scgPullCreds {
		username = this.interpolater.Interpolate(scgPullCreds.Username, scope)
		password = this.interpolater.Interpolate(scgPullCreds.Password, scope)
	}

	pkgPath, ok := this.pkg.Resolve(pkgRef, pkgBasePath, this.pkgCachePath)
	if !ok {
		// pkg not resolved; attempt to pull it
		err := this.pkg.Pull(this.pkgCachePath, pkgRef, &pkg.PullOpts{Username: username, Password: password})
		if nil != err {
			return "", err
		}

		// resolve pulled pkg
		pkgPath, ok = this.pkg.Resolve(pkgRef, pkgBasePath, this.pkgCachePath)
		if !ok {
			return "", fmt.Errorf("Unable to resolve pulled pkg '%v' from '%v'", pkgRef, this.pkgCachePath)
		}
	}
	return pkgPath, nil
}
