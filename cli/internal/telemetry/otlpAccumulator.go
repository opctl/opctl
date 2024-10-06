// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// original reference: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.111.0/exporter/prometheusexporter/accumulator.go
// statement of changes: modified to operate on OTLP metrics directly in a synchronous manner

package telemetry

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/collector/semconv/v1.25.0"
)

var (
	separatorString = string([]byte{model.SeparatorByte})
)

type accumulatedValue struct {
	value         pmetric.Metric // value contains a metric with exactly one aggregated datapoint.
	resourceAttrs pcommon.Map    // resourceAttrs contain the resource attributes. They are used to output instance and job labels.
	scope         pcommon.InstrumentationScope
}

func accumulate(otlpMetrics []pmetric.Metrics) []pmetric.Metrics {
	accumulatedValues := map[string]*accumulatedValue{}
	for _, metrics := range consolidate(otlpMetrics) {
		resourceMetrics := metrics.ResourceMetrics()
		for i := 0; i < resourceMetrics.Len(); i++ {
			rm := resourceMetrics.At(i)
			ilms := rm.ScopeMetrics()
			resourceAttrs := rm.Resource().Attributes()
			for i := 0; i < ilms.Len(); i++ {
				ilm := ilms.At(i)
				metrics := ilm.Metrics()
				for j := 0; j < metrics.Len(); j++ {
					addMetric(accumulatedValues, metrics.At(j), ilm.Scope(), resourceAttrs)
				}
			}
		}
	}
	return collect(accumulatedValues)
}

func consolidate(otlpMetrics []pmetric.Metrics) (result []pmetric.Metrics) {
	groupedByResourceAndScope := map[string]pmetric.Metrics{}

	for _, otlpMetric := range otlpMetrics {
		otlpMetric.MarkReadOnly()
		rm := otlpMetric.ResourceMetrics()
		for rmIdx := 0; rmIdx < rm.Len(); rmIdx++ {
			rm := rm.At(rmIdx)
			resourceAttrKey := attributesToKey(rm.Resource().Attributes())

			scopeMetrics := rm.ScopeMetrics()
			for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
				scopeMetrics := scopeMetrics.At(smIdx)
				scope := scopeMetrics.Scope()
				scopeInfo := fmt.Sprintf("%s:%s", scope.Name(), scope.Version())

				ms := scopeMetrics.Metrics()
				for msIdx := 0; msIdx < ms.Len(); msIdx++ {
					metric := ms.At(msIdx)

					key := strings.Join([]string{resourceAttrKey, scopeInfo}, ":")

					// create a new pmetric.Metrics with the same identifier
					if _, ok := groupedByResourceAndScope[key]; !ok {
						grouped := pmetric.NewMetrics()
						gRM := grouped.ResourceMetrics().AppendEmpty()
						rm.Resource().Attributes().CopyTo(gRM.Resource().Attributes())
						gScope := gRM.ScopeMetrics().AppendEmpty().Scope()
						gScope.SetName(scope.Name())
						gScope.SetVersion(scope.Version())
						groupedByResourceAndScope[key] = grouped
					}

					grouped := groupedByResourceAndScope[key]

					switch metric.Type() {
					case pmetric.MetricTypeHistogram:
						appendHistogramDataPoints(grouped, metric)
					case pmetric.MetricTypeSum:
						appendSumDataPoints(grouped, metric)
					case pmetric.MetricTypeGauge:
						appendGaugeDataPoints(grouped, metric)
					case pmetric.MetricTypeSummary:
						appendSummaryDataPoints(grouped, metric)
					}
				}
			}
		}
	}

	consolidatedMetrics := make([]pmetric.Metrics, 0, len(groupedByResourceAndScope))
	for _, m := range groupedByResourceAndScope {
		consolidatedMetrics = append(consolidatedMetrics, m)
	}
	return consolidatedMetrics
}

