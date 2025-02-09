package resolvercfg

import (
	"context"
	"fmt"
	"os"
	"strings"
)

var etcResolvConfPath = "/etc/resolv.conf"

// Apply to the current system
func Apply(
	ctx context.Context,
	domain,
	nsIPAddress string,
) error {
	rc, err := os.ReadFile(etcResolvConfPath)
	if err != nil {
		return err
	}

	nsPrefix := fmt.Sprintf(
		"nameserver %s",
		nsIPAddress,
	)

	rcString := string(rc)
	if !strings.HasPrefix(rcString, nsPrefix) {
		return os.WriteFile(
			etcResolvConfPath,
			[]byte(fmt.Sprintf(
				"%s # do not edit; managed by opctl\n%s",
				nsPrefix,
				rcString,
			),
			),
			0600,
		)
	}
	return nil
}
