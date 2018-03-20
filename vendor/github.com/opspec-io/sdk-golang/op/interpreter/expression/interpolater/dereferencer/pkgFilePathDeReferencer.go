package dereferencer

//go:generate counterfeiter -o ./fakePkgFilePathDeReferencer.go --fake-name fakePkgFilePathDeReferencer ./ pkgFilePathDeReferencer

import (
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"strings"
)

// pkgFilePathDeReferencer de references pkg file path refs, i.e. refs of the form: $(/name/sub/file.ext)
type pkgFilePathDeReferencer interface {
	DeReferencePkgFilePath(
		ref string,
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (string, bool, error)
}

func newPkgFilePathDeReferencer() pkgFilePathDeReferencer {
	return _pkgFilePathDeReferencer{}
}

type _pkgFilePathDeReferencer struct{}

func (pfd _pkgFilePathDeReferencer) DeReferencePkgFilePath(
	ref string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (string, bool, error) {
	if strings.HasPrefix(ref, "/") {
		// pkg content ref
		contentReadSeekCloser, err := opHandle.GetContent(context.TODO(), ref)
		if nil != err {
			return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		defer contentReadSeekCloser.Close()

		contentBytes, err := ioutil.ReadAll(contentReadSeekCloser)
		if nil != err {
			return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		return string(contentBytes), true, nil
	}
	return ref, false, nil
}