func addMetric(accumulatedValues map[string]*accumulatedValue, metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map) {
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		accumulateGauge(accumulatedValues, metric, il, resourceAttrs)
	case pmetric.MetricTypeSum:
		accumulateSum(accumulatedValues, metric, il, resourceAttrs)
	case pmetric.MetricTypeHistogram:
		accumulateHistogram(accumulatedValues, metric, il, resourceAttrs)
	case pmetric.MetricTypeSummary:
		accumulateSummary(accumulatedValues, metric, il, resourceAttrs)
	}
}

func accumulateGauge(accumulatedValues map[string]*accumulatedValue, metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map) {
	dps := metric.Gauge().DataPoints()
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)

		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		if ip.Flags().NoRecordedValue() {
			delete(accumulatedValues, signature)
			return
		}

		mv, ok := accumulatedValues[signature]
		if !ok {
			m := copyMetricMetadata(metric)
			ip.CopyTo(m.SetEmptyGauge().DataPoints().AppendEmpty())
			accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
			continue
		}

		if ip.Timestamp().AsTime().Before(mv.value.Gauge().DataPoints().At(0).Timestamp().AsTime()) {
			// only keep datapoint with latest timestamp
			continue
		}

		m := copyMetricMetadata(metric)
		ip.CopyTo(m.SetEmptyGauge().DataPoints().AppendEmpty())
		accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
	}
	return
}

func accumulateSum(accumulatedValues map[string]*accumulatedValue, metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map) {
	doubleSum := metric.Sum()

	// Drop metrics with unspecified aggregations
	if doubleSum.AggregationTemporality() == pmetric.AggregationTemporalityUnspecified {
		return
	}

	// Drop non-monotonic and non-cumulative metrics
	if doubleSum.AggregationTemporality() == pmetric.AggregationTemporalityDelta && !doubleSum.IsMonotonic() {
		return
	}

	dps := doubleSum.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)

		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		if ip.Flags().NoRecordedValue() {
			delete(accumulatedValues, signature)
			return
		}

		mv, ok := accumulatedValues[signature]
		if !ok {
			m := copyMetricMetadata(metric)
			m.SetEmptySum().SetIsMonotonic(metric.Sum().IsMonotonic())
			m.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
			ip.CopyTo(m.Sum().DataPoints().AppendEmpty())
			accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
			continue
		}

		if ip.Timestamp().AsTime().Before(mv.value.Sum().DataPoints().At(0).Timestamp().AsTime()) {
			// only keep datapoint with latest timestamp
			continue
		}

		// Delta-to-Cumulative
		if doubleSum.AggregationTemporality() == pmetric.AggregationTemporalityDelta && ip.StartTimestamp() == mv.value.Sum().DataPoints().At(0).Timestamp() {
			ip.SetStartTimestamp(mv.value.Sum().DataPoints().At(0).StartTimestamp())
			switch ip.ValueType() {
			case pmetric.NumberDataPointValueTypeInt:
				ip.SetIntValue(ip.IntValue() + mv.value.Sum().DataPoints().At(0).IntValue())
			case pmetric.NumberDataPointValueTypeDouble:
				ip.SetDoubleValue(ip.DoubleValue() + mv.value.Sum().DataPoints().At(0).DoubleValue())
			}
		}

		m := copyMetricMetadata(metric)
		m.SetEmptySum().SetIsMonotonic(metric.Sum().IsMonotonic())
		m.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
		ip.CopyTo(m.Sum().DataPoints().AppendEmpty())
		accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
	}
	return
}

