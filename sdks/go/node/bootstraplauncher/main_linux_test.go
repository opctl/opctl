package main

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractAuth(t *testing.T) {
	for _, tc := range []struct {
		name, file string
		kind       byte
		want       string
	}{
		{"regular", "helper", tar.TypeReg, ""},
		{"traversal", "../helper", tar.TypeReg, "Invalid Go bootstrap auth file"},
		{"unknown", "archive", tar.TypeReg, "Invalid Go bootstrap auth file"},
		{"symlink", "helper", tar.TypeSymlink, "Go bootstrap auth files must be regular"},
		{"hardlink", "helper", tar.TypeLink, "Go bootstrap auth files must be regular"},
		{"directory", "helper", tar.TypeDir, "Go bootstrap auth files must be regular"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var archive bytes.Buffer
			w := tar.NewWriter(&archive)
			h := &tar.Header{Name: tc.file, Mode: 0777, Typeflag: tc.kind}
			if tc.kind == tar.TypeReg {
				h.Size = 6
			}
			if err := w.WriteHeader(h); err != nil {
				t.Fatal(err)
			}
			if tc.kind == tar.TypeReg {
				w.Write([]byte("secret"))
			}
			w.Close()
			dir := t.TempDir()
			err := extractAuth(&archive, dir, os.Getuid(), os.Getgid())
			if tc.want != "" {
				if err == nil || err.Error() != tc.want {
					t.Fatalf("got %v, want %s", err, tc.want)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(filepath.Join(dir, tc.file))
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(dir, tc.file))
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != 0600 || string(data) != "secret" {
				t.Fatal("unsafe extraction")
			}
		})
	}
	for _, tc := range []struct {
		data []byte
		want string
	}{
		{[]byte("invalid archive"), "Invalid Go bootstrap auth archive"},
		{make([]byte, (1<<20)+1), "Go bootstrap private input exceeds 1 MiB"},
	} {
		if err := extractAuth(bytes.NewReader(tc.data), t.TempDir(), os.Getuid(), os.Getgid()); err == nil || err.Error() != tc.want {
			t.Fatalf("got %v, want %s", err, tc.want)
		}
	}
}

func TestExtractAuthRejectsDuplicatesAndTruncation(t *testing.T) {
	var archive bytes.Buffer
	w := tar.NewWriter(&archive)
	for i := 0; i < 2; i++ {
		w.WriteHeader(&tar.Header{Name: "helper", Mode: 0600, Size: 6})
		w.Write([]byte("secret"))
	}
	w.Close()
	for _, data := range [][]byte{archive.Bytes(), archive.Bytes()[:515]} {
		if err := extractAuth(bytes.NewReader(data), t.TempDir(), os.Getuid(), os.Getgid()); err == nil {
			t.Fatal("invalid archive accepted")
		}
	}
}

func TestDiagnosticCaptureIsBounded(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "diagnostic"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	capture := &cappedFile{file: file, cancel: cancel}
	if _, err = io.Copy(capture, strings.NewReader(strings.Repeat("x", 1<<20))); err == nil {
		t.Fatal("overflow accepted")
	}
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 65537 || ctx.Err() == nil {
		t.Fatal("diagnostics exceeded the tmpfs cap or child was not cancelled")
	}
}

func TestAuthRequiresTmpfs(t *testing.T) {
	if err := requireTmpfs(t.TempDir()); err == nil || err.Error() != "Go bootstrap auth requires tmpfs" {
		t.Fatalf("got %v", err)
	}
}
