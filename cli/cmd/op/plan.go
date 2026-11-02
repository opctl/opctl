package op

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/opspec/sdk/refscan"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"sort"
	"strings"
)

func newPlanCmd(satisfier cliparamsatisfier.CLIParamSatisfier) *cobra.Command {
	var args []string
	var argFile string
	cmd := &cobra.Command{Use: "plan OP_REF", Short: "Inspect an op graph without running its calls", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, refs []string) error {
		handle, err := dataresolver.New(satisfier, node).Resolve(cmd.Context(), refs[0], nil)
		if err != nil {
			return err
		}
		spec, err := node.DescribeOp(cmd.Context(), model.DescribeOpReq{Op: model.StartOpReqOp{Ref: handle.Ref()}})
		if err != nil {
			return err
		}
		params := goop.PlanParams(spec.Inputs)
		noteIgnoredArgs(cmd.ErrOrStderr(), args, params)
		if argFile != "" {
			data, err := os.ReadFile(argFile)
			if err != nil {
				return err
			}
			var doc yaml.Node
			if err = yaml.Unmarshal(data, &doc); err != nil {
				return err
			}
			if len(doc.Content) > 0 {
				root := doc.Content[0]
				for i := 0; i+1 < len(root.Content); i += 2 {
					if params[root.Content[i].Value] == nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "note: input %s isn't passed to plan; ignored\n", root.Content[i].Value)
					}
				}
			}
		}
		file, err := satisfier.NewYMLFileInputSrc(argFile)
		if err != nil {
			return err
		}
		provided := cliparamsatisfier.NewInputSourcer(satisfier.NewSliceInputSrc(args, "="), file, satisfier.NewEnvVarInputSrc())
		defaults, prompt := satisfier.NewParamDefaultInputSrc(params), satisfier.NewCliPromptInputSrc(params)
		source := cliparamsatisfier.NewInputSourcer(satisfier.NewSliceInputSrc(args, "="), file, satisfier.NewEnvVarInputSrc(), defaults, prompt)
		isGo, err := opfile.IsGoOp(cmd.Context(), handle)
		if err != nil {
			return err
		}
		if isGo {
			source = cliparamsatisfier.NewGoInputSourcer(provided, defaults, prompt)
		}
		values, err := satisfier.Satisfy(source, params)
		if err != nil {
			return err
		}
		plan, err := node.PlanOp(cmd.Context(), model.PlanOpReq{Args: values, Op: model.StartOpReqOp{Ref: handle.Ref()}})
		if err != nil {
			return err
		}
		return renderPlan(cmd.OutOrStdout(), plan, values)
	}}
	cmd.Flags().StringArrayVarP(&args, "args", "a", nil, "Set a non-secret value input (NAME=VALUE)")
	cmd.Flags().StringVar(&argFile, "arg-file", "", "Read input values from a YAML file")
	return cmd
}

func noteIgnoredArgs(w io.Writer, args []string, params map[string]*model.ParamSpec) {
	seen := map[string]bool{}
	for _, arg := range args {
		name, _, _ := strings.Cut(arg, "=")
		if params[name] == nil && !seen[name] {
			fmt.Fprintf(w, "note: input %s isn't passed to plan; ignored\n", name)
			seen[name] = true
		}
	}
}

func renderPlan(w io.Writer, plan *model.OpPlan, inputs map[string]*model.Value) error {
	inputs = planHeaderInputs(plan.Inputs, inputs)
	data, err := json.Marshal(plan.OpSpec)
	if err != nil {
		return err
	}
	var raw map[string]any
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}
	delete(raw, "version")
	renamePlanOutputs(raw, plan.OutputBindings)
	data, err = json.Marshal(raw)
	if err != nil {
		return err
	}
	var doc yaml.Node
	if err = yaml.Unmarshal(data, &doc); err != nil {
		return err
	}
	root := doc.Content[0]
	ordered := []*yaml.Node{}
	for _, key := range []string{"name", "description", "inputs", "outputs", "run"} {
		for i := 0; i+1 < len(root.Content); i += 2 {
			if root.Content[i].Value == key {
				ordered = append(ordered, root.Content[i], root.Content[i+1])
			}
		}
	}
	root.Content = ordered
	redactPlanDefaults(root)
	fmt.Fprintf(w, "# opctl op plan: snapshot of the op graph that %s plans for the inputs below.\n", strings.Join(strings.Fields(plan.Name), " "))
	fmt.Fprintln(w, "# For inspection and debugging only. It isn't a conversion: other inputs can plan a different graph,")
	fmt.Fprintln(w, "# and Go branches, helpers and op stubs don't appear in it.")
	fmt.Fprintln(w, "# Secret, dir, file and socket inputs aren't passed to plan, so they stay references.")
	names := make([]string, 0, len(inputs))
	for name := range inputs {
		names = append(names, name)
	}
	sort.Strings(names)
	params := goop.PlanParams(plan.Inputs)
	for _, name := range names {
		if params[name] == nil {
			continue
		}
		v, err := inputs[name].Unbox()
		if err != nil {
			return err
		}
		text, err := json.Marshal(v)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "#   %s: %s\n", name, text)
	}
	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err = encoder.Encode(&doc); err != nil {
		return err
	}
	if err = encoder.Close(); err != nil {
		return err
	}
	_, err = w.Write(out.Bytes())
	return err
}

