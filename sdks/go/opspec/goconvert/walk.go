package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"strings"
)

func (a *Analysis) walk(c *model.CallSpec, path string, scope map[string]*binding, state walkState, root bool) *graph {
	if c == nil {
		if root {
			a.refuse(path, "Go ops can't plan an empty op graph")
		}
		return nil
	}
	g := &graph{Spec: c, Path: path}
	if c.Name != nil {
		a.refuse(path+".name", "call name isn't supported by the Go op SDK")
	}
	if len(c.Needs) > 0 {
		a.refuse(path+".needs", "needs isn't supported by the Go op SDK")
	}
	if c.Description != "" {
		a.Notes = append(a.Notes, Note{path + ".description", "call description becomes a Go comment"})
	}
	if c.If != nil {
		g.Condition = *c.If
		data, _ := json.Marshal(c.If)
		var ps []map[string]any
		_ = json.Unmarshal(data, &ps)
		constant := true
		for i, p := range ps {
			for kind, raw := range p {
				ppath := fmt.Sprintf("%s.if[%d].%s", path, i, kind)
				if kind == "exists" || kind == "notExists" {
					a.refuse(ppath, "exists and notExists predicates need runtime data")
					constant = false
					continue
				}
				operands, _ := raw.([]any)
				for j, x := range operands {
					v := a.value(x, fmt.Sprintf("%s[%d]", ppath, j), "predicate", "string", scope, state)
					if len(v.Refs) > 0 {
						constant = false
					}
				}
			}
		}
		if constant {
			yes, err := predicates.Interpret(*c.If, nil)
			if err != nil {
				a.refuse(path+".if", err.Error())
			} else {
				g.Fold = &yes
				a.Notes = append(a.Notes, Note{path + ".if", fmt.Sprintf("constant predicate folded to %t", yes)})
				if !yes {
					if root {
						a.refuse(path, "Go ops can't plan an empty op graph")
					}
					return g
				}
			}
		} else {
			if root {
				a.refuse(path+".if", "a root if can plan an empty op graph")
			}
			state.branch = path
		}
	}
	switch {
	case c.Container != nil:
		g.Kind = "container"
		a.container(c.Container, path+".container", scope, state)
	case c.Op != nil:
		g.Kind = "op"
		a.op(c.Op, path, scope, state)
	case c.Serial != nil || c.Parallel != nil:
		children := c.Serial
		g.Kind = "serial"
		parallel := false
		if c.Parallel != nil {
			children = c.Parallel
			g.Kind = "parallel"
			parallel = true
		}
		merged := cloneScope(scope)
		owners := map[string]string{}
		for i, child := range *children {
			childScope := scope
			if parallel {
				childScope = cloneScope(scope)
			}
			childPath := fmt.Sprintf("%s.%s[%d]", path, g.Kind, i)
			g.Children = append(g.Children, a.walk(child, childPath, childScope, state, false))
			if parallel {
				for name, b := range childScope {
					if b != scope[name] {
						if previous := owners[name]; previous != "" {
							a.refuse(childPath, "name "+name+" is produced by two parallel siblings (also "+previous+")")
						}
						owners[name] = childPath
						merged[name] = b
					}
				}
			}
		}
		if parallel {
			for name, b := range merged {
				scope[name] = b
			}
		}
		if root && len(*children) == 0 {
			a.refuse(path, "Go ops can't plan an empty op graph")
		}
	case c.SerialLoop != nil || c.ParallelLoop != nil:
		g.Kind = "serialLoop"
		var rangeValue any
		var vars *model.LoopVarsSpec
		var body *model.CallSpec
		if c.SerialLoop != nil {
			l := c.SerialLoop
			rangeValue, vars, body = l.Range, l.Vars, &l.Run
			if len(l.Until) > 0 {
				a.refuse(path+".serialLoop.until", "until needs runtime data and can't be converted")
			}
		}
		if c.ParallelLoop != nil {
			g.Kind = "parallelLoop"
			l := c.ParallelLoop
			rangeValue, vars, body = l.Range, l.Vars, &l.Run
		}
		lp := path + "." + g.Kind
		if rangeValue == nil {
			a.refuse(lp+".range", "a loop needs an input or literal range")
		}
		g.Range = a.value(rangeValue, lp+".range", "range", "", scope, state)
		rangeKind, scalar := literalRange(rangeValue)
		if len(g.Range.Refs) > 0 {
			if !g.Range.Whole || len(g.Range.Refs) != 1 || g.Range.Refs[0].Binding == nil || g.Range.Refs[0].Suffix != "" || g.Range.Refs[0].Nested {
				a.refuse(lp+".range", "loop range must be a literal or a whole non-secret input")
			} else {
				b := g.Range.Refs[0].Binding
				rangeKind = b.Kind
				if b.Source != "input" || secret(b.Param) {
					a.refuse(lp+".range", "loop range can't read results or secrets")
				}
				if rangeKind == "array" {
					scalar = a.loopInputs[b.Name]
					if scalar == "" {
						scalar = "any"
					}
				} else if rangeKind == "object" {
					scalar = "any"
				}
			}
		}
		if rangeKind == "number" {
			a.refuse(lp+".range", "numeric loop ranges aren't executable by the YAML interpreter")
		}
		if rangeKind != "array" && rangeKind != "object" && rangeKind != "number" {
			a.refuse(lp+".range", "loop range must be an array or object")
		}
		if root {
			empty := true
			switch r := rangeValue.(type) {
			case []any:
				empty = len(r) == 0
			case map[string]any:
				empty = len(r) == 0
			}
			if empty {
				a.refuse(lp, "a root loop can plan an empty op graph")
			}
		}
		loopScope := cloneScope(scope)
		g.Vars = map[string]*binding{}
		next := state
		next.loop = path
		next.parallelLoop = g.Kind == "parallelLoop"
		next.loopGraph = g
		if vars != nil {
			for _, entry := range []struct {
				which  string
				name   *string
				scalar string
			}{{"index", vars.Index, "int"}, {"key", vars.Key, "string"}, {"value", vars.Value, scalar}} {
				if entry.name == nil {
					continue
				}
				name, ok := wholeName(*entry.name)
				if !ok {
					a.refuse(lp+".vars."+entry.which, "loop variable must bind a whole name")
					continue
				}
				kind := "string"
				if entry.scalar == "int" || entry.scalar == "float64" {
					kind = "number"
				} else if entry.scalar == "bool" {
					kind = "boolean"
				}
				b := a.newBinding(name, kind, "loop", lp+".vars."+entry.which, next)
				b.Scalar = entry.scalar
				if b.Scalar == "" {
					b.Scalar = "any"
				}
				loopScope[name] = b
				g.Vars[entry.which] = b
			}
		}
		g.Body = a.walk(body, lp+".run", loopScope, next, false)
		for name, b := range loopScope {
			if b != scope[name] {
				scope[name] = b
			}
		}
	default:
		if root {
			a.refuse(path, "Go ops can't plan an empty op graph")
		}
	}
	return g
}
func literalRange(v any) (string, string) {
	switch v := v.(type) {
	case []any:
		scalar := ""
		for i, x := range v {
			k := "any"
			switch x.(type) {
			case string:
				k = "string"
			case float64:
				k = "float64"
			case bool:
				k = "bool"
			}
			if i == 0 {
				scalar = k
			} else if scalar != k {
				scalar = "any"
			}
		}
		if scalar == "" {
			scalar = "any"
		}
		return "array", scalar
	case map[string]any:
		return "object", "any"
	case float64:
		return "number", "int"
	}
	return "", ""
}

