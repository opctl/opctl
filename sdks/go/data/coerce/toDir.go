package coerce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// ToDir attempts to coerce value to a dir
func ToDir(
	value *model.Value,
	scratchDir string,
) (*model.Value, error) {
	switch {
	case value == nil:
		return nil, errors.New("unable to coerce null to dir")
	case value.Array != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce array to dir")
	case value.Boolean != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce boolean to dir")
	case value.Dir != nil:
		return value, nil
	case value.File != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce file to dir")
	case value.Number != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce number to dir")
	case value.Object != nil:
		uniqueStr, err := uniquestring.Construct()
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to dir")
		}

		rootDirPath := filepath.Join(scratchDir, uniqueStr)
		err = rCreateFileItem(rootDirPath, "", *value.Object)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to dir")
		}

		return &model.Value{Dir: &rootDirPath}, nil
	case value.Socket != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce socket to dir")
	case value.String != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce string to dir")
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to dir", value)
	}
}

func rCreateFileItem(
	rootPath,
	relParentPath string,
	children map[string]interface{},
) error {
	itemPath := filepath.Join(rootPath, relParentPath)

	if fileData, ok := children["data"]; ok && len(children) == 1 {

		// handle file
		dataString, ok := fileData.(string)
		if !ok {
			return fmt.Errorf("%s .data not string", relParentPath)
		}

		// ensure parent exists
		err := os.MkdirAll(
			filepath.Dir(itemPath),
			0777,
		)
		if err != nil {
			return errors.Wrap(err, "error creating "+itemPath)
		}

		err = ioutil.WriteFile(itemPath, []byte(dataString), 0777)
		if err != nil {
			return errors.Wrap(err, "error creating "+itemPath)
		}

		return nil
	}

	// ensure dir exists
	err := os.MkdirAll(
		itemPath,
		0777,
	)
	if err != nil {
		return errors.Wrap(err, "error creating "+itemPath)
	}

	for k, v := range children {
		if !strings.HasPrefix(k, "") {
			return fmt.Errorf(`%s %s must start with "/"`, relParentPath, k)
		}

		relChildPath := filepath.Join(relParentPath, k)

		switch v := v.(type) {
		case map[string]interface{}:
			if err := rCreateFileItem(rootPath, relChildPath, v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s not a valid file/dir initializer", relChildPath)
		}
	}

	return nil
}
