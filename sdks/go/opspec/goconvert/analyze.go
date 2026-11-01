// Package goconvert analyzes YAML ops before rendering and verifying Go ops.
package goconvert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"golang.org/x/mod/module"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Describe obtains a callee's current interface. Local refs are absolute paths.
type Describe func(context.Context, string) (*model.OpSpec, error)

// Refusal identifies an unsupported construct in the original op file.
type Refusal struct {
	Path, Reason string
	Line, Column int
}

func (r Refusal) String() string {
	return fmt.Sprintf("op.yml:%d:%d  %s  %s", r.Line, r.Column, r.Path, r.Reason)
}

// Note describes an intentional metadata change or a verbatim reference.
type Note struct{ Path, Reason string }

// Analysis records a single symbolic scope walk, shared by all rendering and verification.
type Analysis struct {
	Dir        string
	Spec       *model.OpSpec
	Refusals   []Refusal
	Notes      []Note
	nodes      map[string]*yaml.Node
	callees    map[string]*model.OpSpec
	refs       map[string]string
	values     map[string]*boundValue
	inputs     map[string]*binding
	produced   map[string][]*binding
	final      map[string]*binding
	graph      *graph
	bindings   []*binding
	loopInputs map[string]string
}
type binding struct {
	Name, Kind, Type, Source, Path, Loop, Branch, Expr, Scalar string
	Param                                                      *model.ParamSpec
	ID                                                         int
	Implicit                                                   bool
}
type reference struct {
	Raw, Body, Suffix string
	Binding           *binding
	Nested            bool
}
type boundValue struct {
	Raw                any
	Refs               []reference
	Whole              bool
	Mode, Kind, Reason string
	Items              []*boundValue
	Members            map[string]*boundValue
}
type graph struct {
	Spec       *model.CallSpec
	Path, Kind string
	Children   []*graph
	Body       *graph
	Condition  []*model.PredicateSpec
	Fold       *bool
	Range      *boundValue
	Vars       map[string]*binding
	Hoists     []*binding
}
type walkState struct {
	loop, branch string
	parallelLoop bool
	loopGraph    *graph
}

var paramName = regexp.MustCompile(`^[-_a-zA-Z0-9]+$`)

// Analyze validates source once, describes its callees, and records every refusal
// before rendering or writing anything. It never changes the op directory.
func Analyze(ctx context.Context, dir string, source []byte, describe Describe) (*Analysis, error) {
	spec, err := opfile.Unmarshal(source)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err = yaml.Unmarshal(source, &doc); err != nil {
		return nil, err
	}
	a := &Analysis{Dir: dir, Spec: spec, nodes: map[string]*yaml.Node{}, callees: map[string]*model.OpSpec{}, refs: map[string]string{}, values: map[string]*boundValue{}, inputs: map[string]*binding{}, produced: map[string][]*binding{}, loopInputs: map[string]string{}}
	if len(doc.Content) > 0 {
		a.index(doc.Content[0], "")
	}
	if err = module.CheckImportPath(filepath.Base(dir)); err != nil {
		a.refuse("name", "op directory name isn't a valid Go module path element")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			a.refuse("name", "op directory already contains "+entry.Name()+"; move Go sources into a subdirectory")
		}
	}
	if spec.Version != "" {
		a.Notes = append(a.Notes, Note{"version", "op version has no Go form; dropped"})
	}
	for _, group := range []string{"inputs", "outputs"} {
		for name, param := range paramsFor(spec, group) {
			a.checkParam(group+"."+name, name, param)
		}
	}
	if err = a.describeCalls(ctx, spec.Run, "run", describe); err != nil {
		return nil, err
	}
	a.findLoopInputs(spec.Run)
	for name, param := range spec.Inputs {
		kind, _ := paramParts(param)
		b := a.newBinding(name, kind, "input", "inputs."+name, walkState{})
		b.Param = param
		if scalar := a.loopInputs[name]; scalar != "" {
			b.Type = "opctl.Array[" + scalar + "]"
		}
		a.inputs[name] = b
	}
	scope := cloneScope(a.inputs)
	a.graph = a.walk(spec.Run, "run", scope, walkState{}, true)
	if !nonempty(a.graph) {
		a.refuse("run", "Go ops can't plan an empty op graph")
	}
	a.final = scope
	for name, param := range spec.Outputs {
		b := scope[name]
		if b == nil {
			if !hasDefault(param) {
				a.refuse("outputs."+name, "declared output is unbound and has no default")
			}
			continue
		}
		if b.Loop != "" {
			a.refuse("outputs."+name, "declared output is produced inside a loop")
		}
		if b.Branch != "" && !hasDefault(param) {
			a.refuse("outputs."+name, "declared output is produced inside an if branch and has no default")
		}
		kind, _ := paramParts(param)
		if b.Kind != kind {
			a.refuse("outputs."+name, "output handle kind differs from its declaration")
		}
	}
	sort.SliceStable(a.Refusals, func(i, j int) bool {
		x, y := a.Refusals[i], a.Refusals[j]
		if x.Line != y.Line {
			return x.Line < y.Line
		}
		if x.Column != y.Column {
			return x.Column < y.Column
		}
		return x.Path < y.Path
	})
	return a, nil
}

