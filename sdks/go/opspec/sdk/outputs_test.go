package sdk

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"
)

func TestOutputVisibility(t *testing.T) {
	for _, tc := range []struct {
		name    string
		build   func(Dir, File) Call
		wantErr string
	}{
		{"serial produce consume", func(d Dir, f File) Call {
			return Serial(Container("producer").Dir("/out", d).File("/out.txt", f), Container("consumer").Dir("/in", d).File("/in.txt", f).EnvFrom("CONTENT", f))
		}, ""},
		{"consume after parallel", func(d Dir, f File) Call {
			return Serial(Parallel(Container("dir").Dir("/out", d), Container("file").File("/out", f)), Container("consumer").Dir("/in", d).File("/in.txt", f))
		}, ""},
		{"parallel sibling consume", func(d Dir, _ File) Call {
			return Parallel(Container("producer").Dir("/out", d), Container("consumer").Dir("/in", d))
		}, "parallel calls can't see each other's outputs"},
		{"two parallel producers", func(d Dir, _ File) Call {
			return Parallel(Container("a").Dir("/out", d), Container("b").Dir("/out", d))
		}, "parallel calls can't see each other's outputs"},
		{"two mounts in one call", func(d Dir, _ File) Call { return Container("a").Dir("/one", d).Dir("/two", d) }, "two mount paths"},
		{"env before producer", func(_ Dir, f File) Call {
			return Serial(Container("consumer").EnvFrom("CONTENT", f), Container("producer").File("/out", f))
		}, "used by EnvFrom"},
		{"env in producing call", func(_ Dir, f File) Call { return Container("a").File("/out", f).EnvFrom("CONTENT", f) }, "before any earlier call mounts it"},
		{"nested invisible producer", func(d Dir, _ File) Call {
			return Parallel(Serial(Container("a").Dir("/out", d)), Serial(Container("b"), Container("c").Dir("/in", d)))
		}, "run.parallel[1].serial[1] (c)"},
		{"shared parallel consumers", func(d Dir, _ File) Call {
			return Serial(Container("producer").Dir("/out", d), Parallel(Container("a").Dir("/in", d), Container("b").Dir("/in", d)))
		}, ""},
		{"visible mounts may alias", func(d Dir, _ File) Call {
			return Serial(Container("producer").Dir("/out", d), Container("consumer").Dir("/one", d).Dir("/two", d))
		}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := Plan(func(struct{}) Call { return tc.build(NewDir(), NewFile()) }, nil)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatal(err)
				}
				if !json.Valid(data) {
					t.Fatal("invalid JSON")
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) || !strings.Contains(err.Error(), "outputs_test.go:") {
				t.Fatalf("error = %v; want %s and creation site", err, tc.wantErr)
			}
		})
	}
}

func TestOutputDiagnosticNamesBothCalls(t *testing.T) {
	_, err := Plan(func(struct{}) Call {
		d := NewDir()
		return Parallel(Container("producer").Dir("/out", d), Container("consumer").Dir("/in", d))
	}, nil)
	if err == nil {
		t.Fatal("invisible output accepted")
	}
	for _, part := range []string{"dir created at", "run.parallel[0] (producer)", "run.parallel[1] (consumer)", "produce it in an earlier serial call"} {
		if !strings.Contains(err.Error(), part) {
			t.Fatalf("error = %v; missing %s", err, part)
		}
	}
}

func TestOutputAllocation(t *testing.T) {
	helper := func() Call { return Container("alpine").Dir("/out", NewDir()).File("/out.txt", NewFile()) }
	build := func(struct{}) Call { return Parallel(helper(), helper()) }
	first, err := Plan(build, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Plan(build, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("successive plans differ")
	}
	for _, name := range []string{"opctl_dir_1", "opctl_file_2", "opctl_dir_3", "opctl_file_4"} {
		if bytes.Count(first, []byte("$("+name+")")) != 1 {
			t.Fatalf("missing unique output %s", name)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := Plan(build, nil)
			if err != nil || !bytes.Equal(first, data) {
				t.Errorf("concurrent plan changed: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestOutputsMustBeCreatedDuringPlan(t *testing.T) {
	d, f := NewDir(), NewFile()
	for _, tc := range []struct {
		name string
		call Call
	}{
		{"NewDir", Container("a").Dir("/out", d)}, {"NewFile", Container("a").File("/out", f)},
	} {
		_, err := Plan(func(struct{}) Call { return tc.call }, nil)
		if err == nil || !strings.Contains(err.Error(), tc.name+" must be called while building the op graph") {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestUnusedOutputIsAllowed(t *testing.T) {
	_, err := Plan(func(struct{}) Call { NewDir(); return Container("a").Env("LITERAL", "$(opctl_dir_1)") }, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeneratedInputPrefixIsReserved(t *testing.T) {
	_, err := Describe[struct {
		Reserved String `opctl:"opctl_dir_1"`
	}]()
	if err == nil || !strings.Contains(err.Error(), "opctl_ is reserved") {
		t.Fatalf("error = %v", err)
	}
}
