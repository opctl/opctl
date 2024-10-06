package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/cli/internal/telemetry"
)

// telemetrySet implements "telemetry set" sub command
func telemetrySet(
	ctx context.Context,
	config *telemetry.TelemetryConfig,
	collector string,
	endpoint string,
	dataDir string,
	pushInterval int,
	tracesEnabled bool,
	metricsEnabled bool,
	logsEnabled bool,
) (string, error) {
	err := config.Update(&telemetry.TelemetryConfig{
		Collector:      collector,
		Endpoint:       endpoint,
		Path:           telemetry.DataFilePath(dataDir),
		PushInterval:   time.Second * time.Duration(pushInterval),
		TracesEnabled:  tracesEnabled,
		MetricsEnabled: metricsEnabled,
		LogsEnabled:    logsEnabled,
	})
	if err != nil {
		return "", fmt.Errorf("failed to write telemetry config: %w", err)
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
