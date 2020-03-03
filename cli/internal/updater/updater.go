package updater

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/equinox-io/equinox"
)

type Update struct {
	Version         string
	equinoxResponse *equinox.Response
}

//counterfeiter:generate -o fakes/updater.go . Updater
type Updater interface {
	GetUpdateIfExists(
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
	return _new(newEquinoxClient())
}

// constructs an Updater w/ provided equinoxClient to enable unit testing
func _new(
	equinoxClient equinoxClient,
) Updater {
	return _updater{
		publicKey: []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEf6wRkF8+1yEmOm/SwMfVUzF+Ouf3JkjH
kLhSsetGPdgalqErhAfgaYXlKKCU9lOinLmmrjm6pTaGgl1vwCOjlj9QjGQMYpOF
NpxANJQvIwvHHHGb1VgCLj2kYOeyIa6D
-----END ECDSA PUBLIC KEY-----
`),
		equinoxClient: equinoxClient,
	}
}

type _updater struct {
	publicKey     []byte
	equinoxClient equinoxClient
}

func (this _updater) GetUpdateIfExists(
	releaseChannel string,
) (
	update *Update,
	err error,
) {

	opts := equinox.Options{Channel: releaseChannel}
	if err = opts.SetPublicKeyPEM(this.publicKey); nil != err {
		return
	}

	// check for an update
	equinoxResponse, err := this.equinoxClient.Check("app_kNrDsPk2bis", opts)
	switch {
	case err == equinox.NotAvailableErr:
		err = nil
		return
	case nil != err:
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
	err = this.equinoxClient.Apply(*update.equinoxResponse)
	if nil != err {
		return
	}

	return

}
