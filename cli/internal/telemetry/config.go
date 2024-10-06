// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Statement of changes: modified from assuming use within an OpenTelemetry Collector to directly integrated with golang sdk

package telemetry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"gopkg.in/yaml.v3"
)

var (
	DefaultPushInterval  = time.Second * 10
	DefaultFlushInterval = time.Second
	defaultRotation      = &Rotation{
		MaxMegabytes: 1,
		MaxDays:      0,
		MaxBackups:   1,
		LocalTime:    false,
	}
)

// Config defines configuration for file exporter.
type TelemetryConfig struct {
	// Collector type
	Collector string `yaml:"collector" json:"collector"`

	// Collector endpoint
	Endpoint string `yaml:"endpoint" json:"endpoint"`

	// Enable traces
	TracesEnabled bool `yaml:"tracesEnabled" json:"tracesEnabled"`

	// Enable Metrics
	MetricsEnabled bool `yaml:"metricsEnabled" json:"metricsEnabled"`

	// Enable Logs
	LogsEnabled bool `yaml:"logsEnabled" json:"logsEnabled"`

	// Path of the file to write to. Path is relative to current directory.
	Path string `yaml:"path" json:"path"`

	// FlushInterval is the duration between flushes.
	// See time.ParseDuration for valid values.
	FlushInterval time.Duration `yaml:"flushInterval" json:"flushInterval"`

	// PushInterval is the duration between pushes.
	// see time.ParseDuration for valid values.
	PushInterval time.Duration `yaml:"pushInterval" json:"pushInterval"`

	// rotation defines an option about rotation of telemetry files. Ignored
	// when GroupByAttribute is used.
	rotation *Rotation // `yaml:"rotation" json:"rotation"`
}

// Rotation an option to rolling log files
type Rotation struct {
	// MaxMegabytes is the maximum size in megabytes of the file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxMegabytes int `yaml:"maxMegabytes" json:"maxMegabytes"`

	// MaxDays is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxDays int `yaml:"maxDays" json:"maxDays"`

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to 100 files.
	MaxBackups int `yaml:"maxBackups" json:"maxBackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `yaml:"localtime" json:"localtime"`
}

var _ component.Config = (*TelemetryConfig)(nil)

func ReadConfigFile(dataDir string) (*TelemetryConfig, error) {
	config := DefaultConfig(dataDir)

	configFilePath, err := xdg.ConfigFile(fmt.Sprintf("%s/telemetry/config.yaml", jobName))
	if err != nil {
		return config, fmt.Errorf("failed to get config file path: %w", err)
	}

	configFileBytes, fileReadErr := os.ReadFile(configFilePath)
	if fileReadErr == nil {
		err = yaml.Unmarshal(configFileBytes, &config)
		if err != nil {
			return config, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	if config.PushInterval <= 0 {
		config.PushInterval = DefaultPushInterval
	}
	if config.Path == "" {
		config.Path = DataFilePath(dataDir)
	}
	if config.Collector == "" {
		config.Collector = os.Getenv(EnvVarOpctlTelemetryCollector)
	}
	if config.Endpoint == "" {
		config.Endpoint = os.Getenv(EnvVarOpctlTelemetryEndpoint)
	}

	if fileReadErr != nil {
		err = config.writeConfigFile()
		if err != nil {
			return config, fmt.Errorf("failed to init config file: %w", err)
		}
	}

	resolvedPath, err := filepath.Abs(filepath.Join(config.Path, "../"))
	if err != nil {
		return config, fmt.Errorf("error initializing telemetry data dir: %w", err)
	}
	if err := os.MkdirAll(resolvedPath, 0775|os.ModeSetgid); err != nil {
		return config, fmt.Errorf("error initializing telemetry data dir: %w", err)
	}

	return config, nil
}

func (cfg *TelemetryConfig) Update(nextConfig *TelemetryConfig) error {
	if cfg.Path != nextConfig.Path {
		err := os.Remove(cfg.Path)
		if err != nil {
			return fmt.Errorf("failed to remove current file: %w", err)
		}
	}
	Merge(cfg, nextConfig)
	return cfg.writeConfigFile()
}

func (cfg *TelemetryConfig) Disabled() bool {
	return (cfg.TracesEnabled == false && cfg.MetricsEnabled == false && cfg.LogsEnabled == false)
}

func Merge(dest *TelemetryConfig, src *TelemetryConfig) {
	collectorEnvVar := os.Getenv(EnvVarOpctlTelemetryCollector)
	if src.Collector == "" && collectorEnvVar != "" {
		src.Collector = collectorEnvVar
	}
	dest.Collector = src.Collector

	endpointEnvVar := os.Getenv(EnvVarOpctlTelemetryEndpoint)
	if src.Endpoint == "" && endpointEnvVar != "" {
		src.Endpoint = endpointEnvVar
	}
	dest.Endpoint = src.Endpoint

	dest.TracesEnabled = src.TracesEnabled
	dest.MetricsEnabled = src.MetricsEnabled
	dest.LogsEnabled = src.LogsEnabled
	dest.Path = src.Path

	if src.FlushInterval > 0 {
		dest.FlushInterval = src.FlushInterval
	} else {
		dest.FlushInterval = DefaultFlushInterval
	}
	if src.PushInterval > 0 {
		dest.PushInterval = src.PushInterval
	} else {
		dest.PushInterval = DefaultPushInterval
	}
	dest.rotation = defaultRotation
	// if src.Rotation != nil {
	// 	dest.Rotation.MaxMegabytes = src.Rotation.MaxMegabytes
	// 	dest.Rotation.MaxDays = src.Rotation.MaxDays
	// 	dest.Rotation.MaxBackups = src.Rotation.MaxBackups
	// 	dest.Rotation.LocalTime = src.Rotation.LocalTime
	// }
}

// Validate checks if the exporter configuration is valid
func (cfg *TelemetryConfig) Validate() error {
	if cfg.Path == "" {
		return errors.New("path must be non-empty")
	}
	if cfg.FlushInterval <= 0 {
		return errors.New("flushInterval must be larger than zero")
	}
	if cfg.PushInterval <= 0 {
		return errors.New("pushInterval must be larger than zero")
	}

	return nil
}

func (cfg *TelemetryConfig) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		return errors.New("empty config for file exporter")
	}
	// first load the config normally
	err := componentParser.Unmarshal(cfg)
	if err != nil {
		return err
	}

	// set flush interval to 1 second if not set.
	if cfg.FlushInterval == 0 {
		cfg.FlushInterval = DefaultFlushInterval
	}
	return nil
}

func DefaultConfig(dataDir string) *TelemetryConfig {
	collector := os.Getenv(EnvVarOpctlTelemetryCollector)
	return &TelemetryConfig{
		Collector:      collector,
		Endpoint:       os.Getenv(EnvVarOpctlTelemetryEndpoint),
		Path:           DataFilePath(dataDir),
		PushInterval:   DefaultPushInterval,
		FlushInterval:  DefaultFlushInterval,
		LogsEnabled:    collector != "prom",
		MetricsEnabled: true,
		TracesEnabled:  collector != "prom",
		rotation:       defaultRotation,
	}
}

func DataFilePath(dataDir string) string {
	return filepath.Join(dataDir, "telemetry", "otlp.jsonl")
}

func (cfg *TelemetryConfig) writeConfigFile() error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	configFilePath, err := xdg.ConfigFile(fmt.Sprintf("%s/telemetry/config.yaml", jobName))
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	yamlFile, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config file: %w", err)
	}

	err = os.WriteFile(configFilePath, yamlFile, 0666)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
