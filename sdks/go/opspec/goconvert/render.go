package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"
)

// RenderOptions identifies the module that will own converted op stubs.
type RenderOptions struct {
	ModuleRoot, ModulePath string
	// StagingDir anchors local callee refs at their original targets during verification.
	StagingDir string
}

// Rendered is formatted Go source and the editor stubs needed to compile it locally.
type Rendered struct {
	Source []byte
	Stubs  map[string][]byte
}
type renderer struct {
	stubCalls               map[string]*convertedStub
	fallbacks               map[string]string
	handleTypes             map[*binding]string
	a                       *Analysis
	options                 RenderOptions
	names                   map[string]bool
	expr                    map[*binding]string
	used, declared, predecl map[*binding]bool
	imports                 map[string]string
	stubs                   map[string][]byte
	declarations            strings.Builder
	outputTypes             int
}

// Render builds Go source from the already analyzed scope walk.
func Render(a *Analysis, options RenderOptions) (*Rendered, error) {
	if len(a.Refusals) > 0 {
		return nil, fmt.Errorf("op can't be converted: %s", a.Refusals[0])
	}
	if options.StagingDir != "" {
		staged := *a
		staged.refs = map[string]string{}
		for site, ref := range a.refs {
			if !strings.Contains(ref, "#") {
				target := ref
				if !filepath.IsAbs(target) {
					target = filepath.Join(a.Dir, target)
				}
				relative, err := filepath.Rel(options.StagingDir, target)
				if err != nil {
					return nil, err
				}
				ref = filepath.ToSlash(relative)
				if !strings.HasPrefix(ref, "../") {
					ref = "./" + ref
				}
			}
			staged.refs[site] = ref
		}
		a = &staged
	}
	r := &renderer{a: a, options: options, names: map[string]bool{}, expr: map[*binding]string{}, used: map[*binding]bool{}, declared: map[*binding]bool{}, predecl: map[*binding]bool{}, imports: map[string]string{}, stubs: map[string][]byte{}}
	for _, name := range []string{"opctl", "plan", "in", "main", "Inputs", "Outputs", "strconv", "maps", "slices"} {
		r.names[name] = true
	}
	if err := r.prepareStubs(); err != nil {
		return nil, err
	}
	fields := fieldNames(keys(a.Spec.Inputs), true)
	for name, b := range a.inputs {
		r.expr[b] = "in." + fields[name]
	}
	for _, v := range a.values {
		for _, ref := range v.Refs {
			if ref.Binding != nil {
				r.used[ref.Binding] = true
			}
		}
	}
	for name := range a.Spec.Outputs {
		for b := a.final[name]; b != nil; b = b.Fallback {
			r.used[b] = true
			if b.Branch != "" {
				r.predecl[b] = true
			}
		}
	}
	for _, b := range a.bindings {
		if b.Source != "input" {
			r.expr[b] = r.variable(b.Name)
		}
	}
	var body strings.Builder
	for _, b := range a.bindings {
		if r.predecl[b] {
			fmt.Fprintf(&body, "var %s %s\n", r.expr[b], b.Type)
			r.declared[b] = true
		}
	}
	expr, err := r.call(a.graph, &body)
	if err != nil {
		return nil, err
	}
	if expr == "" {
		return nil, fmt.Errorf("Go ops can't plan an empty op graph")
	}
	fmt.Fprintf(&body, "return %s", expr)
	if len(a.Spec.Outputs) > 0 {
		fields := fieldNames(keys(a.Spec.Outputs), true)
		body.WriteString(", Outputs{\n")
		for _, name := range a.orderedKeys("outputs", a.Spec.Outputs) {
			if b := a.final[name]; b != nil {
				fmt.Fprintf(&body, "%s: %s,\n", fields[name], r.expr[b])
			}
		}
		body.WriteString("}")
	}
	body.WriteByte('\n')
	source, err := a.renderSource(r.imports, r.declarations.String(), body.String())
	if err != nil {
		return nil, err
	}
	return &Rendered{Source: source, Stubs: r.stubs}, nil
}
func (r *renderer) variable(raw string) string {
	base := fieldNames([]string{raw}, true)[raw]
	base = strings.ToLower(base[:1]) + base[1:]
	if token.Lookup(base).IsKeyword() || types.Universe.Lookup(base) != nil {
		base += "2"
	}
	name := base
	for n := 2; r.names[name]; n++ {
		name = base + strconv.Itoa(n)
	}
	r.names[name] = true
	return name
}
func (r *renderer) call(g *graph, out *strings.Builder) (string, error) {
	if g == nil {
		return "", nil
	}
	if g.Fold != nil && !*g.Fold {
		out.WriteString(comment("Constant predicate folded to false."))
		return "", nil
	}
	out.WriteString(r.a.comments(g.Path))
	if g.Spec.Description != "" {
		out.WriteString(comment(g.Spec.Description))
	}
	if g.Fold != nil {
		out.WriteString(comment("Constant predicate folded to true."))
	}
	if len(g.Condition) > 0 && g.Fold == nil {
		return r.conditional(g, out)
	}
	switch g.Kind {
	case "container":
		return r.container(g, out)
	case "op":
		return r.op(g, out)
	case "serial", "parallel":
		children := []string{}
		for _, child := range g.Children {
			expr, err := r.call(child, out)
			if err != nil {
				return "", err
			}
			if expr != "" {
				children = append(children, expr)
			}
		}
		if len(children) == 0 {
			return "", nil
		}
		if len(children) == 1 {
			return children[0], nil
		}
		name := "Serial"
		if g.Kind == "parallel" {
			name = "Parallel"
		}
		return "opctl." + name + "(" + strings.Join(children, ",\n") + ")", nil
	case "serialLoop", "parallelLoop":
		return r.loop(g, out)
	}
	return "", nil
}
func (r *renderer) container(g *graph, out *strings.Builder) (string, error) {
	c := g.Spec.Container
	p := g.Path + ".container"
	for _, group := range []struct {
		key    string
		values map[string]any
	}{{"dirs", c.Dirs}, {"files", c.Files}} {
		for _, name := range r.a.mapKeys(p+"."+group.key, keys(group.values)) {
			for _, b := range r.a.produced[p+"."+group.key+"."+name] {
				if b.Implicit && !r.declared[b] {
					r.assign(out, b, "opctl.New"+strings.ToUpper(b.Kind[:1])+b.Kind[1:]+"()")
				}
			}
		}
	}
	image, err := r.text(r.a.values[p+".image.ref"])
	if err != nil {
		return "", err
	}
	expression := "opctl.Container(" + image + ")"
	methods := 0
	add := func(path, method string) {
		comments := r.a.comments(path)
		if v := r.a.values[path]; v != nil && v.Mode == "raw" {
			comments += comment("TODO(opctl op convert): " + v.Reason + ".")
		}
		expression += ".\n" + comments + method
		methods++
	}
	for _, key := range r.a.mapKeys(p, []string{"dirs", "files", "envVars", "sockets", "name", "workDir", "cmd"}) {
		switch key {
		case "dirs", "files":
			values := c.Dirs
			method := "Dir"
			if key == "files" {
				values = c.Files
				method = "File"
			}
			for _, name := range r.a.mapKeys(p+"."+key, keys(values)) {
				path := p + "." + key + "." + name
				value, err := r.value(r.a.values[path])
				if err != nil {
					return "", err
				}
				add(path, method+"("+literal(name)+", "+value+")")
			}
		case "envVars":
			if values, ok := c.EnvVars.(map[string]any); ok {
				for _, name := range r.a.mapKeys(p+".envVars", keys(values)) {
					path := p + ".envVars." + name
					v := r.a.values[path]
					var value, method string
					if v.Mode == "literal" || v.Mode == "text" {
						value, err = r.text(v)
						method = "Env"
					} else {
						value, err = r.value(v)
						method = "EnvFrom"
					}
					if err != nil {
						return "", err
					}
					add(path, method+"("+literal(name)+", "+value+")")
				}
			}
		case "sockets":
			for _, name := range r.a.mapKeys(p+".sockets", keys(c.Sockets)) {
				path := p + ".sockets." + name
				value, err := r.value(r.a.values[path])
				if err != nil {
					return "", err
				}
				add(path, "Socket("+literal(name)+", "+value+")")
			}
		case "name":
			if c.Name != nil {
				value, err := r.text(r.a.values[p+".name"])
				if err != nil {
					return "", err
				}
				add(p+".name", "Name("+value+")")
			}
		case "workDir":
			if c.WorkDir != "" {
				value, err := r.text(r.a.values[p+".workDir"])
				if err != nil {
					return "", err
				}
				add(p+".workDir", "WorkDir("+value+")")
			}
		case "cmd":
			if items, ok := c.Cmd.([]any); ok {
				values := []string{}
				for i := range items {
					value, err := r.text(r.a.values[fmt.Sprintf("%s.cmd[%d]", p, i)])
					if err != nil {
						return "", err
					}
					values = append(values, value)
				}
				add(p+".cmd", "Cmd("+strings.Join(values, ", ")+")")
			}
		}
	}
	produces := len(c.Outputs) > 0 || c.Stdout != "" || c.Stderr != ""
	if !produces && methods <= 1 {
		return expression, nil
	}
	imageName := filepath.Base(c.Image.Ref)
	if index := strings.IndexAny(imageName, ":@"); index >= 0 {
		imageName = imageName[:index]
	}
	call := r.variable(imageName)
	if len(c.OutputParams) > 0 {
		r.outputTypes++
		typeName := "ContainerOutputs"
		if r.outputTypes > 1 {
			typeName += strconv.Itoa(r.outputTypes)
		}
		r.declarations.WriteString(comment(typeName + " declares the outputs file of " + g.Path + "."))
		r.declarations.WriteString(r.a.structSource(typeName, p+".outputParams", c.OutputParams, false))
		result := "_"
		for name := range c.OutputParams {
			for _, b := range r.a.produced[p+".outputs."+name] {
				if r.used[b] {
					if result == "_" {
						result = r.variable("results")
					}
				}
			}
		}
		fmt.Fprintf(out, "%s, %s := opctl.Outputs[%s](%s)\n", call, result, typeName, expression)
		fields := fieldNames(keys(c.OutputParams), true)
		for _, name := range keys(c.OutputParams) {
			for _, b := range r.a.produced[p+".outputs."+name] {
				if r.used[b] {
					if r.predecl[b] {
						r.assign(out, b, result+"."+fields[name])
					} else {
						r.expr[b] = result + "." + fields[name]
					}
				}
			}
		}
	} else {
		fmt.Fprintf(out, "%s := %s\n", call, expression)
	}
	for _, name := range keys(c.Outputs) {
		if c.OutputParams[name] != nil {
			continue
		}
		for _, b := range r.a.produced[p+".outputs."+name] {
			r.produce(out, call, b, "opctl.ContainerOutput["+b.Type+"]("+call+", "+literal(name)+")")
		}
	}
	for _, stream := range []struct{ key, method string }{{"stdout", "Stdout"}, {"stderr", "Stderr"}} {
		for _, b := range r.a.produced[p+"."+stream.key] {
			r.produce(out, call, b, call+"."+stream.method+"()")
		}
	}
	return call, nil
}
func (r *renderer) op(g *graph, out *strings.Builder) (string, error) {
	if stub := r.stubCalls[g.Path]; stub != nil {
		return r.stubOp(g, out, stub)
	}
	if reason := r.fallbacks[g.Path]; reason != "" {
		out.WriteString(comment("CallOp, not a stub: " + reason + "."))
	}
	c := g.Spec.Op
	p := g.Path + ".op"
	ref := r.a.refs[g.Path]
	expression := "opctl.CallOp(" + literal(ref) + ")"
	if len(c.Inputs) > 0 {
		parts := []string{}
		for _, name := range r.a.mapKeys(p+".inputs", keys(c.Inputs)) {
			value, err := r.value(r.a.values[p+".inputs."+name])
			if err != nil {
				return "", err
			}
			parts = append(parts, r.a.comments(p+".inputs."+name)+rawTodos(r.a.values[p+".inputs."+name])+literal(name)+": "+value+",")
		}
		expression += ".Inputs(map[string]any{\n" + strings.Join(parts, "\n") + "\n})"
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
		expression += ".PullCreds(" + username + ", " + password + ")"
	}
	if len(c.Outputs) == 0 {
		return expression, nil
	}
	base := filepath.Base(ref)
	if i := strings.IndexByte(base, '#'); i >= 0 {
		base = base[:i]
	}
	call := r.variable(base + "Call")
	fmt.Fprintf(out, "%s := %s\n", call, expression)
	for _, name := range keys(c.Outputs) {
		for _, b := range r.a.produced[p+".outputs."+name] {
			r.produce(out, call, b, "opctl.OpOutput["+b.Type+"]("+call+", "+literal(name)+")")
		}
	}
	return call, nil
}
func (r *renderer) assign(out *strings.Builder, b *binding, expression string) {
	op := ":="
	if r.declared[b] {
		op = "="
	}
	fmt.Fprintf(out, "%s %s %s\n", r.expr[b], op, expression)
	r.declared[b] = true
}
func (r *renderer) produce(out *strings.Builder, call string, b *binding, expression string) {
	name := "_"
	op := "="
	if r.used[b] {
		name = r.expr[b]
		if !r.declared[b] {
			op = ":="
		}
		r.declared[b] = true
	}
	fmt.Fprintf(out, "%s, %s %s %s\n", call, name, op, expression)
}
func (a *Analysis) mapKeys(path string, fallback []string) []string {
	wanted := map[string]bool{}
	for _, k := range fallback {
		wanted[k] = true
	}
	seen := map[string]bool{}
	out := []string{}
	if n := a.nodes[path]; n != nil {
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := sourceKey(n.Content[i])
			if wanted[k] {
				out = append(out, k)
				seen[k] = true
			}
		}
	}
	for _, k := range fallback {
		if !seen[k] {
			out = append(out, k)
		}
	}
	return out
}
func (r *renderer) value(v *boundValue) (string, error) {
	if v == nil {
		return "nil", nil
	}
	if v.Items != nil {
		parts := []string{}
		for _, item := range v.Items {
			x, err := r.value(item)
			if err != nil {
				return "", err
			}
			parts = append(parts, x)
		}
		return "[]any{" + strings.Join(parts, ", ") + "}", nil
	}
	if v.Members != nil {
		parts := []string{}
		for _, name := range keys(v.Members) {
			x, err := r.value(v.Members[name])
			if err != nil {
				return "", err
			}
			parts = append(parts, literal(name)+": "+x)
		}
		return "map[string]any{" + strings.Join(parts, ", ") + "}", nil
	}
	switch v.Mode {
	case "raw":
		return "opctl.RawRef[" + handleType(v.Kind) + "](" + literal(rawText(v.Raw)) + ")", nil
	case "path":
		ref := v.Refs[0]
		kind := v.Kind
		if kind != "file" {
			kind = "dir"
		}
		method := strings.ToUpper(kind[:1]) + kind[1:]
		if ref.Binding != nil {
			return r.expr[ref.Binding] + "." + method + "(" + literal(ref.Suffix[1:]) + ")", nil
		}
		path := ref.Body
		if strings.HasPrefix(path, "/") {
			path = "." + path
		}
		return "opctl.Op" + method + "(" + literal(path) + ")", nil
	case "handle":
		if !v.Whole || len(v.Refs) != 1 || v.Refs[0].Binding == nil {
			return "", fmt.Errorf("unresolved handle %v", v.Raw)
		}
		b := v.Refs[0].Binding
		expression := r.expr[b]
		if b.Source == "loop" {
			if b.Scalar == "any" {
				return expression, nil
			}
			switch b.Kind {
			case "number":
				expression = "opctl.NumberOf(float64(" + expression + "))"
			case "boolean":
				expression = "opctl.BooleanOf(" + expression + ")"
			default:
				expression = "opctl.StringOf(" + expression + ")"
			}
		}
		return expression, nil
	case "text":
		return r.text(v)
	default:
		return goValue(v.Raw), nil
	}
}
func goValue(v any) string {
	switch v := v.(type) {
	case nil:
		return "nil"
	case string:
		var out strings.Builder
		_ = interpolater.Scan(v, func(s string) error { out.WriteString(s); return nil }, func(ref string) error { out.WriteString(ref); return nil })
		return literal(out.String())
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case []any:
		parts := []string{}
		for _, x := range v {
			parts = append(parts, goValue(x))
		}
		return "[]any{" + strings.Join(parts, ", ") + "}"
	case map[string]any:
		parts := []string{}
		for _, key := range keys(v) {
			parts = append(parts, literal(key)+": "+goValue(v[key]))
		}
		return "map[string]any{" + strings.Join(parts, ", ") + "}"
	}
	return "nil"
}
func (r *renderer) text(v *boundValue) (string, error) {
	if v == nil {
		return literal(""), nil
	}
	raw := v.Raw
	text, ok := raw.(string)
	if !ok {
		data, _ := json.Marshal(raw)
		text = string(data)
	}
	parts := []string{}
	idx := 0
	err := interpolater.Scan(text, func(s string) error {
		if s != "" {
			parts = append(parts, literal(s))
		}
		return nil
	}, func(ref string) error {
		if idx >= len(v.Refs) {
			return fmt.Errorf("missing reference analysis for %s", ref)
		}
		b := v.Refs[idx].Binding
		idx++
		if b == nil || !planReadable(b) {
			return fmt.Errorf("reference %s isn't plan-readable", ref)
		}
		parts = append(parts, r.scalarText(b))
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(parts) == 0 {
		return literal(""), nil
	}
	return strings.Join(parts, " + "), nil
}
func (r *renderer) scalarValue(b *binding) string {
	if b.Source == "input" {
		return r.expr[b] + ".Value()"
	}
	return r.expr[b]
}
func (r *renderer) scalarText(b *binding) string {
	expr := r.scalarValue(b)
	switch {
	case b.Source == "loop" && b.Scalar == "int":
		r.imports["strconv"] = ""
		return "strconv.Itoa(" + expr + ")"
	case b.Kind == "number":
		r.imports["strconv"] = ""
		return "strconv.FormatFloat(" + expr + ", 'f', -1, 64)"
	case b.Kind == "boolean":
		r.imports["strconv"] = ""
		return "strconv.FormatBool(" + expr + ")"
	default:
		return expr
	}
}

func rawText(v any) string {
	if text, ok := v.(string); ok {
		return text
	}
	data, _ := json.Marshal(v)
	return string(data)
}

func rawTodos(v *boundValue) string {
	if v == nil {
		return ""
	}
	out := ""
	if v.Mode == "raw" {
		out += comment("TODO(opctl op convert): " + v.Reason + ".")
	}
	for _, item := range v.Items {
		out += rawTodos(item)
	}
	for _, key := range keys(v.Members) {
		out += rawTodos(v.Members[key])
	}
	return out
}
