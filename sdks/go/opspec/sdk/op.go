package sdk

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
)

// OpCall builds a call by ref. Methods return independent copies.
type OpCall struct {
	op  opCallSpec
	err error
}
type handleUse struct {
	ref string
	raw bool
}
type opCallSpec struct {
	Ref            string            `json:"ref"`
	Inputs         map[string]any    `json:"inputs,omitempty"`
	PullCreds      *credsSpec        `json:"pullCreds,omitempty"`
	Outputs        map[string]string `json:"outputs,omitempty"`
	outputTypes    map[string]reflect.Type
	uses, credUses []handleUse
}
type credsSpec struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CallOp calls a local op or a git ref such as repo#tag/path. It does no lookup during plan.
func CallOp(ref string) OpCall {
	c := OpCall{op: opCallSpec{Ref: ref}}
	if ref == "" || wholeReference.MatchString(ref) {
		c.err = fmt.Errorf("CallOp needs a nonempty literal op ref")
	}
	return c
}

// Inputs replaces the named arguments with a snapshot. Omitted arguments use callee defaults.
func (c OpCall) Inputs(values map[string]any) (result OpCall) {
	result = c
	defer recoverBuilder(&result.err)
	result.op.Inputs = map[string]any{}
	result.op.uses = nil
	for _, key := range sortedInputKeys(values) {
		encoded, uses, err := encodeArgument(reflect.ValueOf(values[key]), map[visit]bool{})
		if err != nil {
			if result.err == nil {
				result.err = fmt.Errorf("op input %q: %w", key, err)
			}
			return
		}
		if encoded == "" {
			encoded = "$(opctl_empty_string)"
		}
		result.op.Inputs[key] = encoded
		result.op.uses = append(result.op.uses, uses...)
	}
	return
}

// PullCreds sets literal strings or String handles for fetching the callee.
func (c OpCall) PullCreds(username, password any) (result OpCall) {
	result = c
	defer recoverBuilder(&result.err)
	creds := &credsSpec{}
	result.op.credUses = nil
	for i, value := range []any{username, password} {
		switch value.(type) {
		case string, String:
		default:
			if result.err == nil {
				result.err = fmt.Errorf("op pull credentials must be strings or String handles")
			}
			return
		}
		encoded, uses, err := encodeArgument(reflect.ValueOf(value), map[visit]bool{})
		if err != nil {
			if result.err == nil {
				result.err = err
			}
			return
		}
		if i == 0 {
			creds.Username = encoded.(string)
		} else {
			creds.Password = encoded.(string)
		}
		result.op.credUses = append(result.op.credUses, uses...)
	}
	result.op.PullCreds = creds
	return
}
func recoverBuilder(err *error) {
	if failure := recover(); failure != nil {
		if input, ok := failure.(inputError); ok {
			if *err == nil {
				*err = input.error
			}
		} else {
			panic(failure)
		}
	}
}
func (OpCall) call() {}
func (c OpCall) render() (*callSpec, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.op.Ref == "" {
		return nil, fmt.Errorf("CallOp needs a nonempty literal op ref")
	}
	return &callSpec{Op: &c.op}, nil
}

type visit struct {
	typ     reflect.Type
	pointer uintptr
}

