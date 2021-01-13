package core

import (
	"fmt"
	"path/filepath"

	"net/url"

	"strings"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/skratchdot/open-golang/open"
)

// UIer exposes the "ui" command
type UIer interface {
	UI(
		mountRef string,
	) error
}

// newUIer returns an initialized "ui" command
func newUIer(
	dataResolver dataresolver.DataResolver,
	core node.OpNode,
) UIer {
	return _uier{
		dataResolver: dataResolver,
		core:         core,
	}
}

type _uier struct {
	dataResolver dataresolver.DataResolver
	core         node.OpNode
}

func (ivkr _uier) UI(
	mountRef string,
) error {
	var resolvedMount string
	var err error
	if strings.HasPrefix(mountRef, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRef)
		if nil != err {
			return err
		}
	} else {
		// otherwise use same resolution as run
		mountHandle, err := ivkr.dataResolver.Resolve(
			mountRef,
			nil,
		)
		if nil != err {
			return err
		}
		resolvedMount = mountHandle.Ref()
	}

	webUIURL := fmt.Sprintf("http://localhost:42224?mount=%s", url.QueryEscape(resolvedMount))

	if err := open.Run(webUIURL); err != nil {
		return err
	}

	return nil
}
