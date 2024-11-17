package dataresolver

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/opspath"
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
	opDirRefs, err := opspath.Get(
		ctx,
		dataRef,
		dtr.node,
	)
	if err != nil {
		return nil, err
	}

	dataRefs := []string{
		dataRef,
	}

	if !filepath.IsAbs(dataRef) {
		for _, opDirRef := range opDirRefs {
			dataRefs = append(dataRefs, opDirRef+"/"+dataRef)
		}
	}

dataRefLoop:
	for _, dataRef := range dataRefs {

		domain := strings.Split(dataRef, "/")[0]

		passwordDescription := fmt.Sprintf("Password for %s.", domain)
		if domain == "github.com" {
			// customize github.com password description...
			passwordDescription = "Personal access token for github.com with 'Repo' permissions."
		}

		credsPromptInputs := map[string]*model.ParamSpec{
			usernameInputName: {
				String: &model.StringParamSpec{
					Description: fmt.Sprintf("Username for %s.", domain),
					Constraints: map[string]interface{}{
						"MinLength": 1,
					},
				},
			},
			passwordInputName: {
				String: &model.StringParamSpec{
					Description: passwordDescription,
					Constraints: map[string]interface{}{
						"MinLength": 1,
					},
					IsSecret: true,
				},
			},
		}

		reattemptedAuth := false

		for {
			opDirHandle, err := data.Resolve(
				ctx,
				dataRef,
				fs.New(),
				dataNode.New(
					dtr.node,
					pullCreds,
				),
			)

			if err == nil {
				return opDirHandle, nil
			}

			if model.IsAuthError(err) && !reattemptedAuth {
				// auth errors can be fixed by supplying correct creds so don't give up; prompt
				cliPromptInputSrc := dtr.cliParamSatisfier.NewCliPromptInputSrc(credsPromptInputs)

				argMap, err := dtr.cliParamSatisfier.Satisfy(
					cliparamsatisfier.NewInputSourcer(cliPromptInputSrc),
					credsPromptInputs,
				)
				if err != nil {
					return nil, err
				}

				// save providedArgs & re-attempt
				pullCreds = &model.Creds{
					Username: *(argMap[usernameInputName].String),
					Password: *(argMap[passwordInputName].String),
				}
				reattemptedAuth = true
				continue
			}
			if errors.As(err, &model.ErrDataUnableToResolve{}) {
				continue dataRefLoop
			}

			return nil, err
		}
	}

	return nil, model.ErrDataUnableToResolve{}
}

const (
	usernameInputName = "username"
	passwordInputName = "password"
)
