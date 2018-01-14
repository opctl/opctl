package dereferencer

//go:generate counterfeiter -o ./fakeScopeFileDeReferencer.go --fake-name fakeScopeFileDeReferencer ./ scopeFileDeReferencer

import (
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

// scopeFileDeReferencer de references scope file refs, i.e. refs of the form: $(name/sub/file.ext)
type scopeFileDeReferencer interface {
	DeReferenceScopeFile(
		ref string,
		scope map[string]*model.Value,
	) (string, bool, error)
}

func newScopeFileDeReferencer() scopeFileDeReferencer {
	return _scopeFileDeReferencer{
		ioutil: iioutil.New(),
	}
}

type _scopeFileDeReferencer struct {
	ioutil iioutil.IIOUtil
}

func (sfd _scopeFileDeReferencer) DeReferenceScopeFile(
	ref string,
	scope map[string]*model.Value,
) (string, bool, error) {
	possiblePathIndex := strings.Index(ref, "/")
	if possiblePathIndex > 0 {
		if scopeValue, isFileRef := scope[ref[:possiblePathIndex]]; isFileRef && nil != scopeValue.Dir {

			filePath := ref[possiblePathIndex+1:]

			contentBytes, err := sfd.ioutil.ReadFile(filepath.Join(*scopeValue.Dir, filePath))
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}
			return string(contentBytes), true, nil
		}
	}
	return ref, false, nil
}
