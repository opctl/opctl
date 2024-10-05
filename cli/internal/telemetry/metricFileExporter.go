package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type MetricFileExporter struct {
	exporter FileExporter

	shutdownOnce sync.Once

	temporalitySelector metric.TemporalitySelector
	aggregationSelector metric.AggregationSelector
}

func (e *MetricFileExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(kind)
	// return metricdata.DeltaTemporality
}

func (e *MetricFileExporter) Aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return metric.DefaultAggregationSelector(kind)
}

func (e *MetricFileExporter) Export(ctx context.Context, metrics *metricdata.ResourceMetrics) error {
	data := convertToCollectorMetrics(metrics)

	// Use the internal exporter (from the OpenTelemetry Collector) to send the marshaled data.
	err := e.exporter.ConsumeMetrics(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *MetricFileExporter) ForceFlush(context.Context) error {
	// exporter holds no state, nothing to flush.
	return nil
}

func (e *MetricFileExporter) Shutdown(ctx context.Context) error {
	// var err error
	// e.shutdownOnce.Do(func() {
	// 	err = e.exporter.Shutdown(ctx)
	// })
	// return err
	return nil
}

func convertToCollectorMetrics(metrics *metricdata.ResourceMetrics) pmetric.Metrics {
	// Initialize a new pmetric.Metrics object
	metricData := pmetric.NewMetrics()

	// Create a new ResourceMetrics in the Collector's format
	resourceMetrics := metricData.ResourceMetrics().AppendEmpty()

	// Set the resource attributes if available
	if metrics.Resource != nil {
		mapAttributes(metrics.Resource.Attributes(), resourceMetrics.Resource().Attributes())
	}

	// Iterate over the scope metrics within the resource
	for _, scopeMetrics := range metrics.ScopeMetrics {
		// Create a new ScopeMetrics in the Collector format
		collectorScopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()

		// Set the scope information
		collectorScopeMetrics.Scope().SetName(scopeMetrics.Scope.Name)
		collectorScopeMetrics.Scope().SetVersion(scopeMetrics.Scope.Version)

		// Iterate over the individual metrics and convert them
		for _, metric := range scopeMetrics.Metrics {
			collectorMetric := collectorScopeMetrics.Metrics().AppendEmpty()
			collectorMetric.SetName(metric.Name)
			collectorMetric.SetDescription(metric.Description)
			collectorMetric.SetUnit(metric.Unit)

			// Convert based on the type of metric (Sum, Gauge, Histogram, etc.)
			switch data := metric.Data.(type) {
			case metricdata.Gauge[int64]:
				convertIntGauge(data, collectorMetric)
			case metricdata.Gauge[float64]:
				convertFloatGauge(data, collectorMetric)
			case metricdata.Sum[int64]:
				convertIntSum(data, collectorMetric)
			case metricdata.Sum[float64]:
				convertFloatSum(data, collectorMetric)
			case metricdata.Histogram[float64]:
				convertHistogram(data, collectorMetric)
				// Add more cases for other metric types as necessary
			}
		}
	}

	return metricData
}

// Helper to convert int64 Gauge data
func convertIntGauge(data metricdata.Gauge[int64], collectorMetric pmetric.Metric) {
	gauge := collectorMetric.SetEmptyGauge()
	for _, dp := range data.DataPoints {
		dataPoint := gauge.DataPoints().AppendEmpty()
		dataPoint.SetStartTimestamp(pcommon.NewTimestampFromTime(dp.StartTime))
		dataPoint.SetTimestamp(pcommon.NewTimestampFromTime(dp.Time))
		dataPoint.SetIntValue(dp.Value)
		mapAttributes(dp.Attributes.ToSlice(), dataPoint.Attributes())
	}
}

// Helper to convert float64 Gauge data
func convertFloatGauge(data metricdata.Gauge[float64], collectorMetric pmetric.Metric) {
	gauge := collectorMetric.SetEmptyGauge()
	for _, dp := range data.DataPoints {
		dataPoint := gauge.DataPoints().AppendEmpty()
		dataPoint.SetStartTimestamp(pcommon.NewTimestampFromTime(dp.StartTime))
		dataPoint.SetTimestamp(pcommon.NewTimestampFromTime(dp.Time))
		dataPoint.SetDoubleValue(dp.Value)
		mapAttributes(dp.Attributes.ToSlice(), dataPoint.Attributes())
	}
}

// Helper to convert int64 Sum data
func convertIntSum(data metricdata.Sum[int64], collectorMetric pmetric.Metric) {
	sum := collectorMetric.SetEmptySum()
	sum.SetIsMonotonic(data.IsMonotonic)
	sum.SetAggregationTemporality(mapAggregationTemporality(data.Temporality))

	for _, dp := range data.DataPoints {
		dataPoint := sum.DataPoints().AppendEmpty()
		dataPoint.SetStartTimestamp(pcommon.NewTimestampFromTime(dp.StartTime))
		dataPoint.SetTimestamp(pcommon.NewTimestampFromTime(dp.Time))
		dataPoint.SetIntValue(dp.Value)
		mapAttributes(dp.Attributes.ToSlice(), dataPoint.Attributes())
	}
}

// Helper to convert float64 Sum data
func convertFloatSum(data metricdata.Sum[float64], collectorMetric pmetric.Metric) {
	sum := collectorMetric.SetEmptySum()
	sum.SetIsMonotonic(data.IsMonotonic)
	sum.SetAggregationTemporality(mapAggregationTemporality(data.Temporality))

	for _, dp := range data.DataPoints {
		dataPoint := sum.DataPoints().AppendEmpty()
		dataPoint.SetStartTimestamp(pcommon.NewTimestampFromTime(dp.StartTime))
		dataPoint.SetTimestamp(pcommon.NewTimestampFromTime(dp.Time))
		dataPoint.SetDoubleValue(dp.Value)
		mapAttributes(dp.Attributes.ToSlice(), dataPoint.Attributes())
	}
}

// Helper to convert float64 Histogram data
func convertHistogram(data metricdata.Histogram[float64], collectorMetric pmetric.Metric) {
	hist := collectorMetric.SetEmptyHistogram()
	hist.SetAggregationTemporality(mapAggregationTemporality(data.Temporality))

	for _, dp := range data.DataPoints {
		dataPoint := hist.DataPoints().AppendEmpty()
		dataPoint.SetStartTimestamp(pcommon.NewTimestampFromTime(dp.StartTime))
		dataPoint.SetTimestamp(pcommon.NewTimestampFromTime(dp.Time))
		dataPoint.SetCount(dp.Count)
		dataPoint.SetSum(dp.Sum)
		dataPoint.BucketCounts().FromRaw(dp.BucketCounts)
		dataPoint.ExplicitBounds().FromRaw(dp.Bounds)
		mapAttributes(dp.Attributes.ToSlice(), dataPoint.Attributes())
	}
}

// Helper to map aggregation temporality
func mapAggregationTemporality(temp metricdata.Temporality) pmetric.AggregationTemporality {
	switch temp {
	case metricdata.CumulativeTemporality:
		return pmetric.AggregationTemporalityCumulative
	case metricdata.DeltaTemporality:
		return pmetric.AggregationTemporalityDelta
	default:
		return pmetric.AggregationTemporalityUnspecified
	}
}
