package goop

import (
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type stubImport struct{ name, site string }
type stubPackage struct {
	imports []string
	stubs   []stubImport
	tests   []stubImport
}
type stubScanner struct {
	source       *Source
	module       string
	declarations map[string]StubUse
	packages     map[string]*stubPackage
}

// ScanStubs returns the stubs reachable through non-test imports of the op's
// main module. Declarations may live anywhere in that module, including tests.
func ScanStubs(source *Source) ([]StubUse, error) {
	scan, err := scanStubModule(source)
	if err != nil {
		return nil, err
	}
	needed := []stubImport{}
	visited := map[string]bool{}
	var visit func(string)
	visit = func(dir string) {
		if visited[dir] {
			return
		}
		visited[dir] = true
		pkg := scan.packages[dir]
		if pkg == nil {
			return
		}
		needed = append(needed, pkg.stubs...)
		for _, path := range pkg.imports {
			if strings.HasPrefix(path, scan.module+"/") {
				visit(filepath.Join(source.ModuleRoot, filepath.FromSlash(strings.TrimPrefix(path, scan.module+"/"))))
			} else if path == scan.module {
				visit(source.ModuleRoot)
			}
		}
	}
	visit(source.OpDir)
	return scan.uses(needed)
}

// ScanModuleStubs returns every declared or imported op stub, including tests,
// without entering nested modules or generated packages.
func ScanModuleStubs(moduleRoot string) ([]StubUse, error) {
	scan, err := scanStubModule(&Source{ModuleRoot: moduleRoot})
	if err != nil {
		return nil, err
	}
	needed := []stubImport{}
	for _, dir := range sortedStubKeys(scan.packages) {
		pkg := scan.packages[dir]
		needed = append(needed, pkg.stubs...)
		needed = append(needed, pkg.tests...)
	}
	return scan.uses(needed)
}
func (scan *stubScanner) uses(imports []stubImport) ([]StubUse, error) {
	uses := map[string]StubUse{}
	for _, imported := range imports {
		use, ok := scan.declarations[imported.name]
		if !ok {
			return nil, fmt.Errorf("%s: no ref is declared for op stub %q; add //opctl:op <ref> after the import, e.g. //opctl:op ../%s", imported.site, imported.name, imported.name)
		}
		uses[use.Name] = use
	}
	result := make([]StubUse, 0, len(uses))
	for _, name := range sortedStubKeys(uses) {
		result = append(result, uses[name])
	}
	return result, nil
}
func scanStubModule(source *Source) (*stubScanner, error) {
	module, err := readModule(source.ModuleRoot)
	if err != nil {
		return nil, err
	}
	if module.Module == nil {
		return nil, fmt.Errorf("%s/go.mod: needs a module directive", source.ModuleRoot)
	}
	scan := &stubScanner{source: source, module: module.Module.Mod.Path, declarations: map[string]StubUse{}, packages: map[string]*stubPackage{}}
	err = filepath.WalkDir(source.ModuleRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != source.ModuleRoot {
				switch entry.Name() {
				case ".git", ".hg", ".svn", "vendor", "node_modules", "opctlgen":
					return filepath.SkipDir
				}
				if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
					return filepath.SkipDir
				} else if !os.IsNotExist(err) {
					return err
				}
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		return scan.file(path)
	})
	return scan, err
}
func (scan *stubScanner) file(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, data, parser.ImportsOnly|parser.ParseComments)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	pkg := scan.packages[dir]
	if pkg == nil {
		pkg = &stubPackage{}
		scan.packages[dir] = pkg
	}
	allowed := map[token.Pos]bool{}
	for _, imported := range file.Imports {
		value, err := strconv.Unquote(imported.Path.Value)
		if err != nil {
			return err
		}
		site := fset.Position(imported.Pos()).String()
		prefix := scan.module + "/opctlgen/"
		isStub := strings.HasPrefix(value, prefix) || value == scan.module+"/opctlgen"
		if !isStub {
			if !strings.HasSuffix(path, "_test.go") {
				pkg.imports = append(pkg.imports, value)
			}
			continue
		}
		name := strings.TrimPrefix(value, prefix)
		if strings.Contains(name, "/") {
			return fmt.Errorf("%s: op stub paths must be one level below opctlgen/", site)
		}
		if !validStubName(name) {
			return fmt.Errorf("%s: op stub name %q must be a lowercase Go package name like mysql", site, name)
		}
		use := stubImport{name: name, site: site}
		if strings.HasSuffix(path, "_test.go") {
			pkg.tests = append(pkg.tests, use)
		} else {
			pkg.stubs = append(pkg.stubs, use)
		}
		if imported.Comment == nil {
			continue
		}
		for _, comment := range imported.Comment.List {
			if !strings.HasPrefix(comment.Text, "//opctl:op") {
				continue
			}
			allowed[comment.Pos()] = true
			tokens := strings.Fields(strings.TrimPrefix(comment.Text, "//opctl:op"))
			if !strings.HasPrefix(comment.Text, "//opctl:op ") || len(tokens) != 1 {
				return fmt.Errorf("%s: //opctl:op needs one ref after it, like //opctl:op ../greet", site)
			}
			declaration, err := scan.declaration(name, tokens[0], path, site)
			if err != nil {
				return err
			}
			if previous, ok := scan.declarations[name]; ok && (previous.Local != declaration.Local || previous.ModuleRef != declaration.ModuleRef || !declaration.Local && previous.Ref != declaration.Ref) {
				return fmt.Errorf("op stub %s is declared with different refs: %s (%s) and %s (%s); rename one import", name, previous.Site, previous.Ref, site, declaration.Ref)
			}
			if _, ok := scan.declarations[name]; !ok {
				scan.declarations[name] = declaration
			}
		}
	}
	// Scan the rest of the file too, to catch misplaced declarations in bodies.
	var comments scanner.Scanner
	comments.Init(fset.File(file.Pos()), data, nil, scanner.ScanComments)
	for {
		pos, tok, text := comments.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.COMMENT || !strings.Contains(text, "opctl:op") {
			continue
		}
		site := fset.Position(pos).String()
		if strings.HasPrefix(text, "// opctl:op") {
			return fmt.Errorf("%s: write //opctl:op without a space", site)
		}
		if !allowed[pos] {
			return fmt.Errorf("%s: put //opctl:op <ref> after an op stub import", site)
		}
	}
	sort.Strings(pkg.imports)
	return nil
}
func (scan *stubScanner) declaration(name, ref, path, site string) (StubUse, error) {
	use := StubUse{Name: name, Ref: ref, Site: site}
	if strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../") {
		target := filepath.Join(filepath.Dir(path), filepath.FromSlash(ref))
		if scan.source.RepositoryRoot != "" {
			if err := inRepository(scan.source.RepositoryRoot, target); err != nil {
				return use, fmt.Errorf("%s: op stub %s: %s: %w", site, name, ref, err)
			}
		}
		target, err := canonicalPath(target)
		if err != nil {
			return use, fmt.Errorf("%s: op stub %s: %w", site, name, err)
		}
		rel, err := filepath.Rel(scan.source.ModuleRoot, target)
		if err != nil {
			return use, err
		}
		use.Local = true
		use.ModuleRef = filepath.ToSlash(rel)
		return use, nil
	}
	if !filepath.IsAbs(ref) && strings.Contains(ref, "#") {
		return use, nil
	}
	return use, fmt.Errorf("%s: op stub %s: use ./ or ../ for a local op, or repo#tag/path for a git op", site, name)
}
