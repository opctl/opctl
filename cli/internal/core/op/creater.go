package op

import (
	"path/filepath"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Creater exposes the "op create" sub command
type Creater interface {
	Create(
		path string,
		description string,
		name string,
	)
}

// newCreater returns an initialized "op create" sub command
func newCreater(
	cliExiter cliexiter.CliExiter,
) Creater {
	return _creater{
		cliExiter: cliExiter,
	}
}

type _creater struct {
	cliExiter cliexiter.CliExiter
}

func (ivkr _creater) Create(
	path string,
	description string,
	name string,
) {
	if err := opspec.Create(
		filepath.Join(path, name),
		name,
		description,
	); nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