func accumulateHistogram(accumulatedValues map[string]*accumulatedValue, metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map) {
	histogram := metric.Histogram()
	dps := histogram.DataPoints()

	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)

		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs) // uniquely identify this time series you are accumulating for
		if ip.Flags().NoRecordedValue() {
			delete(accumulatedValues, signature)
			return
		}

		mv, ok := accumulatedValues[signature] // a accumulates metric values for all times series. Get value for particular time series
		if !ok {
			// first data point
			m := copyMetricMetadata(metric)
			ip.CopyTo(m.SetEmptyHistogram().DataPoints().AppendEmpty())
			m.Histogram().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
			accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
			continue
		}

		m := copyMetricMetadata(metric)
		m.SetEmptyHistogram().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)

		switch histogram.AggregationTemporality() {
		case pmetric.AggregationTemporalityDelta:
			pp := mv.value.Histogram().DataPoints().At(0) // previous aggregated value for time range
			if ip.StartTimestamp().AsTime() != pp.Timestamp().AsTime() {
				// treat misalignment as restart and reset, or violation of single-writer principle and drop
				if ip.StartTimestamp().AsTime().After(pp.Timestamp().AsTime()) {
					ip.CopyTo(m.Histogram().DataPoints().AppendEmpty())
				} else {
					continue
				}
			} else {
				accumulateHistogramValues(pp, ip, m.Histogram().DataPoints().AppendEmpty())
			}
		case pmetric.AggregationTemporalityCumulative:
			if ip.Timestamp().AsTime().Before(mv.value.Histogram().DataPoints().At(0).Timestamp().AsTime()) {
				// only keep datapoint with latest timestamp
				continue
			}

			ip.CopyTo(m.Histogram().DataPoints().AppendEmpty())
		default:
			// unsupported temporality
			continue
		}
		accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
	}
	return
}

func accumulateSummary(accumulatedValues map[string]*accumulatedValue, metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map) {
	dps := metric.Summary().DataPoints()
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)

		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		if ip.Flags().NoRecordedValue() {
			delete(accumulatedValues, signature)
			return
		}

		mv, ok := accumulatedValues[signature]
		stalePoint := ok &&
			ip.Timestamp().AsTime().Before(mv.value.Summary().DataPoints().At(0).Timestamp().AsTime())

		if stalePoint {
			// Only keep this datapoint if it has a later timestamp.
			continue
		}

		m := copyMetricMetadata(metric)
		ip.CopyTo(m.SetEmptySummary().DataPoints().AppendEmpty())
		accumulatedValues[signature] = &accumulatedValue{value: m, resourceAttrs: resourceAttrs, scope: il}
	}

	return
}

func appendHistogramDataPoints(metrics pmetric.Metrics, metric pmetric.Metric) {
	// find the resource metric with scoped metrics that has the specified type
	rms := metrics.ResourceMetrics()
	for rmIdx := 0; rmIdx < rms.Len(); rmIdx++ {
		rm := rms.At(rmIdx)
		scopeMetrics := rm.ScopeMetrics()
		for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
			scopeMetrics := scopeMetrics.At(smIdx)
			ms := scopeMetrics.Metrics()
			for msIdx := 0; msIdx < ms.Len(); msIdx++ {
				m := ms.At(msIdx)
				if m.Name() == metric.Name() {
					histogramDataPoints := m.Histogram().DataPoints()
					for dpIdx := 0; dpIdx < metric.Histogram().DataPoints().Len(); dpIdx++ {
						metric.Histogram().DataPoints().At(dpIdx).CopyTo(histogramDataPoints.AppendEmpty())
					}
					return
				}
			}
		}
	}

	// create a new resource metric with scoped metrics
	// resourceMetrics and scopeMetrics are already set correctly
	metric.CopyTo(metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty())
}

func appendSumDataPoints(metrics pmetric.Metrics, metric pmetric.Metric) {
	// find the resource metric with scoped metrics that has the specified type
	rms := metrics.ResourceMetrics()
	for rmIdx := 0; rmIdx < rms.Len(); rmIdx++ {
		rm := rms.At(rmIdx)
		scopeMetrics := rm.ScopeMetrics()
		for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
			scopeMetrics := scopeMetrics.At(smIdx)
			ms := scopeMetrics.Metrics()
			for msIdx := 0; msIdx < ms.Len(); msIdx++ {
				m := ms.At(msIdx)
				if m.Name() == metric.Name() {
					sumDataPoints := m.Sum().DataPoints()
					for dpIdx := 0; dpIdx < metric.Sum().DataPoints().Len(); dpIdx++ {
						metric.Sum().DataPoints().At(dpIdx).CopyTo(sumDataPoints.AppendEmpty())
					}
					return
				}
			}
		}
	}

	// create a new resource metric with scoped metrics
	// resourceMetrics and scopeMetrics are already set correctly
	metric.CopyTo(metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty())
}

