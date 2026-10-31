package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Handle is a typed input binding in an op graph.
type Handle interface{ handle() string }
type inputHandle interface {
	Handle
	inputType() (string, reflect.Type)
	bind(inputValue) Handle
}
type inputValue struct {
	name   string
	value  json.RawMessage
	secret bool
}
type inputError struct{ error }

func (v inputValue) handle() string {
	if v.name == "" {
		panic(inputError{fmt.Errorf("handle is not bound to an input")})
	}
	return "$(" + v.name + ")"
}
func readValue[T any](v inputValue) T {
	v.handle()
	if v.secret {
		panic(inputError{fmt.Errorf("input %q is secret; its value can't be read during plan, bind it with EnvFrom", v.name)})
	}
	var value T
	if err := json.Unmarshal(v.value, &value); err != nil {
		panic(inputError{fmt.Errorf("input %q has an invalid value for its Go type", v.name)})
	}
	return value
}

// String is a string input handle.
type String struct{ inputValue }

func (String) inputType() (string, reflect.Type) { return "string", reflect.TypeOf("") }
func (String) bind(v inputValue) Handle          { return String{v} }
func (v String) Value() string                   { return readValue[string](v.inputValue) }

// Number is a number input handle.
type Number struct{ inputValue }

func (Number) inputType() (string, reflect.Type) { return "number", reflect.TypeOf(float64(0)) }
func (Number) bind(v inputValue) Handle          { return Number{v} }
func (v Number) Value() float64                  { return readValue[float64](v.inputValue) }

// Boolean is a boolean input handle.
type Boolean struct{ inputValue }

func (Boolean) inputType() (string, reflect.Type) { return "boolean", reflect.TypeOf(false) }
func (Boolean) bind(v inputValue) Handle          { return Boolean{v} }
func (v Boolean) Value() bool                     { return readValue[bool](v.inputValue) }

// Array is an array input handle with items of type T.
type Array[T any] struct{ inputValue }

func (Array[T]) inputType() (string, reflect.Type) {
	return "array", reflect.TypeOf((*[]T)(nil)).Elem()
}
func (Array[T]) bind(v inputValue) Handle { return Array[T]{v} }
func (v Array[T]) Value() []T             { return readValue[[]T](v.inputValue) }

// Object is an object input handle decoded into a struct or string-keyed map T.
type Object[T any] struct{ inputValue }

func (Object[T]) inputType() (string, reflect.Type) {
	return "object", reflect.TypeOf((*T)(nil)).Elem()
}
func (Object[T]) bind(v inputValue) Handle { return Object[T]{v} }
func (v Object[T]) Value() T               { return readValue[T](v.inputValue) }

// Dir is a directory input handle. Its contents cannot be read during plan.
type Dir struct{ inputValue }

func (Dir) inputType() (string, reflect.Type) { return "dir", nil }
func (Dir) bind(v inputValue) Handle          { return Dir{v} }

// File is a file input handle. Its contents cannot be read during plan.
type File struct{ inputValue }

func (File) inputType() (string, reflect.Type) { return "file", nil }
func (File) bind(v inputValue) Handle          { return File{v} }

// Socket is a socket input handle.
type Socket struct{ inputValue }

func (Socket) inputType() (string, reflect.Type) { return "socket", nil }
func (Socket) bind(v inputValue) Handle          { return Socket{v} }
