// The launcher is embedded and compiled by the bootstrap toolchain. It has no module dependencies.
package main

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

const authDir = "/run/opctl-auth"

func requireTmpfs(path string) error {
	var stat syscall.Statfs_t
	if syscall.Statfs(path, &stat) != nil || stat.Type != 0x01021994 {
		return errors.New("Go bootstrap auth requires tmpfs")
	}
	return nil
}

func extractAuth(input io.Reader, dir string, uid, gid int) error {
	data, err := io.ReadAll(io.LimitReader(input, (1<<20)+1))
	if err != nil {
		return errors.New("Invalid Go bootstrap auth archive")
	}
	if len(data) > 1<<20 {
		return errors.New("Go bootstrap private input exceeds 1 MiB")
	}
	files := map[string][]byte{}
	reader := tar.NewReader(bytes.NewReader(data))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Invalid Go bootstrap auth archive")
		}
		switch header.Name {
		case "helper", "gitconfig", ".netrc":
		default:
			return errors.New("Invalid Go bootstrap auth file")
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			return errors.New("Go bootstrap auth files must be regular")
		}
		if _, ok := files[header.Name]; ok {
			return errors.New("Cannot extract Go bootstrap auth archive")
		}
		content, err := io.ReadAll(reader)
		if err != nil {
			return errors.New("Invalid Go bootstrap auth archive")
		}
		files[header.Name] = content
	}
	for name, content := range files {
		file, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return errors.New("Cannot extract Go bootstrap auth archive")
		}
		_, writeErr := file.Write(content)
		ownerErr := file.Chown(uid, gid)
		closeErr := file.Close()
		if writeErr != nil || ownerErr != nil || closeErr != nil {
			return errors.New("Cannot extract Go bootstrap auth archive")
		}
	}
	return os.Chown(dir, uid, gid)
}

