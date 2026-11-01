package node

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
)

//go:embed bootstraplauncher/main_linux.go
var goBootstrapLauncher []byte

func prepareGoCaches(spec *model.ContainerCallSpec, scope map[string]*model.Value, dataDir, scratch string) error {
	for destination, name := range map[string]string{"/opctl/mod": "mod", "/opctl/cache": "build", "/opctl/launcher-cache": "launcher-build"} {
		path := filepath.Join(dataDir, "go", name)
		if err := unsudo.CreateDir(path); err != nil {
			return err
		}
		spec.Dirs[destination] = model.Value{Dir: &path}
	}
	if err := prepareGoLauncher(spec, scope, scratch); err != nil {
		return err
	}
	return nil
}

func prepareGoLauncher(spec *model.ContainerCallSpec, scope map[string]*model.Value, scratch string) error {
	path := filepath.Join(scratch, "launcher.go")
	if err := unsudo.CreateFile(path, goBootstrapLauncher); err != nil {
		return err
	}
	if err := os.Chmod(path, 0600); err != nil {
		return err
	}
	scope["launcher"] = &model.Value{File: &path}
	spec.Files["/opctl/launcher.go"] = "$(launcher)"
	spec.ReadOnlyMounts["/opctl/launcher.go"] = true
	return nil
}

func goCacheEnvironment() map[string]string {
	return map[string]string{"GOMODCACHE": "/opctl/mod", "GOCACHE": "/opctl/cache", "GOWORK": "off", "GOENV": "off", "GOOS": "linux", "GOFLAGS": "-mod=readonly -buildvcs=false -trimpath"}
}
