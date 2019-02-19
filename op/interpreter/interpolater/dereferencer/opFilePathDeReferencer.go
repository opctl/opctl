package dereferencer

//go:generate counterfeiter -o ./fakeOpFilePathDeReferencer.go --fake-name fakeOpFilePathDeReferencer ./ opFilePathDeReferencer

import (
	"context"
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"io/ioutil"
	"strings"
)

// opFilePathDeReferencer de references op file path refs, i.e. refs of the form: $(/name/sub/file.ext)
type opFilePathDeReferencer interface {
	DeReferenceOpFilePath(
		ref string,
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (string, bool, error)
}

func newOpFilePathDeReferencer() opFilePathDeReferencer {
	return _opFilePathDeReferencer{}
}

type _opFilePathDeReferencer struct{}

func (pfd _opFilePathDeReferencer) DeReferenceOpFilePath(
	ref string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (string, bool, error) {
	if strings.HasPrefix(ref, "/") {
		// op content ref
		contentReadSeekCloser, err := opHandle.GetContent(context.TODO(), ref)
		if nil != err {
			return "", true, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		defer contentReadSeekCloser.Close()

		contentBytes, err := ioutil.ReadAll(contentReadSeekCloser)
		if nil != err {
			return "", true, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		return string(contentBytes), true, nil
	}
	return ref, false, nil
}
