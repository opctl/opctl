package sdk

import (
	"fmt"
	"reflect"
)

// ContainerOutput requests a named JSON output as a ref-only handle. Keep the
// returned producer in the graph. T expresses the author's expected use; strict
// runtime type declarations are separate from these bindings.
func ContainerOutput[T Handle](c ContainerCall, name string) (ContainerCall, T) {
	var value T
	typ := reflect.TypeOf((*T)(nil)).Elem()
	h, _ := any(value).(inputHandle)
	result, v := registerContainerOutput(c, name, typ, h)
	if h == nil || typ.Kind() == reflect.Pointer || reflect.TypeOf(h.bind(inputValue{})) != typ {
		return result, value
	}
	return result, h.bind(v).(T)
}

func registerContainerOutput(c ContainerCall, name string, typ reflect.Type, h inputHandle) (ContainerCall, inputValue) {
	fail := func(err error) (ContainerCall, inputValue) {
		if c.err == nil {
			c.err = err
		}
		return c, inputValue{err: c.err}
	}
	if h == nil || typ.Kind() == reflect.Pointer || reflect.TypeOf(h.bind(inputValue{})) != typ {
		return fail(fmt.Errorf("ContainerOutput needs a String, Number, Boolean, Array, or Object handle value"))
	}
	kind, _ := h.inputType()
	switch kind {
	case "string", "number", "boolean", "array", "object":
	default:
		return fail(fmt.Errorf("ContainerOutput does not support %s handles", kind))
	}
	if c.err != nil {
		return fail(c.err)
	}
	if activePlan == nil {
		return fail(fmt.Errorf("ContainerOutput must be called while building the op graph"))
	}
	if !outputName.MatchString(name) {
		return fail(fmt.Errorf("invalid container output name %q", name))
	}
	if ref, exists := c.container.Outputs[name]; exists {
		if c.container.outputTypes[name] != typ {
			return fail(fmt.Errorf("container output %q requested with different handle types", name))
		}
		return c, inputValue{name: ref[2 : len(ref)-1], result: true}
	}
	v := newCallResult(kind)
	bindings := make(map[string]string, len(c.container.Outputs)+1)
	types := make(map[string]reflect.Type, len(c.container.outputTypes)+1)
	for key, ref := range c.container.Outputs {
		bindings[key] = ref
		types[key] = c.container.outputTypes[key]
	}
	bindings[name], types[name] = v.handle(), typ
	c.container.Outputs, c.container.outputTypes = bindings, types
	return c, v
}