func redactPlanDefaults(n *yaml.Node) {
	if n.Kind == yaml.MappingNode {
		secret := false
		for i := 0; i+1 < len(n.Content); i += 2 {
			if n.Content[i].Value == "isSecret" && n.Content[i+1].Value == "true" {
				secret = true
			}
		}
		if secret {
			out := []*yaml.Node{}
			for i := 0; i+1 < len(n.Content); i += 2 {
				if n.Content[i].Value == "default" {
					n.HeadComment = "has a secret default"
					continue
				}
				out = append(out, n.Content[i], n.Content[i+1])
			}
			n.Content = out
		}
	}
	n.Style = 0
	for _, c := range n.Content {
		redactPlanDefaults(c)
	}
}
func renamePlanOutputs(raw map[string]any, bindings map[string]string) {
	used := map[string]bool{}
	if params, ok := raw["inputs"].(map[string]any); ok {
		for n := range params {
			used[n] = true
		}
	}
	var collect func(any)
	collect = func(v any) {
		switch v := v.(type) {
		case map[string]any:
			for k, x := range v {
				if k == "stdout" || k == "stderr" {
					if s, ok := x.(string); ok {
						used[strings.TrimSuffix(strings.TrimPrefix(s, "$("), ")")] = true
					}
				}
				if k == "outputs" {
					if m, ok := x.(map[string]any); ok {
						for _, x := range m {
							if s, ok := x.(string); ok {
								used[strings.TrimSuffix(strings.TrimPrefix(s, "$("), ")")] = true
							}
						}
					}
				}
				collect(x)
			}
		case []any:
			for _, x := range v {
				collect(x)
			}
		}
	}
	collect(raw["run"])
	names := make([]string, 0, len(bindings))
	for n := range bindings {
		names = append(names, n)
	}
	sort.Strings(names)
	rename := map[string]string{}
	for _, n := range names {
		ref := bindings[n]
		if strings.HasPrefix(ref, "$(opctl_") && strings.HasSuffix(ref, ")") && !used[n] {
			old := ref[2 : len(ref)-1]
			if _, ok := rename[old]; !ok {
				rename[old] = n
				used[n] = true
			}
		}
	}
	var rewrite func(any) any
	rewrite = func(v any) any {
		switch v := v.(type) {
		case string:
			v = strings.ReplaceAll(v, "$(opctl_empty_string)", "")
			for old, n := range rename {
				for _, suffix := range []string{")", "/", ".", "["} {
					v = strings.ReplaceAll(v, "$("+old+suffix, "$("+n+suffix)
				}
			}
			return v
		case map[string]any:
			for k, x := range v {
				v[k] = rewrite(x)
			}
		case []any:
			for i, x := range v {
				v[i] = rewrite(x)
			}
		}
		return v
	}
	rewrite(raw["run"])
}

// Defaults containing references are visible in the interface, rather than
// claimed as resolved values by a CLI that cannot read remote op directories.
func planHeaderInputs(params map[string]*model.ParamSpec, provided map[string]*model.Value) map[string]*model.Value {
	out := map[string]*model.Value{}
	for name, v := range provided {
		out[name] = v
	}
	for name, param := range goop.PlanParams(params) {
		if out[name] != nil {
			continue
		}
		data, _ := json.Marshal(param)
		var raw map[string]map[string]json.RawMessage
		_ = json.Unmarshal(data, &raw)
		for kind, body := range raw {
			d, ok := body["default"]
			if !ok {
				continue
			}
			var text string
			if json.Unmarshal(d, &text) == nil && refscan.Contains(text) {
				continue
			}
			boxed, _ := json.Marshal(map[string]json.RawMessage{kind: d})
			var value model.Value
			if json.Unmarshal(boxed, &value) == nil {
				out[name] = &value
			}
		}
	}
	return out
}