func appendGaugeDataPoints(metrics pmetric.Metrics, metric pmetric.Metric) {
	// find the resource metric with scoped metrics that has the specified type
	rms := metrics.ResourceMetrics()
	for rmIdx := 0; rmIdx < rms.Len(); rmIdx++ {
		rm := rms.At(rmIdx)
		scopeMetrics := rm.ScopeMetrics()
		for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
			scopeMetrics := scopeMetrics.At(smIdx)
			ms := scopeMetrics.Metrics()
			for msIdx := 0; msIdx < ms.Len(); msIdx++ {
				m := ms.At(msIdx)
				if m.Name() == metric.Name() {
					gaugeDataPoints := m.Gauge().DataPoints()
					for dpIdx := 0; dpIdx < metric.Gauge().DataPoints().Len(); dpIdx++ {
						metric.Gauge().DataPoints().At(dpIdx).CopyTo(gaugeDataPoints.AppendEmpty())
					}
					return
				}
			}
		}
	}

	// create a new resource metric with scoped metrics
	// resourceMetrics and scopeMetrics are already set correctly
	metric.CopyTo(metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty())
}

func appendSummaryDataPoints(metrics pmetric.Metrics, metric pmetric.Metric) {
	// find the resource metric with scoped metrics that has the specified type
	rms := metrics.ResourceMetrics()
	for rmIdx := 0; rmIdx < rms.Len(); rmIdx++ {
		rm := rms.At(rmIdx)
		scopeMetrics := rm.ScopeMetrics()
		for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
			scopeMetrics := scopeMetrics.At(smIdx)
			ms := scopeMetrics.Metrics()
			for msIdx := 0; msIdx < ms.Len(); msIdx++ {
				m := ms.At(msIdx)
				if m.Name() == metric.Name() {
					summaryDataPoints := m.Summary().DataPoints()
					for dpIdx := 0; dpIdx < metric.Summary().DataPoints().Len(); dpIdx++ {
						metric.Summary().DataPoints().At(dpIdx).CopyTo(summaryDataPoints.AppendEmpty())
					}
					return
				}
			}
		}
	}

	// create a new resource metric with scoped metrics
	// resourceMetrics and scopeMetrics are already set correctly
	metric.CopyTo(metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty())
}

func timeseriesSignature(ilmName string, metric pmetric.Metric, attributes pcommon.Map, resourceAttrs pcommon.Map) string {
	var b strings.Builder
	b.WriteString(metric.Type().String())
	b.WriteString("*" + ilmName)
	b.WriteString("*" + metric.Name())
	attrs := make([]string, 0, attributes.Len())
	attributes.Range(func(k string, v pcommon.Value) bool {
		attrs = append(attrs, k+"*"+v.AsString())
		return true
	})
	sort.Strings(attrs)
	b.WriteString("*" + strings.Join(attrs, "*"))
	if job, ok := extractJob(resourceAttrs); ok {
		b.WriteString("*" + model.JobLabel + "*" + job)
	}
	if instance, ok := extractInstance(resourceAttrs); ok {
		b.WriteString("*" + model.InstanceLabel + "*" + instance)
	}
	return b.String()
}

func copyMetricMetadata(metric pmetric.Metric) pmetric.Metric {
	m := pmetric.NewMetric()
	m.SetName(metric.Name())
	m.SetDescription(metric.Description())
	m.SetUnit(metric.Unit())

	return m
}

