package goop

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// Toolchain is the immutable image identity and explicit Linux target used to build.
type Toolchain struct{ ID, OS, Arch, Variant string }

// CacheFormatVersion identifies the on-disk metadata format, not the build recipe.
const CacheFormatVersion = "v1"

// BuildRecipe contains the launcher and the exact command/environment sent to the builder.
type BuildRecipe struct {
	Launcher []byte
	Command  []interface{}
	Env      map[string]interface{}
}

func record(h hash.Hash, fields ...string) {
	var length [8]byte
	for _, field := range fields {
		binary.BigEndian.PutUint64(length[:], uint64(len(field)))
		h.Write(length[:])
		h.Write([]byte(field))
	}
}

// Fingerprint streams source content; it never reads the mutable Go caches or input values.
func (source *Source) Fingerprint(ctx context.Context, excluded []string, image Toolchain, recipe BuildRecipe) (string, error) {
	roots := []string{}
	for _, root := range source.Roots {
		resolved, err := filepath.EvalSymlinks(root)
		if err != nil {
			return "", err
		}
		roots = append(roots, resolved)
	}
	resolvedExclusions := make([]string, 0, len(excluded))
	for _, path := range excluded {
		resolved, err := canonicalPath(path)
		if err != nil {
			return "", err
		}
		resolvedExclusions = append(resolvedExclusions, resolved)
		for _, replace := range source.Replacements {
			if within(resolved, replace.Root) {
				return "", fmt.Errorf("Go op replacement %q uses excluded node cache %q", replace.Root, resolved)
			}
		}
	}
	walker := sourceHasher{ctx: ctx, roots: roots, excluded: resolvedExclusions, active: map[string]bool{}, done: map[string]string{}}
	main, err := walker.digest(source.ModuleRoot)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	recipeBytes, err := json.Marshal(recipe)
	if err != nil {
		return "", err
	}
	record(h, string(recipeBytes), main, source.Package, source.Name, source.GoVersion, image.ID, image.OS, image.Arch, image.Variant)
	if _, err := os.Stat(filepath.Join(source.ModuleRoot, "go.sum")); os.IsNotExist(err) {
		record(h, "missing go.sum")
	} else if err != nil {
		return "", err
	} else {
		record(h, "present go.sum")
	}
	for _, replace := range source.Replacements {
		digest := "missing"
		if !replace.Missing {
			digest, err = walker.digest(replace.Root)
			if err != nil {
				return "", err
			}
		}
		record(h, replace.OldPath, replace.OldVersion, replace.NewPath, digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

type sourceHasher struct {
	ctx      context.Context
	roots    []string
	excluded []string
	active   map[string]bool
	done     map[string]string
}

func (walker *sourceHasher) digest(path string) (string, error) {
	if err := walker.ctx.Err(); err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	for _, excluded := range walker.excluded {
		if within(excluded, path) {
			return "", fmt.Errorf("Go op source %q enters excluded node cache", path)
		}
	}
	if walker.active[path] {
		return "", fmt.Errorf("Go op source symlink cycle at %q", path)
	}
	if digest, ok := walker.done[path]; ok {
		return digest, nil
	}
	walker.active[path] = true
	defer delete(walker.active, path)
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		text, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			return "", fmt.Errorf("Go op source symlink %q: %w", path, err)
		}
		allowed := false
		for _, root := range walker.roots {
			if within(root, target) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("Go op source symlink %q escapes declared source roots", path)
		}
		digest, err := walker.digest(target)
		if err != nil {
			return "", err
		}
		record(h, "symlink", text, digest)
	case info.IsDir():
		record(h, "directory")
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
	children:
		for _, entry := range entries {
			if entry.Name() == ".git" || entry.Name() == ".hg" || entry.Name() == ".svn" {
				continue
			}
			child := filepath.Join(path, entry.Name())
			for _, excluded := range walker.excluded {
				if within(path, excluded) && within(excluded, child) {
					continue children
				}
			}
			digest, err := walker.digest(child)
			if err != nil {
				return "", err
			}
			record(h, entry.Name(), digest)
		}
	case info.Mode().IsRegular():
		record(h, "file", fmt.Sprint(info.Mode().Perm()&0111), fmt.Sprint(info.Size()))
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		_, copyErr := io.Copy(h, file)
		closeErr := file.Close()
		if copyErr != nil {
			return "", copyErr
		}
		if closeErr != nil {
			return "", closeErr
		}
	default:
		return "", fmt.Errorf("unsupported Go op source file %q (%s)", path, info.Mode().Type())
	}
	result := hex.EncodeToString(h.Sum(nil))
	walker.done[path] = result
	return result, nil
}
