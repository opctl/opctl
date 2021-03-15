package model

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

// DataHandle is a provider agnostic interface for interacting w/ data
// @TODO: merge Value and DataHandle
//counterfeiter:generate -o fakes/dataHandle.go . DataHandle
type DataHandle interface {
	// ListDescendants lists descendant of the data node pointed to by the current handle
	ListDescendants(
		ctx context.Context,
	) (
		[]*DirEntry,
		error,
	)

	// GetContent gets data from the current handle
	GetContent(
		ctx context.Context,
		contentPath string,
	) (
		ReadSeekCloser,
		error,
	)

	// Path returns the local path of the data; may or may not be same as Ref
	// returns nil if data doesn't exist locally
	Path() *string

	// Ref returns a ref to the data; may or may not be same as Path
	Ref() string
}

// DirEntry represents an entry in a directory (a sub directory or file)
type DirEntry struct {
	Path string
	Size int64
	Mode os.FileMode
}

// Value represents a typed value
type Value struct {
	Array   *[]interface{}          `json:"array,omitempty"`
	Boolean *bool                   `json:"boolean,omitempty"`
	Dir     *string                 `json:"dir,omitempty"`
	File    *string                 `json:"file,omitempty"`
	Number  *float64                `json:"number,omitempty"`
	Object  *map[string]interface{} `json:"object,omitempty"`
	Socket  *string                 `json:"socket,omitempty"`
	String  *string                 `json:"string,omitempty"`
}

// // ValueNode interfaces with a value of any type
// type ValueNode interface {
// 	// ID returns a globally unique id for this value
// 	ID() string
// }

// // FileValue interfaces with a file typed value
// type FileValue interface {
// 	GetReadSeekCloser() (ReadSeekCloser, error)

// 	ValueNode
// }

// // DirValue interfaces with a directory typed value
// type DirValue interface {
// 	GetDescendant(
// 		ref string,
// 	) (
// 		Value,
// 		error,
// 	)

// 	ValueNode
// }

// Unbox unboxes a Value into a native go type
func (value Value) Unbox() (interface{}, error) {
	if nil != value.Array {
		nativeArray := []interface{}{}
		for itemKey, itemValue := range *value.Array {
			switch typedItemValue := itemValue.(type) {
			case Value:
				nativeItem, err := typedItemValue.Unbox()
				if nil != err {
					return nil, errors.Wrap(err, fmt.Sprintf("unable to unbox item '%v' from array", itemKey))
				}

				nativeArray = append(nativeArray, nativeItem)
			default:
				nativeArray = append(nativeArray, itemValue)
			}
		}
		return nativeArray, nil
	} else if nil != value.Boolean {
		return *value.Boolean, nil
	} else if nil != value.Dir {
		return *value.Dir, nil
	} else if nil != value.File {
		return *value.File, nil
	} else if nil != value.Number {
		return *value.Number, nil
	} else if nil != value.Object {
		nativeObject := map[string]interface{}{}
		for propKey, propValue := range *value.Object {
			switch typedPropValue := propValue.(type) {
			case Value:
				var err error
				if nativeObject[propKey], err = typedPropValue.Unbox(); nil != err {
					return nil, errors.Wrap(err, fmt.Sprintf("unable to unbox property '%v' from object", propKey))
				}
			default:
				nativeObject[propKey] = propValue
			}
		}
		return nativeObject, nil
	} else if nil != value.Socket {
		return *value.Socket, nil
	} else if nil != value.String {
		return *value.String, nil
	}
	return nil, fmt.Errorf("unable to unbox value '%+v'", value)
}
