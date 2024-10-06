package main

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/telemetry"
)

// telemetryGet implements "telemetry" command
func telemetryGet(ctx context.Context, cliOutput clioutput.CliOutput, dataDir string) (string, error) {
	config, err := telemetry.ReadConfigFile(dataDir)
	if err != nil {
		return "", err
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
