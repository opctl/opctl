package dataresolver

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/data/node"
	"github.com/opctl/opctl/sdks/go/model"
)

// DataResolver resolves packages
//counterfeiter:generate -o fakes/dataResolver.go . DataResolver
type DataResolver interface {
	Resolve(
		dataRef string,
		pullCreds *model.Creds,
	) (model.DataHandle, error)
}

func New(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeProvider nodeprovider.NodeProvider,
) DataResolver {
	return _dataResolver{
		cliParamSatisfier: cliParamSatisfier,
		nodeProvider:      nodeProvider,
		os:                ios.New(),
	}
}

type _dataResolver struct {
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	nodeProvider      nodeprovider.NodeProvider
	os                ios.IOS
}

func (dtr _dataResolver) Resolve(
	dataRef string,
	pullCreds *model.Creds,
) (model.DataHandle, error) {
	cwd, err := dtr.os.Getwd()
	if nil != err {
		return nil, err
	}

	fsProvider := fs.New(
		filepath.Join(cwd, ".opspec"),
		cwd,
	)

	nodeHandle, err := dtr.nodeProvider.CreateNodeIfNotExists()
	if nil != err {
		return nil, err
	}

	for {
		opDirHandle, err := data.Resolve(
			context.TODO(),
			dataRef,
			fsProvider,
			node.New(
				nodeHandle.APIClient(),
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
			return opDirHandle, err
		case isAuthError:
			// auth errors can be fixed by supplying correct creds so don't give up; prompt
			argMap, err := dtr.cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					dtr.cliParamSatisfier.NewCliPromptInputSrc(credsPromptInputs),
				),
				credsPromptInputs,
			)
			if nil != err {
				return nil, err
			}

			// save providedArgs & re-attempt
			pullCreds = &model.Creds{
				Username: *(argMap[usernameInputName].String),
				Password: *(argMap[passwordInputName].String),
			}
			continue
		default:
			// uncorrectable error.. give up
			return nil, fmt.Errorf("Unable to resolve pkg '%v'; error was %v", dataRef, err.Error())
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
