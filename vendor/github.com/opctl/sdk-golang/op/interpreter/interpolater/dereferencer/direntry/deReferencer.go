package direntry

//go:generate counterfeiter -o ./fakeDeReferencer.go --fake-name FakeDeReferencer ./ DeReferencer

import (
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"path/filepath"
	"strings"
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
	return _deReferencer{}
}

type _deReferencer struct {
}

func (dr _deReferencer) DeReference(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to deReference '%v'; expected '/'", ref)
	}

	fileValue := filepath.Join(*data.Dir, ref)

	return "", &model.Value{File: &fileValue}, nil
}
