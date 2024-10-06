package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

type TraceFileExporter struct {
	exporter  FileExporter
	stoppedMu sync.RWMutex
	stopped   bool
}

// ExportSpans is called by the OpenTelemetry Go SDK to export traces.
func (e *TraceFileExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if len(spans) == 0 {
		return nil
	}

	// Convert the spans to ptrace.Traces using the marshaller logic provided
	traceData := convertToCollectorTraces(spans)

	// Use the internal exporter (from the OpenTelemetry Collector) to send the marshaled data.
	err := e.exporter.ConsumeTraces(ctx, traceData)
	if err != nil {
		return err
	}
	return nil
}

func (e *TraceFileExporter) Shutdown(ctx context.Context) error {
	return nil
}

func convertToCollectorTraces(spans []sdktrace.ReadOnlySpan) ptrace.Traces {
	// Initialize a new ptrace.Traces object
	traceData := ptrace.NewTraces()

	// Create a new resource spans slice
	resourceSpansSlice := traceData.ResourceSpans()

	// You need to create a new ResourceSpans for each unique resource (optional)
	resourceSpans := resourceSpansSlice.AppendEmpty()

	// Create a new scope spans slice (each instrumentation scope can have multiple spans)
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()

	// Iterate over the spans and convert each one
	for _, span := range spans {
		// Append a new span to the scopeSpans
		collectorSpan := scopeSpans.Spans().AppendEmpty()

		// Set span attributes
		collectorSpan.SetName(span.Name())
		collectorSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(span.StartTime()))
		collectorSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(span.EndTime()))

		// Set the span ID and trace ID
		collectorSpan.SetSpanID(pcommon.SpanID(span.SpanContext().SpanID()))
		collectorSpan.SetTraceID(pcommon.TraceID(span.SpanContext().TraceID()))

		// Map Span Kind
		collectorSpan.SetKind(mapSpanKind(span.SpanKind()))

		// Map Status
		collectorSpan.Status().SetCode(mapSpanStatus(span.Status().Code))

		// Map Attributes
		mapAttributes(span.Attributes(), collectorSpan.Attributes())

		// Map Events
		mapEvents(span.Events(), collectorSpan.Events())

		// Map Links
		mapLinks(span.Links(), collectorSpan.Links())
	}

	return traceData
}

// Helper function to map span kind
func mapSpanKind(kind trace.SpanKind) ptrace.SpanKind {
	switch kind {
	case trace.SpanKindClient:
		return ptrace.SpanKindClient
	case trace.SpanKindServer:
		return ptrace.SpanKindServer
	case trace.SpanKindProducer:
		return ptrace.SpanKindProducer
	case trace.SpanKindConsumer:
		return ptrace.SpanKindConsumer
	case trace.SpanKindInternal:
		return ptrace.SpanKindInternal
	default:
		return ptrace.SpanKindUnspecified
	}
}

// Helper function to map status
func mapSpanStatus(status codes.Code) ptrace.StatusCode {
	switch status {
	case codes.Ok:
		return ptrace.StatusCodeOk
	case codes.Error:
		return ptrace.StatusCodeError
	case codes.Unset:
		return ptrace.StatusCodeUnset
	default:
		return ptrace.StatusCodeUnset
	}
}

// Helper function to map attributes from SDK to Collector format
func mapAttributes(srcAttributes []attribute.KeyValue, dest pcommon.Map) {
	for _, kv := range srcAttributes {
		key := string(kv.Key)
		switch kv.Value.Type() {
		case attribute.BOOL:
			dest.PutBool(key, kv.Value.AsBool())
		case attribute.INT64:
			dest.PutInt(key, kv.Value.AsInt64())
		case attribute.FLOAT64:
			dest.PutDouble(key, kv.Value.AsFloat64())
		case attribute.STRING:
			dest.PutStr(key, kv.Value.AsString())
		case attribute.BOOLSLICE:
			vSlice := pcommon.Slice{}
			for _, v := range kv.Value.AsBoolSlice() {
				vSlice.AppendEmpty().SetBool(v)
			}
			vSlice.MoveAndAppendTo(dest.PutEmptySlice(key))
		case attribute.INT64SLICE:
			vSlice := pcommon.Slice{}
			for _, v := range kv.Value.AsInt64Slice() {
				vSlice.AppendEmpty().SetInt(v)
			}
			vSlice.MoveAndAppendTo(dest.PutEmptySlice(key))
		case attribute.FLOAT64SLICE:
			vSlice := pcommon.Slice{}
			for _, v := range kv.Value.AsFloat64Slice() {
				vSlice.AppendEmpty().SetDouble(v)
			}
			vSlice.MoveAndAppendTo(dest.PutEmptySlice(key))
		case attribute.STRINGSLICE:
			vSlice := pcommon.Slice{}
			for _, v := range kv.Value.AsStringSlice() {
				vSlice.AppendEmpty().SetStr(v)
			}
			vSlice.MoveAndAppendTo(dest.PutEmptySlice(key))
		}
	}
}

// Helper function to map events from SDK to Collector format
func mapEvents(srcEvents []sdktrace.Event, dest ptrace.SpanEventSlice) {
	for _, event := range srcEvents {
		eventEntry := dest.AppendEmpty()
		eventEntry.SetName(event.Name)
		eventEntry.SetTimestamp(pcommon.NewTimestampFromTime(event.Time))
		mapAttributes(event.Attributes, eventEntry.Attributes())
	}
}

// Helper function to map links from SDK to Collector format
func mapLinks(srcLinks []sdktrace.Link, dest ptrace.SpanLinkSlice) {
	for _, link := range srcLinks {
		linkEntry := dest.AppendEmpty()
		linkEntry.SetTraceID(pcommon.TraceID(link.SpanContext.TraceID()))
		linkEntry.SetSpanID(pcommon.SpanID(link.SpanContext.SpanID()))
		mapAttributes(link.Attributes, linkEntry.Attributes())
	}
}
