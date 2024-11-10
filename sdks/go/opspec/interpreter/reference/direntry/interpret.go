package direntry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
)

// Interpret a dir entry ref i.e. refs of the form name/sub/file.ext
// it's an error if ref doesn't start with '/'
// returns ref remainder, dereferenced data, and error if one occurred
func Interpret(
	ref string,
	data *ipld.Node,
	opts *string,
) (string, *ipld.Node, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: expected '/'", ref)
	}

	valuePath := filepath.Join(*data.Dir, ref)

	fileInfo, err := os.Stat(valuePath)
	if err == nil {
		if fileInfo.IsDir() {
			return "", &ipld.Node{Dir: &valuePath}, nil
		}

		return "", &ipld.Node{File: &valuePath}, nil
	} else if opts != nil && os.IsNotExist(err) {

		if *opts == "Dir" {
			err := os.MkdirAll(valuePath, 0770)
			if err != nil {
				return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)
			}

			return "", &ipld.Node{Dir: &valuePath}, nil
		}

		// handle file ref
		err := os.MkdirAll(filepath.Dir(valuePath), 0770)
		if err != nil {
			return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)
		}

		file, err := os.Create(valuePath)
		file.Close()
		if err != nil {
			return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)
		}

		return "", &ipld.Node{File: &valuePath}, nil

	}

	return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref: %w", ref, err)

}
