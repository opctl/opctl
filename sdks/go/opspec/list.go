package opspec

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/pkg/errors"
)

// List ops recursively within a directory, returning discovered op files by path.
func List(
	ctx context.Context,
	dirHandle model.DataHandle,
) (map[string]*model.OpSpec, error) {

	contents, err := dirHandle.ListDescendants(ctx)
	if nil != err {
		return nil, err
	}

	opsByPath := map[string]*model.OpSpec{}
	for _, content := range contents {
		if filepath.Base(content.Path) == opfile.FileName {

			opFileReader, err := dirHandle.GetContent(ctx, content.Path)
			if nil != err {
				return nil, errors.Wrap(err, fmt.Sprintf("error opening %s%s", dirHandle.Ref(), content.Path))
			}

			opFileBytes, err := ioutil.ReadAll(opFileReader)
			opFileReader.Close()
			if nil != err {
				return nil, errors.Wrap(err, fmt.Sprintf("error reading %s%s", dirHandle.Ref(), content.Path))
			}

			opFile, err := opfile.Unmarshal(
				opFileBytes,
			)
			if nil != err {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshalling %s%s", dirHandle.Ref(), content.Path))
			}

			opsByPath[filepath.Dir(content.Path)] = opFile
		}

	}

	return opsByPath, nil
}
