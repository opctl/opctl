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
	nsIPAddress,
	nsPort string,
) error {
	if nsPort != "53" {
		return fmt.Errorf(
			"%s:%s invalid; linux DNS resolver requires servers listen on port 53",
			nsIPAddress,
			nsPort,
		)
	}

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
