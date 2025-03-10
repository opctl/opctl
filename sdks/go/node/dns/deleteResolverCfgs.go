package dns

import (
	"context"
	"strings"

	"github.com/goodhosts/hostsfile"
)

// DeleteResolverCfgs we've made to the OS
func DeleteResolverCfgs(
	ctx context.Context,
) error {
	h, err := hostsfile.NewHosts()
	if err != nil {
		return err
	}

	for _, l := range h.Lines {
		if strings.Contains(l.Comment, managedByOpctlComment) {
			h.Remove(l.IP, l.Hosts...)
		}
	}

	return h.Flush()
}
