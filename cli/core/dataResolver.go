package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeDataResolver.go --fake-name fakeDataResolver ./ dataResolver

import (
	"context"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/cli/util/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
	"net/url"
)

// dataResolver resolves packages
type dataResolver interface {
	Resolve(
		dataRef string,
		pullCreds *model.PullCreds,
	) model.DataHandle
}

func newDataResolver(
	cliExiter cliexiter.CliExiter,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeURL url.URL,
) dataResolver {
	return _dataResolver{
		cliExiter:         cliExiter,
		cliParamSatisfier: cliParamSatisfier,
		data:              data.New(),
		nodeReachabilityEnsurer: newNodeReachabilityEnsurer(cliExiter),
		nodeURL:                 nodeURL,
		os:                      ios.New(),
	}
}

type _dataResolver struct {
	cliExiter               cliexiter.CliExiter
	cliParamSatisfier       cliparamsatisfier.CLIParamSatisfier
	data                    data.Data
	nodeURL                 url.URL
	nodeReachabilityEnsurer nodeReachabilityEnsurer
	os                      ios.IOS
}

func (this _dataResolver) Resolve(
	dataRef string,
	pullCreds *model.PullCreds,
) model.DataHandle {

	// ensure node reachable
	this.nodeReachabilityEnsurer.EnsureNodeReachable()

	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return nil // support fake exiter
	}

	fsProvider := this.data.NewFSProvider(cwd)

	for {
		opDirHandle, err := this.data.Resolve(
			context.TODO(),
			dataRef,
			fsProvider,
			this.data.NewNodeProvider(
				this.nodeURL,
				pullCreds,
			),
		)

		var isAuthError bool
		switch err.(type) {
		case model.ErrDataProviderAuthorization:
			isAuthError = true
		case model.ErrDataProviderAuthentication:
			isAuthError = true
		}

		switch {
		case nil == err:
			return opDirHandle
		case isAuthError:
			// auth errors can be fixed by supplying correct creds so don't give up; prompt
			argMap := this.cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					this.cliParamSatisfier.NewCliPromptInputSrc(credsPromptInputs),
				),
				credsPromptInputs,
			)

			// save providedArgs & re-attempt
			pullCreds = &model.PullCreds{
				Username: *(argMap[usernameInputName].String),
				Password: *(argMap[passwordInputName].String),
			}
			continue
		default:
			// uncorrectable error.. give up
			this.cliExiter.Exit(
				cliexiter.ExitReq{
					Message: fmt.Sprintf("Unable to resolve pkg '%v'; error was %v", dataRef, err.Error()),
					Code:    1,
				},
			)
			return nil // support fake exiter
		}

	}

}

const (
	usernameInputName = "username"
	passwordInputName = "password"
)

var (
	credsPromptInputs = map[string]*model.Param{
		usernameInputName: {
			String: &model.StringParam{
				Description: "username used to auth w/ the pkg source",
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
			},
		},
		passwordInputName: {
			String: &model.StringParam{
				Description: "password used to auth w/ the pkg source",
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
				IsSecret: true,
			},
		},
	}
)
