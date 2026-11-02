package goconvert

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
	"strings"
)

const sdkModule = "github.com/opctl/opctl/sdks/go/opspec/sdk"

// StageFiles prepares a sibling Go op that the invoking user's node can read.
// The caller removes the returned directory after verification.
func StageFiles(a *Analysis, files *FileSet) (string, error) {
	stage, err := os.MkdirTemp(filepath.Dir(a.Dir), "."+filepath.Base(a.Dir)+".opctl-convert-")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(stage)
		}
	}()
	if err = unsudo.Chown(stage); err != nil {
		return "", err
	}
	staged, err := Render(a, RenderOptions{ModuleRoot: files.ModuleRoot, ModulePath: files.ModulePath, StagingDir: stage})
	if err != nil {
		return "", err
	}
	for name, content := range files.Files {
		if name == "op.go" {
			content = staged.Source
		}
		if err = writeNewFile(filepath.Join(stage, name), content); err != nil {
			return "", err
		}
	}
	// Existing per-op modules keep their replace directives at the same depth.
	if !files.NewModule && files.ModuleRoot == a.Dir {
		for _, name := range []string{"go.mod", "go.sum"} {
			var content []byte
			content, err = os.ReadFile(filepath.Join(a.Dir, name))
			if os.IsNotExist(err) {
				err = nil
				continue
			}
			if err != nil {
				return "", err
			}
			if err = writeNewFile(filepath.Join(stage, name), content); err != nil {
				return "", err
			}
		}
	}
	return stage, nil
}

// FileSet identifies the module and files needed to replace one op file.
type FileSet struct {
	Files                  map[string][]byte
	ModuleRoot, ModulePath string
	NewModule              bool
}

// PrepareFiles applies the module rule without changing the filesystem.
func PrepareFiles(dir string, sdk opspec.GoSDK) (*FileSet, error) {
	root, err := goop.ModuleRoot(dir)
	if err == nil {
		data, e := os.ReadFile(filepath.Join(root, "go.mod"))
		if e != nil {
			return nil, e
		}
		mod, e := modfile.Parse("go.mod", data, nil)
		if e != nil {
			return nil, e
		}
		for _, req := range mod.Require {
			if req.Mod.Path == sdkModule {
				if mod.Module == nil {
					return nil, fmt.Errorf("go.mod needs a module directive")
				}
				return &FileSet{Files: map[string][]byte{}, ModuleRoot: root, ModulePath: mod.Module.Mod.Path}, nil
			}
		}
		if root == dir {
			return nil, fmt.Errorf("op dir has a go.mod without the Go op SDK; add require %s or remove it", sdkModule)
		}
	} else if !strings.Contains(err.Error(), "needs a go.mod") {
		return nil, err
	}
	scaffold, err := opspec.GoOpFiles(filepath.Base(dir), "", sdk)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, file := range scaffold {
		if file.Name != "op.go" {
			files[file.Name] = []byte(file.Content)
		}
	}
	ignore := filepath.Join(dir, ".gitignore")
	if info, e := os.Lstat(ignore); e == nil {
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf(".gitignore must be a regular file")
		}
		existing, e := os.ReadFile(ignore)
		if e != nil {
			return nil, e
		}
		found := false
		for _, line := range strings.Split(string(existing), "\n") {
			if line == "/opctlgen/" {
				found = true
			}
		}
		if found {
			delete(files, ".gitignore")
		} else {
			if len(existing) > 0 && existing[len(existing)-1] != '\n' {
				existing = append(existing, '\n')
			}
			files[".gitignore"] = append(existing, []byte("/opctlgen/\n")...)
		}
	} else if !os.IsNotExist(e) {
		return nil, e
	}
	return &FileSet{Files: files, ModuleRoot: dir, ModulePath: "ops.local/" + filepath.Base(dir), NewModule: true}, nil
}

// WriteFiles writes module files and atomically installs op.go before removing op.yml.
func WriteFiles(dir string, files map[string][]byte) error { return writeFiles(dir, files, nil) }
func writeFiles(dir string, files map[string][]byte, beforeRemove func() error) error {
	for _, name := range []string{"op.go", "op.go.tmp"} {
		if _, err := os.Lstat(filepath.Join(dir, name)); err == nil {
			return fmt.Errorf("%s already exists", name)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	createdModule := false
	for _, name := range []string{"go.mod", "go.sum", ".gitignore"} {
		data, ok := files[name]
		if !ok {
			continue
		}
		path := filepath.Join(dir, name)
		if name == ".gitignore" {
			if info, err := os.Lstat(path); err == nil {
				if !info.Mode().IsRegular() {
					return fmt.Errorf(".gitignore must be a regular file")
				}
				existing, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if !strings.HasPrefix(string(data), string(existing)) {
					return fmt.Errorf(".gitignore changed during conversion")
				}
				file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
				if err != nil {
					return err
				}
				_, writeErr := file.Write(data[len(existing):])
				closeErr := file.Close()
				if writeErr != nil {
					return writeErr
				}
				if closeErr != nil {
					return closeErr
				}
				continue
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		if err := writeNewFile(path, data); err != nil {
			if name == "go.sum" && createdModule {
				if removeErr := os.Remove(filepath.Join(dir, "go.mod")); removeErr != nil {
					return fmt.Errorf("%w; remove incomplete go.mod: %v", err, removeErr)
				}
			}
			return err
		}
		createdModule = createdModule || name == "go.mod"
	}
	tmp := filepath.Join(dir, "op.go.tmp")
	if err := writeNewFile(tmp, files["op.go"]); err != nil {
		return err
	}
	if _, err := os.Lstat(filepath.Join(dir, "op.go")); err == nil {
		return fmt.Errorf("op.go appeared during conversion")
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(tmp, filepath.Join(dir, "op.go")); err != nil {
		return err
	}
	if beforeRemove != nil {
		if err := beforeRemove(); err != nil {
			return err
		}
	}
	return os.Remove(filepath.Join(dir, "op.yml"))
}
func writeNewFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(data)
	syncErr := f.Sync()
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	if syncErr != nil {
		return syncErr
	}
	if closeErr != nil {
		return closeErr
	}
	return unsudo.Chown(path)
}
