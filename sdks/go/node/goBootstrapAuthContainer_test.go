package node

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	dockerRuntime "github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
)

type bootstrapTestOutput struct{ io.Writer }

func (bootstrapTestOutput) Close() error { return nil }

// Opt in with a socket and the same host-visible temp path mounted into the test container.
func TestBootstrapAuthContainer(t *testing.T) {
	temp := os.Getenv("OPCTL_TEST_HOST_TEMP")
	if temp == "" {
		t.Skip("requires opt-in Docker daemon and shared host temp path")
	}
	for _, imageTag := range []string{"golang:1.21", "golang:1.25"} {
		for _, outcome := range []string{"success", "download-failure", "diagnostic-overflow", "invalid-archive"} {
			failDownload := outcome != "success"
			t.Run(imageTag+"/"+outcome, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
				defer cancel()
				root, err := os.MkdirTemp(temp, "private-bootstrap-")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(root)
				if err := os.Chmod(root, 0755); err != nil {
					t.Fatal(err)
				}
				write := func(name, text string, mode os.FileMode) {
					t.Helper()
					if err := os.WriteFile(filepath.Join(root, name), []byte(text), mode); err != nil {
						t.Fatal(err)
					}
				}
				for _, name := range []string{"mod", "cache", "out", "src"} {
					if err := os.Mkdir(filepath.Join(root, name), 0755); err != nil {
						t.Fatal(err)
					}
					if err := os.Chown(filepath.Join(root, name), 1000, 1000); err != nil {
						t.Fatal(err)
					}
				}
				write("launcher", string(goBootstrapLauncher), 0644)
				if err := os.WriteFile(filepath.Join(root, "src", "go.mod"), []byte("module example.com/check\ngo 1.21\n"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "src", "main.go"), []byte("package main\nfunc main(){}\n"), 0644); err != nil {
					t.Fatal(err)
				}
				marker := "synthetic-container-secret"
				script := `#!/bin/sh
set -eu
case "$1" in
mod)
 [ "$(stat -f -c %T /run/opctl-auth)" = tmpfs ]
 [ "$(stat -c %a /run/opctl-auth/helper)" = 600 ]
 [ "$(stat -c %u /run/opctl-auth/helper)" = "$(id -u)" ]
 test "$HOME" = /run/opctl-auth
 result=$(printf 'protocol=https\nhost=example.com\npath=acme/repo\n\n' | git credential fill)
 case "$result" in *'password=synthetic-container-secret'*) ;; *) exit 21;; esac
 mkdir -p /opctl/mod/cache/vcs/test
 git init --bare -q /opctl/mod/cache/vcs/test
 git --git-dir=/opctl/mod/cache/vcs/test config remote.origin.url https://example.com/acme/repo.git
 head -c 262144 /dev/zero > /opctl/mod/downloaded-module.zip
 [ "$(wc -c < /opctl/mod/downloaded-module.zip)" -eq 262144 ]
 printf 'synthetic-container-secret' >&2
`
				if outcome == "diagnostic-overflow" {
					script += "head -c 65537 /dev/zero; head -c 65537 /dev/zero >&2; exit 0\n"
				} else if failDownload {
					script += "exit 19\n"
				} else {
					script += "exit 0\n"
				}
				script += `;;
build)
 [ -z "$(find /run/opctl-auth -mindepth 1 -print -quit)" ]
 [ -z "${GIT_CONFIG_GLOBAL:-}" ]
 [ "$HOME" != /run/opctl-auth ]
 [ "$GOPROXY" = off ]
 [ "$GONOPROXY" = none ]
 [ "${GOSUMDB:-}" != off ]
 ;;
esac
exec /usr/local/go/bin/go "$@"
`
				write("go", script, 0755)
				runtime, err := dockerRuntime.New(ctx, "unix:///var/run/docker.sock")
				if err != nil {
					t.Fatal(err)
				}
				identity, err := runtime.ResolveImage(ctx, &model.ContainerCallImage{Ref: &imageTag}, "", discardBootstrapEvents{})
				if err != nil {
					t.Fatal(err)
				}
				auth, err := prepareBootstrapAuth([]model.Auth{{Resources: "example.com/acme", Creds: model.Creds{Username: "synthetic", Password: marker}}}, "", nil)
				if err != nil {
					t.Fatal(err)
				}
				if outcome == "invalid-archive" {
					auth.archive = []byte("invalid archive")
				}
				ctx = containerruntime.WithBootstrapInput(ctx, bytes.NewReader(auth.archive))
				req := &model.ContainerCall{ContainerID: strings.ReplaceAll(filepath.Base(root), "-", ""), Image: &model.ContainerCallImage{Ref: &identity.Ref, Prepared: true}, WorkDir: "/src", Cmd: []string{"sh", "/opctl/launcher", "sh", "-ce", "go version > /opctl/build/go-version; go build -mod=readonly -buildvcs=false -trimpath -o /opctl/build/program ."},
					Dirs:           map[string]string{"/opctl/mod": filepath.Join(root, "mod"), "/opctl/cache": filepath.Join(root, "cache"), "/opctl/build": filepath.Join(root, "out"), "/src": filepath.Join(root, "src")},
					Files:          map[string]string{"/opctl/launcher": filepath.Join(root, "launcher"), "/usr/local/bin/go": filepath.Join(root, "go")},
					ReadOnlyMounts: map[string]bool{"/src": true, "/opctl/launcher": true, "/usr/local/bin/go": true},
					EnvVars:        map[string]string{"PATH": "/usr/local/bin:/usr/local/go/bin:/usr/bin:/bin", "GOMODCACHE": "/opctl/mod", "GOCACHE": "/opctl/cache", "GOPRIVATE": auth.private, "GOWORK": "off", "GOENV": "off", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0"}}
				var out, diagnostic bytes.Buffer
				exit, err := runtime.RunContainer(ctx, req, "", discardBootstrapEvents{}, bootstrapTestOutput{&out}, bootstrapTestOutput{&diagnostic})
				if err != nil {
					t.Fatal(err)
				}
				if exit == nil || (*exit == 0) == failDownload {
					t.Fatalf("unexpected container exit: %v", exit)
				}
				if bytes.Contains(out.Bytes(), []byte(marker)) || bytes.Contains(diagnostic.Bytes(), []byte(marker)) {
					t.Fatal("credential download diagnostics leaked")
				}
				if failDownload {
					want := map[string]string{"download-failure": "Go module download failed", "diagnostic-overflow": "diagnostics exceed 64 KiB", "invalid-archive": "Invalid Go bootstrap auth archive"}[outcome]
					if !strings.Contains(diagnostic.String(), want) {
						t.Fatal("missing controlled download error")
					}
				} else {
					if _, err := os.Stat(filepath.Join(root, "out", "program")); err != nil {
						t.Fatalf("compiler did not run after auth cleanup: %v; diagnostics=%s", err, diagnostic.String())
					}
				}
				err = filepath.Walk(filepath.Join(root, "mod"), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.Mode().IsRegular() {
						return nil
					}
					data, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					if bytes.Contains(data, []byte(marker)) {
						t.Error("credentials entered Go VCS/module cache")
					}
					return nil
				})
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
