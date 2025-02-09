package resolvercfg

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
)

var resolverPrefix = "opctl_"
var resolverDir = "/etc/resolver"

// Apply to the current system
func Apply(
	ctx context.Context,
	domain,
	nsIPAddress string,
) error {
	var buf bytes.Buffer
	buf.WriteString(
		"# do not edit; managed by opctl\n",
	)

	buf.WriteString(
		fmt.Sprintf(
			"domain %s\n",
			domain,
		),
	)

	buf.WriteString(
		fmt.Sprintf(
			"nameserver %s\n",
			nsIPAddress,
		),
	)

	buf.WriteString(
		"nameserver 8.8.8.8\n",
	)

	if err := os.MkdirAll(
		resolverDir,
		0755,
	); err != nil {
		return err
	}

	if err := os.WriteFile(
		path.Join(
			resolverDir,
			resolverPrefix+domain,
		),
		buf.Bytes(),
		0644,
	); err != nil {
		return err
	}

	return clearCaches(ctx)
}
