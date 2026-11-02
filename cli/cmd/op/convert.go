package op

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/opstubs"
	"github.com/opctl/opctl/sdks/go/data"
	datafs "github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/goconvert"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func newConvertCmd(satisfier cliparamsatisfier.CLIParamSatisfier) *cobra.Command {
	var goFlag, dryRun, force bool
	var args []string
	var argFile string
	cmd := &cobra.Command{Use: "convert --go OP_REF", Short: "Convert a local op file to a verified Go op", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, refs []string) error {
		if !goFlag {
			return fmt.Errorf("pass --go; op convert only converts op.yml to Go")
		}
		if strings.Contains(refs[0], "#") {
			return fmt.Errorf("op convert needs a local op")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		handle, err := data.Resolve(cmd.Context(), refs[0], datafs.New(filepath.Join(cwd, ".opspec"), cwd))
		if err != nil {
			return err
		}
		dir, err := filepath.Abs(handle.Ref())
		if err != nil {
			return err
		}
		source, err := os.ReadFile(filepath.Join(dir, "op.yml"))
		if err != nil {
			return err
		}
		a, err := goconvert.Analyze(cmd.Context(), dir, source, func(ctx context.Context, ref string) (*model.OpSpec, error) {
			return node.DescribeOp(ctx, model.DescribeOpReq{Op: model.StartOpReqOp{Ref: ref}})
		})
		if err != nil {
			return err
		}
		if len(a.Refusals) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "opctl op convert: %s can't be converted to Go; nothing was written.\n\n", refs[0])
			for _, r := range a.Refusals {
				fmt.Fprintf(cmd.ErrOrStderr(), "  %s\n", r)
			}
			fmt.Fprintln(cmd.ErrOrStderr(), "\nKeep these calls in YAML: move each into its own YAML op, call that op from here, and rerun convert.\nSee the converting-opctl-ops skill for the split.")
			return fmt.Errorf("conversion refused")
		}
		sdk := opspec.GoSDK{Version: sdkVersion, Sum: sdkSum, GoModSum: sdkGoModSum}
		files, err := goconvert.PrepareFiles(dir, sdk)
		if err != nil {
			return fmt.Errorf("op.yml:1:1  name  %w; nothing was written", err)
		}
		rendered, err := goconvert.Render(a, goconvert.RenderOptions{ModuleRoot: files.ModuleRoot, ModulePath: files.ModulePath})
		if err != nil {
			return err
		}
		files.Files["op.go"] = rendered.Source
		supplied, err := convertInputValues(cmd, satisfier, args, argFile, a.Spec.Inputs)
		if err != nil {
			return err
		}
		sets, verifyErr := goconvert.InputSets(a.Spec, dir, supplied)
		if verifyErr == nil && files.NewModule && !sdk.Published() {
			verifyErr = fmt.Errorf("the Go op SDK isn't published in this opctl build; add a replace directive to go.mod or convert inside a module that requires the SDK")
		}
		if verifyErr == nil {
			verifyErr = verifyConversion(cmd.Context(), a, files, sets)
		}
		if dryRun {
			names := make([]string, 0, len(files.Files))
			for name := range files.Files {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Fprintf(cmd.OutOrStdout(), "// ==> %s <==\n%s\n", name, files.Files[name])
			}
		}
		for _, note := range a.Notes {
			fmt.Fprintf(cmd.ErrOrStderr(), "note: %s: %s\n", note.Path, note.Reason)
		}
		todos := []string{}
		for i, line := range strings.Split(string(rendered.Source), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "// TODO(opctl op convert):") {
				todos = append(todos, fmt.Sprintf("op.go:%d %s", i+1, strings.TrimSpace(strings.TrimPrefix(line, "// TODO(opctl op convert):"))))
			}
		}
		if len(todos) > 0 {
			fmt.Fprintf(cmd.ErrOrStderr(), "%d TODO: %s\n", len(todos), strings.Join(todos, "; "))
		}
		if verifyErr != nil && (!force || dryRun) {
			return fmt.Errorf("verification failed: %w\nNothing was written. Inspect op.go with --dry-run, or write it anyway with --force", verifyErr)
		}
		if dryRun {
			fmt.Fprintf(cmd.ErrOrStderr(), "verified: the Go op plans the op.yml graph for %d input sets; nothing was written\n", len(sets))
			return nil
		}
		if err = goconvert.WriteFiles(dir, files.Files); err != nil {
			return err
		}
		if verifyErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: wrote unverified op.go: %s\n", verifyErr)
		} else {
			fmt.Fprintf(cmd.ErrOrStderr(), "verified: the Go op plans the op.yml graph for %d input sets\n", len(sets))
		}
		if files.NewModule && !sdk.Published() {
			fmt.Fprintln(cmd.ErrOrStderr(), "This opctl build has no published Go op SDK. Add a replace directive for github.com/opctl/opctl/sdks/go/opspec/sdk to go.mod before running the op.")
		}
		if err = opstubs.RefreshModule(cmd.Context(), files.ModuleRoot, node); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: opctlgen/ not refreshed: %s; run opctl op gen for details\n", err)
		}
		names := make([]string, 0, len(files.Files))
		for name := range files.Files {
			names = append(names, name)
		}
		sort.Strings(names)
		fmt.Fprintf(cmd.ErrOrStderr(), "converted %s to a Go op: wrote %s; removed op.yml\n", refs[0], strings.Join(names, ", "))
		return nil
	}}
	cmd.Flags().BoolVar(&goFlag, "go", false, "Convert op.yml to a Go op")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Verify and print files without replacing op.yml")
	cmd.Flags().BoolVar(&force, "force", false, "Write despite verification failure; refusals still stop conversion")
	cmd.Flags().StringArrayVarP(&args, "args", "a", nil, "Set a non-secret verification input (NAME=VALUE)")
	cmd.Flags().StringVar(&argFile, "arg-file", "", "Read verification inputs from a YAML file")
	return cmd
}

