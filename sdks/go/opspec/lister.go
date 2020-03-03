package opspec

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"path/filepath"
)

//counterfeiter:generate -o fakes/lister.go . Lister
type Lister interface {
	// List recursively lists ops within a directory
	List(
		ctx context.Context,
		dirHandle model.DataHandle,
	) ([]*model.OpFile, error)
}

// NewLister returns an initialized Lister instance
func NewLister() Lister {
	return _lister{
		ioUtil:             iioutil.New(),
		opFileUnmarshaller: opfile.NewUnmarshaller(),
	}
}

type _lister struct {
	ioUtil             iioutil.IIOUtil
	opFileUnmarshaller opfile.Unmarshaller
}

func (ls _lister) List(
	ctx context.Context,
	dirHandle model.DataHandle,
) ([]*model.OpFile, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if nil != err {
		return nil, err
	}

	var ops []*model.OpFile
	for _, content := range contents {
		if filepath.Base(content.Path) == opfile.FileName {

			opFileReader, err := dirHandle.GetContent(ctx, content.Path)
			if nil != err {
				// ignore errors for now;
				continue
			}

			opFileBytes, err := ls.ioUtil.ReadAll(opFileReader)
			opFileReader.Close()
			if nil != err {
				// ignore errors for now;
				continue
			}

			if opFile, err := ls.opFileUnmarshaller.Unmarshal(
				opFileBytes,
			); nil == err {
				// ignore invalid ops
				ops = append(ops, opFile)
			}
		}

	}

	return ops, nil
}
