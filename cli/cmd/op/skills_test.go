package op

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func TestConversionSkillDocumentsRefusalsAndRewrites(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	data, err := os.ReadFile(filepath.Join(root, ".agents", "skills", "converting-opctl-ops", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ReplaceAll(string(data), `\|`, "|")
	for _, want := range []string{"Rewrite instead of a split", "Each rewrite also works in op.yml. Make it there", "of the original op before and after should differ only by the new", "go.sum with a published SDK", "op refs can't contain whitespace or control characters", "op ref must be a literal local or git ref", "op ref is empty", "loop body reads and rebinds outer name", "verification coverage incomplete"} {
		if !strings.Contains(text, want) {
			t.Errorf("skill is missing %q", want)
		}
	}
	paths, err := filepath.Glob(filepath.Join(root, "sdks", "go", "opspec", "goconvert", "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 2 {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != "refuse" {
				return true
			}
			literal, ok := call.Args[1].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return true
			}
			reason, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(text, reason) {
				t.Errorf("skill is missing refusal %q", reason)
			}
			return true
		})
	}
}
