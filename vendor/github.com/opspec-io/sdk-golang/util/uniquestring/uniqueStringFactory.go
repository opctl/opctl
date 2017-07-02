package uniquestring

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ UniqueStringFactory

import (
	"github.com/golang-interfaces/satori-go.uuid"
	"strings"
)

type UniqueStringFactory interface {
	Construct() (uniqueString string)
}

func NewUniqueStringFactory() UniqueStringFactory {
	return uniqueStringFactory{
		uuid: iuuid.New(),
	}
}

type uniqueStringFactory struct {
	uuid iuuid.IUUID
}

func (this uniqueStringFactory) Construct() (uniqueString string) {
	uniqueString = strings.Replace(
		this.uuid.NewV4().String(),
		"-",
		"",
		-1,
	)

	return
}
