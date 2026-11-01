package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/sdk/refscan"
	"github.com/opctl/opctl/sdks/go/opspec/snapshot"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Compare checks the interface and independently expanded YAML graph against a Go plan.
func Compare(a *Analysis, described *model.OpSpec, planned *goop.Plan, inputs map[string]any) error {
	if err := difference("interface", canonicalInterface(a.Spec), canonicalInterface(described)); err != nil {
		return err
	}
	expanded, err := snapshot.Expand(a.Spec, inputs)
	if err != nil {
		return err
	}
	left, err := canonicalGraph(a, expanded.Run, nil, inputs)
	if err != nil {
		return err
	}
	right, err := canonicalGraph(a, planned.Run, planned.OutputBindings, inputs)
	if err != nil {
		return err
	}
	return difference("", left, right)
}
func jsonValue(v any) any {
	b, _ := json.Marshal(v)
	var out any
	_ = json.Unmarshal(b, &out)
	return out
}
func canonicalInterface(spec *model.OpSpec) any {
	raw := jsonValue(spec).(map[string]any)
	delete(raw, "run")
	delete(raw, "version")
	for _, group := range []string{"inputs", "outputs"} {
		params, _ := raw[group].(map[string]any)
		for _, v := range params {
			p := v.(map[string]any)
			for _, kind := range []string{"array", "object", "string", "number", "boolean", "dir", "file", "socket"} {
				if body, ok := p[kind].(map[string]any); ok {
					if p["description"] == nil && body["description"] != nil {
						p["description"] = body["description"]
					}
					delete(body, "description")
					if c, ok := body["constraints"].(map[string]any); ok && len(c) == 0 {
						delete(body, "constraints")
					}
					if body["isSecret"] == false {
						delete(body, "isSecret")
					}
				}
			}
		}
	}
	return raw
}
func difference(p string, a, b any) error {
	if reflect.DeepEqual(a, b) {
		return nil
	}
	if x, ok := a.(map[string]any); ok {
		if y, ok := b.(map[string]any); ok {
			all := map[string]bool{}
			for k := range x {
				all[k] = true
			}
			for k := range y {
				all[k] = true
			}
			for _, k := range keys(all) {
				next := k
				if p != "" {
					next = p + "." + k
				}
				if err := difference(next, x[k], y[k]); err != nil {
					return err
				}
			}
		}
	}
	if x, ok := a.([]any); ok {
		if y, ok := b.([]any); ok && len(x) == len(y) {
			for i := range x {
				if err := difference(fmt.Sprintf("%s[%d]", p, i), x[i], y[i]); err != nil {
					return err
				}
			}
		}
	}
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return fmt.Errorf("%s\n  op.yml: %s\n  op.go:  %s", p, x, y)
}

type refToken struct {
	Ref string `json:"ref"`
}
type textTokens struct {
	Parts []any `json:"text"`
}

type canonicalizer struct {
	a      *Analysis
	inputs map[string]any
	next   int
}