func accumulateHistogramValues(prev, current, dest pmetric.HistogramDataPoint) {
	dest.SetStartTimestamp(prev.StartTimestamp())

	older := prev
	newer := current
	if current.Timestamp().AsTime().Before(prev.Timestamp().AsTime()) {
		older = current
		newer = prev
	}

	newer.Attributes().CopyTo(dest.Attributes())
	dest.SetTimestamp(newer.Timestamp())

	// checking for bucket boundary alignment, optionally re-aggregate on newer boundaries
	match := older.ExplicitBounds().Len() == newer.ExplicitBounds().Len()
	for i := 0; match && i < newer.ExplicitBounds().Len(); i++ {
		match = older.ExplicitBounds().At(i) == newer.ExplicitBounds().At(i)
	}

	if match {

		dest.SetCount(newer.Count() + older.Count())
		dest.SetSum(newer.Sum() + older.Sum())

		counts := make([]uint64, newer.BucketCounts().Len())
		for i := 0; i < newer.BucketCounts().Len(); i++ {
			counts[i] = newer.BucketCounts().At(i) + older.BucketCounts().At(i)
		}
		dest.BucketCounts().FromRaw(counts)
	} else {
		// use new value if bucket bounds do not match
		dest.SetCount(newer.Count())
		dest.SetSum(newer.Sum())
		dest.BucketCounts().FromRaw(newer.BucketCounts().AsRaw())
	}

	dest.ExplicitBounds().FromRaw(newer.ExplicitBounds().AsRaw())
}

func resourceSignature(attributes pcommon.Map) string {
	job, _ := extractJob(attributes)
	instance, _ := extractInstance(attributes)
	if job == "" || instance == "" {
		return ""
	}

	return job + separatorString + instance
}

func extractInstance(attributes pcommon.Map) (string, bool) {
	// Map service.instance.id to instance
	if inst, ok := attributes.Get(conventions.AttributeServiceInstanceID); ok {
		return inst.AsString(), true
	}
	return "", false
}

func extractJob(attributes pcommon.Map) (string, bool) {
	// Map service.name + service.namespace to job
	if serviceName, ok := attributes.Get(conventions.AttributeServiceName); ok {
		job := serviceName.AsString()
		if serviceNamespace, ok := attributes.Get(conventions.AttributeServiceNamespace); ok {
			job = fmt.Sprintf("%s/%s", serviceNamespace.AsString(), job)
		}
		return job, true
	}
	return "", false
}

// Collect returns a slice with relevant aggregated metrics and their resource attributes.
func collect(accumulatedValues map[string]*accumulatedValue) []pmetric.Metrics {
	// group the accumulatedValues by resource attributes and scope (identifiers)
	groups := map[pcommon.Map]map[pcommon.InstrumentationScope][]*accumulatedValue{}
	for _, v := range accumulatedValues {
		if _, ok := groups[v.resourceAttrs]; !ok {
			groups[v.resourceAttrs] = map[pcommon.InstrumentationScope][]*accumulatedValue{}
		}
		if _, ok := groups[v.resourceAttrs][v.scope]; !ok {
			groups[v.resourceAttrs][v.scope] = []*accumulatedValue{}
		}
		groups[v.resourceAttrs][v.scope] = append(groups[v.resourceAttrs][v.scope], v)
	}

	// convert the accumulatedValues to pmetric.Metrics again
	metrics := []pmetric.Metrics{}
	for resourceAttrs, scopeValues := range groups {
		for scope, values := range scopeValues {
			metric := pmetric.NewMetrics()
			resourceMetrics := metric.ResourceMetrics().AppendEmpty()
			resourceAttrs.CopyTo(resourceMetrics.Resource().Attributes())
			scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()
			scope.CopyTo(scopeMetrics.Scope())
			for _, v := range values {
				v.value.CopyTo(scopeMetrics.Metrics().AppendEmpty())
			}
			metrics = append(metrics, metric)
		}
	}

	return metrics
}

type withFormatTime struct {
	t time.Time
	f fs.DirEntry
}

type byFormatTime []withFormatTime

func (b byFormatTime) Less(i, j int) bool {
	return b[i].t.After(b[j].t)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
