package opspec

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

//counterfeiter:generate -o fakes/lister.go . Lister
type Lister interface {
	// List recursively lists ops within a directory, returning discovered op files by path.
	List(
		ctx context.Context,
		dirHandle model.DataHandle,
	) (map[string]*model.OpFile, error)
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
) (map[string]*model.OpFile, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if nil != err {
		return nil, err
	}

	opsByPath := map[string]*model.OpFile{}
	for _, content := range contents {
		if filepath.Base(content.Path) == opfile.FileName {

			opFileReader, err := dirHandle.GetContent(ctx, content.Path)
			if nil != err {
				return nil, fmt.Errorf("error opening %s%s; %s", dirHandle.Ref(), content.Path, err)
			}

			opFileBytes, err := ls.ioUtil.ReadAll(opFileReader)
			opFileReader.Close()
			if nil != err {
				return nil, fmt.Errorf("error reading %s%s; %s", dirHandle.Ref(), content.Path, err)
			}

			opFile, err := ls.opFileUnmarshaller.Unmarshal(
				opFileBytes,
			)
			if nil != err {
				return nil, fmt.Errorf("error unmarshalling %s%s; %s", dirHandle.Ref(), content.Path, err)
			}

			opsByPath[filepath.Dir(content.Path)] = opFile
		}

	}

	return opsByPath, nil
}
