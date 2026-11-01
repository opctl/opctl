package sdk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type typedFixtureOutputs struct {
	Count Number `opctl:"count,minimum=1"`
	Text  String `opctl:"text,description='Output text'"`
}

func typedFixtureGraph() Call {
	producer, out := Outputs[typedFixtureOutputs](Container("alpine:3").Cmd("true"))
	return Serial(producer, Container("alpine:3").EnvFrom("COUNT", out.Count).EnvFrom("TEXT", out.Text).Cmd("true"))
}
func TestTypedContainerOutputGolden(t *testing.T) {
	fixture, _ := filepath.Abs("testdata/plan.typed-container-outputs.json")
	helloDir(t)
	expected, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		data, err := Plan(func(struct{}) Call { return typedFixtureGraph() }, nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != string(expected) {
			t.Fatalf("typed fixture mismatch: %s", data)
		}
	}
}
func TestTypedContainerOutputDeclarations(t *testing.T) {
	type Nested struct {
		Revision string `json:"revision"`
	}
	type Out struct {
		Count    Number         `opctl:"count,format=integer,minimum=1"`
		Ready    Boolean        `opctl:"ready"`
		Tags     Array[string]  `opctl:"tags,minItems=1"`
		Metadata Object[Nested] `opctl:"metadata,required=revision,secret"`
	}
	data, err := Plan(func(struct{}) Call {
		p, out := Outputs[Out](Container("alpine"))
		return Serial(p, Container("alpine").EnvFrom("N", out.Count).EnvFrom("M", out.Metadata))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"opspecSchema":4`, `"outputParams"`, `"items":{"type":"string"}`, `"required":["revision"]`, `"isSecret":true`, `"minimum":1`} {
		if !strings.Contains(string(data), want) {
			t.Fatal("missing declaration " + want)
		}
	}
	desc, err := DescribeWithOutputs[struct{}, Out]()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(desc), "outputParams") || !strings.Contains(string(desc), `"opspecSchema":1`) {
		t.Fatal("describe acquired a graph feature")
	}
}
func TestTypedContainerOutputReuseAndCopies(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		base := Container("x")
		p, original := ContainerOutput[Number](base, "count")
		p, out := Outputs[struct {
			Count Number `opctl:"count,minimum=1"`
		}](p)
		again, more := Outputs[struct {
			Count Number `opctl:"count,minimum=1"`
		}](p)
		again, untyped := ContainerOutput[Number](again, "count")
		if original.handle() != out.Count.handle() || out.Count.handle() != more.Count.handle() || out.Count.handle() != untyped.handle() {
			t.Fatal("same result was allocated twice")
		}
		if len(base.container.Outputs) != 0 || len(base.container.OutputParams) != 0 {
			t.Fatal("input builder mutated")
		}
		copied, _ := Outputs[struct {
			Other String `opctl:"other"`
		}](p)
		if len(p.container.OutputParams) != 1 || len(copied.container.OutputParams) != 2 {
			t.Fatal("declaration maps shared")
		}
		return Serial(again, Container("x").EnvFrom("COUNT", out.Count))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
func TestInvalidTypedContainerOutputs(t *testing.T) {
	failures := []func() Call{
		func() Call { p, _ := Outputs[int](Container("x")); return p },
		func() Call { p, _ := Outputs[struct{}](Container("x")); return p },
		func() Call {
			p, _ := Outputs[struct {
				X string `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X *String `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X Dir `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X File `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X Socket `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X String `opctl:"x,default=''"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X Boolean `opctl:"x,secret"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X String `opctl:"opctl_x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X String `opctl:"x"`
				Y Number `opctl:"x"`
			}](Container("x"))
			return p
		},
		func() Call {
			p, _ := Outputs[struct {
				X Number `opctl:"x,minimum=1"`
			}](Container("x"))
			p, _ = Outputs[struct {
				X Number `opctl:"x,minimum=2"`
			}](p)
			return p
		},
		func() Call {
			p, _ := ContainerOutput[Array[string]](Container("x"), "x")
			p, _ = Outputs[struct {
				X Array[int] `opctl:"x"`
			}](p)
			return p
		},
		func() Call {
			_, out := Outputs[struct {
				Good String `opctl:"good"`
				Bad  Number `opctl:"bad,default=1"`
			}](Container("x"))
			return Container("x").EnvFrom("V", out.Good)
		},
		func() Call {
			p, out := Outputs[typedFixtureOutputs](Container("x"))
			return Parallel(p, Container("x").EnvFrom("V", out.Text))
		},
		func() Call { p, out := Outputs[typedFixtureOutputs](Container("x")); return p.EnvFrom("V", out.Text) },
		func() Call { p, _ := Outputs[typedFixtureOutputs](Container("x")); return Serial(p, p.Cmd("other")) },
		func() Call {
			_, out := Outputs[typedFixtureOutputs](Container("x"))
			return Container("x").EnvFrom("V", out.Text)
		},
		func() Call { p, out := Outputs[typedFixtureOutputs](Container("x")); out.Count.Value(); return p },
	}
	for i, fn := range failures {
		if _, err := Plan(func(struct{}) Call { return fn() }, nil); err == nil {
			t.Fatalf("accepted invalid case %d", i)
		}
	}
	p, out := Outputs[typedFixtureOutputs](Container("x"))
	if _, err := Plan(func(struct{}) Call { return p.EnvFrom("V", out.Text) }, nil); err == nil {
		t.Fatal("accepted output declaration outside plan")
	}
}
func TestTypedContainerOutputSchemaRange(t *testing.T) {
	for _, max := range []int{1, 2, 3, 4} {
		request := fmt.Sprintf(`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":%d},"inputs":{},"knownOutputs":{}}`, max)
		path := filepath.Join(t.TempDir(), "request.json")
		if err := os.WriteFile(path, []byte(request), 0600); err != nil {
			t.Fatal(err)
		}
		command := exec.Command(os.Args[0], "-test.run=^TestStreamMainHelper$")
		command.Env = append(os.Environ(), "OPCTL_STREAM_HELPER=1", "OPCTL_TYPED_HELPER=1", "OPCTL_STREAM_PHASE=plan", "OPCTL_BOOTSTRAP_REQUEST="+path)
		var stderr strings.Builder
		command.Stderr = &stderr
		out, err := command.Output()
		if max < 4 {
			if err == nil || len(out) != 0 || !strings.Contains(stderr.String(), "this op uses typed container outputs (opspec >= 4)") {
				t.Fatal("unsupported typed graph did not fail clearly before emission")
			}
		} else if err != nil || !strings.Contains(string(out), `"opspecSchema":4`) {
			t.Fatalf("supported typed graph failed: %v", err)
		}
	}
}
func TestTypedResultCreationSite(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		_, out := Outputs[typedFixtureOutputs](Container("x"))
		return Container("x").EnvFrom("COUNT", out.Count)
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "typed_outputs_test.go:") {
		t.Fatalf("lost author creation site: %v", err)
	}
}
