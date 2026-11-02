package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"reflect"
	"sort"
	"strings"
)

// InputSets returns at most sixteen deterministic, validated public input sets.
func InputSets(spec *model.OpSpec, opDir string, supplied map[string]any) ([]map[string]any, error) {
	publicParams := goop.PlanParams(spec.Inputs)
	defaults, err := params.ApplyDefaults(nil, publicParams, opDir, "")
	if err != nil {
		return nil, fmt.Errorf("verification input defaults: %w", err)
	}
	base := map[string]any{}
	candidates := map[string][]any{}
	for _, name := range keys(publicParams) {
		p := publicParams[name]
		kind, body := paramParts(p)
		cs, _ := body["constraints"].(map[string]any)
		values := []any{}
		if v, ok := supplied[name]; ok {
			values = append(values, v)
		}
		if boxed := defaults[name]; boxed != nil {
			v, err := boxed.Unbox()
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		}
		// Literal operands supply boundary and opposite-side candidates for predicates.
		var visit func(any)
		visit = func(v any) {
			switch x := v.(type) {
			case map[string]any:
				for _, key := range keys(x) {
					visit(x[key])
				}
			case []any:
				if len(x) == 2 && (x[0] == "$("+name+")" || x[1] == "$("+name+")") {
					v := x[0]
					if v == "$("+name+")" {
						v = x[1]
					}
					switch y := v.(type) {
					case float64:
						values = append(values, y, y-1, y+1)
					case string:
						values = append(values, y, y+"-x")
					case bool:
						values = append(values, y, !y)
					}
				}
				for _, y := range x {
					visit(y)
				}
			}
		}
		visit(jsonValue(spec.Run))
		values = append(values, sampleValues(kind, cs)...)
		seen := map[string]bool{}
		for _, v := range values {
			if validCandidate(p, kind, v) {
				b, _ := json.Marshal(v)
				if !seen[string(b)] {
					candidates[name] = append(candidates[name], v)
					seen[string(b)] = true
				}
			}
		}
		if len(candidates[name]) == 0 {
			return nil, fmt.Errorf("verification needs a value for input %s; pass -a %s=…", name, name)
		}
		base[name] = candidates[name][0]
	}
	sets := []map[string]any{}
	seen := map[string]bool{}
	limit := 15 // Reserve the final slot for the all-flipped set.
	add := func(set map[string]any) {
		b, _ := json.Marshal(set)
		if !seen[string(b)] && len(sets) < limit {
			sets = append(sets, set)
			seen[string(b)] = true
		}
	}
	add(base)
	flipped := copyValues(base)
	run, _ := json.Marshal(spec.Run)
	names := keys(publicParams)
	control := controlInputNames(spec.Run)
	sort.SliceStable(names, func(i, j int) bool { return control[names[i]] && !control[names[j]] })
	for _, name := range names {
		if !strings.Contains(string(run), "$("+name) {
			continue
		}
		n := 0
		for _, v := range candidates[name] {
			if n == 4 {
				break
			}
			set := copyValues(base)
			set[name] = v
			add(set)
			n++
		}
		for _, v := range candidates[name] {
			if !reflect.DeepEqual(v, base[name]) {
				flipped[name] = v
				break
			}
		}
	}
	limit = 16
	add(flipped)
	return sets, nil
}
func copyValues(in map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func validCandidate(p *model.ParamSpec, kind string, v any) bool {
	b, _ := json.Marshal(map[string]any{kind: v})
	var boxed model.Value
	if json.Unmarshal(b, &boxed) != nil {
		return false
	}
	return len(param.Validate(&boxed, p, true)) == 0
}
func sampleValues(kind string, cs map[string]any) []any {
	if cs == nil {
		cs = map[string]any{}
	}
	values := []any{}
	if e, ok := cs["enum"].([]any); ok {
		values = append(values, e...)
	}
	switch kind {
	case "string":
		samples := map[string]string{"semver": "1.2.3", "uri": "https://example.com", "email": "user@example.com", "hostname": "example.com", "ipv4": "192.0.2.1", "ipv6": "2001:db8::1", "date-time": "2026-01-02T03:04:05Z", "uuid": "123e4567-e89b-12d3-a456-426614174000"}
		if f, ok := cs["format"].(string); ok {
			if s := samples[f]; s != "" {
				values = append(values, s)
			}
		}
		values = append(values, "opctl-sample", "")
	case "number":
		values = append(values, float64(0), float64(1), float64(2))
		for _, key := range []string{"minimum", "maximum"} {
			if v, ok := cs[key]; ok {
				values = append(values, v)
			}
		}
	case "boolean":
		values = append(values, true, false)
	case "array":
		values = append(values, []any{})
		items, _ := cs["items"].(map[string]any)
		typ, _ := items["type"].(string)
		if typ == "integer" {
			typ = "number"
		}
		if typ == "" {
			typ = "string"
		}
		itemsSamples := sampleValues(typ, items)
		if len(itemsSamples) > 0 {
			values = append(values, []any{itemsSamples[0], itemsSamples[0]})
		}
	case "object":
		values = append(values, map[string]any{})
	}
	return values
}

func controlInputNames(run *model.CallSpec) map[string]bool {
	names := map[string]bool{}
	var collect func(any)
	collect = func(v any) {
		switch x := v.(type) {
		case string:
			_ = interpolater.Scan(x, func(string) error { return nil }, func(ref string) error { name, _ := rootRef(ref[2 : len(ref)-1]); names[name] = true; return nil })
		case []any:
			for _, v := range x {
				collect(v)
			}
		case map[string]any:
			for _, v := range x {
				collect(v)
			}
		}
	}
	var visit func(*model.CallSpec)
	visit = func(c *model.CallSpec) {
		if c == nil {
			return
		}
		collect(jsonValue(c.If))
		if c.SerialLoop != nil {
			collect(c.SerialLoop.Range)
			visit(&c.SerialLoop.Run)
		}
		if c.ParallelLoop != nil {
			collect(c.ParallelLoop.Range)
			visit(&c.ParallelLoop.Run)
		}
		for _, group := range []*[]*model.CallSpec{c.Serial, c.Parallel} {
			if group != nil {
				for _, child := range *group {
					visit(child)
				}
			}
		}
	}
	visit(run)
	return names
}
