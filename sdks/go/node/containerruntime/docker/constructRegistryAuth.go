package docker

import (
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types/registry"
)

func constructRegistryAuth(
	identity string,
	secret string,
) (registryAuth string, err error) {

	// EncodeAuthToBase64 serializes the auth configuration as JSON base64 payload
	buf, err := json.Marshal(
		registry.AuthConfig{
			Username: identity,
			Password: secret,
		})
	if err != nil {
		return
	}
	registryAuth = base64.URLEncoding.EncodeToString(buf)

	return
}
