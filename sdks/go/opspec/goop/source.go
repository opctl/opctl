package goop

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Source describes the nearest module and the main module's effective local replacements.
type Source struct {
	RepositoryRoot                              string
	ModuleRoot, OpDir, Package, Name, GoVersion string
	Roots                                       []string
	Replacements                                []Replacement
}

type Replacement struct {
	OldPath, OldVersion, NewPath, Root string
	Missing                            bool
}

func Discover(opDir, repositoryRoot string) (*Source, error) {
	opDir, err := filepath.Abs(opDir)
	if err != nil {
		return nil, err
	}
	if strings.Contains(opDir, "$(") {
		return nil, fmt.Errorf("Go op paths containing $( are not supported yet")
	}
	root, err := moduleRoot(opDir, repositoryRoot)
	if err != nil {
		return nil, err
	}
	if err := validateMount(root); err != nil {
		return nil, err
	}
	module, err := readModule(root)
	if err != nil {
		return nil, err
	}
	if module.Go == nil {
		return nil, fmt.Errorf("Go op module needs a go directive")
	}
	relative, err := filepath.Rel(root, opDir)
	if err != nil {
		return nil, err
	}
	source := &Source{ModuleRoot: root, OpDir: opDir, Package: filepath.ToSlash(relative), Name: filepath.Base(opDir), GoVersion: module.Go.Version, Roots: []string{root}}
	for _, replace := range module.Replace {
		if replace.New.Version != "" {
			continue
		}
		target := replace.New.Path
		if !filepath.IsAbs(target) {
			target = filepath.Join(root, target)
		}
		target = filepath.Clean(target)
		if repositoryRoot != "" {
			if err := inRepository(repositoryRoot, target); err != nil {
				return nil, fmt.Errorf("Go op local replace %q must resolve within fetched repository root %q: %w", replace.New.Path, repositoryRoot, err)
			}
		}
		if !within(root, target) {
			if err := validateMountDestination(target); err != nil {
				return nil, err
			}
		}
		replacement := Replacement{OldPath: replace.Old.Path, OldVersion: replace.Old.Version, NewPath: replace.New.Path, Root: target}
		if _, err := os.Stat(target); os.IsNotExist(err) {
			replacement.Missing = true
		} else if err != nil {
			return nil, err
		} else {
			if !within(root, target) {
				if err := validateMount(target); err != nil {
					return nil, err
				}
			}
			source.Roots = append(source.Roots, target)
		}
		source.Replacements = append(source.Replacements, replacement)
	}
	sort.Slice(source.Replacements, func(i, j int) bool {
		a, b := source.Replacements[i], source.Replacements[j]
		if a.OldPath == b.OldPath {
			return a.OldVersion < b.OldVersion
		}
		return a.OldPath < b.OldPath
	})
	// One read-only mount covers all nested roots without mounting a common ancestor.
	sort.Slice(source.Roots, func(i, j int) bool {
		if len(source.Roots[i]) == len(source.Roots[j]) {
			return source.Roots[i] < source.Roots[j]
		}
		return len(source.Roots[i]) < len(source.Roots[j])
	})
	roots := []string{}
	for _, target := range source.Roots {
		covered := false
		for _, parent := range roots {
			if within(parent, target) {
				covered = true
				break
			}
		}
		if !covered {
			roots = append(roots, target)
		}
	}
	source.Roots = roots
	source.RepositoryRoot = repositoryRoot
	return source, nil
}

// canonicalPath resolves existing ancestors too, so a missing leaf cannot hide a symlink escape.
func canonicalPath(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	if info, statErr := os.Lstat(path); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", err
	}
	parent := filepath.Dir(path)
	if parent == path {
		return "", err
	}
	resolved, err = canonicalPath(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(resolved, filepath.Base(path)), nil
}
