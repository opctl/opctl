package sdk

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLiteralHandles(t *testing.T) {
	type config struct {
		Text string `json:"text"`
	}
	v := config{Text: "$(not-a-ref)"}
	h := ObjectOf(v)
	v.Text = "changed"
	if h.Value().Text != "$(not-a-ref)" || StringOf("").Value() != "" || NumberOf(2).Value() != 2 || !BooleanOf(true).Value() || ArrayOf([]string{"x"}).Value()[0] != "x" {
		t.Fatal("literal values or snapshot changed")
	}
	data, err := Plan(func(struct{}) Call {
		return Serial(CallOp("../callee").Inputs(map[string]any{"empty": StringOf(""), "config": h, "number": NumberOf(2), "array": ArrayOf([]string{"$(literal)"}), "bool": BooleanOf(true)}).PullCreds(StringOf("user"), StringOf("$(password)")), Container("alpine:3").EnvFrom("TEXT", StringOf("$(literal)")).EnvFrom("CONFIG", h).EnvFrom("EMPTY", StringOf("")))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join("testdata", "plan.literal-handles.json")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, want) {
		t.Fatalf("got %s; want %s", data, want)
	}
}

func TestInvalidLiteralHandles(t *testing.T) {
	for _, h := range []Handle{NumberOf(math.Inf(1)), NumberOf(math.NaN()), ObjectOf(make(chan int)), ArrayOf[string](nil)} {
		_, err := Plan(func(struct{}) Call { return CallOp("../callee").Inputs(map[string]any{"value": h}) }, nil)
		if err == nil {
			t.Fatalf("accepted invalid literal %T", h)
		}
	}
	type outputs struct {
		Text String `opctl:"text"`
	}
	_, err := PlanWithOutputs(func(struct{}) (Call, outputs) { return Container("alpine:3"), outputs{StringOf("constant")} }, nil)
	if err == nil || !strings.Contains(err.Error(), "use a default tag for constants") {
		t.Fatalf("got %v", err)
	}
}
