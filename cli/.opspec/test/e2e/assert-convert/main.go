// assert-convert copies fixtures before conversion and checks inspection-only plans.
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	"golang.org/x/mod/modfile"
)

type conversion struct {
	Root, Op, Expect string
	ExpectMessages   []string
	PlanOnly         bool
	Snapshot         string
}

const sdkModule = "github.com/opctl/opctl/sdks/go/opspec/sdk"

func tree(root string) (string, error) {
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		fmt.Fprintln(hash, rel, entry.Type())
		if !entry.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash.Write(data)
		}
		return nil
	})
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

func prepareCopy(source, target, sdk string) error {
	if err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == "opctlgen" {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(target, rel)
		if entry.IsDir() {
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
			return os.Chmod(dest, 0755)
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("fixture must contain regular files: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, info.Mode().Perm())
	}); err != nil {
		return err
	}
	path := filepath.Join(target, "go.mod")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		data = []byte("module opctl.test/convert\ngo 1.25\n")
	} else if err != nil {
		return err
	}
	mod, err := modfile.Parse(path, data, nil)
	if err != nil {
		return err
	}
	if err = mod.AddRequire(sdkModule, "v0.0.0"); err != nil {
		return err
	}
	if err = mod.AddReplace(sdkModule, "", sdk, ""); err != nil {
		return err
	}
	data, err = mod.Format()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func opctl(args ...string) ([]byte, []byte, error) {
	binary := os.Getenv("OPCTL_BINARY")
	if binary == "" {
		binary = "opctl"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func noSecrets(data []byte) error {
	for _, marker := range []string{"SECRET-MARKER", "PRIVATE-DEFAULT"} {
		if bytes.Contains(data, []byte(marker)) {
			return fmt.Errorf("secret marker leaked")
		}
	}
	return nil
}

func checkPlanEvent(event model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err = noSecrets(data); err != nil {
		return err
	}
	if out := event.ContainerStdOutWrittenTo; out != nil {
		if err = noSecrets(out.Data); err != nil {
			return err
		}
	}
	if out := event.ContainerStdErrWrittenTo; out != nil {
		if err = noSecrets(out.Data); err != nil {
			return err
		}
	}
	if start := event.CallStarted; start != nil && (start.Call.Op != nil || start.Call.Container != nil && (start.Call.Name == nil || !strings.HasPrefix(*start.Call.Name, "bootstrap "))) {
		return fmt.Errorf("plan executed an op call")
	}
	return nil
}

func plan(ref, argsFile, snapshot string) error {
	since := time.Now().UTC()
	stdout, stderr, err := opctl("op", "plan", ref, "--arg-file", argsFile, "-a", "token=SECRET-MARKER")
	if e := noSecrets(append(append([]byte{}, stdout...), stderr...)); e != nil {
		return e
	}
	fmt.Fprint(os.Stderr, string(stderr), string(stdout))
	if err != nil {
		return fmt.Errorf("plan failed: %w", err)
	}
	if !bytes.HasPrefix(stdout, []byte("# opctl op plan: snapshot")) {
		return fmt.Errorf("missing plan header")
	}
	if snapshot != "" {
		want, err := os.ReadFile(snapshot)
		if err != nil {
			return err
		}
		if !bytes.Equal(stdout, want) {
			return fmt.Errorf("plan differs from golden %s", snapshot)
		}
	} else if !bytes.Contains(stdout, []byte("has a secret default")) || !bytes.Contains(stdout, []byte("$(token)")) {
		return fmt.Errorf("secret default or symbolic token missing")
	}
	api := os.Getenv("OPCTL_PLAN_API_URL")
	if api == "" {
		api = "http://127.0.0.1:42224/api"
	}
	base, err := url.Parse(api)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stream, err := client.New(*base, nil).GetEventStream(ctx, &model.GetEventStreamReq{Filter: model.EventFilter{Since: &since}})
	if err != nil {
		return err
	}
	quiet := time.NewTimer(500 * time.Millisecond)
	defer quiet.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-quiet.C:
			return nil
		case event, ok := <-stream:
			if !ok {
				return fmt.Errorf("event stream closed")
			}
			if event.Timestamp.Before(since) {
				continue
			}
			if err = checkPlanEvent(event); err != nil {
				return err
			}
			if !quiet.Stop() {
				select {
				case <-quiet.C:
				default:
				}
			}
			quiet.Reset(500 * time.Millisecond)
		}
	}
}

func convert(config conversion, caller, target, argsFile string) (string, error) {
	sdk := os.Getenv("sdkPath")
	if sdk == "" {
		sdk = "/repo/sdks/go/opspec/sdk"
	}
	if err := prepareCopy(config.Root, target, sdk); err != nil {
		return "", err
	}
	before, err := tree(target)
	if err != nil {
		return "", err
	}
	ref := filepath.Join(target, config.Op)
	callerRel, err := filepath.Rel(config.Root, caller)
	if err != nil || callerRel == ".." || strings.HasPrefix(callerRel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("caller is outside conversion root")
	}
	caller = filepath.Join(target, callerRel)
	var callerSource []byte
	var callerPath string
	if caller != ref {
		for _, name := range []string{"op.go", "op.yml"} {
			path := filepath.Join(caller, name)
			if data, e := os.ReadFile(path); e == nil {
				callerSource, callerPath = data, path
				break
			}
		}
		if callerPath == "" {
			return "", fmt.Errorf("caller declaration missing")
		}
	}
	stdout, stderr, convertErr := opctl("op", "convert", "--go", ref, "--arg-file", argsFile)
	fmt.Fprint(os.Stderr, string(stdout), string(stderr))
	if config.Expect == "refused" {
		if convertErr == nil {
			return "", fmt.Errorf("expected conversion refusal")
		}
		after, err := tree(target)
		if err != nil || before != after {
			return "", fmt.Errorf("refused conversion changed files: %v", err)
		}
		for _, message := range config.ExpectMessages {
			if !bytes.Contains(stderr, []byte(message)) {
				return "", fmt.Errorf("missing refusal: %s", message)
			}
		}
		return "", nil
	}
	if config.Expect != "written" {
		return "", fmt.Errorf("unknown convert expectation: %s", config.Expect)
	}
	if convertErr != nil {
		return "", fmt.Errorf("convert failed: %w", convertErr)
	}
	if bytes.Contains(stderr, []byte("unverified")) {
		return "", fmt.Errorf("conversion was not verified")
	}
	if _, err := os.Stat(filepath.Join(ref, "op.go")); err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(ref, "op.yml")); !os.IsNotExist(err) {
		return "", fmt.Errorf("op.yml still exists")
	}
	if callerPath != "" {
		data, err := os.ReadFile(callerPath)
		if err != nil || !bytes.Equal(data, callerSource) {
			return "", fmt.Errorf("caller changed")
		}
	}
	if config.PlanOnly {
		if err := plan(ref, argsFile, filepath.Join(target, config.Snapshot)); err != nil {
			return "", err
		}
		return "", nil
	}
	return caller, nil
}

func run() error {
	if len(os.Args) == 4 && os.Args[1] == "plan" {
		return plan(os.Args[2], os.Args[3], "")
	}
	if len(os.Args) != 6 || os.Args[1] != "convert" {
		return fmt.Errorf("usage: assert-convert convert CONFIG_JSON CALLER COPY_DEST ARGS_FILE | plan OP_REF ARGS_FILE")
	}
	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		return err
	}
	var config conversion
	if err = json.Unmarshal(data, &config); err != nil {
		return err
	}
	ref, err := convert(config, os.Args[3], os.Args[4], os.Args[5])
	if err == nil && ref != "" {
		fmt.Println(ref)
	}
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
