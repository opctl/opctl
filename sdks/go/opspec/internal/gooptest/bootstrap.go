// Package gooptest runs bootstrap phases locally for integration tests.
package gooptest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"os/exec"
	"path/filepath"
)

// Phase runs the same program contract as node bootstrap, without a node.
func Phase(opDir, phase string, inputs map[string]any) ([]byte, error) {
	request, err := json.Marshal(goop.NewRequest(phase, inputs))
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp("", "goop-request-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	requestPath := filepath.Join(dir, "request.json")
	if err = os.WriteFile(requestPath, request, 0600); err != nil {
		return nil, err
	}
	cmd := exec.Command("go", "run", ".", phase)
	cmd.Dir = opDir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly", "GOTOOLCHAIN=local", "CGO_ENABLED=0", "GOWORK=off", "OPCTL_BOOTSTRAP_REQUEST="+requestPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", phase, err, stderr.String())
	}
	return stdout.Bytes(), nil
}

// Bootstrap runs describe and plan against an existing Go op.
func Bootstrap(opDir string, inputs map[string]any) (*model.OpSpec, *goop.Plan, error) {
	data, err := Phase(opDir, "describe", nil)
	if err != nil {
		return nil, nil, err
	}
	desc, err := goop.ParseDescribe(data)
	if err != nil {
		return nil, nil, err
	}
	data, err = Phase(opDir, "plan", inputs)
	if err != nil {
		return nil, nil, err
	}
	plan, err := goop.ParsePlan(data)
	return desc, plan, err
}
