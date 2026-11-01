package node

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

func authFiles(t *testing.T, auth bootstrapAuth) map[string][]byte {
	t.Helper()
	reader := tar.NewReader(bytes.NewReader(auth.archive))
	files := map[string][]byte{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		if header.Mode != 0600 {
			t.Fatal("credential archive file mode")
		}
		files[header.Name] = data
	}
	return files
}
func TestBootstrapAuthScopesAndProtocol(t *testing.T) {
	records := []model.Auth{
		{Resources: "github.com", Creds: model.Creds{Username: "host", Password: "host-token"}},
		{Resources: "https://github.com/acme", Creds: model.Creds{Username: "org", Password: "org-token"}},
		{Resources: "github.com/acme/repo", Creds: model.Creds{Username: " ", Password: "literal 'quote\" $x \\ token"}},
		{Resources: "git.example.com/team", Creds: model.Creds{Username: "transitive", Password: "transitive-token"}},
		{Resources: "docker.io/library/golang:1.21", Creds: model.Creds{Username: "ignored", Password: "ignored"}},
	}
	auth, err := prepareBootstrapAuth(records, "github.com/acme/repo", &model.Creds{Username: "explicit", Password: "explicit-token"})
	if err != nil {
		t.Fatal(err)
	}
	if auth.private != "git.example.com/team,github.com,github.com/acme,github.com/acme/repo" {
		t.Fatalf("GOPRIVATE: %s", auth.private)
	}
	files := authFiles(t, auth)
	if string(files[".netrc"]) != "machine github.com login host password host-token\n" {
		t.Fatal("netrc widened a repository or whitespace credential")
	}
	for _, tc := range []struct{ protocol, host, path, wantUser, wantPass string }{
		{"https", "github.com", "acme/repo", "explicit", "explicit-token"},
		{"https", "github.com", "acme/repo/sub", "explicit", "explicit-token"},
		{"https", "github.com", "acme/repository", "org", "org-token"},
		{"https", "github.com", "acme-other/repo", "host", "host-token"},
		{"https", "github.com.evil", "acme/repo", "", ""},
		{"https", "github.com", "acme/repo/../other", "", ""},
		{"https", "github.com", "acme/repo/%2e%2e/other", "", ""},
		{"http", "github.com", "acme/repo", "", ""},
		{"https", "git.example.com", "team/dependency", "transitive", "transitive-token"},
	} {
		command := exec.Command("sh", "-c", string(files["helper"]), "helper", "get")
		command.Stdin = strings.NewReader("protocol=" + tc.protocol + "\nhost=" + tc.host + "\npath=" + tc.path + "\n\n")
		output, err := command.Output()
		if err != nil {
			t.Fatal("helper execution failed")
		}
		want := ""
		if tc.wantUser != "" {
			want = "username=" + tc.wantUser + "\npassword=" + tc.wantPass + "\n"
		}
		if string(output) != want {
			t.Fatal("helper released credentials outside selected HTTPS scope")
		}
	}
	auth, err = prepareBootstrapAuth(records, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	files = authFiles(t, auth)
	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, files["helper"], 0600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("git", "-c", "credential.helper=", "-c", "credential.helper=!sh "+helper, "-c", "credential.useHttpPath=true", "credential", "fill")
	command.Dir = dir
	command.Env = append(os.Environ(), "HOME="+dir, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0")
	command.Stdin = strings.NewReader("protocol=https\nhost=github.com\npath=acme/repo\n\n")
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	output, err := command.Output()
	if err != nil {
		message := diagnostic.String()
		for _, record := range records {
			message = strings.ReplaceAll(message, record.Password, "<redacted>")
		}
		t.Fatalf("Git credential protocol failed: %s", message)
	}
	if !bytes.Contains(output, []byte("password=literal 'quote\" $x \\ token\n")) {
		t.Fatal("Git protocol changed literal credentials")
	}
	for _, operation := range []string{"store", "erase"} {
		command := exec.Command("sh", helper, operation)
		output, err := command.Output()
		if err != nil || len(output) != 0 {
			t.Fatal("helper persisted credentials")
		}
	}
}
func TestBootstrapAuthRejectsUnsupportedCredentials(t *testing.T) {
	for _, password := range []string{"line\nbreak", "nul\x00value", "cr\rvalue"} {
		if _, err := prepareBootstrapAuth([]model.Auth{{Resources: "example.com/repo", Creds: model.Creds{Password: password}}}, "", nil); err == nil {
			t.Fatal("invalid credential accepted")
		}
	}
	for _, resource := range []string{"http://example.com/repo", "image:1.21", "example.com/repo#version", "example.com/../other", "example.com/acme*"} {
		auth, err := prepareBootstrapAuth([]model.Auth{{Resources: resource}}, "", nil)
		if err != nil || len(auth.archive) != 0 {
			t.Fatalf("unsupported resource %q accepted", resource)
		}
	}
	auth, err := prepareBootstrapAuth(nil, "/local/path", &model.Creds{Username: "local", Password: "must-not-infer"})
	if err != nil || len(auth.archive) != 0 {
		t.Fatal("local op inferred remote credential scope")
	}
}
func TestBootstrapHostNetrcDiscoveryMapping(t *testing.T) {
	user, token := "synthetic-user", "synthetic-token"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != token {
			w.Header().Set("WWW-Authenticate", `Basic realm="modules"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte("go-import discovery"))
	}))
	defer server.Close()
	auth, err := prepareBootstrapAuth([]model.Auth{{Resources: "127.0.0.1", Creds: model.Creds{Username: user, Password: token}}}, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(string(authFiles(t, auth)[".netrc"]))
	if len(fields) != 6 {
		t.Fatal("host discovery netrc missing")
	}
	request, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL+"?go-get=1", nil)
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 401 {
		t.Fatal("discovery did not challenge")
	}
	request.SetBasicAuth(fields[3], fields[5])
	response, err = server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatal("host netrc mapping failed discovery")
	}
}

type bootstrapSnapshotStore struct {
	stateStore
	auths []model.Auth
}

func (store *bootstrapSnapshotStore) TryGetAuth(string) *model.Auth { return nil }
func (store *bootstrapSnapshotStore) AuthSnapshot() ([]model.Auth, error) {
	return append([]model.Auth(nil), store.auths...), nil
}

func TestBootstrapAuthDeliveryAndRotation(t *testing.T) {
	bootstrap, cache, op, builds := cachedBootstrapFixture(t)
	marker := "synthetic-delivery-secret"
	store := &bootstrapSnapshotStore{auths: []model.Auth{{Resources: "example.com/team", Creds: model.Creds{Username: "user", Password: marker}}}}
	cache.stateStore = store
	caller := bootstrap.(_bootstrapper).caller.(*fakes.FakeCaller)
	original := caller.CallStub
	var delivered bool
	caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
		input := containerruntime.BootstrapInputFrom(ctx)
		if *spec.Name == "bootstrap build" {
			if input == nil {
				t.Fatal("build received no private auth archive")
			}
			data, err := io.ReadAll(input)
			if err != nil {
				t.Fatal(err)
			}
			delivered = bytes.Contains(data, []byte(marker))
			if spec.Container.EnvVars.(map[string]interface{})["GOPRIVATE"] != "example.com/team" {
				t.Fatal("private scope mapping")
			}
		} else if input != nil {
			t.Fatal("describe/plan received auth")
		}
		data, err := json.Marshal([]interface{}{op, spec, scope})
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte(marker)) {
			t.Fatal("credentials entered serialized call or scope")
		}
		return original(ctx, id, scope, spec, path, parent, root)
	}
	if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := bootstrap.Plan(context.Background(), op, "", map[string]interface{}{"public": "fresh-plan-marker"}); err != nil {
		t.Fatal(err)
	}
	store.auths[0].Password = "rotated-secret"
	if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	if !delivered || builds.Load() != 1 {
		t.Fatal("auth rotation changed compiled identity or auth was not delivered")
	}
	err := filepath.Walk(cache.directory, func(path string, info os.FileInfo, err error) error {
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
		if bytes.Contains(data, []byte(marker)) || bytes.Contains(data, []byte("fresh-plan-marker")) {
			t.Fatal("credentials or plan response entered program cache")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBootstrapAuthMetadataStaysPrivate(t *testing.T) {
	marker := "synthetic-provenance-secret"
	creds := &model.Creds{Username: "u", Password: marker}
	data, err := json.Marshal([]interface{}{model.OpCall{BootstrapSource: "example.com/team/repo", BootstrapCreds: creds}, model.OpCallSpec{BootstrapSource: "example.com/team/repo", BootstrapCreds: creds}})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte(marker)) || bytes.Contains(data, []byte("BootstrapSource")) {
		t.Fatal("bootstrap provenance was serialized")
	}
}
