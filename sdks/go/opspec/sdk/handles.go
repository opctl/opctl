package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// Handle is a typed input or output binding in an op graph.
type Handle interface{ handle() string }
type inputHandle interface {
	Handle
	inputType() (string, reflect.Type)
	bind(inputValue) Handle
}
type inputValue struct {
	raw    string
	name   string
	value  json.RawMessage
	result bool
	secret bool
	err    error
}
type inputError struct{ error }

func (v inputValue) handle() string {
	if v.err != nil {
		panic(inputError{v.err})
	}
	if v.raw != "" {
		return v.raw
	}
	if v.name == "" {
		panic(inputError{fmt.Errorf("handle is not bound to an input")})
	}
	return "$(" + v.name + ")"
}
func readValue[T any](v inputValue) T {
	v.handle()
	if v.result {
		panic(inputError{fmt.Errorf("call outputs can't be read during plan")})
	}
	if v.isRaw() {
		panic(inputError{fmt.Errorf("RawRef values can't be read during plan")})
	}
	if v.secret {
		panic(inputError{fmt.Errorf("input %q is secret; its value can't be read during plan, bind it with EnvFrom", v.name)})
	}
	var value T
	if err := json.Unmarshal(v.value, &value); err != nil {
		panic(inputError{fmt.Errorf("input %q has an invalid value for its Go type", v.name)})
	}
	return value
}

// String is a string handle.
type String struct{ inputValue }

func (String) inputType() (string, reflect.Type) { return "string", reflect.TypeOf("") }
func (String) bind(v inputValue) Handle          { return String{v} }
func (v String) Value() string                   { return readValue[string](v.inputValue) }

// Number is a number handle.
type Number struct{ inputValue }

func (Number) inputType() (string, reflect.Type) { return "number", reflect.TypeOf(float64(0)) }
func (Number) bind(v inputValue) Handle          { return Number{v} }
func (v Number) Value() float64                  { return readValue[float64](v.inputValue) }

// Boolean is a boolean handle.
type Boolean struct{ inputValue }

func (Boolean) inputType() (string, reflect.Type) { return "boolean", reflect.TypeOf(false) }
func (Boolean) bind(v inputValue) Handle          { return Boolean{v} }
func (v Boolean) Value() bool                     { return readValue[bool](v.inputValue) }

// Array is an array handle with items of type T.
type Array[T any] struct{ inputValue }

func (Array[T]) inputType() (string, reflect.Type) {
	return "array", reflect.TypeOf((*[]T)(nil)).Elem()
}
func (Array[T]) bind(v inputValue) Handle { return Array[T]{v} }
func (v Array[T]) Value() []T             { return readValue[[]T](v.inputValue) }

// Object is an object handle decoded into a struct or string-keyed map T.
type Object[T any] struct{ inputValue }

func (Object[T]) inputType() (string, reflect.Type) {
	return "object", reflect.TypeOf((*T)(nil)).Elem()
}
func (Object[T]) bind(v inputValue) Handle { return Object[T]{v} }
func (v Object[T]) Value() T               { return readValue[T](v.inputValue) }

// Dir is a directory handle. Its contents cannot be read during plan.
type Dir struct{ inputValue }

func (Dir) inputType() (string, reflect.Type) { return "dir", nil }
func (Dir) bind(v inputValue) Handle          { return Dir{v} }

// File is a file handle. Its contents cannot be read during plan.
type File struct{ inputValue }

func (File) inputType() (string, reflect.Type) { return "file", nil }
func (File) bind(v inputValue) Handle          { return File{v} }

// Socket is a socket handle.
type Socket struct{ inputValue }

func (Socket) inputType() (string, reflect.Type) { return "socket", nil }
func (Socket) bind(v inputValue) Handle          { return Socket{v} }

func (v inputValue) bindingValue() inputValue { return v }

// OpDir is a handle to a directory in the op's source, relative to the op
// directory. The path must start with ./ or ../.
func OpDir(path string) Dir { return Dir{opPath("OpDir", path)} }

// OpFile is a handle to a file in the op's source, relative to the op
// directory. The path must start with ./ or ../.
func OpFile(path string) File { return File{opPath("OpFile", path)} }

func opPath(fn, path string) inputValue {
	if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		return inputValue{err: fmt.Errorf("%s needs a path starting with ./ or ../, got %q", fn, path)}
	}
	if strings.ContainsAny(path, "$()") {
		return inputValue{err: fmt.Errorf("%s path %q must not contain $, ( or )", fn, path)}
	}
	return inputValue{name: path}
}

// Dir returns a handle to a directory beneath an input or op directory.
// The path is relative and can't contain .. elements.
func (v Dir) Dir(path string) Dir { return Dir{v.subpath(path)} }

// File returns a handle to a file beneath an input or op directory.
// The path is relative and can't contain .. elements.
func (v Dir) File(path string) File { return File{v.subpath(path)} }

func (v inputValue) subpath(path string) inputValue {
	switch {
	case v.err != nil:
		return v
	case v.isRaw():
		return inputValue{err: errors.New("subpaths of RawRef handles are not supported; pass the whole reference to RawRef")}
	case v.result || strings.HasPrefix(v.name, "opctl_"):
		return inputValue{err: errors.New("subpaths of call outputs are not supported yet")}
	case v.name == "":
		return inputValue{err: errors.New("handle is not bound to an input")}
	}
	if path == "" || strings.HasPrefix(path, "/") || strings.ContainsAny(path, "$()") || slices.Contains(strings.Split(path, "/"), "..") {
		return inputValue{err: fmt.Errorf("subpath %q must be relative, without .. elements, $, ( or )", path)}
	}
	return inputValue{name: strings.TrimSuffix(v.name, "/") + "/" + path}
}
