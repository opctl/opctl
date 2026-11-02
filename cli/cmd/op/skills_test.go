package op

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConversionSkillsAreIdenticalRealCopies(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	var canonical []byte
	for _, agent := range []string{".agents", ".claude"} {
		dir := filepath.Join(root, agent, "skills", "converting-opctl-ops")
		for _, path := range []string{dir, filepath.Join(dir, "SKILL.md")} {
			info, err := os.Lstat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode()&os.ModeSymlink != 0 {
				t.Fatalf("skill must be a real copy: %s", path)
			}
		}
		data, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Fatal("empty skill")
		}
		if canonical == nil {
			canonical = data
		} else if !bytes.Equal(canonical, data) {
			t.Fatal(".agents and .claude conversion skills have drifted")
		}
	}
}
