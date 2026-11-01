package node_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Run against a separately started isolated Docker node. The ordinary unit gate
// does not start a daemon or change the developer's running node.
func TestTypedOutputsDocker(t *testing.T) {
	binary, repo, data := os.Getenv("OPCTL_TYPED_E2E_BINARY"), os.Getenv("OPCTL_TYPED_E2E_REPO"), os.Getenv("OPCTL_TYPED_E2E_DATA")
	if binary == "" || repo == "" || data == "" {
		t.Skip("set OPCTL_TYPED_E2E_BINARY, OPCTL_TYPED_E2E_REPO, and OPCTL_TYPED_E2E_DATA for a live isolated Docker node")
	}
	base, _ := url.Parse("http://127.0.0.1:42999/api")
	api := client.New(*base, nil)
	for _, language := range []string{"YAML", "Go"} {
		cases := []struct{ name, diagnostic string }{
			{"pass", ""},
			{"missing", `"count": required output is missing`},
			{"wrong-type", `"count": expected number, received string`},
			{"wrong-string", `"version": expected string, received number`},
			{"bad-constraint", `"count": Must be greater than or equal to 1`},
			{"nested-constraint", "revision"},
			{"secret", ""},
			{"secret-constraint", "sensitive"},
		}
		if language == "YAML" {
			cases = append(cases, struct{ name, diagnostic string }{"declared-unbound", `"extra": required output is missing`})
		}
		for _, tc := range cases {
			t.Run(language+"/"+tc.name, func(t *testing.T) {
				dir := filepath.Join(repo, "test-suite/run/container/object/typed-outputs", tc.name)
				if language == "Go" {
					dir = filepath.Join(repo, "test-suite/run/go/typed-outputs", tc.name)
				}
				since := time.Now().UTC()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
				defer cancel()
				cmd := exec.CommandContext(ctx, binary, "--api-listen-address", "127.0.0.1:42999", "--dns-listen-address", "127.0.0.2:53", "--data-dir", data, "--container-runtime", "docker", "--no-color", "run", "--no-progress", dir)
				cmd.Dir = repo
				// Process diagnostics are kept in memory and deliberately omitted on
				// failure; event assertions below identify the failing boundary.
				output, runErr := cmd.CombinedOutput()
				if ctx.Err() != nil {
					t.Fatal("CLI timed out")
				}
				if (runErr != nil) != (tc.diagnostic != "") {
					t.Fatalf("CLI success=%t, want success=%t", runErr == nil, tc.diagnostic == "")
				}
				ctxEvents, stop := context.WithTimeout(context.Background(), 20*time.Second)
				defer stop()
				events, err := api.GetEventStream(ctxEvents, &model.GetEventStreamReq{Filter: model.EventFilter{Since: &since}})
				if err != nil {
					t.Fatal(err)
				}
				var root, producer string
				var producerEnd, rootEnd *model.CallEnded
				var containers int
				private := strings.HasPrefix(tc.name, "secret")
				marker := fmt.Sprintf("%032d", 42)
				assertPrivate := func(raw []byte) {
					for _, forbidden := range []string{marker, base64.StdEncoding.EncodeToString([]byte(marker)), base64.StdEncoding.EncodeToString([]byte(marker + "\n"))} {
						if strings.Contains(string(raw), forbidden) {
							t.Fatal("generated secret appeared in public diagnostics or events")
						}
					}
				}
				if private {
					assertPrivate(output)
				}
			collect:
				for {
					select {
					case <-ctxEvents.Done():
						t.Fatal("root completion not observed")
					case event, ok := <-events:
						if !ok {
							t.Fatal("event stream closed before root completion")
						}
						if start := event.CallStarted; start != nil && start.Call.ParentID == nil && start.Call.Op != nil && filepath.Clean(start.Call.Op.OpPath) == filepath.Clean(dir) {
							root = start.Call.RootID
						}
						if private {
							raw, err := json.Marshal(event)
							if err != nil {
								t.Fatal(err)
							}
							assertPrivate(raw)
							if out := event.ContainerStdOutWrittenTo; out != nil {
								assertPrivate(out.Data)
							}
							if out := event.ContainerStdErrWrittenTo; out != nil {
								assertPrivate(out.Data)
							}
						}
						if start := event.CallStarted; start != nil && root != "" && start.Call.RootID == root && start.Call.Container != nil && (start.Call.Name == nil || !strings.HasPrefix(*start.Call.Name, "bootstrap ")) {
							containers++
							if producer == "" {
								producer = start.Call.ID
							}
						}
						if end := event.CallEnded; end != nil && root != "" && end.Call.RootID == root {
							if end.Call.Container != nil && (end.Outcome == model.OpOutcomeFailed || (tc.diagnostic == "" && end.Call.ID == producer)) {
								producerEnd = end
							}
							if end.Call.ID == root {
								rootEnd = end
								break collect
							}
						}
					}
				}
				if rootEnd == nil || producerEnd == nil {
					t.Fatal("missing root or graph producer completion")
				}
				if tc.diagnostic != "" {
					if rootEnd.Outcome != model.OpOutcomeFailed {
						t.Fatal("root succeeded despite invalid produced value")
					}
					if producerEnd.Error == nil || !strings.Contains(producerEnd.Error.Message, tc.diagnostic) || producerEnd.Outcome != model.OpOutcomeFailed {
						t.Fatal("producer diagnostic did not identify the invalid member")
					}
					if containers != 1 {
						t.Fatal("consumer started after invalid producer")
					}
					if len(producerEnd.Outputs) != 0 || len(rootEnd.Outputs) != 0 {
						t.Fatal("failed producer published partial outputs")
					}
				} else if rootEnd.Outcome != model.OpOutcomeSucceeded || containers < 2 {
					t.Fatal("successful producer values did not reach the consumer")
				}
				if private && len(rootEnd.Outputs) != 0 {
					t.Fatal("public parent declaration exposed private result")
				}
			})
		}
	}
}
