package goconvert

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (r *renderer) conditional(g *graph, out *strings.Builder) (string, error) {
	condition, err := r.predicate(g)
	if err != nil {
		return "", err
	}
	calls := r.variable("calls")
	fmt.Fprintf(out, "%s := []opctl.Call{}\n", calls)
	for _, b := range r.a.bindings {
		if b.Branch == g.Path && r.predecl[b] && b.Fallback != nil {
			r.assign(out, b, r.expr[b.Fallback])
		}
	}
	fmt.Fprintf(out, "if %s {\n", condition)
	body := *g
	body.Condition = nil
	expr, err := r.call(&body, out)
	if err != nil {
		return "", err
	}
	if expr != "" {
		fmt.Fprintf(out, "%s = append(%s, %s)\n", calls, calls, expr)
	}
	out.WriteString("}\n")
	return "opctl.Serial(" + calls + "...)", nil
}
func (r *renderer) predicate(g *graph) (string, error) {
	data, _ := json.Marshal(g.Condition)
	var predicates []map[string]any
	_ = json.Unmarshal(data, &predicates)
	clauses := []string{}
	for i, predicate := range predicates {
		for _, kind := range keys(predicate) {
			operands, _ := predicate[kind].([]any)
			values := []string{}
			for j := range operands {
				v := r.a.values[fmt.Sprintf("%s.if[%d].%s[%d]", g.Path, i, kind, j)]
				var text string
				var err error
				if kind == "eq" || kind == "ne" {
					text, err = r.text(v)
				} else {
					text, err = r.number(v)
				}
				if err != nil {
					return "", err
				}
				values = append(values, text)
			}
			operator := map[string]string{"eq": "==", "ne": "!=", "lt": "<", "lte": "<=", "gt": ">", "gte": ">="}[kind]
			if operator == "" {
				return "", fmt.Errorf("unsupported predicate %s", kind)
			}
			if kind == "ne" {
				for j := range values {
					for k := 0; k < j; k++ {
						clauses = append(clauses, "("+values[k]+" != "+values[j]+")")
					}
				}
			} else {
				for j := 1; j < len(values); j++ {
					clauses = append(clauses, "("+values[j-1]+" "+operator+" "+values[j]+")")
				}
			}
		}
	}
	if len(clauses) == 0 {
		return "true", nil
	}
	return strings.Join(clauses, " && "), nil
}
func (r *renderer) number(v *boundValue) (string, error) {
	if len(v.Refs) > 0 {
		b := v.Refs[0].Binding
		if b == nil || b.Kind != "number" {
			return "", fmt.Errorf("ordered predicate needs a number")
		}
		expr := r.scalarValue(b)
		if b.Source == "loop" && b.Scalar == "int" {
			expr = "float64(" + expr + ")"
		}
		return expr, nil
	}
	n, err := strconv.ParseFloat(rawText(v.Raw), 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(n, 'g', -1, 64), nil
}
func (r *renderer) loop(g *graph, out *strings.Builder) (string, error) {
	for _, b := range g.Hoists {
		if !r.declared[b] {
			r.assign(out, b, "opctl.New"+strings.ToUpper(b.Kind[:1])+b.Kind[1:]+"()")
		}
	}
	calls := r.variable("calls")
	fmt.Fprintf(out, "%s := []opctl.Call{}\n", calls)
	kind, scalar := literalRange(g.Range.Raw)
	rangeExpr := ""
	input := false
	if g.Range.Whole && len(g.Range.Refs) > 0 && g.Range.Refs[0].Binding != nil {
		b := g.Range.Refs[0].Binding
		rangeExpr = r.expr[b] + ".Value()"
		kind = b.Kind
		input = true
	} else {
		switch v := g.Range.Raw.(type) {
		case []any:
			parts := []string{}
			for _, x := range v {
				parts = append(parts, goValue(x))
			}
			rangeExpr = "[]" + scalar + "{" + strings.Join(parts, ", ") + "}"
		case map[string]any:
			parts := []string{}
			for _, key := range keys(v) {
				parts = append(parts, literal(key)+": "+goValue(v[key]))
			}
			rangeExpr = "map[string]" + scalar + "{" + strings.Join(parts, ", ") + "}"
		}
	}
	if rangeExpr == "" {
		return "", fmt.Errorf("unsupported loop range")
	}
	index, value, key := g.Vars["index"], g.Vars["value"], g.Vars["key"]
	indexName := "_"
	if index != nil && r.used[index] {
		indexName = r.expr[index]
	}
	valueName := "_"
	if value != nil && r.used[value] {
		valueName = r.expr[value]
	}
	if kind == "object" {
		r.imports["maps"] = ""
		r.imports["slices"] = ""
		objects := r.variable("items")
		fmt.Fprintf(out, "%s := %s\n", objects, rangeExpr)
		keyName := r.variable("key")
		if key != nil && r.used[key] {
			keyName = r.expr[key]
		}
		if valueName == "_" && (key == nil || !r.used[key]) {
			keyName = "_"
		}
		if keyName == "_" && indexName == "_" {
			fmt.Fprintf(out, "for range slices.Sorted(maps.Keys(%s)) {\n", objects)
		} else {
			fmt.Fprintf(out, "for %s, %s := range slices.Sorted(maps.Keys(%s)) {\n", indexName, keyName, objects)
		}
		if valueName != "_" {
			expr := objects + "[" + keyName + "]"
			if input && value.Scalar != "any" {
				if value.Scalar == "int" {
					expr = "int(" + expr + ".(float64))"
				} else {
					expr += ".(" + value.Scalar + ")"
				}
			}
			fmt.Fprintf(out, "%s := %s\n", valueName, expr)
		}
	} else {
		switch {
		case indexName == "_" && valueName == "_":
			fmt.Fprintf(out, "for range %s {\n", rangeExpr)
		case valueName == "_":
			fmt.Fprintf(out, "for %s := range %s {\n", indexName, rangeExpr)
		default:
			fmt.Fprintf(out, "for %s, %s := range %s {\n", indexName, valueName, rangeExpr)
		}
	}
	expr, err := r.call(g.Body, out)
	if err != nil {
		return "", err
	}
	if expr != "" {
		fmt.Fprintf(out, "%s = append(%s, %s)\n", calls, calls, expr)
	}
	out.WriteString("}\n")
	group := "Serial"
	if g.Kind == "parallelLoop" {
		group = "Parallel"
	}
	return "opctl." + group + "(" + calls + "...)", nil
}
