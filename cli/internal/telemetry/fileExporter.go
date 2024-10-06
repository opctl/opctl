// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Statement of changes: modified from assuming use within an OpenTelemetry Collector to directly integrated with golang sdk

package telemetry // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"

import (
	"context"

	"github.com/gofrs/flock"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// fileExporter is the implementation of file exporter that writes telemetry data to a file
type fileExporter struct {
	conf       *TelemetryConfig
	marshaller *marshaller
	writer     *fileWriter
	fileLock   *flock.Flock
}

func (e *fileExporter) ConsumeTraces(_ context.Context, td ptrace.Traces) error {
	if e.marshaller == nil {
		return nil
	}

	buf, err := e.marshaller.traces(td)
	if err != nil {
		return err
	}

	return e.export(buf)
}

func (e *fileExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	if e.marshaller == nil {
		return nil
	}

	buf, err := e.marshaller.metrics(md)
	if err != nil {
		return err
	}

	return e.export(buf)
}

func (e *fileExporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	if e.marshaller == nil {
		return nil
	}

	buf, err := e.marshaller.logs(ld)
	if err != nil {
		return err
	}

	return e.export(buf)
}

// Start starts the flush timer if set.
func (e *fileExporter) Start(_ context.Context) error {
	var err error
	e.marshaller = newMarshaller()
	e.writer, err = newFileWriter(e.conf)
	if err != nil {
		return err
	}
	e.writer.start()
	return nil
}

// Shutdown stops the exporter and is invoked during shutdown.
// It stops the flush ticker if set.
func (e *fileExporter) Shutdown(context.Context) error {
	if e.writer == nil {
		return nil
	}
	w := e.writer
	e.writer = nil
	return w.shutdown()
}

func (e *fileExporter) export(buf []byte) error {
	locked, err := e.fileLock.TryLock()
	if err != nil {
		return err
	}
	if locked {
		err = e.writer.export(buf)
		if err != nil {
			e.fileLock.Unlock()
			return err
		}
		return e.fileLock.Unlock()
	}

	return nil
}
