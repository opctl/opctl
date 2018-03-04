package core

//go:generate counterfeiter -o ./fakePkgResolver.go --fake-name fakePkgResolver ./ pkgResolver

import (
	"context"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"net/url"
)

// pkgResolver resolves packages
type pkgResolver interface {
	Resolve(
		pkgRef string,
		pullCreds *model.PullCreds,
	) model.PkgHandle
}

func newPkgResolver(
	cliExiter cliexiter.CliExiter,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeURL url.URL,
) pkgResolver {
	return _pkgResolver{
		cliExiter:               cliExiter,
		cliParamSatisfier:       cliParamSatisfier,
		nodeReachabilityEnsurer: newNodeReachabilityEnsurer(cliExiter),
		nodeURL:                 nodeURL,
		pkg:                     pkg.New(),
		os:                      ios.New(),
	}
}

type _pkgResolver struct {
	cliExiter               cliexiter.CliExiter
	cliParamSatisfier       cliparamsatisfier.CLIParamSatisfier
	nodeURL                 url.URL
	nodeReachabilityEnsurer nodeReachabilityEnsurer
	os                      ios.IOS
	pkg                     pkg.Pkg
}

func (this _pkgResolver) Resolve(
	pkgRef string,
	pullCreds *model.PullCreds,
) model.PkgHandle {

	// ensure node reachable
	this.nodeReachabilityEnsurer.EnsureNodeReachable()

	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return nil // support fake exiter
	}

	fsProvider := this.pkg.NewFSProvider(cwd)

	for {
		pkgHandle, err := this.pkg.Resolve(
			context.TODO(),
			pkgRef,
			fsProvider,
			this.pkg.NewNodeProvider(
				this.nodeURL,
				pullCreds,
			),
		)

		var isAuthError bool
		switch err.(type) {
		case model.ErrPkgPullAuthorization:
			isAuthError = true
		case model.ErrPkgPullAuthentication:
			isAuthError = true
		}

		switch {
		case nil == err:
			return pkgHandle
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
					Message: fmt.Sprintf("Unable to resolve pkg '%v'; error was %v", pkgRef, err.Error()),
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
				Constraints: &model.StringConstraints{
					MinLength: 1,
				},
			},
		},
		passwordInputName: {
			String: &model.StringParam{
				Description: "password used to auth w/ the pkg source",
				Constraints: &model.StringConstraints{
					MinLength: 1,
				},
				IsSecret: true,
			},
		},
	}
)
