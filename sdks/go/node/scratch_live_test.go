//go:build livescratch

// Proves scratch cleanup against a real Docker runtime, driving the node core as the API does:
//
//	go test -tags livescratch -run TestScratchLive -v ./node/
//
// Containers are unnamed, so no DNS resolver entries are written.
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
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
)

func TestScratchLive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runtime, err := docker.New(ctx, "unix:///var/run/docker.sock")
	if err != nil {
		t.Fatal(err)
	}

	dataDir := t.TempDir()
	c, err := New(ctx, runtime, dataDir, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	op, _ := filepath.Abs("testdata/scratchproof")

	start := func(ref, secret string, clean bool) string {
		id, err := c.StartOp(ctx, model.StartOpReq{
			Args:         map[string]*model.Value{"secret": {String: &secret}},
			CleanScratch: clean,
			Op:           model.StartOpReqOp{Ref: ref},
		})
		if err != nil {
			t.Fatal(err)
		}

		return id
	}

	// ended waits for the root's CallEnded, as forge-k8s-connector does.
	ended := func(rootID string) model.CallEnded {
		events, _ := c.GetEventStream(ctx, &model.GetEventStreamReq{Filter: model.EventFilter{Roots: []string{rootID}}})
		for {
			select {
			case e := <-events:
				if e.CallEnded != nil && e.CallEnded.Call.ID == rootID {
					// CallEnded is published before StartOp's goroutine records the end; give it a moment.
					time.Sleep(500 * time.Millisecond)
					return *e.CallEnded
				}
			case <-time.After(2 * time.Minute):
				t.Fatalf("%s never ended", rootID)
			}
		}
	}

	// onDisk reports whether secret is in any file under the data dir.
	onDisk := func(secret string) bool {
		found := false
		_ = filepath.Walk(filepath.Join(dataDir, "dcg"), func(p string, info os.FileInfo, err error) error {
			if err == nil && info.Mode().IsRegular() && !strings.Contains(p, string(filepath.Separator)+"events"+string(filepath.Separator)) {
				if b, _ := os.ReadFile(p); bytes.Contains(b, []byte(secret)) {
					found = true
				}
			}
			return nil
		})
		return found
	}

	t.Run("a CLI-style run is cleaned when it ends, after its steps shared outputs", func(t *testing.T) {
		secret := "cli-secret-7f3a"
		end := ended(start(op, secret, true))

		if end.Outcome != model.OpOutcomeSucceeded {
			t.Fatalf("outcome %s: %v", end.Outcome, end.Error)
		}
		if onDisk(secret) {
			t.Fatal("secret still in scratch after a CleanScratch run ended")
		}
	})

	t.Run("an API-style run keeps a dir output readable after it ends, until the TTL", func(t *testing.T) {
		secret := "api-secret-91bc"
		end := ended(start(op, secret, false))

		if end.Outcome != model.OpOutcomeSucceeded {
			t.Fatalf("outcome %s: %v", end.Outcome, end.Error)
		}

		// forge-k8s-connector's pull: list the dir output, then fetch each file through the node.
		dir := *end.Outputs["out"].Dir
		entries, err := c.ListDescendants(ctx, model.ListDescendantsReq{DataRef: dir})
		if err != nil || len(entries) == 0 {
			t.Fatalf("listing %s after CallEnded: %v (%d entries)", dir, err, len(entries))
		}
		r, err := c.GetData(ctx, model.GetDataReq{DataRef: filepath.Join(dir, entries[0].Path)})
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(r)
		r.Close()
		if strings.TrimSpace(string(b)) != secret {
			t.Fatalf("read %q from the dir output", b)
		}

		time.Sleep(8 * time.Second) // past the 3s TTL and its sweep
		if onDisk(secret) {
			t.Fatal("secret still in scratch after the TTL")
		}
	})

	t.Run("a killed CLI-style run is cleaned too", func(t *testing.T) {
		secret := "killed-secret-c4d2"
		slow, _ := filepath.Abs("testdata/scratchproof/slow")
		rootID := start(slow, secret, true)

		deadline := time.Now().Add(time.Minute)
		for !onDisk(secret) {
			if time.Now().After(deadline) {
				t.Fatal("secret never reached scratch")
			}
			time.Sleep(200 * time.Millisecond)
		}

		if err := c.KillOp(ctx, model.KillOpReq{OpID: rootID, RootCallID: rootID}); err != nil {
			t.Fatal(err)
		}
		if end := ended(rootID); end.Outcome != model.OpOutcomeKilled {
			t.Fatalf("outcome %s", end.Outcome)
		}
		if onDisk(secret) {
			t.Fatal("secret still in scratch after a killed CleanScratch run")
		}
	})
}
