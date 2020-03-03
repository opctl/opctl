package uniquestring

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/golang-interfaces/satori-go.uuid"
	"strings"
)

//counterfeiter:generate -o fakes/uniqueStringFactory.go . UniqueStringFactory
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
