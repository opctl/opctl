package goconvert

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"path/filepath"
	"strconv"
	"strings"
)

type convertedStub struct {
	use             goop.StubUse
	spec            *model.OpSpec
	inputs, outputs map[string]string
}

func (r *renderer) prepareStubs() error {
	r.stubCalls = map[string]*convertedStub{}
	r.fallbacks = map[string]string{}
	r.handleTypes = map[*binding]string{}
	for _, b := range r.a.bindings {
		r.handleTypes[b] = b.Type
	}
	if r.options.ModulePath == "" {
		return nil
	}
	existing := map[string]string{}
	reuse := map[string]string{}
	if r.options.ModuleRoot != "" {
		uses, _ := goop.ScanModuleStubs(r.options.ModuleRoot)
		for _, use := range uses {
			target := goop.StubModuleRef(r.options.ModuleRoot, use)
			existing[use.Name] = target
			reuse[target] = use.Name
		}
	}
	calls := []*graph{}
	var visit func(*graph)
	visit = func(g *graph) {
		if g == nil {
			return
		}
		if g.Kind == "op" {
			calls = append(calls, g)
		}
		for _, c := range g.Children {
			visit(c)
		}
		visit(g.Body)
	}
	visit(r.a.graph)
	eligible := map[string]*convertedStub{}
	for _, g := range calls {
		spec := r.a.callees[g.Path]
		inputs, outputs, err := goop.StubHandleTypes(spec)
		if err != nil {
			return err
		}
		reason := ""
		c := g.Spec.Op
		p := g.Path + ".op"
		for _, name := range keys(spec.Inputs) {
			if _, ok := c.Inputs[name]; !ok && !hasDefault(spec.Inputs[name]) {
				reason = "required input " + strconv.Quote(name) + " isn't bound"
				break
			}
		}
		if reason == "" {
			for _, name := range keys(c.Inputs) {
				v := r.a.values[p+".inputs."+name]
				if _, ok := r.stubValue(v, inputs[name]); !ok {
					actual := r.valueType(v)
					reason = "input " + strconv.Quote(name) + " is " + shortType(actual) + " here and " + shortType(inputs[name]) + " in the stub"
					if (v.Items != nil || v.Members != nil) && containsReferences(v) {
						reason = "input " + strconv.Quote(name) + " contains references in an array or object"
					}
					break
				}
			}
		}
		if reason == "" {
			for _, name := range keys(c.Outputs) {
				for _, b := range r.a.produced[p+".outputs."+name] {
					for output := range r.a.Spec.Outputs {
						for final := r.a.final[output]; final != nil; final = final.Fallback {
							if final == b {
								kind, _ := paramParts(r.a.Spec.Outputs[output])
								if outputs[name] != handleType(kind) {
									reason = "output " + strconv.Quote(name) + " is " + shortType(outputs[name]) + " in the stub and " + shortType(handleType(kind)) + " in the op's outputs"
								}
							}
						}
					}
				}
			}
		}
		if reason != "" {
			r.fallbacks[g.Path] = reason
			continue
		}
		stub := &convertedStub{spec: spec, inputs: inputs, outputs: outputs}
		eligible[g.Path] = stub
		for name := range c.Outputs {
			for _, b := range r.a.produced[p+".outputs."+name] {
				r.handleTypes[b] = outputs[name]
			}
		}
	}
	targets := map[string]string{}
	for _, g := range calls {
		if eligible[g.Path] != nil {
			ref := r.a.refs[g.Path]
			targets[ref] = r.target(ref)
		}
	}
	labels := map[string]string{}
	for _, ref := range keys(targets) {
		target := targets[ref]
		if name := reuse[target]; name != "" {
			labels[ref] = name
			continue
		}
		parts := stubParts(ref)
		name := parts[0]
		for depth := 1; depth < len(parts) && r.labelCollision(name, ref, targets, parts, depth); depth++ {
			name = parts[depth] + name
		}
		base := name
		for n := 2; !goop.ValidStubName(name) || existing[name] != "" && existing[name] != target || r.names[name] || labelUsed(labels, name, ref); n++ {
			name = base + strconv.Itoa(n)
		}
		labels[ref] = name
		existing[name] = target
	}
	for _, g := range calls {
		stub := eligible[g.Path]
		if stub == nil {
			continue
		}
		ref := r.a.refs[g.Path]
		name := labels[ref]
		if _, err := NormalizeRef(r.a.Dir, ref); err != nil {
			return err
		}
		use := goop.StubUse{Name: name, Ref: ref, Local: !strings.Contains(ref, "#")}
		if use.Local {
			relative, err := filepath.Rel(r.options.ModuleRoot, r.target(ref))
			if err != nil {
				return err
			}
			use.ModuleRef = filepath.ToSlash(relative)
		}
		stub.use = use
		r.stubCalls[g.Path] = stub
		r.names[name] = true
		r.imports[strings.TrimSuffix(r.options.ModulePath, "/")+"/opctlgen/"+name] = "//opctl:op " + ref
		data, err := goop.RenderStub(use, stub.spec)
		if err != nil {
			return err
		}
		r.stubs[name] = data
	}
	return nil
}
func (r *renderer) target(ref string) string {
	if strings.Contains(ref, "#") {
		return ref
	}
	if filepath.IsAbs(ref) {
		return filepath.Clean(ref)
	}
	return filepath.Clean(filepath.Join(r.a.Dir, ref))
}
func stubParts(ref string) []string {
	parts := strings.Split(strings.Trim(ref, "/"), "/")
	out := []string{}
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part == ".opspec" || part == "." || part == ".." || part == "" {
			continue
		}
		if index := strings.IndexByte(part, '#'); index >= 0 {
			part = part[:index]
		}
		part = strings.Map(func(r rune) rune {
			if r >= 'A' && r <= 'Z' {
				return r + 'a' - 'A'
			}
			if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, part)
		if part != "" {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		return []string{"op"}
	}
	if out[0][0] < 'a' || out[0][0] > 'z' {
		out[0] = "op" + out[0]
	}
	return out
}
func (r *renderer) labelCollision(name, ref string, targets map[string]string, parts []string, depth int) bool {
	for other := range targets {
		if other == ref || targets[other] == targets[ref] {
			continue
		}
		otherParts := stubParts(other)
		candidate := otherParts[0]
		for i := 1; i < depth && i < len(otherParts); i++ {
			candidate = otherParts[i] + candidate
		}
		if candidate == name {
			return true
		}
	}
	return false
}
func labelUsed(labels map[string]string, name, ref string) bool {
	for other, value := range labels {
		if other != ref && value == name {
			return true
		}
	}
	return false
}
func shortType(typ string) string { return strings.ReplaceAll(typ, "opctl.", "") }
func containsReferences(v *boundValue) bool {
	if len(v.Refs) > 0 {
		return true
	}
	for _, x := range v.Items {
		if containsReferences(x) {
			return true
		}
	}
	for _, x := range v.Members {
		if containsReferences(x) {
			return true
		}
	}
	return false
}
func (r *renderer) valueType(v *boundValue) string {
	if v == nil {
		return "nil"
	}
	if v.Mode == "handle" && v.Whole && len(v.Refs) > 0 && v.Refs[0].Binding != nil {
		b := v.Refs[0].Binding
		if b.Source == "loop" && b.Scalar == "any" {
			return "any"
		}
		return r.handleTypes[b]
	}
	if v.Mode == "path" {
		if v.Kind == "file" {
			return "opctl.File"
		}
		return "opctl.Dir"
	}
	if v.Mode == "raw" {
		return handleType(v.Kind)
	}
	if v.Mode == "text" {
		return "string"
	}
	switch v.Raw.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	}
	return "nil"
}
func (r *renderer) stubValue(v *boundValue, want string) (string, bool) {
	typ := r.valueType(v)
	if typ == want {
		expr, err := r.value(v)
		return expr, err == nil
	}
	if v == nil {
		return "", false
	}
	if typ == "string" && want == "opctl.String" {
		expr, err := r.text(v)
		return "opctl.StringOf(" + expr + ")", err == nil
	}
	if typ == "number" && want == "opctl.Number" {
		return "opctl.NumberOf(" + goValue(v.Raw) + ")", true
	}
	if typ == "boolean" && want == "opctl.Boolean" {
		return "opctl.BooleanOf(" + goValue(v.Raw) + ")", true
	}
	if containsReferences(v) {
		return "", false
	}
	if list, ok := v.Raw.([]any); ok && strings.HasPrefix(want, "opctl.Array[") {
		itemType := strings.TrimSuffix(strings.TrimPrefix(want, "opctl.Array["), "]")
		if itemType != "any" && itemType != "string" && itemType != "int" && itemType != "float64" && itemType != "bool" {
			return "", false
		}
		parts := []string{}
		for _, item := range list {
			actual := literalScalar(item)
			if itemType != "any" && actual != itemType && !(itemType == "int" && actual == "float64") {
				return "", false
			}
			if itemType == "int" {
				number := item.(float64)
				if number != float64(int(number)) {
					return "", false
				}
			}
			parts = append(parts, goValue(item))
		}
		return "opctl.ArrayOf([]" + itemType + "{" + strings.Join(parts, ", ") + "})", true
	}
	if _, ok := v.Raw.(map[string]any); ok && want == "opctl.Object[map[string]any]" {
		return "opctl.ObjectOf(" + goValue(v.Raw) + ")", true
	}
	return "", false
}
func (r *renderer) stubOp(g *graph, out *strings.Builder, stub *convertedStub) (string, error) {
	p := g.Path + ".op"
	c := g.Spec.Op
	fields := goop.StubFieldNames(stub.spec.Inputs)
	parts := []string{}
	for _, name := range r.a.mapKeys(p+".inputs", keys(c.Inputs)) {
		value, ok := r.stubValue(r.a.values[p+".inputs."+name], stub.inputs[name])
		if !ok {
			return "", fmt.Errorf("stub type changed while rendering input %s", name)
		}
		parts = append(parts, r.a.comments(p+".inputs."+name)+rawTodos(r.a.values[p+".inputs."+name])+fields[name]+": "+value+",")
	}
	call := r.variable(stub.use.Name + "Call")
	result := "_"
	for name := range c.Outputs {
		for _, b := range r.a.produced[p+".outputs."+name] {
			if r.used[b] && result == "_" {
				result = r.variable(stub.use.Name + "Outputs")
			}
		}
	}
	fmt.Fprintf(out, "%s, %s := %s.Call(%s.Inputs{\n%s\n})\n", call, result, stub.use.Name, stub.use.Name, strings.Join(parts, "\n"))
	outputFields := goop.StubFieldNames(stub.spec.Outputs)
	for _, name := range keys(c.Outputs) {
		for _, b := range r.a.produced[p+".outputs."+name] {
			if r.used[b] {
				expression := result + "." + outputFields[name]
				if r.predecl[b] {
					r.assign(out, b, expression)
				} else {
					r.expr[b] = expression
				}
			}
		}
	}
	if c.PullCreds != nil {
		username, err := r.value(r.a.values[p+".pullCreds.username"])
		if err != nil {
			return "", err
		}
		password, err := r.value(r.a.values[p+".pullCreds.password"])
		if err != nil {
			return "", err
		}
		fmt.Fprintf(out, "%s = %s.PullCreds(%s, %s)\n", call, call, username, password)
	}
	return call, nil
}
