package goop

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func generatedStub(data []byte) bool {
	line, _, _ := bytes.Cut(data, []byte("\n"))
	return bytes.HasPrefix(line, []byte(GeneratedStubHeader)) && bytes.HasSuffix(line, []byte(". DO NOT EDIT."))
}

// WriteStubs atomically refreshes opctlgen and removes stale generated packages.
// A nil file preserves a declared stub whose describe failed. User files and
// symlinks are reported and left untouched, including during stale cleanup.
func WriteStubs(moduleRoot string, files map[string][]byte) error {
	root := filepath.Join(moduleRoot, "opctlgen")
	if info, err := os.Lstat(root); err == nil {
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s must be a directory, not a file or symlink", root)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	ignore, err := os.OpenFile(filepath.Join(root, ".gitignore"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err == nil {
		_, writeErr := ignore.WriteString("*\n")
		closeErr := ignore.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return err
		}
	} else if !os.IsExist(err) {
		return err
	}
	var failures []error
	for _, name := range sortedStubKeys(files) {
		if !validStubName(name) {
			failures = append(failures, fmt.Errorf("invalid op stub name %q", name))
			continue
		}
		if files[name] == nil {
			continue
		}
		if err := writeStub(root, name, files[name]); err != nil {
			failures = append(failures, err)
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return errors.Join(append(failures, err)...)
	}
	for _, entry := range entries {
		if entry.Name() == ".gitignore" {
			continue
		}
		if _, needed := files[entry.Name()]; needed {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			failures = append(failures, fmt.Errorf("leaving user file %s", path))
			continue
		}
		if err := filepath.WalkDir(path, func(file string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("leaving symlink %s", file)
			}
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if !generatedStub(data) {
				return fmt.Errorf("leaving user file %s", file)
			}
			return nil
		}); err != nil {
			failures = append(failures, err)
			continue
		}
		if err := os.RemoveAll(path); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}
func writeStub(root, name string, data []byte) error {
	dir := filepath.Join(root, name)
	if info, err := os.Lstat(dir); err == nil {
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing op stub directory %s: not a directory or is a symlink", dir)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "stub.go")
	if info, err := os.Lstat(path); err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("refusing to overwrite user file %s", path)
		}
		existing, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !generatedStub(existing) {
			return fmt.Errorf("refusing to overwrite user file %s", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if !generatedStub(data) {
		return fmt.Errorf("op stub %s is missing its generated header", name)
	}
	temporary, err := os.CreateTemp(dir, ".stub-")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())
	_, writeErr := temporary.Write(data)
	chmodErr := temporary.Chmod(0644)
	closeErr := temporary.Close()
	if err := errors.Join(writeErr, chmodErr, closeErr); err != nil {
		return err
	}
	return os.Rename(temporary.Name(), path)
}

// StubModuleRef returns the absolute local target or the unchanged git ref.
func StubModuleRef(moduleRoot string, use StubUse) string {
	if use.Local {
		return filepath.Join(moduleRoot, filepath.FromSlash(use.ModuleRef))
	}
	return strings.TrimSpace(use.Ref)
}
