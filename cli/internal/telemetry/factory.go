// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Statement of changes: modified from assuming use within an OpenTelemetry Collector to directly integrated with golang sdk

package telemetry

import (
	"context"

	"github.com/gofrs/flock"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TracesStability  = component.StabilityLevelAlpha
	MetricsStability = component.StabilityLevelAlpha
	LogsStability    = component.StabilityLevelAlpha
)

type FileExporter interface {
	ConsumeTraces(context.Context, ptrace.Traces) error
	ConsumeMetrics(context.Context, pmetric.Metrics) error
	ConsumeLogs(context.Context, plog.Logs) error
	Start(context.Context) error
	Shutdown(context.Context) error
}

func NewFileExporter(conf *TelemetryConfig, fileLock *flock.Flock) FileExporter {
	return &fileExporter{
		conf:     conf,
		fileLock: fileLock,
	}
}

func newFileWriter(conf *TelemetryConfig) (*fileWriter, error) {
	return &fileWriter{
		path: conf.Path,
		file: &lumberjack.Logger{
			Filename:   conf.Path,
			MaxSize:    conf.rotation.MaxMegabytes,
			MaxAge:     conf.rotation.MaxDays,
			MaxBackups: conf.rotation.MaxBackups,
			LocalTime:  conf.rotation.LocalTime,
		},
		flushInterval: conf.FlushInterval,
	}, nil
}
