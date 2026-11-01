package goconvert

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// RenderInterface renders a compilable declaration for inspecting the converted
// interface. Render supplies the original op graph as the plan body.
func RenderInterface(a *Analysis) ([]byte, error) {
	body := "return opctl.Container(\"alpine:3\")"
	if len(a.Spec.Outputs) > 0 {
		body += ", Outputs{}"
	}
	return a.renderSource(nil, "", body)
}
func (a *Analysis) renderSource(imports map[string]string, declarations, body string) ([]byte, error) {
	if len(a.Refusals) > 0 {
		return nil, fmt.Errorf("op can't be rendered: %s", a.Refusals[0])
	}
	var out bytes.Buffer
	docComment, override := description(a.Spec.Description)
	out.WriteString(docComment)
	out.WriteString("package main\n\nimport (\nopctl \"github.com/opctl/opctl/sdks/go/opspec/sdk\"\n")
	for _, path := range keys(imports) {
		fmt.Fprintf(&out, "%s %s\n", literal(path), imports[path])
	}
	out.WriteString(")\n\n")
	if len(a.Spec.Inputs) > 0 {
		out.WriteString(a.structSource("Inputs", "inputs", a.Spec.Inputs, true))
	}
	if len(a.Spec.Outputs) > 0 {
		out.WriteString(a.structSource("Outputs", "outputs", a.Spec.Outputs, false))
	}
	out.WriteString(declarations)
	input := "struct{}"
	if len(a.Spec.Inputs) > 0 {
		input = "Inputs"
	}
	result := "opctl.Call"
	entry := "Main"
	if len(a.Spec.Outputs) > 0 {
		result = "(opctl.Call, Outputs)"
		entry = "MainWithOutputs"
	}
	fmt.Fprintf(&out, "func plan(in %s) %s {\n%s\n}\n\n", input, result, body)
	metadata := []string{}
	if a.Spec.Name != filepath.Base(a.Dir) {
		metadata = append(metadata, "Name: "+literal(a.Spec.Name))
	}
	if override {
		metadata = append(metadata, "Description: "+literal(a.Spec.Description))
	}
	option := ""
	if len(metadata) > 0 {
		option = ", opctl.Op{" + strings.Join(metadata, ", ") + "}"
	}
	fmt.Fprintf(&out, "func main() { opctl.%s(plan%s) }\n", entry, option)
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format converted Go op: %w\n%s", err, &out)
	}
	for _, line := range strings.Split(string(formatted), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build") {
			return nil, fmt.Errorf("converted source contains a build constraint")
		}
	}
	return formatted, nil
}
func (a *Analysis) structSource(name, path string, params map[string]*model.ParamSpec, input bool) string {
	var out strings.Builder
	fmt.Fprintf(&out, "type %s struct {\n", name)
	names := fieldNames(keys(params), true)
	for _, key := range a.orderedKeys(path, params) {
		param := params[key]
		kind, _ := paramParts(param)
		typ := handleType(kind)
		if input && a.inputs[key] != nil {
			typ = a.inputs[key].Type
		}
		out.WriteString(a.comments(path + "." + key))
		fmt.Fprintf(&out, "%s %s %s\n", names[key], typ, paramTag(key, param, typ))
	}
	out.WriteString("}\n\n")
	return out.String()
}
func (a *Analysis) orderedKeys(path string, params map[string]*model.ParamSpec) []string {
	out := []string{}
	seen := map[string]bool{}
	if n := a.nodes[path]; n != nil {
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := sourceKey(n.Content[i])
			if _, ok := params[k]; ok {
				out = append(out, k)
				seen[k] = true
			}
		}
	}
	for _, k := range keys(params) {
		if !seen[k] {
			out = append(out, k)
		}
	}
	return out
}
func fieldNames(names []string, human bool) map[string]string {
	base := func(raw string) string {
		if human && strings.ToUpper(raw) == raw {
			raw = strings.ToLower(raw)
		}
		parts := strings.FieldsFunc(raw, func(r rune) bool { return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') })
		for i, p := range parts {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
		s := strings.Join(parts, "")
		if s == "" {
			s = "X" + shortHash(raw)
		}
		if s[0] >= '0' && s[0] <= '9' {
			s = "X" + s
		}
		return s
	}
	count := map[string]int{}
	for _, n := range names {
		count[base(n)]++
	}
	out := map[string]string{}
	for _, n := range names {
		b := base(n)
		if count[b] > 1 {
			b += "_" + shortHash(n)
		}
		out[n] = b
	}
	return out
}
func shortHash(s string) string { sum := sha256.Sum256([]byte(s)); return hex.EncodeToString(sum[:4]) }
func literal(s string) string {
	if strings.Contains(s, "\n") && !strings.ContainsAny(s, "`\r") {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}
func tagValue(s string) string {
	if s != "" && !strings.ContainsAny(s, ",'\\") {
		return s
	}
	return "'" + strings.NewReplacer("\\", "\\\\", "'", "\\'").Replace(s) + "'"
}
func paramTag(name string, param *model.ParamSpec, typ string) string {
	kind, body := paramParts(param)
	items := []string{name}
	description := param.Description
	if description == "" {
		description, _ = body["description"].(string)
	}
	if description != "" {
		items = append(items, "description="+tagValue(description))
	}
	if d, ok := body["default"]; ok {
		var text string
		if s, ok := d.(string); ok {
			text = s
		} else {
			b, _ := json.Marshal(d)
			text = string(b)
		}
		items = append(items, "default="+tagValue(text))
	}
	if body["isSecret"] == true {
		items = append(items, "secret")
	}
	constraints, _ := body["constraints"].(map[string]any)
	remaining := map[string]any{}
	for _, key := range keys(constraints) {
		value := constraints[key]
		if key == "items" && kind == "array" && typ != "opctl.Array[any]" {
			continue
		}
		if text, ok := flatConstraint(kind, key, value); ok {
			items = append(items, key+"="+tagValue(text))
		} else {
			remaining[key] = value
		}
	}
	if len(remaining) > 0 {
		b, _ := json.Marshal(remaining)
		items = append(items, "constraints="+tagValue(string(b)))
	}
	tag := "opctl:" + strconv.Quote(strings.Join(items, ","))
	if strings.Contains(tag, "`") {
		return strconv.Quote(tag)
	}
	return "`" + tag + "`"
}
func flatConstraint(kind, key string, value any) (string, bool) {
	switch key {
	case "format":
		s, ok := value.(string)
		if !ok {
			return "", false
		}
		if kind == "number" && s == "integer" || kind == "string" && strings.Contains("|date-time|docker-image-ref|email|hostname|ipv4|ipv6|semver|uri|", "|"+s+"|") {
			return s, true
		}
		return "", false
	case "enum", "required":
		list, ok := value.([]any)
		if !ok || len(list) == 0 {
			return "", false
		}
		parts := []string{}
		for _, v := range list {
			switch v := v.(type) {
			case string:
				if key == "enum" && kind != "string" || strings.Contains(v, "|") {
					return "", false
				}
				parts = append(parts, v)
			case float64:
				if key != "enum" || kind != "number" {
					return "", false
				}
				parts = append(parts, strconv.FormatFloat(v, 'g', -1, 64))
			default:
				return "", false
			}
		}
		return strings.Join(parts, "|"), true
	case "pattern":
		s, ok := value.(string)
		return s, ok && kind == "string"
	case "minLength", "maxLength", "minItems", "maxItems", "minProperties", "maxProperties", "minimum", "maximum", "multipleOf", "uniqueItems":
		b, err := json.Marshal(value)
		return string(b), err == nil
	}
	return "", false
}
func comment(text string) string {
	var out strings.Builder
	text = strings.Map(func(r rune) rune {
		if r != '\n' && (unicode.IsControl(r) || unicode.Is(unicode.Cf, r)) {
			return -1
		}
		return r
	}, text)
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "+build") || strings.HasPrefix(trimmed, "go:build") {
			line = "Comment: " + line
		}
		out.WriteString("// " + line + "\n")
	}
	return out.String()
}
func description(text string) (string, bool) {
	if text == "" {
		return "", false
	}
	candidate := comment(text)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "op.go", candidate+"package main\n", parser.PackageClauseOnly|parser.ParseComments)
	if err != nil {
		return "", true
	}
	docs, err := doc.NewFromFiles(fset, []*ast.File{f}, "main")
	if err != nil {
		return "", true
	}
	if strings.TrimRight(string(docs.Markdown(docs.Doc)), "\n") == text {
		return candidate, false
	}
	return "", true
}
func (a *Analysis) comments(path string) string {
	n := a.nodes[path]
	if n == nil {
		return ""
	}
	parts := []string{}
	for _, text := range []string{n.HeadComment, n.LineComment} {
		if text != "" {
			for _, line := range strings.Split(text, "\n") {
				parts = append(parts, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "#")))
			}
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return comment(strings.Join(parts, "\n"))
}
