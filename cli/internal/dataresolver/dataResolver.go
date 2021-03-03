package dataresolver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	dataNode "github.com/opctl/opctl/sdks/go/data/node"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
)

// DataResolver resolves packages
type DataResolver interface {
	Resolve(
		ctx context.Context,
		dataRef string,
		pullCreds *model.Creds,
	) (model.DataHandle, error)
}

func New(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	node node.Node,
) DataResolver {
	return _dataResolver{
		cliParamSatisfier: cliParamSatisfier,
		node:              node,
	}
}

type _dataResolver struct {
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	node              node.Node
}

func (dtr _dataResolver) Resolve(
	ctx context.Context,
	dataRef string,
	pullCreds *model.Creds,
) (model.DataHandle, error) {
	cwd, err := os.Getwd()
	if nil != err {
		return nil, err
	}

	fsProvider := fs.New(
		filepath.Join(cwd, ".opspec"),
		cwd,
	)

	domain := strings.Split(dataRef, "/")[0]

	passwordDescription := fmt.Sprintf("Password for %s.", domain)
	if domain == "github.com" {
		// customize github.com password description...
		passwordDescription = "Personal access token for github.com with 'Repo' permissions."
	}

	credsPromptInputs := map[string]*model.Param{
		usernameInputName: {
			String: &model.StringParam{
				Description: fmt.Sprintf("Username for %s.", domain),
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
			},
		},
		passwordInputName: {
			String: &model.StringParam{
				Description: passwordDescription,
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
				IsSecret: true,
			},
		},
	}

	for {
		opDirHandle, err := data.Resolve(
			ctx,
			dataRef,
			fsProvider,
			dataNode.New(
				dtr.node,
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
			return nil, fmt.Errorf("Unable to resolve '%v'; error was %v", dataRef, err.Error())
		}

	}

}

const (
	usernameInputName = "username"
	passwordInputName = "password"
)
