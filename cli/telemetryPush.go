package main

import (
	"github.com/opctl/opctl/cli/internal/telemetry"
)

// telemetryPush implements "telemetry push" sub command
func telemetryPush(t *telemetry.Telemetry, dataDir string) error {
	return telemetry.Push(dataDir)
}
