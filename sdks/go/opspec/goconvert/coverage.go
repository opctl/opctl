package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"sort"
	"strings"
)

// CheckCoverage reports every branch and sampled loop length the capped sets miss.
func CheckCoverage(a *Analysis, sets []map[string]any) error {
	branches := map[string]map[bool]bool{}
	lengths := map[string]map[int]bool{}
	var visit func(*graph, map[string]*model.Value) error
	visit = func(g *graph, scope map[string]*model.Value) error {
		if g == nil {
			return nil
		}
		if len(g.Condition) > 0 {
			yes, err := predicates.Interpret(g.Condition, scope)
			if err != nil {
				return err
			}
			if branches[g.Path] == nil {
				branches[g.Path] = map[bool]bool{}
			}
			branches[g.Path][yes] = true
			if !yes {
				return nil
			}
		}
		if g.Body != nil {
			ranged, err := loopable.Interpret(g.Range.Raw, scope)
			if err != nil {
				return err
			}
			count := 0
			if ranged.Array != nil {
				count = len(*ranged.Array)
			} else if ranged.Object != nil {
				count = len(*ranged.Object)
			}
			if lengths[g.Path] == nil {
				lengths[g.Path] = map[int]bool{}
			}
			lengths[g.Path][count] = true
			if count > 1000 {
				return fmt.Errorf("verification loop at %s exceeds 1000 iterations", g.Path)
			}
			var vars *model.LoopVarsSpec
			if g.Spec.SerialLoop != nil {
				vars = g.Spec.SerialLoop.Vars
			} else {
				vars = g.Spec.ParallelLoop.Vars
			}
			for i := 0; i < count; i++ {
				local, err := iteration.Scope(i, scope, g.Range.Raw, vars)
				if err != nil {
					return err
				}
				if err = visit(g.Body, local); err != nil {
					return err
				}
			}
		}
		for _, child := range g.Children {
			if err := visit(child, scope); err != nil {
				return err
			}
		}
		return nil
	}
	for _, set := range sets {
		scope := map[string]*model.Value{}
		for name, raw := range set {
			kind, _ := paramParts(a.Spec.Inputs[name])
			data, err := json.Marshal(map[string]any{kind: raw})
			if err != nil {
				return err
			}
			var boxed model.Value
			if err = json.Unmarshal(data, &boxed); err != nil {
				return err
			}
			scope[name] = &boxed
		}
		if err := visit(a.graph, scope); err != nil {
			return err
		}
	}
	missing := []string{}
	var check func(*graph)
	check = func(g *graph) {
		if g == nil || g.Fold != nil && !*g.Fold {
			return
		}
		if len(g.Condition) > 0 && g.Fold == nil {
			for _, side := range []bool{false, true} {
				if !branches[g.Path][side] {
					missing = append(missing, fmt.Sprintf("%s.if: %t branch wasn't reached", g.Path, side))
				}
			}
		}
		if g.Body != nil {
			candidates := []any{g.Range.Raw}
			if g.Range.Whole && len(g.Range.Refs) == 1 {
				b := g.Range.Refs[0].Binding
				kind, body := paramParts(b.Param)
				cs, _ := body["constraints"].(map[string]any)
				candidates = nil
				for _, set := range sets {
					candidates = append(candidates, set[b.Name])
				}
				for _, v := range sampleValues(kind, cs) {
					if validCandidate(b.Param, kind, v) {
						candidates = append(candidates, v)
					}
				}
			}
			wanted := map[int]bool{}
			for _, v := range candidates {
				switch x := v.(type) {
				case []any:
					wanted[len(x)] = true
				case map[string]any:
					wanted[len(x)] = true
				}
			}
			for n := range wanted {
				if !lengths[g.Path][n] {
					missing = append(missing, fmt.Sprintf("%s.%s.range: loop length %d wasn't reached", g.Path, g.Kind, n))
				}
			}
			check(g.Body)
		}
		for _, child := range g.Children {
			check(child)
		}
	}
	check(a.graph)
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("verification coverage incomplete:\n  %s", strings.Join(missing, "\n  "))
	}
	return nil
}
