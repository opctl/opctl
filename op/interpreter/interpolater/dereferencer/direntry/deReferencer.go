package direntry

//go:generate counterfeiter -o ./fakeDeReferencer.go --fake-name FakeDeReferencer ./ DeReferencer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
)

// DeReferencer de references a dir entry ref i.e. refs of the form name/sub/file.ext
// it's an error if ref doesn't start with '/'
// returns ref remainder, dereferenced data, and error if one occurred
type DeReferencer interface {
	DeReference(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
}

func NewDeReferencer() DeReferencer {
	return _deReferencer{
		ioutil: iioutil.New(),
	}
}

type _deReferencer struct {
	ioutil iioutil.IIOUtil
}

func (dr _deReferencer) DeReference(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to deReference '%v'; expected '/'", ref)
	}

	contentBytes, err := dr.ioutil.ReadFile(filepath.Join(*data.Dir, ref))
	if nil != err {
		return "", nil, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
	}

	contentString := string(contentBytes)

	return "", &model.Value{File: &contentString}, nil
}
