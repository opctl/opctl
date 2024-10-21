package dns

import (
	"bytes"
	"fmt"
	"os"
	"path"
)

var etcResolverPath = "/etc/resolver"

func setResolver(
	domain string,
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
		"nameserver 127.0.0.1\n",
	)
	buf.WriteString(
		"nameserver 8.8.8.8\n",
	)

	if err := os.MkdirAll(
		etcResolverPath,
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		path.Join(
			etcResolverPath,
			fmt.Sprintf(
				"opctl_%s",
				domain,
			),
		),
		buf.Bytes(),
		0644,
	)
}