func cleanAuth() error {
	entries, err := os.ReadDir(authDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(authDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

type cappedFile struct {
	file   *os.File
	size   int
	cancel context.CancelFunc
}

func (w *cappedFile) Write(p []byte) (int, error) {
	remaining := 65537 - w.size
	if len(p) > remaining {
		p = p[:remaining]
	}
	n, err := w.file.Write(p)
	w.size += n
	if w.size > 65536 {
		w.cancel()
		return n, errors.New("Go module download diagnostics exceed 64 KiB")
	}
	return n, err
}

func run(ctx context.Context, uid, gid int, env []string, stdout, stderr io.Writer, args ...string) error {
	command := exec.CommandContext(ctx, "setpriv", append([]string{"--reuid", strconv.Itoa(uid), "--regid", strconv.Itoa(gid), "--clear-groups", "--no-new-privs"}, args...)...)
	command.Env = env
	command.WaitDelay = time.Second
	// Child processes never inherit the private Docker input stream.
	command.Stdout, command.Stderr = stdout, stderr
	return command.Run()
}

func download(ctx context.Context, uid, gid int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	out, err := os.OpenFile(filepath.Join(authDir, "download.stdout"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("Go module download failed before diagnostic capture")
	}
	defer out.Close()
	diagnostic, err := os.OpenFile(filepath.Join(authDir, "download.stderr"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("Go module download failed before diagnostic capture")
	}
	defer diagnostic.Close()
	stdout, stderr := &cappedFile{file: out, cancel: cancel}, &cappedFile{file: diagnostic, cancel: cancel}
	env := append(os.Environ(), "HOME="+authDir, "GIT_CONFIG_GLOBAL="+authDir+"/gitconfig", "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0", "GIT_TRACE=0", "GIT_TRACE_CURL=0", "GIT_CURL_VERBOSE=0", "GIT_TRACE_PACKET=0", "GIT_TRACE_SETUP=0", "GIT_TRACE_PERFORMANCE=0", "GIT_SSH_COMMAND=false")
	err = run(ctx, uid, gid, env, stdout, stderr, "go", "mod", "download")
	if stdout.size > 65536 || stderr.size > 65536 {
		return errors.New("Go module download diagnostics exceed 64 KiB")
	}
	if err != nil {
		data, _ := os.ReadFile(diagnostic.Name())
		if regexp.MustCompile(`(?i)401|403|unauthorized|authentication|terminal prompts disabled`).Match(data) {
			return errors.New("Go module download/auth failed. Repository-scoped credentials support HTTPS Git; authenticated vanity discovery on Go 1.21 needs host-scoped auth.")
		}
		return errors.New("Go module download failed. Check auth and update go.mod/go.sum outside bootstrap; private download diagnostics are suppressed.")
	}
	return nil
}

func launch(ctx context.Context, args []string) error {
	owner := "/opctl/mod"
	if _, err := os.Stat(owner); os.IsNotExist(err) {
		owner = "/opctl/program"
	}
	if owner == "/opctl/mod" {
		if err := os.Setenv("GOCACHE", "/opctl/cache"); err != nil {
			return err
		}
	}
	info, err := os.Stat(owner)
	if err != nil {
		return errors.New("Go bootstrap cache must be shared at the same host path with the container daemon")
	}
	stat := info.Sys().(*syscall.Stat_t)
	uid, gid := int(stat.Uid), int(stat.Gid)
	if _, err := exec.LookPath("setpriv"); err != nil {
		return errors.New("Go bootstrap requires setpriv to drop privileges")
	}
	scratch, err := os.MkdirTemp("/tmp", "opctl-bootstrap-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratch)
	defer cleanAuth()
	if err := os.Chown(scratch, uid, gid); err != nil {
		return err
	}
	for _, name := range []string{"home", "tmp", "gopath"} {
		path := filepath.Join(scratch, name)
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
		if err := os.Chown(path, uid, gid); err != nil {
			return err
		}
	}
	for name, dir := range map[string]string{"HOME": "home", "TMPDIR": "tmp", "GOPATH": "gopath"} {
		if err := os.Setenv(name, filepath.Join(scratch, dir)); err != nil {
			return err
		}
	}
	_, buildErr := os.Stat("/opctl/build")
	building := buildErr == nil
	if building {
		if err := os.Chown("/opctl/build", uid, gid); err != nil {
			return err
		}
	}
	if _, err := os.Stat(authDir); err == nil {
		if err := requireTmpfs(authDir); err != nil {
			return err
		}
		err := extractAuth(os.Stdin, authDir, uid, gid)
		os.Stdin.Close()
		if err != nil {
			return err
		}
		if err := download(ctx, uid, gid); err != nil {
			return err
		}
		if err := cleanAuth(); err != nil {
			return err
		}
	} else {
		os.Stdin.Close()
		if building {
			if err := run(ctx, uid, gid, os.Environ(), os.Stdout, os.Stderr, "go", "mod", "download"); err != nil {
				return err
			}
		}
	}
	if building {
		os.Setenv("GOPROXY", "off")
		os.Setenv("GONOPROXY", "none")
	}
	if len(args) == 1 && args[0] == "build" {
		// Cache the launcher beside the program so describe and plan need no compiler.
		path, err := os.Executable()
		if err != nil {
			return err
		}
		binary, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile("/opctl/build/launcher", binary, 0700); err != nil {
			return err
		}
		if err := os.Chown("/opctl/build/launcher", uid, gid); err != nil {
			return err
		}
		version, err := os.OpenFile("/opctl/build/go-version", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		err = run(ctx, uid, gid, os.Environ(), version, os.Stderr, "go", "version")
		closeErr := version.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		return run(ctx, uid, gid, os.Environ(), os.Stdout, os.Stderr, "go", "build", "-mod=readonly", "-buildvcs=false", "-trimpath", "-o", "/opctl/build/program", ".")
	}
	if len(args) == 0 {
		return errors.New("Go bootstrap requires a program")
	}
	return run(ctx, uid, gid, os.Environ(), os.Stdout, os.Stderr, args...)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	if err := launch(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
