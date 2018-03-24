package dereferencer

//go:generate counterfeiter -o ./fakeScopeFilePathDeReferencer.go --fake-name fakeScopeFilePathDeReferencer ./ scopeFilePathDeReferencer

import (
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

// scopeFilePathDeReferencer de references scope file path refs, i.e. refs of the form: $(name/sub/file.ext)
type scopeFilePathDeReferencer interface {
	DeReferenceScopeFilePath(
		ref string,
		scope map[string]*model.Value,
	) (string, bool, error)
}

func newScopeFilePathDeReferencer() scopeFilePathDeReferencer {
	return _scopeFilePathDeReferencer{
		ioutil: iioutil.New(),
	}
}

type _scopeFilePathDeReferencer struct {
	ioutil iioutil.IIOUtil
}

func (sfd _scopeFilePathDeReferencer) DeReferenceScopeFilePath(
	ref string,
	scope map[string]*model.Value,
) (string, bool, error) {
	possiblePathIndex := strings.Index(ref, "/")
	if possiblePathIndex > 0 {
		if scopeValue, isFileRef := scope[ref[:possiblePathIndex]]; isFileRef && nil != scopeValue.Dir {

			filePath := ref[possiblePathIndex+1:]

			contentBytes, err := sfd.ioutil.ReadFile(filepath.Join(*scopeValue.Dir, filePath))
			if nil != err {
				return "", true, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}
			return string(contentBytes), true, nil
		}
	}
	return ref, false, nil
}
