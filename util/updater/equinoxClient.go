package updater

//go:generate counterfeiter -o ./fakeEquinoxClient.go --fake-name fakeEquinoxClient ./ equinoxClient

import "github.com/equinox-io/equinox"

// client interface for equinox
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
