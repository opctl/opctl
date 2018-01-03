package core

import (
	"errors"

	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (string, error) {
	if nil == req.Pkg {
		return "", errors.New("pkg required")
	}

	var pullCreds *model.PullCreds
	if nil != req.Pkg.PullCreds {
		pullCreds = &model.PullCreds{
			Username: req.Pkg.PullCreds.Username,
			Password: req.Pkg.PullCreds.Password,
		}
	}

	pkgHandle, err := this.pkg.Resolve(
		req.Pkg.Ref,
		this.pkg.NewFSProvider(),
		this.pkg.NewGitProvider(this.pkgCachePath, pullCreds),
	)
	if nil != err {
		return "", err
	}

	opId := this.uniqueStringFactory.Construct()

	// construct scgOpCall
	scgOpCall := &model.SCGOpCall{
		Pkg: &model.SCGOpCallPkg{
			Ref: pkgHandle.Ref(),
		},
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}
	for name := range req.Args {
		// implicitly bind
		scgOpCall.Inputs[name] = ""
	}

	pkgManifest, err := this.pkg.GetManifest(pkgHandle)
	if nil != err {
		return "", err
	}
	for name := range pkgManifest.Outputs {
		// implicitly bind
		scgOpCall.Outputs[name] = ""
	}

	inboundScope := map[string]*model.Value{}
	for name, param := range pkgManifest.Inputs {
		if arg, isArg := req.Args[name]; isArg {
			switch {
			case nil != param.Array:
				if arrayArg, isArray := arg.([]interface{}); isArray {
					inboundScope[name] = &model.Value{Array: arrayArg}
				}
			case nil != param.Dir:
				if dirArg, isDir := arg.(string); isDir {
					inboundScope[name] = &model.Value{Dir: &dirArg}
				}
			case nil != param.File:
				if fileArg, isFile := arg.(string); isFile {
					inboundScope[name] = &model.Value{File: &fileArg}
				}
			case nil != param.Number:
				if numberArg, isNumber := arg.(float64); isNumber {
					inboundScope[name] = &model.Value{Number: &numberArg}
				}
			case nil != param.Object:
				if objectArg, isObject := arg.(map[string]interface{}); isObject {
					inboundScope[name] = &model.Value{Object: objectArg}
				}
			case nil != param.Socket:
				if socketArg, isSocket := arg.(string); isSocket {
					inboundScope[name] = &model.Value{Socket: &socketArg}
				}
			case nil != param.String:
				if stringArg, isString := arg.(string); isString {
					inboundScope[name] = &model.Value{String: &stringArg}
				}
			}
		}
	}

	go func() {
		this.opCaller.Call(
			inboundScope,
			opId,
			pkgHandle,
			opId,
			scgOpCall,
		)
	}()

	return opId, nil

}
