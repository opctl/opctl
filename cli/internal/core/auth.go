package core

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/core/auth"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Auther exposes the "auth" sub command
type Auther interface {
	Auth() auth.Auth
}

// newAuther returns an initialized "auth" sub command
func newAuther(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
	nodeProvider nodeprovider.NodeProvider,
) Auther {
	return _auther{
		auth: auth.New(
			cliExiter,
			dataResolver,
			nodeProvider,
		),
	}
}

type _auther struct {
	auth auth.Auth
}

func (ivkr _auther) Auth() auth.Auth {
	return ivkr.auth
}
