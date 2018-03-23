package iyaml

import (
	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter -o fake.go --fake-name Fake ./ IYAML

type IYAML interface {
	// Marshal serializes the value provided into a YAML document. The structure
	// of the generated document will reflect the structure of the value itself.
	// Maps and pointers (to struct, string, int, etc) are accepted as the in value.
	//
	// Struct fields are only unmarshalled if they are exported (have an upper case
	// first letter), and are unmarshalled using the field name lowercased as the
	// default key. Custom keys may be defined via the "yaml" name in the field
	// tag: the content preceding the first comma is used as the key, and the
	// following comma-separated options are used to tweak the marshalling process.
	// Conflicting names result in a runtime error.
	//
	// The field tag format accepted is:
	//
	//     `(...) yaml:"[<key>][,<flag1>[,<flag2>]]" (...)`
	//
	// The following flags are currently supported:
	//
	//     omitempty    Only include the field if it's not set to the zero
	//                  value for the type or to empty slices or maps.
	//                  Does not apply to zero valued structs.
	//
	//     flow         Marshal using a flow style (useful for structs,
	//                  sequences and maps).
	//
	//     inline       Inline the field, which must be a struct or a map,
	//                  causing all of its fields or keys to be processed as if
	//                  they were part of the outer struct. For maps, keys must
	//                  not conflict with the yaml keys of other struct fields.
	//
	// In addition, if the key is "-", the field is ignored.
	//
	// For example:
	//
	//     type T struct {
	//         F int "a,omitempty"
	//         B int
	//     }
	//     yaml.Marshal(&T{B: 2}) // Returns "b: 2\n"
	//     yaml.Marshal(&T{F: 1}} // Returns "a: 1\nb: 0\n"
	//
	Marshal(v interface{}) ([]byte, error)

	// Unmarshal decodes the first document found within the in byte slice
	// and assigns decoded values into the out value.
	//
	// Maps and pointers (to a struct, string, int, etc) are accepted as out
	// values. If an internal pointer within a struct is not initialized,
	// the yaml package will initialize it if necessary for unmarshalling
	// the provided data. The out parameter must not be nil.
	//
	// The type of the decoded values should be compatible with the respective
	// values in out. If one or more values cannot be decoded due to a type
	// mismatches, decoding continues partially until the end of the YAML
	// content, and a *yaml.TypeError is returned with details for all
	// missed values.
	//
	// Struct fields are only unmarshalled if they are exported (have an
	// upper case first letter), and are unmarshalled using the field name
	// lowercased as the default key. Custom keys may be defined via the
	// "yaml" name in the field tag: the content preceding the first comma
	// is used as the key, and the following comma-separated options are
	// used to tweak the marshalling process (see Marshal).
	// Conflicting names result in a runtime error.
	//
	// For example:
	//
	//     type T struct {
	//         F int `yaml:"a,omitempty"`
	//         B int
	//     }
	//     var t T
	//     yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)
	//
	// See the documentation of Marshal for the format of tags and a list of
	// supported tag options.
	//
	Unmarshal(data []byte, v interface{}) error
}

func New() IYAML {
	return _IJSON{}
}

type _IJSON struct{}

func (ijson _IJSON) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (ijson _IJSON) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
