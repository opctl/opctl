package uniquestring

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ UniqueStringFactory

import (
	"github.com/golang-interfaces/satori-go.uuid"
	"strings"
)

type UniqueStringFactory interface {
	Construct() (string, error)
}

func NewUniqueStringFactory() UniqueStringFactory {
	return uniqueStringFactory{
		uuid: iuuid.New(),
	}
}

type uniqueStringFactory struct {
	uuid iuuid.IUUID
}

func (this uniqueStringFactory) Construct() (string, error) {
	uuid, err := this.uuid.NewV4()
	if nil != err {
		return "", err
	}

	return strings.Replace(
		uuid.String(),
		"-",
		"",
		-1,
	), nil
}
