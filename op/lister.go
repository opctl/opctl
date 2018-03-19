package op

//go:generate counterfeiter -o ./fakeLister.go --fake-name FakeLister ./ Lister

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"path"
)

type Lister interface {
	// List recursively lists ops within a directory
	List(
		ctx context.Context,
		dirHandle model.DataHandle,
	) ([]*model.PkgManifest, error)
}

// NewLister returns an initialized Lister instance
func NewLister() Lister {
	return _lister{
		ioUtil:             iioutil.New(),
		dotYmlUnmarshaller: dotyml.NewUnmarshaller(),
	}
}

type _lister struct {
	ioUtil             iioutil.IIOUtil
	dotYmlUnmarshaller dotyml.Unmarshaller
}

func (ls _lister) List(
	ctx context.Context,
	dirHandle model.DataHandle,
) ([]*model.PkgManifest, error) {

	contents, err := dirHandle.ListContents(ctx)
	if nil != err {
		return nil, err
	}

	var ops []*model.PkgManifest
	for _, content := range contents {
		if path.Base(content.Path) == dotyml.FileName {

			opDotYmlReader, err := dirHandle.GetContent(ctx, content.Path)
			if nil != err {
				// ignore errors for now;
				continue
			}

			opDotYmlBytes, err := ls.ioUtil.ReadAll(opDotYmlReader)
			opDotYmlReader.Close()
			if nil != err {
				// ignore errors for now;
				continue
			}

			if opDotYml, err := ls.dotYmlUnmarshaller.Unmarshal(
				opDotYmlBytes,
			); nil == err {
				// ignore err'd pkgs
				ops = append(ops, opDotYml)
			}
		}

	}

	return ops, nil
}
