package dataresolver

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ DataResolver

import (
	"context"
	"fmt"
	"net/url"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/apireachabilityensurer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
)

// DataResolver resolves packages
type DataResolver interface {
	Resolve(
		dataRef string,
		pullCreds *model.PullCreds,
	) model.DataHandle
}

func New(
	cliExiter cliexiter.CliExiter,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeURL url.URL,
) DataResolver {
	return _dataResolver{
		apiReachabilityEnsurer: apireachabilityensurer.New(cliExiter),
		cliExiter:              cliExiter,
		cliParamSatisfier:      cliParamSatisfier,
		data:                   data.New(),
		nodeURL:                nodeURL,
		os:                     ios.New(),
	}
}

type _dataResolver struct {
	cliExiter              cliexiter.CliExiter
	cliParamSatisfier      cliparamsatisfier.CLIParamSatisfier
	data                   data.Data
	nodeURL                url.URL
	apiReachabilityEnsurer apireachabilityensurer.APIReachabilityEnsurer
	os                     ios.IOS
}

func (dtr _dataResolver) Resolve(
	dataRef string,
	pullCreds *model.PullCreds,
) model.DataHandle {

	// ensure node reachable
	dtr.apiReachabilityEnsurer.Ensure()

	cwd, err := dtr.os.Getwd()
	if nil != err {
		dtr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return nil // support fake exiter
	}

	fsProvider := dtr.data.NewFSProvider(cwd)

	for {
		opDirHandle, err := dtr.data.Resolve(
			context.TODO(),
			dataRef,
			fsProvider,
			dtr.data.NewNodeProvider(
				dtr.nodeURL,
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
			argMap := dtr.cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					dtr.cliParamSatisfier.NewCliPromptInputSrc(credsPromptInputs),
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
			dtr.cliExiter.Exit(
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
