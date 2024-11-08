package dns

import (
	"context"

	"github.com/opctl/opctl/sdks/go/node/dns/internal/resolvercfg"
)

// DeleteResolverCfgs we've made to the OS
func DeleteResolverCfgs(
	ctx context.Context,
) error {
	return resolvercfg.Delete(ctx)
}
