package sdk

import (
	"strings"
	"testing"
)

func TestOpOutputVisibility(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(OpCall, Dir) Call
		want  string
	}{
		{"serial", func(c OpCall, d Dir) Call { return Serial(c, Container("alpine").Dir("/out", d)) }, ""},
		{"omitted", func(c OpCall, d Dir) Call { return Container("alpine").Dir("/out", d) }, "before"},
		{"earlier", func(c OpCall, d Dir) Call { return Serial(Container("alpine").Dir("/out", d), c) }, "before"},
		{"parallel", func(c OpCall, d Dir) Call { return Parallel(c, Container("alpine").Dir("/out", d)) }, "before"},
		{"after parallel", func(c OpCall, d Dir) Call {
			return Serial(Parallel(c, Container("alpine")), Container("alpine").Dir("/out", d))
		}, ""},
		{"duplicate", func(c OpCall, d Dir) Call { return Serial(c, c) }, "produced"},
		{"self input", func(c OpCall, d Dir) Call { return c.Inputs(map[string]any{"self": d}) }, "before"},
		{"raw", func(c OpCall, d Dir) Call { return Container("alpine").Dir("/out", RawRef[Dir](d.handle())) }, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Plan(func(struct{}) Call { c, d := OpOutput[Dir](CallOp("./yaml"), "artifact"); return tc.build(c, d) }, nil)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err=%v want %s", err, tc.want)
			}
		})
	}
}
func TestOpOutputTypesAndCopies(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		base := CallOp("./yaml")
		a, x := OpOutput[String](base, "text")
		a, y := OpOutput[String](a, "text")
		if x.handle() != y.handle() {
			t.Fatal("same output allocated twice")
		}
		b, z := OpOutput[String](base, "text")
		if x.handle() == z.handle() {
			t.Fatal("copy reused output")
		}
		return Serial(a, b)
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Plan(func(struct{}) Call {
		c, _ := OpOutput[String](CallOp("x"), "out")
		c, _ = OpOutput[Number](c, "out")
		return c
	}, nil)
	if err == nil {
		t.Fatal("accepted conflicting output types")
	}
	_, err = Plan(func(struct{}) Call { c, v := OpOutput[String](CallOp("x"), "out"); v.Value(); return c }, nil)
	if err == nil || !strings.Contains(err.Error(), "op call outputs can't be read during plan") {
		t.Fatalf("%v", err)
	}
	c, _ := OpOutput[Dir](CallOp("x"), "out")
	if _, err := Plan(func(struct{}) Call { return c }, nil); err == nil {
		t.Fatal("accepted accessor outside plan")
	}
}
func TestOpOutputAllKinds(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		c := CallOp("x")
		c, _ = OpOutput[String](c, "string")
		c, _ = OpOutput[Number](c, "number")
		c, _ = OpOutput[Boolean](c, "boolean")
		c, _ = OpOutput[Array[string]](c, "array")
		c, _ = OpOutput[Object[map[string]int]](c, "object")
		c, _ = OpOutput[Dir](c, "dir")
		c, _ = OpOutput[File](c, "file")
		c, _ = OpOutput[Socket](c, "socket")
		return c
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range []func() Call{
		func() Call { c, _ := OpOutput[*String](CallOp("x"), "out"); return c },
		func() Call { c, _ := OpOutput[Handle](CallOp("x"), "out"); return c },
		func() Call { c, _ := OpOutput[String](CallOp("x"), ""); return c },
	} {
		if _, err := Plan(func(struct{}) Call { return fn() }, nil); err == nil {
			t.Fatal("accepted invalid output accessor")
		}
	}
}
