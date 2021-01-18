package coerce

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
)

// ToDir attempts to coerce value to a dir
func ToDir(
	value *model.Value,
	scratchDir string,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, fmt.Errorf("unable to coerce null to dir")
	case nil != value.Array:
		return nil, fmt.Errorf("unable to coerce array to dir; incompatible types")
	case nil != value.Boolean:
		return nil, fmt.Errorf("unable to coerce boolean to dir; incompatible types")
	case nil != value.Dir:
		return value, nil
	case nil != value.File:
		return nil, fmt.Errorf("unable to coerce file to dir; incompatible types")
	case nil != value.Number:
		return nil, fmt.Errorf("unable to coerce number to dir; incompatible types")
	case nil != value.Object:
		uniqueStr, err := uniquestring.Construct()
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to dir; error was %v", err.Error())
		}

		rootDirPath := filepath.Join(scratchDir, uniqueStr)
		err = rCreateFileItem(rootDirPath, "", *value.Object)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to dir; error was %v", err.Error())
		}

		return &model.Value{Dir: &rootDirPath}, nil
	case nil != value.Socket:
		return nil, fmt.Errorf("unable to coerce socket to dir; incompatible types")
	case nil != value.String:
		return nil, fmt.Errorf("unable to coerce string to dir; incompatible types")
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
		if nil != err {
			return fmt.Errorf("error creating %s; error was %s", itemPath, err.Error())
		}

		err = ioutil.WriteFile(itemPath, []byte(dataString), 0777)
		if nil != err {
			return fmt.Errorf("error creating %s; error was %s", itemPath, err.Error())
		}

		return nil
	}

	// ensure dir exists
	err := os.MkdirAll(
		itemPath,
		0777,
	)
	if nil != err {
		return fmt.Errorf("error creating %s; error was %s", itemPath, err.Error())
	}

	for k, v := range children {
		if !strings.HasPrefix(k, "") {
			return fmt.Errorf("%s %s must start with '/'", relParentPath, k)
		}

		relChildPath := filepath.Join(relParentPath, k)

		switch v := v.(type) {
		case map[string]interface{}:
			if err := rCreateFileItem(rootPath, relChildPath, v); nil != err {
				return err
			}
		default:
			return fmt.Errorf("%s not a valid file/dir initializer", relChildPath)
		}
	}

	return nil
}
