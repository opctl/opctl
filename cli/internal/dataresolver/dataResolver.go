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
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
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
	if err != nil {
		return nil, err
	}

	fsProvider := fs.New(
		filepath.Join(cwd, ".opspec"),
		cwd,
	)

	domain := strings.Split(dataRef, "/")[0]

	passwordDescription := fmt.Sprintf("Password for %s.", domain)

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
			fsProvider,
			dataNode.New(
				dtr.node,
				pullCreds,
			),
		)

		if err == nil {
			return opDirHandle, nil
		}

		if model.IsAuthError(err) && !reattemptedAuth {
			if domain == "github.com" {
				t, err := authenticateWithGitHub(ctx)
				if err != nil {
					return nil, err
				}

				pullCreds = &model.Creds{
					Username: "n/a",
					Password: t.AccessToken,
				}
				reattemptedAuth = true
				continue
			}
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

		return nil, err
	}
}

const (
	usernameInputName = "username"
	passwordInputName = "password"
)

func authenticateWithGitHub(
	ctx context.Context,
) (*oauth2.Token, error) {
	config := &oauth2.Config{
		// safe to embed per https://github.com/cli/cli/blob/7d558edd12d6b4101a0b3d3df3c689385a1c830a/internal/authflow/flow.go#L21-L25
		ClientID:     "Ov23liVQSZYxmPWGLNGe",
		ClientSecret: "48d696d29d06eab46ef562535d64775f852e343b",
		Endpoint:     endpoints.GitHub,
		Scopes:       []string{"repo"},
	}

	// Initiate device flow
	deviceCode, err := config.DeviceAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device code: %w", err)
	}

	// Display instructions to the user
	fmt.Printf("Please visit %s and enter the code: %s\n", deviceCode.VerificationURI, deviceCode.UserCode)
	fmt.Println("Waiting for authorization...")

	err = open.Run(
		deviceCode.VerificationURI,
	)
	if err != nil {
		return nil, err
	}

	// Poll for the access token
	return config.DeviceAccessToken(ctx, deviceCode)
}
