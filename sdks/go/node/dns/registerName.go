package dns

import (
	"context"
	"strings"

	"github.com/opctl/opctl/sdks/go/node/dns/internal/resolvercfg"
)

func RegisterName(
	ctx context.Context,
	name,
	ipAddress string,
) error {
	ips := ipsByHostname[name]
	if ips == nil {
		ipsByHostname[name] = map[string]interface{}{}
	}

	ipsByHostname[name][ipAddress] = nil

	nameParts := strings.Split(name, ".")

	return resolvercfg.Apply(
		ctx,
		nameParts[len(nameParts)-1],
		nsIPAddress,
		nsPort,
	)
}
