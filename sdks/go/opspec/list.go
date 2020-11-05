package opspec

import (
	"io/ioutil"
	"context"
	"fmt"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
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
				return nil, fmt.Errorf("error opening %s%s; %s", dirHandle.Ref(), content.Path, err)
			}

			opFileBytes, err := ioutil.ReadAll(opFileReader)
			opFileReader.Close()
			if nil != err {
				return nil, fmt.Errorf("error reading %s%s; %s", dirHandle.Ref(), content.Path, err)
			}

			opFile, err := opfile.Unmarshal(
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