func nonempty(g *graph) bool {
	if g == nil || g.Fold != nil && !*g.Fold || g != nil && len(g.Condition) > 0 && g.Fold == nil {
		return false
	}
	switch g.Kind {
	case "container", "op":
		return true
	case "serialLoop", "parallelLoop":
		return nonempty(g.Body)
	default:
		for _, child := range g.Children {
			if nonempty(child) {
				return true
			}
		}
	}
	return false
}
func paramsFor(s *model.OpSpec, key string) map[string]*model.ParamSpec {
	if key == "inputs" {
		return s.Inputs
	}
	return s.Outputs
}
func (a *Analysis) index(n *yaml.Node, path string) {
	a.nodes[path] = n
	if n.Kind == yaml.AliasNode {
		if n.Alias != nil {
			a.index(n.Alias, path)
		}
		return
	}
	switch n.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(n.Content); i += 2 {
			key := n.Content[i].Value
			p := key
			if path != "" {
				p = path + "." + key
			}
			a.index(n.Content[i+1], p)
			if n.Content[i].HeadComment != "" && n.Content[i+1].HeadComment == "" {
				n.Content[i+1].HeadComment = n.Content[i].HeadComment
			}
			if n.Content[i].LineComment != "" && n.Content[i+1].LineComment == "" {
				n.Content[i+1].LineComment = n.Content[i].LineComment
			}
		}
	case yaml.SequenceNode:
		for i, c := range n.Content {
			a.index(c, fmt.Sprintf("%s[%d]", path, i))
		}
	}
}
func (a *Analysis) refuse(path, reason string) {
	for _, r := range a.Refusals {
		if r.Path == path && r.Reason == reason {
			return
		}
	}
	n := a.nodes[path]
	p := path
	for n == nil && p != "" {
		if i := strings.LastIndexAny(p, ".["); i >= 0 {
			p = p[:i]
		} else {
			p = ""
		}
		n = a.nodes[p]
	}
	line, col := 1, 1
	if n != nil {
		line, col = n.Line, n.Column
	}
	a.Refusals = append(a.Refusals, Refusal{path, reason, line, col})
}
func (a *Analysis) newBinding(name, kind, source, path string, state walkState) *binding {
	b := &binding{Name: name, Kind: kind, Type: handleType(kind), Source: source, Path: path, Loop: state.loop, Branch: state.branch, ID: len(a.bindings) + 1}
	a.bindings = append(a.bindings, b)
	return b
}
func cloneScope(in map[string]*binding) map[string]*binding {
	out := map[string]*binding{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func paramParts(param *model.ParamSpec) (string, map[string]any) {
	data, _ := json.Marshal(param)
	var raw map[string]any
	_ = json.Unmarshal(data, &raw)
	for _, kind := range []string{"string", "number", "boolean", "dir", "file", "socket", "array", "object"} {
		if body, ok := raw[kind].(map[string]any); ok {
			return kind, body
		}
	}
	return "", nil
}
func handleType(kind string) string {
	switch kind {
	case "array":
		return "opctl.Array[any]"
	case "object":
		return "opctl.Object[map[string]any]"
	case "":
		return "opctl.String"
	default:
		return "opctl." + strings.ToUpper(kind[:1]) + kind[1:]
	}
}
func hasDefault(param *model.ParamSpec) bool {
	_, body := paramParts(param)
	_, ok := body["default"]
	return ok || param != nil && param.DefaultPresent
}
func secret(param *model.ParamSpec) bool {
	_, body := paramParts(param)
	return body["isSecret"] == true
}
func (a *Analysis) checkParam(path, name string, param *model.ParamSpec) {
	if !paramName.MatchString(name) || strings.HasPrefix(name, "opctl_") {
		a.refuse(path, "parameter names must contain only letters, numbers, - or _, and can't start with opctl_")
	}
	kind, body := paramParts(param)
	if body == nil {
		return
	}
	if d, present := body["default"]; present || param.DefaultPresent {
		valid := false
		switch kind {
		case "string", "dir", "file":
			_, valid = d.(string)
		case "number":
			_, valid = d.(float64)
		case "boolean":
			_, valid = d.(bool)
		case "array":
			_, valid = d.([]any)
		case "object":
			_, valid = d.(map[string]any)
		}
		if !valid {
			a.refuse(path+"."+kind+".default", "null defaults and defaults of the wrong type can't be converted")
		}
	}
	if nested, _ := body["description"].(string); nested != "" {
		a.Notes = append(a.Notes, Note{path, "deprecated nested description moved to the parameter"})
	}
}
func (a *Analysis) describeCalls(ctx context.Context, c *model.CallSpec, path string, describe Describe) error {
	if c == nil {
		return nil
	}
	if c.Op != nil {
		ref, err := NormalizeRef(a.Dir, c.Op.Ref)
		if err != nil {
			a.refuse(path+".op.ref", err.Error())
		} else {
			a.refs[path] = ref
			if describe == nil {
				return fmt.Errorf("describing callee %s needs a node", ref)
			}
			target := ref
			if !strings.Contains(ref, "#") {
				target = filepath.Join(a.Dir, ref)
			}
			spec, err := describe(ctx, target)
			if err != nil {
				return fmt.Errorf("describe %s: %w", ref, err)
			}
			a.callees[path] = spec
		}
	}
	for _, group := range []struct {
		name  string
		calls *[]*model.CallSpec
	}{{"serial", c.Serial}, {"parallel", c.Parallel}} {
		if group.calls != nil {
			for i, child := range *group.calls {
				if err := a.describeCalls(ctx, child, fmt.Sprintf("%s.%s[%d]", path, group.name, i), describe); err != nil {
					return err
				}
			}
		}
	}
	if c.SerialLoop != nil {
		if err := a.describeCalls(ctx, &c.SerialLoop.Run, path+".serialLoop.run", describe); err != nil {
			return err
		}
	}
	if c.ParallelLoop != nil {
		if err := a.describeCalls(ctx, &c.ParallelLoop.Run, path+".parallelLoop.run", describe); err != nil {
			return err
		}
	}
	return nil
}

// NormalizeRef makes local op refs explicit while preserving git refs.
func NormalizeRef(dir, ref string) (string, error) {
	if strings.HasPrefix(ref, "$(") && strings.HasSuffix(ref, ")") {
		body := ref[2 : len(ref)-1]
		if strings.HasPrefix(body, "./") || strings.HasPrefix(body, "../") {
			ref = body
		} else if strings.HasPrefix(body, "/") {
			ref = "." + body
		} else {
			return "", fmt.Errorf("op ref must be a literal local or git ref")
		}
	}
	if strings.Contains(ref, "$(") {
		return "", fmt.Errorf("op ref must be a literal local or git ref")
	}
	if strings.Contains(ref, "#") || filepath.IsAbs(ref) || strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../") {
		return ref, nil
	}
	if ref == "" {
		return "", fmt.Errorf("op ref is empty")
	}
	if info, err := os.Stat(filepath.Join(dir, ref)); err == nil && info.IsDir() {
		return "./" + ref, nil
	}
	return "../" + ref, nil
}
func (a *Analysis) findLoopInputs(c *model.CallSpec) {
	if c == nil {
		return
	}
	var r any
	if c.SerialLoop != nil {
		r = c.SerialLoop.Range
		a.findLoopInputs(&c.SerialLoop.Run)
	}
	if c.ParallelLoop != nil {
		r = c.ParallelLoop.Range
		a.findLoopInputs(&c.ParallelLoop.Run)
	}
	if ref, ok := r.(string); ok {
		if name, ok := wholeName(ref); ok {
			if p := a.Spec.Inputs[name]; p != nil && p.Array != nil {
				schema := p.Array.Constraints
				items, ok := schema["items"].(map[string]any)
				if ok && len(items) == 1 {
					if scalar := scalarType(items["type"]); scalar != "" {
						a.loopInputs[name] = scalar
					}
				}
			}
		}
	}
	for _, children := range []*[]*model.CallSpec{c.Serial, c.Parallel} {
		if children != nil {
			for _, child := range *children {
				a.findLoopInputs(child)
			}
		}
	}
}
func scalarType(kind any) string {
	switch kind {
	case "string":
		return "string"
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	}
	return ""
}
func wholeName(ref string) (string, bool) {
	if len(ref) > 3 && strings.HasPrefix(ref, "$(") && strings.HasSuffix(ref, ")") {
		name := ref[2 : len(ref)-1]
		if paramName.MatchString(name) {
			return name, true
		}
	}
	return "", false
}
func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
