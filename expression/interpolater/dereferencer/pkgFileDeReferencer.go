package dereferencer

//go:generate counterfeiter -o ./fakePkgFileDeReferencer.go --fake-name fakePkgFileDeReferencer ./ pkgFileDeReferencer

import (
	"context"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"strings"
)

// pkgFileDeReferencer de references pkg file refs, i.e. refs of the form: $(/name/sub/file.ext)
type pkgFileDeReferencer interface {
	DeReferencePkgFile(
		ref string,
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) (string, bool, error)
}

func newPkgFileDeReferencer() pkgFileDeReferencer {
	return _pkgFileDeReferencer{}
}

type _pkgFileDeReferencer struct{}

func (pfd _pkgFileDeReferencer) DeReferencePkgFile(
	ref string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (string, bool, error) {
	if strings.HasPrefix(ref, "/") {
		// pkg content ref
		contentReadSeekCloser, err := pkgHandle.GetContent(context.TODO(), ref)
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