func canonicalGraph(a *Analysis, run *model.CallSpec, outputs map[string]string, inputs map[string]any) (any, error) {
	c := &canonicalizer{a: a, inputs: inputs}
	scope := map[string]string{}
	for name := range a.Spec.Inputs {
		scope[name] = "in:" + name
	}
	raw, _ := jsonValue(run).(map[string]any)
	tree, err := c.call(raw, scope)
	if err != nil {
		return nil, err
	}
	bound := map[string]any{}
	for _, name := range keys(a.Spec.Outputs) {
		ref := "$(" + name + ")"
		if outputs != nil {
			ref = outputs[name]
		}
		if ref == "" {
			bound[name] = refToken{"unbound:" + name}
		} else {
			bound[name] = c.value(ref, scope, true)
		}
	}
	return normalizeAliases(map[string]any{"run": tree, "outputs": bound}), nil
}
func cloneNames(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func rootRef(body string) (string, string) {
	i := strings.IndexAny(body, "/.[")
	if i < 0 {
		return body, ""
	}
	return body[:i], body[i:]
}
func (c *canonicalizer) reference(body string, scope map[string]string) any {
	// Resolve inner known values before resolving the enclosing path reference.
	var nested strings.Builder
	_ = refscan.ScanBody(body, func(s string) error { nested.WriteString(s); return nil }, func(ref string) error {
		v := c.reference(ref[2:len(ref)-1], scope)
		if s, ok := v.(string); ok {
			nested.WriteString(s)
		} else if m, ok := v.(refToken); ok {
			nested.WriteString("$(" + m.Ref + ")")
		} else {
			b, _ := json.Marshal(v)
			nested.Write(b)
		}
		return nil
	})
	body = nested.String()
	if body == "opctl_empty_string" {
		return ""
	}
	if strings.HasPrefix(body, "/") {
		body = "." + body
	}
	if body == "." || strings.HasPrefix(body, "./") || strings.HasPrefix(body, "../") {
		clean := path.Clean(body)
		if clean != ".." && !strings.HasPrefix(clean, "../") {
			clean = "./" + clean
		}
		return refToken{clean}
	}
	root, suffix := rootRef(body)
	if scope[root] == "in:"+root {
		if value, ok := c.inputs[root]; ok && suffix == "" {
			return value
		}
	}
	name := scope[root]
	if name == "" {
		name = "unbound:" + root
	}
	return refToken{name + suffix}
}
func textValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}
func (c *canonicalizer) value(v any, scope map[string]string, typed bool) any {
	switch v := v.(type) {
	case string:
		tokens := []any{}
		literal := ""
		flush := func() {
			if literal != "" {
				tokens = append(tokens, literal)
				literal = ""
			}
		}
		_ = interpolater.Scan(v, func(s string) error { literal += s; return nil }, func(ref string) error {
			x := c.reference(ref[2:len(ref)-1], scope)
			if _, ok := x.(refToken); ok {
				flush()
				tokens = append(tokens, x)
			} else {
				if typed && ref == v {
					flush()
					tokens = append(tokens, x)
				} else {
					literal += textValue(x)
				}
			}
			return nil
		})
		flush()
		if len(tokens) == 0 {
			return ""
		}
		if len(tokens) == 1 {
			return tokens[0]
		}
		return textTokens{tokens}
	case map[string]any:
		out := map[string]any{}
		for k, x := range v {
			out[k] = c.value(x, scope, typed)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, x := range v {
			out[i] = c.value(x, scope, typed)
		}
		return out
	default:
		if !typed && v != nil {
			return textValue(v)
		}
		return v
	}
}
func (c *canonicalizer) producer(ref string, scope map[string]string) any {
	name := strings.TrimSuffix(strings.TrimPrefix(ref, "$("), ")")
	c.next++
	id := fmt.Sprintf("v%d", c.next)
	scope[name] = id
	return refToken{id}
}
func (c *canonicalizer) call(raw map[string]any, scope map[string]string) (any, error) {
	if raw == nil {
		return nil, nil
	}
	for _, kind := range []string{"serial", "parallel"} {
		if children, ok := raw[kind].([]any); ok {
			out := []any{}
			merged := cloneNames(scope)
			for _, child := range children {
				childScope := scope
				if kind == "parallel" {
					childScope = cloneNames(scope)
				}
				x, err := c.call(child.(map[string]any), childScope)
				if err != nil {
					return nil, err
				}
				if x != nil {
					out = append(out, x)
				}
				if kind == "parallel" {
					for k, v := range childScope {
						if v != scope[k] {
							merged[k] = v
						}
					}
				}
			}
			if kind == "parallel" {
				for k, v := range merged {
					scope[k] = v
				}
			}
			if len(out) == 0 {
				return nil, nil
			}
			if len(out) == 1 {
				return out[0], nil
			}
			return map[string]any{kind: out}, nil
		}
	}
	if ctr, ok := raw["container"].(map[string]any); ok {
		out := map[string]any{}
		for _, kind := range []string{"dirs", "files"} {
			if mounts, ok := ctr[kind].(map[string]any); ok {
				for _, key := range keys(mounts) {
					v := mounts[key]
					if v == nil {
						v = "$(./" + strings.TrimPrefix(key, "/") + ")"
						mounts[key] = v
					}
					if s, ok := v.(string); ok && strings.HasPrefix(s, "$(") && strings.HasSuffix(s, ")") {
						body := s[2 : len(s)-1]
						if root, suffix := rootRef(body); root != "" && suffix == "" && scope[root] == "" {
							c.producer(s, scope)
						}
					}
				}
			}
		}
		for _, key := range keys(ctr) {
			v := ctr[key]
			if key == "stdout" || key == "stderr" || key == "outputs" {
				continue
			}
			if key == "envVars" {
				if env, ok := v.(map[string]any); ok {
					for k, x := range env {
						if x == nil {
							env[k] = "$(" + k + ")"
						} else if _, ok := x.(string); !ok {
							env[k] = textValue(x)
						}
					}
				}
			}
			out[key] = c.value(v, scope, key == "outputParams")
		}
		for _, key := range []string{"stdout", "stderr"} {
			if s, ok := ctr[key].(string); ok && s != "" {
				out[key] = c.producer(s, scope)
			}
		}
		if outputs, ok := ctr["outputs"].(map[string]any); ok {
			m := map[string]any{}
			for _, key := range keys(outputs) {
				m[key] = c.producer(outputs[key].(string), scope)
			}
			out["outputs"] = m
		}
		return map[string]any{"container": out}, nil
	}
	if op, ok := raw["op"].(map[string]any); ok {
		out := map[string]any{}
		for _, key := range keys(op) {
			if key == "outputs" {
				continue
			}
			v := op[key]
			if key == "ref" {
				normalized, err := NormalizeRef(c.a.Dir, v.(string))
				if err != nil {
					return nil, err
				}
				out[key] = normalized
				continue
			}
			if key == "inputs" {
				if args, ok := v.(map[string]any); ok {
					for k, x := range args {
						if x == nil {
							args[k] = "$(" + k + ")"
						}
					}
				}
			}
			out[key] = c.value(v, scope, true)
		}
		if outputs, ok := op["outputs"].(map[string]any); ok {
			m := map[string]any{}
			for _, key := range keys(outputs) {
				s, _ := outputs[key].(string)
				if s == "" {
					s = "$(" + key + ")"
				}
				m[key] = c.producer(s, scope)
			}
			out["outputs"] = m
		}
		return map[string]any{"op": out}, nil
	}
	return nil, fmt.Errorf("unexpanded call in verification")
}

// Stub calls bind every declared output. Unread op-call aliases have no effect;
// container output declarations still require their values and remain intact.
func normalizeAliases(doc map[string]any) any {

	reads := map[string]bool{}
	var collect func(any)
	collect = func(v any) {
		switch x := v.(type) {
		case refToken:
			root, _ := rootRef(x.Ref)
			reads[root] = true
		case textTokens:
			for _, p := range x.Parts {
				collect(p)
			}
		case []any:
			for _, p := range x {
				collect(p)
			}
		case map[string]any:
			for _, p := range x {
				collect(p)
			}
		}
	}
	var readCall func(any)
	readCall = func(v any) {
		call, ok := v.(map[string]any)
		if !ok {
			return
		}
		for _, kind := range []string{"serial", "parallel"} {
			if children, ok := call[kind].([]any); ok {
				for _, child := range children {
					readCall(child)
				}
			}
		}
		for _, kind := range []string{"container", "op"} {
			if body, ok := call[kind].(map[string]any); ok {
				for key, value := range body {
					if key != "outputs" && !(kind == "container" && (key == "stdout" || key == "stderr")) {
						collect(value)
					}
				}
			}
		}
	}
	readCall(doc["run"])
	collect(doc["outputs"])
	var prune func(any)
	prune = func(v any) {
		call, ok := v.(map[string]any)
		if !ok {
			return
		}
		if op, ok := call["op"].(map[string]any); ok {
			if outputs, ok := op["outputs"].(map[string]any); ok {
				for name, p := range outputs {
					if ref, ok := p.(refToken); ok && !reads[ref.Ref] {
						delete(outputs, name)
					}
				}
				if len(outputs) == 0 {
					delete(op, "outputs")
				}
			}
		}
		for _, kind := range []string{"serial", "parallel"} {
			if children, ok := call[kind].([]any); ok {
				for _, child := range children {
					prune(child)
				}
			}
		}
	}
	prune(doc["run"])
	ids := map[int]bool{}
	var visit func(any, func(refToken) refToken) any
	visit = func(v any, f func(refToken) refToken) any {
		switch x := v.(type) {
		case refToken:
			return f(x)
		case textTokens:
			for i, p := range x.Parts {
				x.Parts[i] = visit(p, f)
			}
			return x
		case []any:
			for i, p := range x {
				x[i] = visit(p, f)
			}
		case map[string]any:
			for k, p := range x {
				x[k] = visit(p, f)
			}
		}
		return v
	}
	visit(doc, func(ref refToken) refToken {
		root, _ := rootRef(ref.Ref)
		if strings.HasPrefix(root, "v") {
			if n, err := strconv.Atoi(root[1:]); err == nil {
				ids[n] = true
			}
		}
		return ref
	})
	ordered := make([]int, 0, len(ids))
	for n := range ids {
		ordered = append(ordered, n)
	}
	sort.Ints(ordered)
	names := map[string]string{}
	for i, n := range ordered {
		names[fmt.Sprintf("v%d", n)] = fmt.Sprintf("v%d", i+1)
	}
	return visit(doc, func(ref refToken) refToken {
		root, suffix := rootRef(ref.Ref)
		if name := names[root]; name != "" {
			ref.Ref = name + suffix
		}
		return ref
	})
}
