package core

import (
	"github.com/opctl/opctl/cli/internal/core/auth"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node"
)

// Auther exposes the "auth" sub command
type Auther interface {
	Auth() auth.Auth
}

// newAuther returns an initialized "auth" sub command
func newAuther(
	dataResolver dataresolver.DataResolver,
	opNode node.OpNode,
) Auther {
	return _auther{
		auth: auth.New(
			dataResolver,
			opNode,
		),
	}
}

type _auther struct {
	auth auth.Auth
}

func (ivkr _auther) Auth() auth.Auth {
	return ivkr.auth
}
