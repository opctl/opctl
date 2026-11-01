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
	dir, err := filepath.Abs(opDir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(dir)
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
	if phase != "plan" && phase != "describe" {
		return nil, "", fmt.Errorf("unsupported bootstrap phase %q", phase)
	}
	opDir, err := filepath.Abs(opDir)
	if err != nil {
		return nil, "", err
	}
	if strings.Contains(opDir, "$(") {
		return nil, "", fmt.Errorf("Go op paths containing $( are not supported yet")
	}
	root, err := ModuleRoot(opDir)
	if err != nil {
		return nil, "", err
	}
	version, err := GoVersion(root)
	if err != nil {
		return nil, "", err
	}
	module, err := readModule(root)
	if err != nil {
		return nil, "", err
	}
	dirs := map[string]interface{}{root: "$(module)"}
	for _, replace := range module.Replace {
		target := replace.New.Path
		if !filepath.IsAbs(target) && !strings.HasPrefix(target, "./") && !strings.HasPrefix(target, "../") && target != "." && target != ".." {
			continue
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(root, target)
		}
		target = filepath.Clean(target)
		if target != root {
			dirs[target] = model.Value{Dir: &target}
		}
	}
	for destination := range dirs {
		for _, protected := range []string{"/usr", "/usr/local", "/go", "/opctl"} {
			if destination == "/" || destination == protected || strings.HasPrefix(protected, destination+"/") {
				return nil, "", fmt.Errorf("Go op mount destination %q would hide required bootstrap image directory %q", destination, protected)
			}
		}
	}
	return &model.ContainerCallSpec{
		Image:   &model.ContainerCallImageSpec{Ref: "golang:" + version},
		Cmd:     []interface{}{"go", "run", ".", phase},
		Dirs:    dirs,
		Files:   map[string]interface{}{"/opctl/request.json": "$(request)"},
		EnvVars: map[string]interface{}{"GOFLAGS": "-mod=readonly", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "OPCTL_BOOTSTRAP_REQUEST": "/opctl/request.json"},
		WorkDir: filepath.ToSlash(opDir),
	}, root, nil
}
