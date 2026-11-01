// assert-outputs checks a scenario's boxed root completion payload.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

func assertOutputs(ctx context.Context, events <-chan model.Event, opPath string, expected []byte) error {
	var want map[string]interface{}
	if err := json.Unmarshal(expected, &want); err != nil || want == nil {
		return fmt.Errorf("expected outputs must be a JSON object")
	}
	canonicalWant, err := json.Marshal(want)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for root outputs: %w", ctx.Err())
		case event, ok := <-events:
			if !ok {
				return fmt.Errorf("event stream closed before root completion")
			}
			end := event.CallEnded
			if end == nil || end.Call.Op == nil || end.Call.ParentID != nil || end.Call.ID != end.Call.RootID || filepath.Clean(end.Call.Op.OpPath) != filepath.Clean(opPath) {
				continue
			}
			if end.Outcome != model.OpOutcomeSucceeded {
				return fmt.Errorf("expected root op to succeed")
			}
			outputs := end.Outputs
			if outputs == nil {
				outputs = map[string]*model.Value{}
			}
			raw, err := json.Marshal(outputs)
			if err != nil {
				return err
			}
			var actual map[string]interface{}
			if err := json.Unmarshal(raw, &actual); err != nil {
				return err
			}
			canonicalActual, err := json.Marshal(actual)
			if err != nil {
				return err
			}
			if !bytes.Equal(canonicalWant, canonicalActual) {
				return fmt.Errorf("root outputs do not match expected boxed JSON")
			}
			return nil
		}
	}
}
func run() error {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		return fmt.Errorf("usage: assert-outputs EXPECTED_JSON OP_PATH [API_URL]")
	}
	expected, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	opPath, err := filepath.Abs(os.Args[2])
	if err != nil {
		return err
	}
	api := "http://127.0.0.1:42224/api"
	if len(os.Args) == 4 {
		api = os.Args[3]
	}
	base, err := url.Parse(api)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	events, err := client.New(*base, nil).GetEventStream(ctx, &model.GetEventStreamReq{})
	if err != nil {
		return err
	}
	return assertOutputs(ctx, events, opPath, expected)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
