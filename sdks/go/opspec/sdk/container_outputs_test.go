package sdk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func namedGraph() Call {
	p := Container("alpine:3").Cmd("sh", "-ce", `printf '%s' '{"version":"v1.2.3","count":3,"ready":true,"targets":["linux"],"metadata":{"revision":"abc123"}}' > "$OPCTL_OUTPUTS"`)
	p, version := ContainerOutput[String](p, "version")
	p, count := ContainerOutput[Number](p, "count")
	p, ready := ContainerOutput[Boolean](p, "ready")
	p, targets := ContainerOutput[Array[string]](p, "targets")
	p, metadata := ContainerOutput[Object[map[string]string]](p, "metadata")
	return Serial(p, Container("alpine:3").EnvFrom("VERSION", version).EnvFrom("COUNT", count).EnvFrom("READY", ready).EnvFrom("TARGETS", targets).EnvFrom("METADATA", metadata).Cmd("true"))
}
func TestContainerOutputGolden(t *testing.T) {
	fixture, _ := filepath.Abs("testdata/plan.container-outputs.json")
	helloDir(t)
	data, err := Plan(func(struct{}) Call { return namedGraph() }, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(expected) {
		t.Fatalf("named fixture mismatch: %s", data)
	}
}
func TestContainerOutputTypesAndCopies(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		base := Container("x")
		p, v := ContainerOutput[String](base, "member")
		again, w := ContainerOutput[String](p, "member")
		if v.handle() != w.handle() {
			t.Fatal("repeated accessor allocated twice")
		}
		independent, other := ContainerOutput[String](base, "member")
		if other.handle() == v.handle() {
			t.Fatal("base builder mutated")
		}
		return Serial(again.Env("A", "B"), independent, Container("x").EnvFrom("V", v).EnvFrom("W", other))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	failures := []func() Call{
		func() Call { p, _ := ContainerOutput[Dir](Container("x"), "x"); return p },
		func() Call { p, _ := ContainerOutput[File](Container("x"), "x"); return p },
		func() Call { p, _ := ContainerOutput[Socket](Container("x"), "x"); return p },
		func() Call { p, _ := ContainerOutput[*String](Container("x"), "x"); return p },
		func() Call { p, _ := ContainerOutput[Handle](Container("x"), "x"); return p },
		func() Call {
			p, _ := ContainerOutput[String](Container("x"), "x")
			p, _ = ContainerOutput[Number](p, "x")
			return p
		},
		func() Call {
			p, _ := ContainerOutput[Array[string]](Container("x"), "x")
			p, _ = ContainerOutput[Array[int]](p, "x")
			return p
		},
		func() Call { p, _ := ContainerOutput[String](Container("x"), "bad/name"); return p },
		func() Call { p, v := ContainerOutput[String](Container("x"), "x"); return p.EnvFrom("SELF", v) },
		func() Call {
			p, v := ContainerOutput[String](Container("x"), "x")
			return Parallel(p, Container("x").EnvFrom("V", v))
		},
		func() Call { p, _ := ContainerOutput[String](Container("x"), "x"); return Serial(p, p.Cmd("other")) },
		func() Call {
			_, v := ContainerOutput[String](Container("x"), "x")
			return Container("x").EnvFrom("V", v)
		},
		func() Call { p, v := ContainerOutput[Number](Container("x"), "x"); v.Value(); return p },
		func() Call { p, v := ContainerOutput[Boolean](Container("x"), "x"); v.Value(); return p },
		func() Call { p, v := ContainerOutput[Array[string]](Container("x"), "x"); v.Value(); return p },
		func() Call { p, v := ContainerOutput[Object[map[string]int]](Container("x"), "x"); v.Value(); return p },
	}
	for i, build := range failures {
		if _, err := Plan(func(struct{}) Call { return build() }, nil); err == nil {
			t.Fatalf("accepted invalid named output case %d", i)
		}
	}
	p, v := ContainerOutput[String](Container("x"), "x")
	_, err = Plan(func(struct{}) Call { return p.EnvFrom("V", v) }, nil)
	if err == nil {
		t.Fatal("accepted named accessor outside plan")
	}
}
func TestContainerOutputSharesAllocation(t *testing.T) {
	data, err := Plan(func(struct{}) Call {
		dir := NewDir()
		op, number := OpOutput[Number](CallOp("./yaml"), "count")
		p, stdout := Container("x").Dir("/out", dir).Stdout()
		p, member := ContainerOutput[String](p, "stdout")
		return Serial(op, p, Container("x").EnvFrom("N", number).EnvFrom("S", stdout).EnvFrom("M", member))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"opctl_dir_1", "opctl_number_2", "opctl_string_3", "opctl_string_4"} {
		if !strings.Contains(string(data), name) {
			t.Fatal("allocator not shared", name)
		}
	}
}

func TestContainerOutputSchemaRequests(t *testing.T) {
	for _, max := range []int{1, 2, 3} {
		request := fmt.Sprintf(`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":%d},"inputs":{},"knownOutputs":{}}`, max)
		path := filepath.Join(t.TempDir(), "request.json")
		if err := os.WriteFile(path, []byte(request), 0600); err != nil {
			t.Fatal(err)
		}
		command := exec.Command(os.Args[0], "-test.run=^TestStreamMainHelper$")
		command.Env = append(os.Environ(), "OPCTL_STREAM_HELPER=1", "OPCTL_NAMED_HELPER=1", "OPCTL_STREAM_PHASE=plan", "OPCTL_BOOTSTRAP_REQUEST="+path)
		var stderr strings.Builder
		command.Stderr = &stderr
		out, err := command.Output()
		if max < 3 {
			if err == nil || len(out) != 0 || !strings.Contains(stderr.String(), "schema 3") {
				t.Fatalf("unsupported named schema: %s %s %v", out, stderr.String(), err)
			}
		} else if err != nil || !strings.Contains(string(out), `"opspecSchema":3`) {
			t.Fatalf("supported named schema: %s %s %v", out, stderr.String(), err)
		}
	}
}
