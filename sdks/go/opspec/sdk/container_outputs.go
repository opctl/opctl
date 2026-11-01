package sdk

import (
	"fmt"
	"reflect"
)

// ContainerOutput requests a named JSON output as a ref-only handle. Keep the
// returned producer in the graph. T expresses the author's expected use; strict
// runtime type declarations are separate from these bindings.
func ContainerOutput[T Handle](c ContainerCall, name string) (result ContainerCall, value T) {
	result = c
	typ := reflect.TypeOf((*T)(nil)).Elem()
	h, ok := any(value).(inputHandle)
	if !ok || typ.Kind() == reflect.Pointer || reflect.TypeOf(h.bind(inputValue{})) != typ {
		if result.err == nil {
			result.err = fmt.Errorf("ContainerOutput needs a String, Number, Boolean, Array, or Object handle value")
		}
		return
	}
	fail := func(err error) (ContainerCall, T) {
		if result.err == nil {
			result.err = err
		}
		return result, h.bind(inputValue{err: result.err}).(T)
	}
	kind, _ := h.inputType()
	switch kind {
	case "string", "number", "boolean", "array", "object":
	default:
		return fail(fmt.Errorf("ContainerOutput does not support %s handles", kind))
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
		return result, h.bind(inputValue{name: ref[2 : len(ref)-1], result: true, err: result.err}).(T)
	}
	v := newCallResult(kind)
	result.container.Outputs = make(map[string]string, len(c.container.Outputs)+1)
	result.container.outputTypes = make(map[string]reflect.Type, len(c.container.outputTypes)+1)
	for key, ref := range c.container.Outputs {
		result.container.Outputs[key] = ref
		result.container.outputTypes[key] = c.container.outputTypes[key]
	}
	result.container.Outputs[name] = v.handle()
	result.container.outputTypes[name] = typ
	if result.err != nil {
		v.err = result.err
	}
	return result, h.bind(v).(T)
}
