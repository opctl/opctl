package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/flock"
	promClient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/global"
	metricAPI "go.opentelemetry.io/otel/metric"
	logSDK "go.opentelemetry.io/otel/sdk/log"
	metricSDK "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	PromPushGatewayCollector = "prompush"
	// OTLPCollector            = "otlp"
	jobName                       = "opctl"
	EnvVarOpctlTelemetryCollector = "OPCTL_TELEMETRY_COLLECTOR"
	EnvVarOpctlTelemetryEndpoint  = "OPCTL_TELEMETRY_ENDPOINT"
)

type Telemetry struct {
	config        *TelemetryConfig
	fileLock      *flock.Flock
	fileExporter  FileExporter
	filesReader   *filesReader
	adapters      *FileExporterAdapters
	pusher        *push.Pusher
	traceProvider *trace.TracerProvider
	meterProvider *metricSDK.MeterProvider
	logProvider   *logSDK.LoggerProvider
}

// InstrumentCommandUsage collects usage metrics for CLI commands.
func InstrumentCommandUsage(ctx context.Context, commandName string) error {
	meter := otel.Meter(jobName)
	commandCounter, err := meter.Int64Counter(fmt.Sprintf("%s.command.usage", jobName))
	if err != nil {
		return fmt.Errorf("Error creating command usage counter: %v", err)
	}
	commandCounter.Add(ctx, 1, metricAPI.WithAttributes(attribute.String("command", commandName)))
	return nil
}

// InstrumentExecutionTime measures and records the execution time of a command.
func InstrumentExecutionTime(ctx context.Context, commandName string, startTime time.Time) error {
	meter := otel.Meter(jobName)
	execTime, err := meter.Float64Histogram(fmt.Sprintf("%s.command.execution_time", jobName))
	if err != nil {
		return fmt.Errorf("Error creating execution time histogram: %v", err)
	}
	duration := float64(time.Since(startTime).Milliseconds())
	execTime.Record(ctx, duration, metricAPI.WithAttributes(attribute.String("command", commandName)))
	return nil
}

// InstrumentErrorMetrics tracks and counts errors that occur during command execution.
func InstrumentErrorMetrics(ctx context.Context, appErr error) error {
	meter := otel.Meter(jobName)
	errorCounter, errCounter := meter.Int64Counter(fmt.Sprintf("%s.command.errors", jobName))
	if errCounter != nil {
		return fmt.Errorf("Error creating error counter: %v", errCounter)
	}
	errorCounter.Add(ctx, 1, metricAPI.WithAttributes(attribute.String("error_type", fmt.Sprintf("%T", appErr))))
	return nil
}

// Setup initializes the metrics provider and sets up the file and Prometheus PushGateway exporters.
func Setup(ctx context.Context, dataDir string) (*Telemetry, error) {
	config, err := ReadConfigFile(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	t := Telemetry{}
	t.config = config

	if t.config.Disabled() {
		// returning early will mean the global telemetry providers will be no-op
		return &t, nil
	}

	t.fileLock = flock.New(config.Path)

	t.fileExporter = NewFileExporter(config, t.fileLock)
	t.filesReader = newFilesReader(config.Path, t.fileLock)
	t.adapters = NewExporterAdapters(ctx, t.fileExporter)

	resource := resource.NewSchemaless(
		attribute.String("service.name", jobName),
	)

	// Set global providers
	if config.TracesEnabled {
		t.traceProvider = trace.NewTracerProvider(
			trace.WithSyncer(&t.adapters.trace),
			trace.WithResource(resource),
		)
		otel.SetTracerProvider(t.traceProvider)
	}

	if config.MetricsEnabled {
		t.meterProvider = metricSDK.NewMeterProvider(
			metricSDK.WithReader(
				metricSDK.NewPeriodicReader(
					&t.adapters.metric,
				),
			),
			metricSDK.WithResource(resource),
		)
		otel.SetMeterProvider(t.meterProvider)
	}

	if config.LogsEnabled {
		t.logProvider = logSDK.NewLoggerProvider(
			logSDK.WithProcessor(
				logSDK.NewSimpleProcessor(&t.adapters.log),
			),
			logSDK.WithResource(resource),
		)
		global.SetLoggerProvider(t.logProvider)
	}

	t.fileExporter.Start(ctx)

	if config.Collector == PromPushGatewayCollector && config.Endpoint != "" {
		t.pusher, err = setupPrometheusPush(config, t.filesReader)
		if err != nil {
			return &t, fmt.Errorf("failed to setup Prometheus PushGateway: %w", err)
		}
	}

	return &t, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	var err error
	if t.traceProvider != nil {
		t.traceProvider.Shutdown(ctx)
	}
	if t.meterProvider != nil {
		t.meterProvider.Shutdown(ctx)
	}
	if t.logProvider != nil {
		t.logProvider.Shutdown(ctx)
	}
	if t.adapters != nil {
		t.adapters.Shutdown(ctx)
	}
	// Attempt to push any remaining metrics before shutting down
	if t.pusher != nil {
		pushPrometheusMetrics(t.pusher, t.filesReader)
	}
	return err
}

func Push(dataDir string) error {
	cfg, err := ReadConfigFile(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if cfg.Disabled() {
		return fmt.Errorf("invalid configuration: collector, endpoint, and at least one telemetry type must be enabled")
	}

	filesReader := newFilesReader(cfg.Path, flock.New(cfg.Path))
	pusher, err := getPrometheusPusher(cfg, filesReader)
	if err != nil {
		return err
	}

	return pushPrometheusMetrics(pusher, filesReader)
}

// setupPrometheusPush sets up metrics pushing to Prometheus PushGateway.
func setupPrometheusPush(cfg *TelemetryConfig, filesReader *filesReader) (*push.Pusher, error) {
	pusher, err := getPrometheusPusher(cfg, filesReader)
	if err != nil {
		return nil, err
	}

	// Start the periodic push in a background goroutine
	go func() {
		for {
			// Push metrics to Prometheus PushGateway
			pushPrometheusMetrics(pusher, filesReader)
			// Sleep for the interval before the next push
			time.Sleep(cfg.PushInterval)
		}
	}()

	return pusher, nil
}

func getPrometheusPusher(cfg *TelemetryConfig, filesReader *filesReader) (*push.Pusher, error) {
	reg := promClient.NewRegistry()
	collector, err := NewPrometheusFileCollector(filesReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus collector: %w", err)
	}
	reg.MustRegister(collector)
	pusher := push.New(cfg.Endpoint, jobName).Gatherer(reg)
	return pusher, nil
}

func pushPrometheusMetrics(pusher *push.Pusher, filesReader *filesReader) error {
	err := pusher.Add()
	if err != nil {
		return err
	}
	filesReader.Clear()
	return nil
}
