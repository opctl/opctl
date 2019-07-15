package op

//go:generate counterfeiter -o ./fakeLister.go --fake-name FakeLister ./ Lister

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/opfile"
	"path/filepath"
)

type Lister interface {
	// List recursively lists ops within a directory
	List(
		ctx context.Context,
		dirHandle model.DataHandle,
	) ([]*model.OpDotYml, error)
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
) ([]*model.OpDotYml, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if nil != err {
		return nil, err
	}

	var ops []*model.OpDotYml
	for _, content := range contents {
		if filepath.Base(content.Path) == dotyml.FileName {

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
				// ignore invalid ops
				ops = append(ops, opDotYml)
			}
		}

	}

	return ops, nil
}
