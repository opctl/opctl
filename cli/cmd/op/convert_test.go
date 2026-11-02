package op

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestConvertCommand(t *testing.T) {
	for _, mode := range []string{"dry-run", "write", "mismatch", "dry-mismatch", "graph-mismatch", "force", "refuse", "child"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "hello")
			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}
			source := []byte("name: hello\ninputs: {token: {string: {isSecret: true}}}\nrun: {container: {image: {ref: alpine}, envVars: {AUTH: 'Bearer $(token)'}}}\n")
			if mode == "child" {
				source = []byte("name: hello\nrun: {op: {ref: ./child}}\n")
			}
			if mode == "refuse" {
				source = []byte("name: hello\nrun: {name: named, container: {image: {ref: alpine}}}\n")
			}
			if err := os.WriteFile(filepath.Join(dir, "op.yml"), source, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module ops.local/test\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n"), 0600); err != nil {
				t.Fatal(err)
			}
			spec, err := opfile.Unmarshal(source)
			if err != nil {
				t.Fatal(err)
			}
			api := new(fakes.FakeNode)
			api.DescribeOpReturns(spec, nil)
			api.PlanOpCalls(func(ctx context.Context, req model.PlanOpReq) (*model.OpPlan, error) {
				if filepath.Dir(req.Op.Ref) != root || !strings.Contains(filepath.Base(req.Op.Ref), ".hello.opctl-convert-") {
					t.Fatalf("wrong staging directory: %s", req.Op.Ref)
				}
				if _, err := os.Stat(filepath.Join(req.Op.Ref, "op.go")); err != nil {
					t.Fatal(err)
				}
				if os.Getuid() == 0 {
					for _, path := range []string{req.Op.Ref, filepath.Join(req.Op.Ref, "op.go")} {
						info, err := os.Stat(path)
						if err != nil || info.Sys().(*syscall.Stat_t).Uid != 65534 {
							t.Fatalf("staged path isn't owned by the sudo caller: %s (%v)", path, err)
						}
					}
				}
				if len(req.Args) != 0 {
					t.Fatal("secret reached plan")
				}
				result := *spec
				result.Name = filepath.Base(req.Op.Ref)
				if mode == "child" {
					data, err := os.ReadFile(filepath.Join(req.Op.Ref, "op.go"))
					if err != nil {
						t.Fatal(err)
					}
					if !strings.Contains(string(data), "../hello/child") {
						t.Fatalf("staged child ref: %s", data)
					}
					result.Run = &model.CallSpec{Op: &model.OpCallSpec{Ref: "../hello/child"}}
				}
				if mode == "graph-mismatch" {
					result.Run = &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "other"}}}
				}
				if mode == "mismatch" || mode == "dry-mismatch" || mode == "force" {
					return nil, fmt.Errorf("test bootstrap unavailable")
				}
				return &model.OpPlan{OpSpec: &result, OutputBindings: map[string]string{}}, nil
			})
			saved := node
			node = api
			defer func() { node = saved }()
			cmd := newConvertCmd(cliparamsatisfier.New(nil))
			if os.Getuid() == 0 {
				t.Setenv("SUDO_UID", "65534")
				t.Setenv("SUDO_GID", "65534")
			}
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			args := []string{"--go", dir, "-a", "token=SECRET-MARKER"}
			if mode == "dry-run" || mode == "dry-mismatch" {
				args = append(args, "--dry-run")
			}
			if mode == "force" || mode == "refuse" {
				args = append(args, "--force")
			}
			cmd.SetArgs(args)
			err = cmd.Execute()
			if (mode == "mismatch" || mode == "dry-mismatch" || mode == "graph-mismatch" || mode == "refuse") != (err != nil) {
				t.Fatalf("%s: %v", mode, err)
			}
			written := mode == "write" || mode == "force" || mode == "child"
			_, goErr := os.Stat(filepath.Join(dir, "op.go"))
			_, yamlErr := os.Stat(filepath.Join(dir, "op.yml"))
			if written != (goErr == nil) || written != (os.IsNotExist(yamlErr)) {
				t.Fatalf("wrong file state: %v %v", goErr, yamlErr)
			}
			matches, _ := filepath.Glob(filepath.Join(root, ".hello.opctl-convert-*"))
			if len(matches) > 0 {
				t.Fatal("staging directory leaked")
			}
			if strings.Contains(stdout.String()+stderr.String(), "SECRET-MARKER") {
				t.Fatal("secret leaked")
			}
			if (mode == "dry-run" || mode == "dry-mismatch") && !strings.Contains(stdout.String(), "// ==> op.go <==") {
				t.Fatal(stdout.String())
			}
			if mode == "write" && !strings.Contains(stderr.String(), "1 TODO") {
				t.Fatal(stderr.String())
			}
			if mode == "force" && !strings.Contains(stderr.String(), "warning: wrote unverified op.go") {
				t.Fatal(stderr.String())
			}
			if mode == "refuse" && api.PlanOpCallCount() != 0 {
				t.Fatal("planned refused op")
			}
			if api.StartOpCallCount() != 0 {
				t.Fatal("conversion ran an op")
			}
		})
	}
}

func TestConvertReportsUnreachedControlFlow(t *testing.T) {
	for _, force := range []bool{false, true} {
		t.Run(fmt.Sprint(force), func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "coverage")
			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}
			source := []byte("name: coverage\ninputs: {zone: {string: {default: dev, constraints: {enum: [dev]}}}, zones: {array: {default: []}}}\nrun:\n  serial:\n  - container: {image: {ref: alpine}}\n  - if: [{eq: ['$(zone)', prod]}]\n    serialLoop: {range: '$(zones)', run: {container: {image: {ref: alpine}}}}\n")
			if err := os.WriteFile(filepath.Join(dir, "op.yml"), source, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module ops.local/test\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n"), 0600); err != nil {
				t.Fatal(err)
			}
			spec, err := opfile.Unmarshal(source)
			if err != nil {
				t.Fatal(err)
			}
			api := new(fakes.FakeNode)
			api.DescribeOpReturns(spec, nil)
			api.PlanOpCalls(func(context.Context, model.PlanOpReq) (*model.OpPlan, error) {
				t.Fatal("incomplete coverage reached node verification")
				return nil, nil
			})
			saved := node
			node = api
			defer func() { node = saved }()
			cmd := newConvertCmd(cliparamsatisfier.New(nil))
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			args := []string{"--go", dir}
			if force {
				args = append(args, "--force")
			}
			cmd.SetArgs(args)
			err = cmd.Execute()
			if (err == nil) != force {
				t.Fatalf("force=%v: %v", force, err)
			}
			diagnostic := stderr.String()
			if err != nil {
				diagnostic += err.Error()
			}
			for _, want := range []string{"true branch wasn't reached", "loop length 0 wasn't reached", "loop length 2 wasn't reached"} {
				if !strings.Contains(diagnostic, want) {
					t.Fatalf("missing %q: %s", want, diagnostic)
				}
			}
			_, goErr := os.Stat(filepath.Join(dir, "op.go"))
			if force != (goErr == nil) {
				t.Fatalf("force=%v: op.go %v", force, goErr)
			}
		})
	}
}
