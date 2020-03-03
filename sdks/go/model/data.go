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

// Unbox unboxes a Value into a native go type
func (vlu Value) Unbox() (interface{}, error) {
	if nil != vlu.Array {
		nativeArray := []interface{}{}
		for itemKey, itemValue := range *vlu.Array {
			switch typedItemValue := itemValue.(type) {
			case Value:
				nativeItem, err := typedItemValue.Unbox()
				if nil != err {
					return nil, fmt.Errorf("unable to unbox array; error unboxing item '%v' was %v", itemKey, err)
				}

				nativeArray = append(nativeArray, nativeItem)
			default:
				nativeArray = append(nativeArray, itemValue)
			}
		}
		return nativeArray, nil
	} else if nil != vlu.Boolean {
		return *vlu.Boolean, nil
	} else if nil != vlu.Dir {
		return *vlu.Dir, nil
	} else if nil != vlu.File {
		return *vlu.File, nil
	} else if nil != vlu.Number {
		return *vlu.Number, nil
	} else if nil != vlu.Object {
		nativeObject := map[string]interface{}{}
		for propKey, propValue := range *vlu.Object {
			switch typedPropValue := propValue.(type) {
			case Value:
				var err error
				if nativeObject[propKey], err = typedPropValue.Unbox(); nil != err {
					return nil, fmt.Errorf("unable to unbox object; error unboxing property '%v' was %v", propKey, err)
				}
			default:
				nativeObject[propKey] = propValue
			}
		}
		return nativeObject, nil
	} else if nil != vlu.Socket {
		return *vlu.Socket, nil
	} else if nil != vlu.String {
		return *vlu.String, nil
	}
	return nil, fmt.Errorf("unable to unbox value %+v; box unknown", vlu)
}

// PullCreds contains authentication attributes for auth'ing w/ a data provider
type PullCreds struct {
	Username,
	Password string
}
