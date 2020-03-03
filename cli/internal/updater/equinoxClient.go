package updater

import "github.com/equinox-io/equinox"

// client interface for equinox
//counterfeiter:generate -o internal/fakes/equinoxClient.go . equinoxClient
type equinoxClient interface {
	Check(appID string, opts equinox.Options) (equinox.Response, error)
	Apply(response equinox.Response) error
}

func newEquinoxClient() equinoxClient {
	return _equinoxClient{}
}

type _equinoxClient struct{}

func (this _equinoxClient) Check(appID string, opts equinox.Options) (equinox.Response, error) {
	return equinox.Check(appID, opts)
}

func (this _equinoxClient) Apply(response equinox.Response) error {
	return response.Apply()
}
