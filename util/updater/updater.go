package updater

//go:generate counterfeiter -o ./fakeUpdater.go --fake-name FakeUpdater ./ Updater

import (
	"github.com/equinox-io/equinox"
)

type Update struct {
	Version         string
	equinoxResponse *equinox.Response
}

type Updater interface {
	TryGetUpdate(
		releaseChannel string,
	) (
		update *Update,
		err error,
	)
	ApplyUpdate(
		update *Update,
	) (
		err error,
	)
}

func New() Updater {
	return _updater{}
}

type _updater struct {
	publicKey []byte
}

func (this _updater) TryGetUpdate(
	releaseChannel string,
) (
	update *Update,
	err error,
) {

	publicKey := []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEf6wRkF8+1yEmOm/SwMfVUzF+Ouf3JkjH
kLhSsetGPdgalqErhAfgaYXlKKCU9lOinLmmrjm6pTaGgl1vwCOjlj9QjGQMYpOF
NpxANJQvIwvHHHGb1VgCLj2kYOeyIa6D
-----END ECDSA PUBLIC KEY-----
`)

	opts := equinox.Options{Channel: releaseChannel}
	if err = opts.SetPublicKeyPEM(publicKey); err != nil {
		return
	}

	// check for the update
	equinoxResponse, err := equinox.Check("app_kNrDsPk2bis", opts)
	switch {
	case err == equinox.NotAvailableErr:
		err = nil
		return
	case err != nil:
		return
	}
	update = &Update{
		equinoxResponse: &equinoxResponse,
		Version:         equinoxResponse.ReleaseVersion,
	}

	return
}

func (this _updater) ApplyUpdate(
	update *Update,
) (err error) {

	// fetch the update and apply it
	err = update.equinoxResponse.Apply()
	if err != nil {
		return
	}

	return

}
