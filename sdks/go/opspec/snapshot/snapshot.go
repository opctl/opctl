// Package snapshot expands the input-only control flow of an op graph without running calls.
package snapshot

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"strings"
)

// Expand copies spec and expands predicates and loops that can be resolved from args.
// Inputs remain references; only iteration variables are substituted. Large or
// runtime-dependent loops and predicates stay in the snapshot.
func Expand(spec *model.OpSpec, args map[string]any) (*model.OpSpec, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	var out model.OpSpec
	if err = json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	scope := map[string]*model.Value{}
	for k, raw := range args {
		v, err := value.Interpret(escapeValue(raw), nil)
		if err != nil {
			return nil, fmt.Errorf("plan input %s: %w", k, err)
		}
		scope[k] = &v
	}
	out.Run, err = expand(out.Run, scope, map[string]*model.Value{})
	return &out, err
}

// escapeValue keeps input text from being interpreted as another reference.
func escapeValue(v any) any {
	switch v := v.(type) {
	case string:
		return strings.ReplaceAll(v, "$", "\\$")
	case []any:
		out := make([]any, len(v))
		for i, x := range v {
			out[i] = escapeValue(x)
		}
		return out
	case map[string]any:
		out := map[string]any{}
		for k, x := range v {
			out[k] = escapeValue(x)
		}
		return out
	default:
		return v
	}
}
func copyScope(in map[string]*model.Value) map[string]*model.Value {
	out := map[string]*model.Value{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func known(v any, scope map[string]*model.Value) bool {
	switch v := v.(type) {
	case string:
		ok := true
		_ = interpolater.Scan(v, func(string) error { return nil }, func(ref string) error {
			body := ref[2 : len(ref)-1]
			root := strings.FieldsFunc(body, func(r rune) bool { return r == '/' || r == '.' || r == '[' })
			if len(root) == 0 || scope[root[0]] == nil || strings.Contains(body, "$(") {
				ok = false
			}
			return nil
		})
		return ok
	case []any:
		for _, x := range v {
			if !known(x, scope) {
				return false
			}
		}
	case map[string]any:
		for _, x := range v {
			if !known(x, scope) {
				return false
			}
		}
	}
	return true
}
func predicatesKnown(ps []*model.PredicateSpec, scope map[string]*model.Value) bool {
	for _, p := range ps {
		if p.Exists != nil || p.NotExists != nil {
			return false
		}
	}
	b, _ := json.Marshal(ps)
	var raw any
	_ = json.Unmarshal(b, &raw)
	return known(raw, scope)
}

func expand(c *model.CallSpec, scope, locals map[string]*model.Value) (*model.CallSpec, error) {
	if c == nil {
		return nil, nil
	}
	if c.If != nil && predicatesKnown(*c.If, scope) {
		yes, err := predicates.Interpret(*c.If, scope)
		if err != nil {
			return nil, err
		}
		if !yes {
			return nil, nil
		}
		c.If = nil
	}
	var rangeValue any
	var vars *model.LoopVarsSpec
	var body *model.CallSpec
	parallel := false
	if c.SerialLoop != nil && len(c.SerialLoop.Until) == 0 {
		rangeValue = c.SerialLoop.Range
		vars = c.SerialLoop.Vars
		body = &c.SerialLoop.Run
	}
	if c.ParallelLoop != nil {
		rangeValue = c.ParallelLoop.Range
		vars = c.ParallelLoop.Vars
		body = &c.ParallelLoop.Run
		parallel = true
	}
	if body != nil && rangeValue != nil && known(rangeValue, scope) {
		ranged, err := loopable.Interpret(rangeValue, scope)
		if err != nil {
			return nil, err
		}
		count := 0
		if ranged.Array != nil {
			count = len(*ranged.Array)
		} else if ranged.Object != nil {
			count = len(*ranged.Object)
		}
		if count <= 1000 {
			children := []*model.CallSpec{}
			loopScope := copyScope(scope)
			for i := 0; i < count; i++ {
				base := scope
				if !parallel {
					base = loopScope
				}
				iterScope, err := iteration.Scope(i, base, rangeValue, vars)
				if err != nil {
					return nil, err
				}
				iterLocals := copyScope(locals)
				if vars != nil {
					for _, name := range []*string{vars.Index, vars.Key, vars.Value} {
						if name != nil {
							key := strings.TrimSuffix(strings.TrimPrefix(*name, "$("), ")")
							if v := iterScope[key]; v != nil {
								iterLocals[key] = v
							}
						}
					}
				}
				b, _ := json.Marshal(body)
				var child model.CallSpec
				_ = json.Unmarshal(b, &child)
				result, err := expand(&child, iterScope, iterLocals)
				if err != nil {
					return nil, err
				}
				if result != nil {
					children = append(children, result)
				}
				if !parallel {
					loopScope = iterScope
				}
			}
			for k := range scope {
				if loopScope[k] == nil {
					delete(scope, k)
				}
			}
			c.SerialLoop = nil
			c.ParallelLoop = nil
			if parallel {
				c.Parallel = &children
			} else {
				c.Serial = &children
			}
		}
	}
	for _, group := range []*[]*model.CallSpec{c.Serial, c.Parallel} {
		if group == nil {
			continue
		}
		children := []*model.CallSpec{}
		merged := copyScope(scope)
		for _, child := range *group {
			childScope := scope
			if c.Parallel != nil {
				childScope = copyScope(scope)
			}
			result, err := expand(child, childScope, locals)
			if err != nil {
				return nil, err
			}
			if result != nil {
				children = append(children, result)
			}
			if c.Parallel != nil {
				for k := range merged {
					if childScope[k] == nil {
						delete(merged, k)
					}
				}
			}
		}
		if c.Parallel != nil {
			for k := range scope {
				if merged[k] == nil {
					delete(scope, k)
				}
			}
		}
		*group = children
		if len(children) == 0 {
			return nil, nil
		}
		if len(children) == 1 && c.If == nil && c.Name == nil && len(c.Needs) == 0 {
			return children[0], nil
		}
	}
	// An unresolved loop's variables have no plan-time values. Leave its body intact.
	if c.SerialLoop == nil && c.ParallelLoop == nil {
		b, _ := json.Marshal(c)
		var raw any
		_ = json.Unmarshal(b, &raw)
		raw, err := substitute(raw, locals, false)
		if err != nil {
			return nil, err
		}
		b, err = json.Marshal(raw)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(b, c); err != nil {
			return nil, err
		}
	}
	if c.Container != nil {
		for _, ref := range []string{c.Container.Stdout, c.Container.Stderr} {
			shadow(ref, scope, locals)
		}
		for _, ref := range c.Container.Outputs {
			shadow(ref, scope, locals)
		}
	}
	if c.Op != nil {
		for _, ref := range c.Op.Outputs {
			shadow(ref, scope, locals)
		}
	}
	return c, nil
}
func shadow(ref string, scope, locals map[string]*model.Value) {
	if strings.HasPrefix(ref, "$(") && strings.HasSuffix(ref, ")") {
		name := ref[2 : len(ref)-1]
		delete(scope, name)
		delete(locals, name)
	}
}

func substitute(v any, locals map[string]*model.Value, typed bool) (any, error) {
	switch v := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, x := range v {
			next := typed || k == "inputs"
			y, err := substitute(x, locals, next)
			if err != nil {
				return nil, err
			}
			out[k] = y
		}
		return out, nil
	case []any:
		out := make([]any, len(v))
		for i, x := range v {
			y, err := substitute(x, locals, typed)
			if err != nil {
				return nil, err
			}
			out[i] = y
		}
		return out, nil
	case string:
		if len(locals) == 0 {
			return v, nil
		}
		if typed && strings.HasPrefix(v, "$(") && strings.HasSuffix(v, ")") {
			if val := locals[v[2:len(v)-1]]; val != nil {
				raw, err := val.Unbox()
				return escapeValue(raw), err
			}
		}
		var out strings.Builder
		err := interpolater.Scan(v, func(s string) error { out.WriteString(strings.ReplaceAll(s, "$", "\\$")); return nil }, func(ref string) error {
			if val := locals[ref[2:len(ref)-1]]; val != nil {
				str, err := coerce.ToString(val)
				if err != nil {
					return err
				}
				out.WriteString(strings.ReplaceAll(*str.String, "$", "\\$"))
				return nil
			}
			// Nested iteration variables can occur inside a path reference.
			body, err := substitute(ref[2:len(ref)-1], locals, false)
			if err != nil {
				return err
			}
			out.WriteString("$(" + body.(string) + ")")
			return nil
		})
		return out.String(), err
	default:
		return v, nil
	}
}
