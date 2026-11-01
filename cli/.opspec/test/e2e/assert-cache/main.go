// assert-cache runs one Go op twice against the same node and checks its bootstrap cache.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func checkCacheEvents(events []model.Event, warm bool) error {
	counts := map[string]int{}
	bootstrapIDs := map[string]bool{}
	for _, event := range events {
		if start := event.CallStarted; start != nil && start.Call.Name != nil {
			counts[*start.Call.Name]++
			if strings.HasPrefix(*start.Call.Name, "bootstrap ") {
				bootstrapIDs[start.Call.ID] = true
			}
		}
	}
	wantBuilds := 1
	if warm {
		wantBuilds = 0
	}
	if counts["bootstrap build"] != wantBuilds || counts["bootstrap describe"] != map[bool]int{false: 1, true: 0}[warm] || counts["bootstrap describe (cached)"] != map[bool]int{false: 1, true: 2}[warm] || counts["bootstrap plan"] != 1 {
		return fmt.Errorf("unexpected bootstrap counts: %v", counts)
	}
	for _, event := range events {
		if log := event.ContainerStdOutWrittenTo; log != nil && bootstrapIDs[log.ContainerID] {
			return fmt.Errorf("bootstrap stdout was published")
		}
	}
	return nil
}
func collect(ctx context.Context, base url.URL, since time.Time, path string) ([]model.Event, error) {
	stream, err := client.New(base, nil).GetEventStream(ctx, &model.GetEventStreamReq{Filter: model.EventFilter{Since: &since}})
	if err != nil {
		return nil, err
	}
	events := []model.Event{}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case event, ok := <-stream:
			if !ok {
				return nil, fmt.Errorf("event stream ended before root completion")
			}
			if event.Timestamp.Before(since) {
				continue
			}
			events = append(events, event)
			if end := event.CallEnded; end != nil && end.Call.Op != nil && end.Call.ParentID == nil && filepath.Clean(end.Call.Op.OpPath) == path {
				if end.Outcome != model.OpOutcomeSucceeded {
					return nil, fmt.Errorf("root op failed")
				}
				return events, nil
			}
		}
	}
}

type executable struct {
	Path, Checksum string
	Modified       time.Time
}

func entry(data string) (executable, map[string]interface{}, error) {
	entries, err := os.ReadDir(filepath.Join(data, "go", "programs", "v1"))
	if err != nil {
		return executable{}, nil, err
	}
	if len(entries) != 1 {
		return executable{}, nil, fmt.Errorf("expected one program entry, got %d", len(entries))
	}
	root := filepath.Join(data, "go", "programs", "v1", entries[0].Name())
	files, err := os.ReadDir(root)
	if err != nil {
		return executable{}, nil, err
	}
	if len(files) != 2 {
		return executable{}, nil, fmt.Errorf("cache retained files beyond executable and metadata")
	}
	path := filepath.Join(root, "program")
	binary, err := os.ReadFile(path)
	if err != nil {
		return executable{}, nil, err
	}
	sum := sha256.Sum256(binary)
	info, err := os.Stat(path)
	if err != nil {
		return executable{}, nil, err
	}
	metadataBytes, err := os.ReadFile(filepath.Join(root, "metadata.json"))
	if err != nil {
		return executable{}, nil, err
	}
	var metadata map[string]interface{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return executable{}, nil, err
	}
	return executable{path, hex.EncodeToString(sum[:]), info.ModTime()}, metadata, nil
}
func durations(events []model.Event) (time.Duration, time.Duration) {
	starts := map[string]time.Time{}
	var build, bootstrap time.Duration
	for _, event := range events {
		if start := event.CallStarted; start != nil && start.Call.Name != nil && strings.HasPrefix(*start.Call.Name, "bootstrap ") {
			starts[start.Call.ID] = event.Timestamp
		}
		if end := event.CallEnded; end != nil {
			if start, ok := starts[end.Call.ID]; ok {
				elapsed := event.Timestamp.Sub(start)
				bootstrap += elapsed
				if end.Call.Name != nil && *end.Call.Name == "bootstrap build" {
					build += elapsed
				}
			}
		}
	}
	return build, bootstrap
}
func run() error {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		return fmt.Errorf("usage: assert-cache OP_PATH ARG_FILE [API_URL]")
	}
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		return err
	}
	data := os.Getenv("OPCTL_DATA_DIR")
	if data == "" {
		return fmt.Errorf("OPCTL_DATA_DIR must select the scenario's isolated node")
	}
	api := "http://127.0.0.1:42224/api"
	if len(os.Args) == 4 {
		api = os.Args[3]
	}
	base, err := url.Parse(api)
	if err != nil {
		return err
	}
	binary := os.Getenv("OPCTL_BINARY")
	if binary == "" {
		binary = "opctl"
	}
	// Image preparation is outside the measurements. Versioned local tags work offline.
	for _, image := range []string{"golang:1.21", "alpine:3"} {
		if err := exec.Command("docker", "image", "inspect", image).Run(); err != nil {
			if output, err := exec.Command("docker", "pull", image).CombinedOutput(); err != nil {
				return fmt.Errorf("prepare fixture image: %w: %s", err, output)
			}
		}
	}
	var first executable
	var cold time.Duration
	var metadata map[string]interface{}
	for index, suffix := range []string{"first", "second"} {
		since := time.Now().UTC()
		start := time.Now()
		command := exec.Command(binary, "run", "--no-progress", "--arg-file", os.Args[2], "-a", "suffix="+suffix, path)
		output, err := command.CombinedOutput()
		elapsed := time.Since(start)
		fmt.Print(string(output))
		if err != nil {
			return fmt.Errorf("cache scenario run: %w", err)
		}
		if !strings.Contains(string(output), "hello from the app module "+suffix) {
			return fmt.Errorf("fresh plan did not consume changed input")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		events, err := collect(ctx, *base, since, path)
		cancel()
		if err != nil {
			return err
		}
		if err := checkCacheEvents(events, index == 1); err != nil {
			return err
		}
		actual, entryMetadata, err := entry(data)
		if err != nil {
			return err
		}
		metadata = entryMetadata
		build, bootstrap := durations(events)
		if index == 0 {
			first, cold = actual, elapsed
			fmt.Printf("COLD total=%s bootstrap=%s build=%s\n", elapsed, bootstrap, build)
		} else {
			if actual != first {
				return fmt.Errorf("warm execution changed executable checksum or mtime")
			}
			fmt.Printf("WARM total=%s bootstrap=%s build=%s ratio=%.2fx\n", elapsed, bootstrap, build, float64(cold)/float64(elapsed))
		}
	}
	source, err := goop.Discover(path, "")
	if err != nil {
		return err
	}
	image := goop.Toolchain{ID: metadata["ImageID"].(string), OS: metadata["OS"].(string), Arch: metadata["Arch"].(string), Variant: metadata["Variant"].(string)}
	start := time.Now()
	key, err := source.Fingerprint(context.Background(), data, image)
	if err != nil {
		return err
	}
	fmt.Printf("HASH duration=%s key=%s\n", time.Since(start), key)
	fmt.Println("PASS: one cold build, zero warm builds, stable executable, memoised describe, fresh plans with changed inputs")
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
