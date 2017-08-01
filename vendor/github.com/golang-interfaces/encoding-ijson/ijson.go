package ijson

import (
	"encoding/json"
	"io"
)

//go:generate counterfeiter -o fake.go --fake-name Fake ./ IJSON

type IJSON interface {
	// Marshal returns the JSON encoding of v.
	//
	// Marshal traverses the value v recursively.
	// If an encountered value implements the Marshaler interface
	// and is not a nil pointer, Marshal calls its MarshalJSON method
	// to produce JSON. If no MarshalJSON method is present but the
	// value implements encoding.TextMarshaler instead, Marshal calls
	// its MarshalText method and encodes the result as a JSON string.
	// The nil pointer exception is not strictly necessary
	// but mimics a similar, necessary exception in the behavior of
	// UnmarshalJSON.
	//
	// Otherwise, Marshal uses the following type-dependent default encodings:
	//
	// Boolean values encode as JSON booleans.
	//
	// Floating point, integer, and Number values encode as JSON numbers.
	//
	// String values encode as JSON strings coerced to valid UTF-8,
	// replacing invalid bytes with the Unicode replacement rune.
	// The angle brackets "<" and ">" are escaped to "\u003c" and "\u003e"
	// to keep some browsers from misinterpreting JSON output as HTML.
	// Ampersand "&" is also escaped to "\u0026" for the same reason.
	// This escaping can be disabled using an Encoder that had SetEscapeHTML(false)
	// called on it.
	//
	// Array and slice values encode as JSON arrays, except that
	// []byte encodes as a base64-encoded string, and a nil slice
	// encodes as the null JSON value.
	//
	// Struct values encode as JSON objects.
	// Each exported struct field becomes a member of the object, using the
	// field name as the object key, unless the field is omitted for one of the
	// reasons given below.
	//
	// The encoding of each struct field can be customized by the format string
	// stored under the "json" key in the struct field's tag.
	// The format string gives the name of the field, possibly followed by a
	// comma-separated list of options. The name may be empty in order to
	// specify options without overriding the default field name.
	//
	// The "omitempty" option specifies that the field should be omitted
	// from the encoding if the field has an empty value, defined as
	// false, 0, a nil pointer, a nil interface value, and any empty array,
	// slice, map, or string.
	//
	// As a special case, if the field tag is "-", the field is always omitted.
	// Note that a field with name "-" can still be generated using the tag "-,".
	//
	// Examples of struct field tags and their meanings:
	//
	//   // Field appears in JSON as key "myName".
	//   Field int `json:"myName"`
	//
	//   // Field appears in JSON as key "myName" and
	//   // the field is omitted from the object if its value is empty,
	//   // as defined above.
	//   Field int `json:"myName,omitempty"`
	//
	//   // Field appears in JSON as key "Field" (the default), but
	//   // the field is skipped if empty.
	//   // Note the leading comma.
	//   Field int `json:",omitempty"`
	//
	//   // Field is ignored by this package.
	//   Field int `json:"-"`
	//
	//   // Field appears in JSON as key "-".
	//   Field int `json:"-,"`
	//
	// The "string" option signals that a field is stored as JSON inside a
	// JSON-encoded string. It applies only to fields of string, floating point,
	// integer, or boolean types. This extra level of encoding is sometimes used
	// when communicating with JavaScript programs:
	//
	//    Int64String int64 `json:",string"`
	//
	// The key name will be used if it's a non-empty string consisting of
	// only Unicode letters, digits, and ASCII punctuation except quotation
	// marks, backslash, and comma.
	//
	// Anonymous struct fields are usually marshaled as if their inner exported fields
	// were fields in the outer struct, subject to the usual Go visibility rules amended
	// as described in the next paragraph.
	// An anonymous struct field with a name given in its JSON tag is treated as
	// having that name, rather than being anonymous.
	// An anonymous struct field of interface type is treated the same as having
	// that type as its name, rather than being anonymous.
	//
	// The Go visibility rules for struct fields are amended for JSON when
	// deciding which field to marshal or unmarshal. If there are
	// multiple fields at the same level, and that level is the least
	// nested (and would therefore be the nesting level selected by the
	// usual Go rules), the following extra rules apply:
	//
	// 1) Of those fields, if any are JSON-tagged, only tagged fields are considered,
	// even if there are multiple untagged fields that would otherwise conflict.
	//
	// 2) If there is exactly one field (tagged or not according to the first rule), that is selected.
	//
	// 3) Otherwise there are multiple fields, and all are ignored; no error occurs.
	//
	// Handling of anonymous struct fields is new in Go 1.1.
	// Prior to Go 1.1, anonymous struct fields were ignored. To force ignoring of
	// an anonymous struct field in both current and earlier versions, give the field
	// a JSON tag of "-".
	//
	// Map values encode as JSON objects. The map's key type must either be a
	// string, an integer type, or implement encoding.TextMarshaler. The map keys
	// are sorted and used as JSON object keys by applying the following rules,
	// subject to the UTF-8 coercion described for string values above:
	//   - string keys are used directly
	//   - encoding.TextMarshalers are marshaled
	//   - integer keys are converted to strings
	//
	// Pointer values encode as the value pointed to.
	// A nil pointer encodes as the null JSON value.
	//
	// Interface values encode as the value contained in the interface.
	// A nil interface value encodes as the null JSON value.
	//
	// Channel, complex, and function values cannot be encoded in JSON.
	// Attempting to encode such a value causes Marshal to return
	// an UnsupportedTypeError.
	//
	// JSON cannot represent cyclic data structures and Marshal does not
	// handle them. Passing cyclic structures to Marshal will result in
	// an infinite recursion.
	//
	Marshal(v interface{}) ([]byte, error)

	// NewEncoder returns a new encoder that writes to w.
	NewEncoder(w io.Writer) *json.Encoder

	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v.
	//
	// Unmarshal uses the inverse of the encodings that
	// Marshal uses, allocating maps, slices, and pointers as necessary,
	// with the following additional rules:
	//
	// To unmarshal JSON into a pointer, Unmarshal first handles the case of
	// the JSON being the JSON literal null. In that case, Unmarshal sets
	// the pointer to nil. Otherwise, Unmarshal unmarshals the JSON into
	// the value pointed at by the pointer. If the pointer is nil, Unmarshal
	// allocates a new value for it to point to.
	//
	// To unmarshal JSON into a value implementing the Unmarshaler interface,
	// Unmarshal calls that value's UnmarshalJSON method, including
	// when the input is a JSON null.
	// Otherwise, if the value implements encoding.TextUnmarshaler
	// and the input is a JSON quoted string, Unmarshal calls that value's
	// UnmarshalText method with the unquoted form of the string.
	//
	// To unmarshal JSON into a struct, Unmarshal matches incoming object
	// keys to the keys used by Marshal (either the struct field name or its tag),
	// preferring an exact match but also accepting a case-insensitive match.
	// Unmarshal will only set exported fields of the struct.
	//
	// To unmarshal JSON into an interface value,
	// Unmarshal stores one of these in the interface value:
	//
	//	bool, for JSON booleans
	//	float64, for JSON numbers
	//	string, for JSON strings
	//	[]interface{}, for JSON arrays
	//	map[string]interface{}, for JSON objects
	//	nil for JSON null
	//
	// To unmarshal a JSON array into a slice, Unmarshal resets the slice length
	// to zero and then appends each element to the slice.
	// As a special case, to unmarshal an empty JSON array into a slice,
	// Unmarshal replaces the slice with a new empty slice.
	//
	// To unmarshal a JSON array into a Go array, Unmarshal decodes
	// JSON array elements into corresponding Go array elements.
	// If the Go array is smaller than the JSON array,
	// the additional JSON array elements are discarded.
	// If the JSON array is smaller than the Go array,
	// the additional Go array elements are set to zero values.
	//
	// To unmarshal a JSON object into a map, Unmarshal first establishes a map to
	// use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal
	// reuses the existing map, keeping existing entries. Unmarshal then stores
	// key-value pairs from the JSON object into the map. The map's key type must
	// either be a string, an integer, or implement encoding.TextUnmarshaler.
	//
	// If a JSON value is not appropriate for a given target type,
	// or if a JSON number overflows the target type, Unmarshal
	// skips that field and completes the unmarshaling as best it can.
	// If no more serious errors are encountered, Unmarshal returns
	// an UnmarshalTypeError describing the earliest such error.
	//
	// The JSON null value unmarshals into an interface, map, pointer, or slice
	// by setting that Go value to nil. Because null is often used in JSON to mean
	// ``not present,'' unmarshaling a JSON null into any other Go type has no effect
	// on the value and produces no error.
	//
	// When unmarshaling quoted strings, invalid UTF-8 or
	// invalid UTF-16 surrogate pairs are not treated as an error.
	// Instead, they are replaced by the Unicode replacement
	// character U+FFFD.
	//
	Unmarshal(data []byte, v interface{}) error
}

func New() IJSON {
	return _IJSON{}
}

type _IJSON struct{}

func (ijson _IJSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (ijson _IJSON) NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

func (ijson _IJSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
