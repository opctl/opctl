package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Outputs declares required JSON members using a nonempty tagged struct of
// String, Number, Boolean, Array, or Object handles. Keep the returned call in
// the graph; its outputs are unavailable to Go control flow during plan.
func Outputs[T any](c ContainerCall) (result ContainerCall, out T) {
	result = c
	target := reflect.ValueOf(&out).Elem()
	defer func() {
		if result.err == nil || target.Kind() != reflect.Struct {
			return
		}
		for i := 0; i < target.NumField(); i++ {
			field := target.Field(i)
			if !field.CanSet() || !field.CanInterface() || field.Kind() == reflect.Pointer {
				continue
			}
			if h, ok := field.Interface().(inputHandle); ok {
				bound := reflect.ValueOf(h.bind(inputValue{err: result.err}))
				if bound.Type() == field.Type() {
					field.Set(bound)
				}
			}
		}
	}()
	fail := func(err error) (ContainerCall, T) {
		if result.err == nil {
			result.err = err
		}
		return result, out
	}
	if result.err != nil {
		return result, out
	}
	specs, fields, err := parameterSpecs[T]("container output")
	if err != nil {
		return fail(err)
	}
	if len(fields) == 0 {
		return fail(fmt.Errorf("container outputs must be a nonempty struct of handles"))
	}
	for _, field := range fields {
		if field.hasDefault {
			return fail(fmt.Errorf("container output %q cannot have a default; every declared member must be produced", field.name))
		}
		switch field.kind {
		case "string", "number", "boolean", "array", "object":
		default:
			return fail(fmt.Errorf("container output %q does not support %s; use mount/socket bindings", field.name, field.kind))
		}
	}
	// Declarations are private snapshots, including nested constraints and schemas.
	encoded, err := json.Marshal(c.container.OutputParams)
	if err != nil {
		return fail(err)
	}
	var declarations map[string]paramSpec
	if err = json.Unmarshal(encoded, &declarations); err != nil {
		return fail(err)
	}
	if declarations == nil {
		declarations = map[string]paramSpec{}
	}
	for _, field := range fields {
		declaration := specs[field.name]
		if previous, exists := declarations[field.name]; exists {
			before, _ := json.Marshal(previous)
			after, err := json.Marshal(declaration)
			if err != nil {
				return fail(err)
			}
			if !bytes.Equal(before, after) {
				return fail(fmt.Errorf("conflicting container output declaration %q", field.name))
			}
		}
		var value inputValue
		result, value = registerContainerOutput(result, field.name, target.Field(field.index).Type(), field.handle)
		if result.err != nil {
			return result, out
		}
		target.Field(field.index).Set(reflect.ValueOf(field.handle.bind(value)))
		declarations[field.name] = declaration
	}
	result.container.OutputParams = declarations
	return result, out
}
