package goop

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSourceFingerprint(t *testing.T) {
	root := t.TempDir()
	write := func(path, text string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	app := filepath.Join(root, "app")
	sdk := filepath.Join(root, "sdk")
	data := filepath.Join(app, "node")
	write(filepath.Join(app, "go.mod"), "module app\ngo 1.21\nreplace example.com/sdk => ../sdk\n")
	write(filepath.Join(app, "ops", "first", "op.go"), "package main\nfunc main(){}\n")
	write(filepath.Join(app, "ops", "second", "op.go"), "package main\nfunc main(){}\n")
	write(filepath.Join(sdk, "go.mod"), "module example.com/sdk\ngo 1.21\n")
	write(filepath.Join(app, "asset.txt"), "asset")
	image := Toolchain{ID: "sha256:a", OS: "linux", Arch: "amd64"}
	fingerprint := func(op string, toolchain Toolchain) string {
		t.Helper()
		source, err := Discover(op, "")
		if err != nil {
			t.Fatal(err)
		}
		key, err := source.Fingerprint(context.Background(), data, toolchain)
		if err != nil {
			t.Fatal(err)
		}
		return key
	}
	op := filepath.Join(app, "ops", "first")
	first := fingerprint(op, image)
	if first == fingerprint(filepath.Join(app, "ops", "second"), image) {
		t.Fatal("different op packages shared a key")
	}
	if err := os.Chtimes(filepath.Join(app, "asset.txt"), time.Now(), time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if first != fingerprint(op, image) {
		t.Fatal("mtime invalidated content key")
	}
	write(filepath.Join(app, ".git", "ignored"), "git metadata")
	write(filepath.Join(data, "programs", "ignored"), "node cache")
	if first != fingerprint(op, image) {
		t.Fatal("excluded metadata invalidated key")
	}
	for _, tc := range []struct{ path, text string }{
		{filepath.Join(app, "asset.txt"), "changed asset"},
		{filepath.Join(sdk, "go.mod"), "module example.com/sdk\ngo 1.21\n// SDK changed\n"},
		{filepath.Join(app, "untracked.go"), "package app\n"},
		{filepath.Join(app, "go.sum"), "new checksum"},
	} {
		before := fingerprint(op, image)
		write(tc.path, tc.text)
		if before == fingerprint(op, image) {
			t.Fatalf("edit ignored: %s", tc.path)
		}
	}
	before := fingerprint(op, image)
	if err := os.Remove(filepath.Join(app, "go.sum")); err != nil {
		t.Fatal(err)
	}
	if before == fingerprint(op, image) {
		t.Fatal("deleted go.sum ignored")
	}
	before = fingerprint(op, image)
	image.ID = "sha256:b"
	if before == fingerprint(op, image) {
		t.Fatal("image identity ignored")
	}
	before = fingerprint(op, image)
	image.Arch = "arm64"
	if before == fingerprint(op, image) {
		t.Fatal("target ignored")
	}
	// An in-root link includes its referent; escaping links and cycles fail.
	if err := os.Symlink("asset.txt", filepath.Join(app, "linked")); err != nil {
		t.Fatal(err)
	}
	before = fingerprint(op, image)
	write(filepath.Join(app, "asset.txt"), "linked edit")
	if before == fingerprint(op, image) {
		t.Fatal("referent ignored")
	}
	if err := os.Symlink(root, filepath.Join(app, "escape")); err != nil {
		t.Fatal(err)
	}
	source, err := Discover(op, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = source.Fingerprint(context.Background(), data, image); err == nil {
		t.Fatal("escaping symlink accepted")
	}
	if err := os.Remove(filepath.Join(app, "escape")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(".", filepath.Join(app, "cycle")); err != nil {
		t.Fatal(err)
	}
	if _, err = source.Fingerprint(context.Background(), data, image); err == nil {
		t.Fatal("source cycle accepted")
	}
}
