package node

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/data"
	datafs "github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/data/git"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

type goStubChainKey struct{}
type goStubSet struct {
	uses  []goop.StubUse
	files map[string][]byte
}

func (bootstrap _bootstrapper) stubs(ctx context.Context, op *model.OpCall, root string, source *goop.Source) (goStubSet, error) {
	set := goStubSet{files: map[string][]byte{}}
	path, err := filepath.EvalSymlinks(op.OpPath)
	if err != nil {
		return set, err
	}
	chain, _ := ctx.Value(goStubChainKey{}).([]string)
	for _, previous := range chain {
		if previous == path {
			return set, fmt.Errorf("op stubs form a cycle: %s. Call one of these ops with opctl.CallOp instead of its stub", strings.Join(append(append([]string{}, chain...), path), " -> "))
		}
	}
	ctx = context.WithValue(ctx, goStubChainKey{}, append(append([]string{}, chain...), path))
	set.uses, err = goop.ScanStubs(source)
	if err != nil {
		return set, err
	}
	for _, use := range set.uses {
		if err := ctx.Err(); err != nil {
			return set, err
		}
		ref := use.Ref
		var creds *model.Creds
		var dir model.DataHandle
		if use.Local {
			ref = filepath.Join(source.ModuleRoot, filepath.FromSlash(use.ModuleRef))
			dir, err = data.Resolve(ctx, ref, datafs.New())
			if err != nil {
				return set, fmt.Errorf("%s: op stub %s: no op at %s (%s): %w", use.Site, use.Name, use.Ref, ref, err)
			}
		} else {
			if bootstrap.cache.stateStore != nil {
				if auth := bootstrap.cache.stateStore.TryGetAuth(ref); auth != nil {
					creds = &auth.Creds
				}
			}
			dir, err = data.Resolve(ctx, ref, git.New(filepath.Join(bootstrap.dataDirPath, "ops"), creds))
			if err != nil {
				return set, fmt.Errorf("%s: op stub %s: fetch %s: %w; add credentials with opctl auth add", use.Site, use.Name, ref, err)
			}
		}
		// Filesystem handles do not carry the fetched caller's repository bound.
		if use.Local && op.RepositoryRoot != "" {
			creds = op.BootstrapCreds
			dir = stubRepositoryDir{DataHandle: dir, repository: op.RepositoryRoot, source: op.BootstrapSource}
		}
		spec, err := describeOpDir(ctx, dir, creds, op.OpID, root, bootstrap, filepath.Join(bootstrap.dataDirPath, "ops"))
		if err != nil {
			return set, fmt.Errorf("%s: op stub %s: describe %s: %w", use.Site, use.Name, use.Ref, err)
		}
		set.files[use.Name], err = goop.RenderStub(use, spec)
		if err != nil {
			return set, fmt.Errorf("%s: %w", use.Site, err)
		}
	}
	return set, nil
}

type stubRepositoryDir struct {
	model.DataHandle
	repository, source string
}

func (d stubRepositoryDir) RepositoryRoot() string { return d.repository }
func (d stubRepositoryDir) SourceResource() string { return d.source }

func (set goStubSet) digests() map[string]string {
	if len(set.files) == 0 {
		return nil
	}
	digests := map[string]string{}
	for name, data := range set.files {
		sum := sha256.Sum256(data)
		digests["opctlgen/"+name+"/stub.go"] = hex.EncodeToString(sum[:])
	}
	return digests
}
func (set goStubSet) materialize(dataDir, moduleRoot string) (string, error) {
	parent := filepath.Join(dataDir, "dcg")
	if err := unsudo.CreateDir(parent); err != nil {
		return "", err
	}
	path, err := os.MkdirTemp(parent, "go-stubs-")
	if err != nil {
		return "", err
	}
	keep := false
	defer func() {
		if !keep {
			os.RemoveAll(path)
		}
	}()
	if err := os.Chmod(path, 0755); err != nil {
		return "", err
	}
	replacements := map[string]string{}
	for name, data := range set.files {
		stubDir := filepath.Join(moduleRoot, "opctlgen", name)
		err := filepath.WalkDir(stubDir, func(file string, entry fs.DirEntry, err error) error {
			if os.IsNotExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				replacements[file] = ""
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		replacements[filepath.Join(stubDir, "stub.go")] = "/opctl/stubs/" + name + ".go"
		if err := os.WriteFile(filepath.Join(path, name+".go"), data, 0644); err != nil {
			return "", err
		}
	}
	data, err := json.Marshal(struct{ Replace map[string]string }{replacements})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(path, "overlay.json"), data, 0644); err != nil {
		return "", err
	}
	keep = true
	return path, nil
}
func (set goStubSet) hint(moduleRoot, diagnostic string) string {
	names := []string{}
	for _, use := range set.uses {
		name := use.Name + " (" + use.Ref + ")"
		if host, err := os.ReadFile(filepath.Join(moduleRoot, "opctlgen", use.Name, "stub.go")); err == nil && !bytes.Equal(host, set.files[use.Name]) {
			name += " [opctlgen/" + use.Name + " is out of date]"
		}
		names = append(names, name)
	}
	hint := "op stubs were generated from the current interfaces of: " + strings.Join(names, ", ") + "."
	if strings.Contains(filepath.ToSlash(diagnostic), "opctlgen/") {
		hint += " The op's code uses inputs or outputs these ops no longer declare."
	}
	hint += " Run \"opctl op gen\" to update opctlgen/ for your editor."
	if strings.Contains(diagnostic, "undefined: opctl.StubCall") {
		hint += " Op stubs need Go op SDK >= this opctl release's SDK; update the SDK requirement in go.mod."
	}
	return hint
}
