package model

import (
	"context"
	"fmt"
	"io"
	"os"
)

type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

// DataHandle is a provider agnostic interface for interacting w/ data
// @TODO: merge Value and DataHandle
//
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

	// Ref returns a ref to the data
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

// Unbox unboxes a Value into a native go type
func (value Value) Unbox() (interface{}, error) {
	if value.Array != nil {
		nativeArray := []interface{}{}
		for itemKey, itemValue := range *value.Array {
			switch typedItemValue := itemValue.(type) {
			case Value:
				nativeItem, err := typedItemValue.Unbox()
				if err != nil {
					return nil, fmt.Errorf("unable to unbox item '%v' from array: %w", itemKey, err)
				}

				nativeArray = append(nativeArray, nativeItem)
			default:
				nativeArray = append(nativeArray, itemValue)
			}
		}
		return nativeArray, nil
	} else if value.Boolean != nil {
		return *value.Boolean, nil
	} else if value.Dir != nil {
		return *value.Dir, nil
	} else if value.File != nil {
		return *value.File, nil
	} else if value.Number != nil {
		return *value.Number, nil
	} else if value.Object != nil {
		nativeObject := map[string]interface{}{}
		for propKey, propValue := range *value.Object {
			switch typedPropValue := propValue.(type) {
			case Value:
				var err error
				if nativeObject[propKey], err = typedPropValue.Unbox(); err != nil {
					return nil, fmt.Errorf("unable to unbox property '%v' from object: %w", propKey, err)
				}
			default:
				nativeObject[propKey] = propValue
			}
		}
		return nativeObject, nil
	} else if value.Socket != nil {
		return *value.Socket, nil
	} else if value.String != nil {
		return *value.String, nil
	}
	return nil, fmt.Errorf("unable to unbox value '%+v'", value)
}
