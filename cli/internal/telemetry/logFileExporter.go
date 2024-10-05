package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/otel/log"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

type LogFileExporter struct {
	exporter  FileExporter
	stoppedMu sync.RWMutex
	stopped   bool
}

func (e *LogFileExporter) Export(ctx context.Context, records []sdkLog.Record) error {
	if len(records) == 0 {
		return nil
	}

	// Convert the logs to pdata.Logs using the marshaller logic provided
	logData := convertToCollectorLogs(records)

	// Use the internal exporter (from the OpenTelemetry Collector) to send the marshaled data.
	err := e.exporter.ConsumeLogs(ctx, logData)
	if err != nil {
		return err
	}
	return nil
}

// Shutdown shuts down the Exporter.
// Calls to Export will perform no operation after this is called.
func (e *LogFileExporter) Shutdown(ctx context.Context) error {
	// var err error
	// e.stoppedMu.Lock()
	// err = e.exporter.Shutdown(ctx)
	// e.stopped = true
	// e.stoppedMu.Unlock()
	// return err
	return nil
}

// ForceFlush performs no action.
func (e *LogFileExporter) ForceFlush(context.Context) error {
	return nil
}

func convertToCollectorLogs(records []sdkLog.Record) plog.Logs {
	// Initialize a new plog.Logs object
	logData := plog.NewLogs()

	// Create a new ResourceLogs in the Collector's format
	resourceLogs := logData.ResourceLogs().AppendEmpty()

	// Optionally set resource attributes, if logs are grouped by resources
	// Example: resourceLogs.Resource().Attributes().InsertString("service.name", "example-service")

	// Create a new ScopeLogs slice to hold logs grouped by instrumentation scope
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	// Optionally set scope name and version (e.g., the name and version of the logger)
	// Example: scopeLogs.Scope().SetName("example-logger")
	// Example: scopeLogs.Scope().SetVersion("v1.0.0")

	// Iterate over each log record and convert to plog.LogRecord
	for _, record := range records {
		// Create a new log record in the Collector format
		logRecord := scopeLogs.LogRecords().AppendEmpty()

		// Set the log message (stored in `Body`)
		logRecord.Body().SetStr(record.Body().AsString())

		// Set the timestamp of the log record
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(record.Timestamp()))

		// Set the severity level
		logRecord.SetSeverityNumber(mapSeverity(record.Severity()))
		logRecord.SetSeverityText(record.SeverityText())

		// Set attributes from the `attribute.Set` in the SDK format to `pcommon.Map`
		destAttributes := logRecord.Attributes()
		record.WalkAttributes(func(kv log.KeyValue) bool {
			mapAttributeLogKV(kv, destAttributes)
			return true
		})

		// Optionally, set additional fields like TraceID, SpanID, etc., if available
		if record.TraceID().IsValid() {
			logRecord.SetTraceID(pcommon.TraceID(record.TraceID()))
		}
		if record.SpanID().IsValid() {
			logRecord.SetSpanID(pcommon.SpanID(record.SpanID()))
		}

		// Optionally, set flags or other fields if your log.Record contains them
	}

	return logData
}

// Helper to map severity from SDK `log.Severity` to Collector's `plog.SeverityNumber`
func mapSeverity(sev log.Severity) plog.SeverityNumber {
	switch sev {
	case log.SeverityTrace:
		return plog.SeverityNumberTrace
	case log.SeverityDebug:
		return plog.SeverityNumberDebug
	case log.SeverityInfo:
		return plog.SeverityNumberInfo
	case log.SeverityWarn:
		return plog.SeverityNumberWarn
	case log.SeverityError:
		return plog.SeverityNumberError
	case log.SeverityFatal:
		return plog.SeverityNumberFatal
	default:
		return plog.SeverityNumberUnspecified
	}
}

func mapAttributeLogKV(kv log.KeyValue, dest pcommon.Map) {
	switch kv.Value.Kind() {
	case log.KindBool:
		dest.PutBool(kv.Key, kv.Value.AsBool())
	case log.KindInt64:
		dest.PutInt(kv.Key, kv.Value.AsInt64())
	case log.KindFloat64:
		dest.PutDouble(kv.Key, kv.Value.AsFloat64())
	case log.KindString:
		dest.PutStr(kv.Key, kv.Value.AsString())
	case log.KindBytes:
		dest.PutEmptyBytes(kv.Key).Append(kv.Value.AsBytes()...)
	case log.KindEmpty:
		dest.PutEmpty(kv.Key)
	case log.KindMap:
		vMap := dest.PutEmptyMap(kv.Key)
		for _, subKV := range kv.Value.AsMap() {
			mapAttributeLogKV(subKV, vMap)
		}
	case log.KindSlice:
		vSlice := pcommon.Slice{}
		for _, v := range kv.Value.AsSlice() {
			vSlice.AppendEmpty().FromRaw(logValueRaw(v))
		}
		vSlice.MoveAndAppendTo(dest.PutEmptySlice(kv.Key))
	}
}

func logValueRaw(value log.Value) any {
	switch value.Kind() {
	case log.KindBool:
		return value.AsBool()
	case log.KindInt64:
		return value.AsInt64()
	case log.KindFloat64:
		return value.AsFloat64()
	case log.KindString:
		return value.AsString()
	case log.KindBytes:
		return value.AsBytes()
	case log.KindEmpty:
		return nil
	case log.KindMap:
		vMap := pcommon.Map{}
		for _, subKV := range value.AsMap() {
			mapAttributeLogKV(subKV, vMap)
		}
		return vMap
	case log.KindSlice:
		vSlice := pcommon.Slice{}
		for _, v := range value.AsSlice() {
			vSlice.AppendEmpty().FromRaw(v)
		}
		return vSlice
	}
	return nil
}
