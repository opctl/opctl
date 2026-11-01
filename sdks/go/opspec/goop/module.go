package goop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"golang.org/x/mod/modfile"
)

// ModuleRoot finds the nearest go.mod in the op dir or its parents.
func ModuleRoot(opDir string) (string, error) {
	return moduleRoot(opDir, "")
}

func moduleRoot(opDir, repositoryRoot string) (string, error) {
	dir, err := filepath.Abs(opDir)
	if err != nil {
		return "", err
	}
	if repositoryRoot != "" {
		if err := inRepository(repositoryRoot, dir); err != nil {
			return "", err
		}
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if repositoryRoot != "" {
				if err := inRepository(repositoryRoot, filepath.Join(dir, "go.mod")); err != nil {
					return "", err
				}
			}
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
		if repositoryRoot != "" && filepath.Clean(dir) == filepath.Clean(repositoryRoot) {
			return "", fmt.Errorf("Go op needs a go.mod within fetched repository root %q", repositoryRoot)
		}
		parent := filepath.Dir(dir)
		if repositoryRoot != "" {
			if err := inRepository(repositoryRoot, parent); err != nil {
				return "", fmt.Errorf("Go op needs a go.mod within fetched repository root %q", repositoryRoot)
			}
		}
		if parent == dir {
			return "", fmt.Errorf("Go op needs a go.mod in the op dir or a parent")
		}
		dir = parent
	}
}

func readModule(root string) (*modfile.File, error) {
	path := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return modfile.Parse(path, data, nil)
}

// GoVersion returns the module's Go directive verbatim.
func GoVersion(moduleRoot string) (string, error) {
	module, err := readModule(moduleRoot)
	if err != nil {
		return "", err
	}
	if module.Go == nil {
		return "", fmt.Errorf("Go op module needs a go directive")
	}
	return module.Go.Version, nil
}

// ContainerSpec prepares a bootstrap and returns the host dir to bind as
// the module scope value. Local replace targets retain their host paths.
func ContainerSpec(opDir, phase string) (*model.ContainerCallSpec, string, error) {
	return ContainerSpecInRepository(opDir, phase, "")
}

// ContainerSpecInRepository bounds git-fetched Go sources and local replacements
// to their resolved repository root. An empty root retains local-op behavior.
func ContainerSpecInRepository(opDir, phase, repositoryRoot string) (*model.ContainerCallSpec, string, error) {
	if phase != "plan" && phase != "describe" {
		return nil, "", fmt.Errorf("unsupported bootstrap phase %q", phase)
	}
	source, err := Discover(opDir, repositoryRoot)
	if err != nil {
		return nil, "", err
	}
	return source.ContainerSpec(phase), source.ModuleRoot, nil
}

// ContainerSpec binds this prepared source identity without rediscovering a changed module.
func (source *Source) ContainerSpec(phase string) *model.ContainerCallSpec {
	root, version := source.ModuleRoot, source.GoVersion
	dirs := map[string]interface{}{}
	for _, target := range source.Roots {
		if target == root {
			dirs[target] = "$(module)"
		} else {
			dirs[target] = model.Value{Dir: &target}
		}
	}

	readOnly := map[string]bool{"/opctl/request.json": true}
	for destination := range dirs {
		readOnly[destination] = true
	}
	return &model.ContainerCallSpec{
		ReadOnlyMounts: readOnly,
		Image:          &model.ContainerCallImageSpec{Ref: "golang:" + version},
		Cmd:            []interface{}{"go", "run", ".", phase},
		Dirs:           dirs,
		Files:          map[string]interface{}{"/opctl/request.json": "$(request)"},
		EnvVars:        map[string]interface{}{"GOFLAGS": "-mod=readonly", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "OPCTL_BOOTSTRAP_REQUEST": "/opctl/request.json"},
		WorkDir:        filepath.ToSlash(source.OpDir),
	}
}

func inRepository(root, path string) error {
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve fetched repository root %q: %w", root, err)
	}
	canonicalPath, err := canonicalPath(path)
	if err != nil {
		return fmt.Errorf("resolve path %q within fetched repository root %q: %w", path, root, err)
	}
	relative, err := filepath.Rel(canonicalRoot, canonicalPath)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("Go op path %q resolves outside fetched repository root %q", path, root)
	}
	return nil
}

func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func validateMountDestination(path string) error {
	for _, protected := range []string{"/usr", "/usr/local", "/go", "/opctl"} {
		if path == "/" || path == protected || strings.HasPrefix(protected, path+"/") || (protected == "/opctl" && within(protected, path)) {
			return fmt.Errorf("Go op mount destination %q would hide required bootstrap image directory %q", path, protected)
		}
	}
	return nil
}

func validateMount(path string) error {
	if err := validateMountDestination(path); err != nil {
		return err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("resolve Go op mount %q: %w", path, err)
	}
	if resolved != filepath.Clean(path) {
		return fmt.Errorf("Go op mount %q follows a host symlink to %q", path, resolved)
	}
	return nil
}
