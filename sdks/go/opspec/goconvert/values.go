package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"strings"
)

func (a *Analysis) value(raw any, path, position, kind string, scope map[string]*binding, state walkState) *boundValue {
	v := &boundValue{Raw: raw, Kind: kind, Mode: "literal"}
	a.values[path] = v
	if position == "input" {
		switch raw := raw.(type) {
		case []any:
			v.Items = make([]*boundValue, len(raw))
			for i, x := range raw {
				v.Items[i] = a.value(x, fmt.Sprintf("%s[%d]", path, i), position, "", scope, state)
			}
			return v
		case map[string]any:
			v.Members = map[string]*boundValue{}
			for _, key := range keys(raw) {
				v.Members[key] = a.value(raw[key], path+"."+key, position, "", scope, state)
			}
			return v
		}
	}
	text, ok := raw.(string)
	if !ok {
		if raw == nil {
			return v
		}
		b, _ := json.Marshal(raw)
		text = string(b)
	}
	literals := 0
	_ = interpolater.Scan(text, func(s string) error { literals += len(s); return nil }, func(ref string) error { v.Refs = append(v.Refs, a.reference(ref, path, scope, state)); return nil })
	if len(v.Refs) == 0 {
		return v
	}
	v.Whole = literals == 0 && len(v.Refs) == 1 && text == v.Refs[0].Raw
	v.Mode = "handle"
	readable := true
	for _, ref := range v.Refs {
		b := ref.Binding
		if b == nil {
			if !(isOpPath(ref.Body) && position != "literal") {
				a.refuse(path, "unknown reference "+ref.Raw)
			}
			readable = false
			continue
		}
		if b.Source == "produced" && (!v.Whole || ref.Suffix != "" || ref.Nested) {
			a.refuse(path, "templates, paths, properties and items over produced values can't be converted")
		}
		if !planReadable(b) || ref.Suffix != "" || ref.Nested {
			readable = false
		}
		if position == "literal" && b.Source == "loop" && b.Scalar == "any" {
			a.refuse(path, "untyped loop item in a string position; add items: {type: string|number|boolean} to the array input")
		}
	}
	if position == "literal" || position == "predicate" {
		if !readable {
			a.refuse(path, "this position needs a plan-readable scalar; results, secrets and filesystem handles can't be read during plan")
		}
		v.Mode = "text"
		return v
	}
	if v.Whole {
		ref := v.Refs[0]
		if ref.Binding != nil && ref.Suffix == "" && !ref.Nested {
			v.Kind = ref.Binding.Kind
			return v
		}
		if isOpPath(ref.Body) && !ref.Nested && !strings.ContainsAny(ref.Body, "$()") {
			if kind == "" {
				v.Kind = "dir"
			}
			v.Mode = "path"
			return v
		}
		if ref.Binding != nil && ref.Binding.Source == "input" && ref.Binding.Kind == "dir" && strings.HasPrefix(ref.Suffix, "/") && !ref.Nested && safeSubpath(ref.Suffix[1:]) {
			if kind == "" {
				v.Kind = "dir"
			}
			v.Mode = "path"
			return v
		}
	}
	if readable {
		v.Mode = "text"
		return v
	}
	v.Mode = "raw"
	if v.Kind == "" {
		v.Kind = "string"
	}
	v.Reason = "input expression kept as RawRef to match op.yml exactly; use typed handles when the expression has no dynamic paths or properties"
	for _, r := range v.Refs {
		if r.Nested {
			v.Reason = "nested reference kept as RawRef to match op.yml exactly; use Dir.Dir or Dir.File when every subpath is relative and contains no .. elements"
		}
	}
	a.Notes = append(a.Notes, Note{path, v.Reason})
	return v
}
func (a *Analysis) reference(raw, path string, scope map[string]*binding, state walkState) reference {
	body := raw[2 : len(raw)-1]
	r := reference{Raw: raw, Body: body, Nested: strings.Contains(body, "$(")}
	if isOpPath(body) {
		return r
	}
	end := len(body)
	for i, c := range body {
		if c == '/' || c == '.' || c == '[' || c == '$' {
			end = i
			break
		}
	}
	name := body[:end]
	r.Suffix = body[end:]
	r.Binding = scope[name]
	if strings.HasPrefix(name, "opctl_") {
		a.refuse(path, "references starting with opctl_ collide with generated names")
	}
	if r.Binding != nil {
		b := r.Binding
		if b.Loop != "" && !strings.HasPrefix(state.loop, b.Loop) {
			a.refuse(path, "loop-body producer is referenced outside its loop")
		}
		if b.Branch != "" && !strings.HasPrefix(state.branch, b.Branch) {
			a.refuse(path, "if-branch producer is referenced outside its branch")
		}
	}
	if r.Nested {
		_ = interpolater.Scan(body, func(string) error { return nil }, func(nested string) error {
			inner := a.reference(nested, path, scope, state)
			if inner.Binding == nil && !isOpPath(inner.Body) {
				a.refuse(path, "unknown reference "+nested)
			}
			if inner.Binding != nil && inner.Binding.Source == "produced" {
				a.refuse(path, "nested references over produced values can't be converted")
			}
			return nil
		})
	}
	return r
}
func isOpPath(body string) bool {
	return strings.HasPrefix(body, "./") || strings.HasPrefix(body, "../") || strings.HasPrefix(body, "/")
}
func safeSubpath(s string) bool {
	if s == "" || strings.HasPrefix(s, "/") || strings.ContainsAny(s, "$()") {
		return false
	}
	for _, part := range strings.Split(s, "/") {
		if part == ".." {
			return false
		}
	}
	return true
}
func planReadable(b *binding) bool {
	if b.Source == "loop" {
		return b.Scalar == "string" || b.Scalar == "int" || b.Scalar == "float64" || b.Scalar == "bool"
	}
	return b.Source == "input" && !secret(b.Param) && (b.Kind == "string" || b.Kind == "number" || b.Kind == "boolean")
}