func (a *Analysis) container(c *model.ContainerCallSpec, path string, scope map[string]*binding, state walkState) {
	if c.Image != nil {
		a.value(c.Image.Ref, path+".image.ref", "literal", "string", scope, state)
		if c.Image.PullCreds != nil {
			a.refuse(path+".image.pullCreds", "image pullCreds aren't supported by the Go op SDK")
		}
		if c.Image.Platform != nil {
			a.refuse(path+".image.platform", "image platform isn't supported by the Go op SDK")
		}
	}
	if len(c.Ports) > 0 {
		a.refuse(path+".ports", "container ports aren't supported by the Go op SDK")
	}
	for _, entry := range []struct {
		key, kind string
		values    map[string]any
	}{{"dirs", "dir", c.Dirs}, {"files", "file", c.Files}} {
		for _, mount := range keys(entry.values) {
			raw := entry.values[mount]
			p := path + "." + entry.key + "." + mount
			if raw == nil {
				raw = "$(./" + strings.TrimPrefix(mount, "/") + ")"
			}
			text, ok := raw.(string)
			if !ok || !strings.HasPrefix(text, "$(") || !strings.HasSuffix(text, ")") {
				a.refuse(p, entry.kind+" initializers aren't supported by the Go op SDK")
				continue
			}
			if name, ok := wholeName(text); ok && scope[name] == nil {
				b := a.newBinding(name, entry.kind, "produced", p, state)
				b.Implicit = true
				if state.loop != "" && !state.parallelLoop {
					b.Loop = ""
					b.Branch = state.branch
					if state.loopGraph != nil {
						state.loopGraph.Hoists = append(state.loopGraph.Hoists, b)
					}
				}
				scope[name] = b
				a.produced[p] = append(a.produced[p], b)
			}
			v := a.value(raw, p, "handle", entry.kind, scope, state)
			if v.Whole && len(v.Refs) > 0 && v.Refs[0].Binding != nil && v.Refs[0].Suffix == "" && v.Refs[0].Binding.Kind != entry.kind {
				a.refuse(p, "mount needs a "+entry.kind+" handle")
			}
		}
	}
	if c.Cmd != nil {
		items, ok := c.Cmd.([]any)
		if !ok {
			a.refuse(path+".cmd", "cmd as one reference isn't supported")
		} else {
			for i, item := range items {
				a.value(item, fmt.Sprintf("%s.cmd[%d]", path, i), "literal", "string", scope, state)
			}
		}
	}
	if c.EnvVars != nil {
		env, ok := c.EnvVars.(map[string]any)
		if !ok {
			a.refuse(path+".envVars", "envVars as one reference isn't supported")
		} else {
			for _, name := range keys(env) {
				raw := env[name]
				if raw == nil {
					raw = "$(" + name + ")"
				}
				a.value(raw, path+".envVars."+name, "handle", "string", scope, state)
			}
		}
	}
	for _, name := range keys(c.Sockets) {
		v := a.value(c.Sockets[name], path+".sockets."+name, "handle", "socket", scope, state)
		if !v.Whole || len(v.Refs) != 1 || v.Refs[0].Binding == nil || v.Refs[0].Binding.Source != "input" || v.Refs[0].Binding.Kind != "socket" {
			a.refuse(path+".sockets."+name, "sockets must bind a whole socket input")
		}
	}
	if c.WorkDir != "" {
		a.value(c.WorkDir, path+".workDir", "literal", "string", scope, state)
	}
	if c.Name != nil {
		a.value(*c.Name, path+".name", "literal", "string", scope, state)
	}
	for _, stream := range []struct{ key, ref string }{{"stdout", c.Stdout}, {"stderr", c.Stderr}} {
		if stream.ref != "" {
			a.producer(stream.ref, stream.key, "string", path+"."+stream.key, scope, state)
		}
	}
	for _, name := range keys(c.Outputs) {
		kind := "string"
		if p := c.OutputParams[name]; p != nil {
			kind, _ = paramParts(p)
			a.checkParam(path+".outputParams."+name, name, p)
		} else {
			if target, ok := wholeName(c.Outputs[name]); ok {
				if p := a.Spec.Outputs[target]; p != nil {
					kind, _ = paramParts(p)
				} else {
					kind = a.consumerKind(target, kind)
				}
			}
		}
		a.producer(c.Outputs[name], name, kind, path+".outputs."+name, scope, state)
	}
	for name := range c.OutputParams {
		if _, ok := c.Outputs[name]; !ok {
			a.refuse(path+".outputParams."+name, "typed container output has no output binding")
		}
	}
}
func (a *Analysis) consumerKind(name, fallback string) string {
	for path, spec := range a.callees {
		var c *model.OpCallSpec
		if n := a.nodes[path+".op"]; n != nil {
			b, _ := json.Marshal(nodeValue(n))
			_ = json.Unmarshal(b, &c)
		}
		if c != nil {
			for input, raw := range c.Inputs {
				if raw == "$("+name+")" {
					if p := spec.Inputs[input]; p != nil {
						kind, _ := paramParts(p)
						return kind
					}
				}
			}
		}
	}
	return fallback
}
func nodeValue(n interface{ Decode(any) error }) any { var out any; _ = n.Decode(&out); return out }
func (a *Analysis) op(c *model.OpCallSpec, path string, scope map[string]*binding, state walkState) {
	callee := a.callees[path]
	if callee == nil {
		return
	}
	for _, name := range keys(c.Inputs) {
		raw := c.Inputs[name]
		if raw == nil {
			raw = "$(" + name + ")"
		}
		kind, _ := paramParts(callee.Inputs[name])
		if kind == "" {
			a.refuse(path+".op.inputs."+name, "callee has no input named "+name)
		}
		a.value(raw, path+".op.inputs."+name, "input", kind, scope, state)
	}
	if c.PullCreds != nil {
		a.value(c.PullCreds.Username, path+".op.pullCreds.username", "handle", "string", scope, state)
		a.value(c.PullCreds.Password, path+".op.pullCreds.password", "handle", "string", scope, state)
	}
	for _, name := range keys(c.Outputs) {
		kind, _ := paramParts(callee.Outputs[name])
		if kind == "" {
			a.refuse(path+".op.outputs."+name, "callee has no output named "+name)
			kind = "string"
		}
		a.producer(c.Outputs[name], name, kind, path+".op.outputs."+name, scope, state)
	}
}
func (a *Analysis) producer(ref, fallback, kind, path string, scope map[string]*binding, state walkState) {
	if ref == "" {
		ref = "$(" + fallback + ")"
	}
	name, ok := wholeName(ref)
	if !ok {
		a.refuse(path, "output binding must be null or a whole name")
		return
	}
	if strings.HasPrefix(name, "opctl_") {
		a.refuse(path, "output binding starts with the reserved opctl_ prefix")
	}
	b := a.newBinding(name, kind, "produced", path, state)
	scope[name] = b
	a.produced[path] = append(a.produced[path], b)
}
