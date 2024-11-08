package direntry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
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

		if *opts == "Dir" {
			if err := unsudo.CreateDir(valuePath); err != nil {
				return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)
			}

			return "", &model.Value{Dir: &valuePath}, nil
		}

		// handle file ref
		if err := unsudo.CreateFile(valuePath, []byte{}); err != nil {
			return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)
		}

		return "", &model.Value{File: &valuePath}, nil

	}

	return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)

}
