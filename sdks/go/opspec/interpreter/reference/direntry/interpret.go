package direntry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// Interpret a dir entry ref i.e. refs of the form name/sub/file.ext
// it's an error if ref doesn't start with '/'
// returns ref remainder, dereferenced data, and error if one occurred
func Interpret(
	ref string,
	data *model.Value,
	opts *string,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: expected '/'", ref)
	}

	valuePath := filepath.Join(*data.Dir, ref)

	fileInfo, err := os.Stat(valuePath)
	if err == nil {
		if fileInfo.IsDir() {
			return "", &model.Value{Dir: &valuePath}, nil
		}

		return "", &model.Value{File: &valuePath}, nil
	} else if opts != nil && os.IsNotExist(err) {

		if "Dir" == *opts {
			err := os.MkdirAll(valuePath, 0700)
			if err != nil {
				return "", nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v' as dir entry ref", ref))
			}

			return "", &model.Value{Dir: &valuePath}, nil
		}

		// handle file ref
		err := os.MkdirAll(filepath.Dir(valuePath), 0700)
		if err != nil {
			return "", nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v' as dir entry ref", ref))
		}

		file, err := os.Create(valuePath)
		file.Close()
		if err != nil {
			return "", nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v' as dir entry ref", ref))
		}

		return "", &model.Value{File: &valuePath}, nil

	}

	return "", nil, errors.Wrap(err, fmt.Sprintf("unable to interpret '%v' as dir entry ref", ref))

}
