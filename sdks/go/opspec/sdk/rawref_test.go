package sdk

import (
	"bytes"
	"strings"
	"testing"
)

func TestRawRefBindings(t *testing.T) {
	refs := []Handle{
		RawRef[String]("$(text)"), RawRef[Number]("$(number)"), RawRef[Boolean]("$(boolean)"),
		RawRef[Array[string]]("$(array)"), RawRef[Object[map[string]int]]("$(object)"),
		RawRef[Dir]("$(dir)"), RawRef[File]("$(file)"), RawRef[Socket]("$(socket)"),
	}
	for _, handle := range refs {
		data, err := Plan(func(struct{}) Call { return Container("alpine").EnvFrom("VALUE", handle) }, nil)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte(`\\$`)) {
			t.Fatalf("escaped RawRef: %s", data)
		}
	}
	data, err := Plan(func(struct{}) Call {
		return Container("alpine").
			Dir("/src", RawRef[Dir]("$(./scripts)")).File("/cfg", RawRef[File]("$(cfg.$(key))")).
			Socket("/sock", RawRef[Socket]("$(socket)")).Env("LITERAL", "$(input)").Cmd("echo", "$(input)")
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{`"/src":"$(./scripts)"`, `"/cfg":"$(cfg.$(key))"`, `"/sock":"$(socket)"`, `"LITERAL":"\\$(input)"`, `"cmd":["echo","\\$(input)"]`} {
		if !bytes.Contains(data, []byte(part)) {
			t.Fatalf("missing %s in %s", part, data)
		}
	}
}

func TestRawRefRejectsPartialReferences(t *testing.T) {
	for _, ref := range []string{"", "text", "prefix $(x)", "$(x) suffix", "$()", "$(x)\n", "$(a\nb)"} {
		_, err := Plan(func(struct{}) Call { return Container("a").File("/x", RawRef[File](ref)) }, nil)
		if err == nil || !strings.Contains(err.Error(), "RawRef needs a whole $(...) reference") {
			t.Fatalf("ref %q: %v", ref, err)
		}
	}
}

func TestRawRefValuesCannotBeRead(t *testing.T) {
	for _, read := range []func(){
		func() { RawRef[String]("$(text)").Value() },
		func() { RawRef[Number]("$(number)").Value() },
		func() { RawRef[Boolean]("$(boolean)").Value() },
		func() { RawRef[Array[string]]("$(array)").Value() },
		func() { RawRef[Object[map[string]int]]("$(object)").Value() },
	} {
		_, err := Plan(func(struct{}) Call { read(); return Container("a") }, nil)
		if err == nil || err.Error() != "RawRef values can't be read during plan" {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestRawRefBypassesVisibility(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		NewFile()
		return Container("a").EnvFrom("CONTENT", RawRef[File]("$(opctl_file_1)"))
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Overwriting an unchecked binding with a typed handle restores validation.
	_, err = Plan(func(struct{}) Call {
		f := NewFile()
		return Container("a").EnvFrom("CONTENT", RawRef[File]("$(opctl_file_1)")).EnvFrom("CONTENT", f)
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "before any earlier call mounts it") {
		t.Fatalf("overwritten RawRef bypassed visibility: %v", err)
	}
}

func TestRawRefCopiesDoNotChangeVisibility(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		f := NewFile()
		base := Container("a").EnvFrom("CONTENT", f)
		_ = base.EnvFrom("CONTENT", RawRef[File]("$(opctl_file_1)"))
		return base
	}, nil)
	if err == nil {
		t.Fatal("derived RawRef mutated original binding")
	}
}
