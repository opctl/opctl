package containerruntime

import "os"

// proxyEnvVarNames are the proxy related environment variables opctl propagates
// from its own process environment into containers it runs. Both the uppercase
// and lowercase forms are propagated, since programs differ in which form they
// honor.
var proxyEnvVarNames = []string{
	"HTTP_PROXY", "http_proxy",
	"HTTPS_PROXY", "https_proxy",
	"NO_PROXY", "no_proxy",
	"ALL_PROXY", "all_proxy",
}

// ProxyEnvVars returns proxy related environment variables present in the opctl
// node's own process environment which are not already set in opEnvVars.
//
// This lets containers reach the network on hosts where the only egress route
// is an HTTP/HTTPS forward proxy (e.g. CI runners behind a proxy). Because
// opctl talks to the container runtime via its SDK rather than the docker CLI,
// the `~/.docker/config.json` proxies convenience does not apply, so opctl must
// inject these itself.
//
// Values explicitly provided by an op are never overridden; only proxy
// variables absent from opEnvVars are returned.
func ProxyEnvVars(opEnvVars map[string]string) map[string]string {
	propagated := map[string]string{}
	for _, name := range proxyEnvVarNames {
		if _, ok := opEnvVars[name]; ok {
			// never clobber an op-declared value
			continue
		}
		if value, ok := os.LookupEnv(name); ok {
			propagated[name] = value
		}
	}
	return propagated
}
