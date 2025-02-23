package dns

import (
	"context"
	"strings"

	"github.com/opctl/opctl/sdks/go/node/dns/internal/resolvercfg"
)

func RegisterNames(
	ctx context.Context,
	names []string,
	ipAddress string,
) error {
	for _, name := range names {
		ips := ipsByHostname[name]
		if ips == nil {
			ipsByHostname[name] = map[string]interface{}{}
		}

		ipsByHostname[name][ipAddress] = nil

		nameParts := strings.Split(name, ".")

		err := resolvercfg.Apply(
			ctx,
			nameParts[len(nameParts)-1],
			nsIPAddress,
			nsPort,
		)
		if err != nil {
			return nil
		}
	}

	return nil
}
