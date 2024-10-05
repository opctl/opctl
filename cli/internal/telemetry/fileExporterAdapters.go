package telemetry

import (
	"context"
	"fmt"
	"sync"
)

// FileExporterAdapters wraps the logic of your custom OpenTelemetry exporter.
type FileExporterAdapters struct {
	trace  TraceFileExporter  // For traces
	metric MetricFileExporter // For metrics
	log    LogFileExporter    // For logs
	fe     FileExporter

	stoppedMu sync.RWMutex
	stopped   bool
}

// NewCustomExporter initializes the custom exporter for traces, metrics, and logs.
func NewExporterAdapters(ctx context.Context, fe FileExporter) *FileExporterAdapters {
	ea := &FileExporterAdapters{
		trace:  TraceFileExporter{exporter: fe},
		metric: MetricFileExporter{exporter: fe},
		log:    LogFileExporter{exporter: fe},
		fe:     fe,
	}

	return ea
}

func (ea *FileExporterAdapters) Start(ctx context.Context) error {
	return ea.fe.Start(ctx)
}

// Shutdown shuts down the exporter and releases resources.
func (ea *FileExporterAdapters) Shutdown(ctx context.Context) error {
	if ea.stopped {
		return nil
	}

	ea.stoppedMu.Lock()
	ea.stopped = true

	var allErrs error
	err := ea.trace.Shutdown(ctx)
	if err != nil {
		allErrs = fmt.Errorf("failed to shutdown trace exporter: %w", err)
	}

	err = ea.metric.Shutdown(ctx)
	if err != nil {
		allErrs = fmt.Errorf("failed to shutdown metric exporter: %w", err)
	}

	err = ea.log.Shutdown(ctx)
	if err != nil {
		allErrs = fmt.Errorf("failed to shutdown log exporter: %w", err)
	}

	ea.fe.Shutdown(ctx)

	ea.stoppedMu.Unlock()

	return allErrs
}
