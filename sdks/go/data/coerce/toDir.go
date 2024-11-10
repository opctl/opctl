package coerce

import (
	"bytes"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
)

// ToDir attempts to coerce value to a dir
func ToDir(
	value ipld.Node,
	scratchDir string,
) (ipld.Node, error) {
	switch {
	case value.IsNull():
		return nil, errors.New("unable to coerce null to dir")
	case value.Kind() == ipld.Kind_List:
		return nil, fmt.Errorf("unable to coerce array to dir: %w", errIncompatibleTypes)
	case value.Kind() == ipld.Kind_Bool:
		return nil, fmt.Errorf("unable to coerce boolean to dir: %w", errIncompatibleTypes)
	case value.Kind() == ipld.Kind_String:
		str, err := value.AsString()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to dir: %w", err)
		}

		strIsPath := filepath.IsAbs(str)
		if strIsPath {
			stat, err := os.Stat(str)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce string to dir: %w", err)
			}

			if stat.IsDir() {
				return value, nil
			}

			if !stat.Mode().IsRegular() {
				return nil, fmt.Errorf("unable to coerce socket to dir: %w", errIncompatibleTypes)
			}

			return nil, fmt.Errorf("unable to coerce file to dir: %w", errIncompatibleTypes)
		}

		return nil, fmt.Errorf("unable to coerce string to dir: %w", errIncompatibleTypes)
	case value.Kind() == ipld.Kind_Int || value.Kind() == ipld.Kind_Float:
		return nil, fmt.Errorf("unable to coerce number to dir: %w", errIncompatibleTypes)
	case value.Kind() == ipld.Kind_Map:
		uniqueStr, err := uniquestring.Construct()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to dir: %w", err)
		}

		rootDirPath := filepath.Join(scratchDir, uniqueStr)

		var buf bytes.Buffer
		if err := json.Encode(value, &buf); err != nil {
			return nil, fmt.Errorf("unable to coerce object to dir: %w", err)
		}
	
		// Deserialize JSON into a map[string]interface{}
		var children map[string]interface{}
		if err := stdjson.Unmarshal(buf.Bytes(), &children); err != nil {
			return nil, fmt.Errorf("unable to coerce object to dir: %w", err)
		}

		err = rCreateFileItem(rootDirPath, "", children)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to dir: %w", err)
		}

		nb := basicnode.Prototype.Bool.NewBuilder()
		nb.AssignString(rootDirPath)
		return nb.Build(), nil
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
			return fmt.Errorf("error creating %s: %w", itemPath, err)
		}

		err = os.WriteFile(itemPath, []byte(dataString), 0777)
		if err != nil {
			return fmt.Errorf("error creating %s: %w", itemPath, err)
		}

		return nil
	}

	// ensure dir exists
	err := os.MkdirAll(
		itemPath,
		0777,
	)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", itemPath, err)
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
