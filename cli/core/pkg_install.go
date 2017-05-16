package core

import (
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

func (this _core) PkgInstall(
	path,
	pkgRef,
	username,
	password string,
) {

	parsedPkgRef, err := this.pkg.ParseRef(pkgRef)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {
		err := this.pkg.Pull(
			path,
			parsedPkgRef,
			&pkg.PullOpts{
				Username: username,
				Password: password,
			},
		)

		_, isAuthError := err.(pkg.ErrAuthenticationFailed)

		switch {
		case nil == err:
			return
		case isAuthError:
			// auth errors can be fixed by supplying correct creds so don't give up; prompt
			argMap := this.cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					cliparamsatisfier.NewCliPromptInputSrc(credsPromptInputs),
				),
				credsPromptInputs,
			)

			// save providedArgs & re-attempt
			username = *(argMap[usernameInputName].String)
			password = *(argMap[passwordInputName].String)
			continue
		default:
			// uncorrectable error.. give up
			this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
			return // support fake exiter
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