func encodeArgument(v reflect.Value, visiting map[visit]bool) (any, []handleUse, error) {
	if !v.IsValid() {
		return nil, nil, fmt.Errorf("null arguments are not supported")
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, nil, fmt.Errorf("null arguments are not supported")
		}
		return encodeArgument(v.Elem(), visiting)
	}
	if h, ok := v.Interface().(Handle); ok {
		concrete, valid := h.(inputHandle)
		if !valid || v.Kind() == reflect.Pointer || reflect.TypeOf(concrete.bind(inputValue{})) != v.Type() {
			return nil, nil, fmt.Errorf("argument must be an opctl handle value")
		}
		ref := h.handle()
		raw, _ := h.(interface{ isRaw() bool })
		return ref, []handleUse{{ref: ref, raw: raw != nil && raw.isRaw()}}, nil
	}
	if v.Type() == reflect.TypeOf(json.Number("")) {
		n, err := v.Interface().(json.Number).Float64()
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, nil, fmt.Errorf("number must be finite")
		}
		return n, nil, nil
	}
	switch v.Kind() {
	case reflect.String:
		return escapeLiteral(v.String()), nil, nil
	case reflect.Bool:
		return v.Bool(), nil, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), nil, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), nil, nil
	case reflect.Float32, reflect.Float64:
		n := v.Float()
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return nil, nil, fmt.Errorf("number must be finite")
		}
		return n, nil, nil
	case reflect.Map, reflect.Slice:
		if v.IsNil() {
			return nil, nil, fmt.Errorf("null arguments are not supported")
		}
		key := visit{v.Type(), v.Pointer()}
		if visiting[key] {
			return nil, nil, fmt.Errorf("cyclic arguments are not supported")
		}
		visiting[key] = true
		defer delete(visiting, key)
	case reflect.Array:
	default:
		return nil, nil, fmt.Errorf("unsupported argument type %s", v.Type())
	}
	var uses []handleUse
	if v.Kind() == reflect.Map {
		if v.Type().Key().Kind() != reflect.String {
			return nil, nil, fmt.Errorf("object arguments need string keys")
		}
		result := map[string]any{}
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
		for _, key := range keys {
			value, nested, err := encodeArgument(v.MapIndex(key), visiting)
			if err != nil {
				return nil, nil, err
			}
			result[escapeLiteral(key.String())] = value
			uses = append(uses, nested...)
		}
		return result, uses, nil
	}
	result := make([]any, v.Len())
	for i := range result {
		value, nested, err := encodeArgument(v.Index(i), visiting)
		if err != nil {
			return nil, nil, err
		}
		result[i] = value
		uses = append(uses, nested...)
	}
	return result, uses, nil
}
func sortedInputKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// OpOutput requests a callee output as a typed, ref-only handle. Keep the returned call in the graph.
// The callee validates the actual value at execution; T describes its use in this graph.
func OpOutput[T Handle](c OpCall, name string) (result OpCall, value T) {
	result = c
	typ := reflect.TypeOf((*T)(nil)).Elem()
	h, ok := any(value).(inputHandle)
	if !ok || typ.Kind() == reflect.Pointer || reflect.TypeOf(h.bind(inputValue{})) != typ {
		if result.err == nil {
			result.err = fmt.Errorf("OpOutput needs an opctl handle value")
		}
		return
	}
	if activePlan == nil {
		if result.err == nil {
			result.err = fmt.Errorf("OpOutput must be called while building the op graph")
		}
		return result, h.bind(inputValue{err: result.err}).(T)
	}

	if !outputName.MatchString(name) {
		if result.err == nil {
			result.err = fmt.Errorf("invalid op output name %q", name)
		}
		return result, h.bind(inputValue{err: result.err}).(T)
	}
	if existing, ok := c.op.Outputs[name]; ok {
		if c.op.outputTypes[name] != typ {
			if result.err == nil {
				result.err = fmt.Errorf("op output %q requested with different handle types", name)
			}
			return
		}
		return result, h.bind(inputValue{name: existing[2 : len(existing)-1], result: true}).(T)
	}
	kind, _ := h.inputType()
	v := newOutput(kind)
	v.result = true
	if v.err != nil {
		if result.err == nil {
			result.err = v.err
		}
		return
	}
	site := activePlan.handles[v.name]
	site.opResult = true
	activePlan.handles[v.name] = site
	result.op.Outputs = map[string]string{}
	result.op.outputTypes = map[string]reflect.Type{}
	for key, ref := range c.op.Outputs {
		result.op.Outputs[key] = ref
		result.op.outputTypes[key] = c.op.outputTypes[key]
	}
	result.op.Outputs[name] = v.handle()
	result.op.outputTypes[name] = typ
	return result, h.bind(v).(T)
}

var outputName = regexp.MustCompile(`^[-_a-zA-Z0-9]+$`)