func convertInputValues(cmd *cobra.Command, satisfier cliparamsatisfier.CLIParamSatisfier, args []string, argFile string, all map[string]*model.ParamSpec) (map[string]any, error) {
	params := goop.PlanParams(all)
	noteIgnoredArgs(cmd.ErrOrStderr(), args, params)
	if argFile != "" {
		data, err := os.ReadFile(argFile)
		if err != nil {
			return nil, err
		}
		var doc yaml.Node
		if err = yaml.Unmarshal(data, &doc); err != nil {
			return nil, err
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
		return nil, err
	}
	source := cliparamsatisfier.NewInputSourcer(satisfier.NewSliceInputSrc(args, "="), file)
	values := map[string]any{}
	names := make([]string, 0, len(params))
	for name := range params {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		raw, ok := source.Source(name)
		if !ok || raw == nil {
			continue
		}
		p := params[name]
		kind := "string"
		var value any = *raw
		switch {
		case p.Array != nil:
			kind = "array"
		case p.Object != nil:
			kind = "object"
		case p.Boolean != nil:
			kind = "boolean"
		case p.Number != nil:
			kind = "number"
		}
		if kind != "string" {
			if err = yaml.Unmarshal([]byte(*raw), &value); err != nil {
				return nil, fmt.Errorf("input %s: %w", name, err)
			}
		}
		b, err := json.Marshal(map[string]any{kind: value})
		if err != nil {
			return nil, err
		}
		var boxed model.Value
		if err = json.Unmarshal(b, &boxed); err != nil {
			return nil, fmt.Errorf("input %s: %w", name, err)
		}
		if errs := param.Validate(&boxed, p, true); len(errs) > 0 {
			return nil, fmt.Errorf("input %s: %v", name, errs)
		}
		values[name], err = boxed.Unbox()
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

func verifyConversion(ctx context.Context, a *goconvert.Analysis, files *goconvert.FileSet, sets []map[string]any) error {
	stage, err := goconvert.StageFiles(a, files)
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	for _, set := range sets {
		if err = ctx.Err(); err != nil {
			return err
		}
		args := map[string]*model.Value{}
		for name, value := range set {
			p := goop.PlanParams(a.Spec.Inputs)[name]
			kind := "string"
			switch {
			case p.Array != nil:
				kind = "array"
			case p.Object != nil:
				kind = "object"
			case p.Boolean != nil:
				kind = "boolean"
			case p.Number != nil:
				kind = "number"
			}
			data, _ := json.Marshal(map[string]any{kind: value})
			var v model.Value
			if err = json.Unmarshal(data, &v); err != nil {
				return err
			}
			args[name] = &v
		}
		planned, err := node.PlanOp(ctx, model.PlanOpReq{Op: model.StartOpReqOp{Ref: stage}, Args: args})
		if err != nil {
			return err
		}
		described := *planned.OpSpec
		described.Run = nil
		// Name comes from the final directory when the source has no explicit override.
		if described.Name == filepath.Base(stage) {
			described.Name = filepath.Base(a.Dir)
		}
		if err = goconvert.Compare(a, &described, &goop.Plan{OpSpec: &model.OpSpec{Run: planned.Run}, OutputBindings: planned.OutputBindings}, set); err != nil {
			inputJSON, _ := json.Marshal(set)
			return fmt.Errorf("inputs %s: %w", inputJSON, err)
		}
	}
	return nil
}
