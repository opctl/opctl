// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Original reference: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.111.0/exporter/prometheusexporter/collector.go
// Statement of changes: modified to collect from otlp metrics stored in a file with superfluous options removed

package telemetry

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/protobuf/proto"
)

type prometheusExporter struct {
	registry *prometheus.Registry
	mu       sync.Mutex
}

const (
	targetInfoMetricName  = "target_info"
	targetInfoDescription = "Target metadata"

	scopeInfoMetricName  = "otel_scope_info"
	scopeInfoDescription = "Instrumentation Scope metadata"

	traceIDExemplarKey = "trace_id"
	spanIDExemplarKey  = "span_id"
)

var (
	scopeInfoKeys   = [2]string{"otel_scope_name", "otel_scope_version"}
	errScopeInvalid = errors.New("invalid scope")
)

// Exporter is a Prometheus Exporter that embeds the OTel metric.Reader
// interface for easy instantiation with a MeterProvider.
type Exporter struct {
	metric.Reader
}

var _ metric.Reader = &Exporter{}

// keyVals is used to store resource attribute key value pairs.
type keyVals struct {
	keys []string
	vals []string
}

// prometheus counters MUST have a _total suffix by default:
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/compatibility/prometheus_and_openmetrics.md
const counterSuffix = "_total"

// collector is used to implement prometheus.Collector.
type collector struct {
	filesReader       *filesReader
	collectMutex      *sync.Mutex
	scopeInfos        map[pcommon.InstrumentationScope]prometheus.Metric
	scopeInfosInvalid map[pcommon.InstrumentationScope]struct{}
	metricFamilies    map[string]*dto.MetricFamily
	resourceKeyVals   keyVals
}

// New returns a Prometheus collector.
func NewPrometheusFileCollector(filesReader *filesReader) (*collector, error) {
	collector := &collector{
		collectMutex:      &sync.Mutex{},
		filesReader:       filesReader,
		scopeInfos:        make(map[pcommon.InstrumentationScope]prometheus.Metric),
		scopeInfosInvalid: make(map[pcommon.InstrumentationScope]struct{}),
		metricFamilies:    make(map[string]*dto.MetricFamily),
	}

	return collector, nil
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// The Opentelemetry SDK doesn't have information on which will exist when the collector
	// is registered. By returning nothing we are an "unchecked" collector in Prometheus,
	// and assume responsibility for consistency of the metrics produced.
	//
	// See https://pkg.go.dev/github.com/prometheus/client_golang@v1.13.0/prometheus#hdr-Custom_Collectors_and_constant_Metrics
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	otlpMetrics, err := c.filesReader.Read()
	if err != nil {
		return
	}

	c.convertOTLPToPrometheus(ch, otlpMetrics)
}

func (c *collector) convertOTLPToPrometheus(ch chan<- prometheus.Metric, allMetrics []pmetric.Metrics) {
	if len(allMetrics) == 0 {
		return
	}

	for _, metrics := range allMetrics {
		metrics.MarkReadOnly()
		resourceMetrics := metrics.ResourceMetrics()
		for rmIdx := 0; rmIdx < resourceMetrics.Len(); rmIdx++ {
			rm := resourceMetrics.At(rmIdx)
			scopeMetrics := rm.ScopeMetrics()
			for smIdx := 0; smIdx < scopeMetrics.Len(); smIdx++ {
				var keys, values [2]string
				scopeMetrics := scopeMetrics.At(smIdx)

				scope := scopeMetrics.Scope()
				scopeInfo, err := c.scopeInfo(scope)
				if errors.Is(err, errScopeInvalid) {
					// Do not report the same error multiple times.
					continue
				}
				if err != nil {
					otel.Handle(err)
					continue
				}

				ch <- scopeInfo

				keys = scopeInfoKeys
				values = [2]string{scope.Name(), scope.Version()}

				smm := scopeMetrics.Metrics()
				for smmIdx := 0; smmIdx < smm.Len(); smmIdx++ {
					m := smm.At(smmIdx)
					typ := c.metricType(m)
					if typ == nil {
						continue
					}
					name := c.getName(m, typ)

					drop, help := c.validateMetrics(name, m.Description(), typ)
					if drop {
						continue
					}

					if help != "" {
						m.SetDescription(help)
					}

					switch m.Type() {
					case pmetric.MetricTypeHistogram:
						addHistogramMetric(ch, m, keys, values, name, c.resourceKeyVals)
					case pmetric.MetricTypeSum:
						addSumMetric(ch, m, keys, values, name, c.resourceKeyVals)
					case pmetric.MetricTypeGauge:
						addGaugeMetric(ch, m, keys, values, name, c.resourceKeyVals)
					case pmetric.MetricTypeSummary:
						addSummaryMetric(ch, m, keys, values, name, c.resourceKeyVals)
					}
				}
			}
		}
	}
}

func addHistogramMetric(ch chan<- prometheus.Metric, m pmetric.Metric, ks, vs [2]string, name string, resourceKV keyVals) {
	histogram := m.Histogram()
	dataPoints := histogram.DataPoints()
	for dpIdx := 0; dpIdx < dataPoints.Len(); dpIdx++ {
		dp := dataPoints.At(dpIdx)
		keys, values := getMapAttrs(dp.Attributes(), ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description(), keys, nil)
		bounds := dp.ExplicitBounds()
		bucketCounts := dp.BucketCounts()
		buckets := make(map[float64]uint64, bounds.Len())

		cumulativeCount := uint64(0)
		for bIdx := 0; bIdx < bounds.Len(); bIdx++ {
			bound := bounds.At(bIdx)
			cumulativeCount += bucketCounts.At(bIdx)
			buckets[bound] = cumulativeCount
		}
		m, err := prometheus.NewConstHistogram(desc, dp.Count(), float64(dp.Sum()), buckets, values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		m = addExemplars(m, dp.Exemplars())
		// m = prometheus.NewMetricWithTimestamp(
		// 	dp.StartTimestamp().AsTime(),
		// 	addExemplars(m, dp.Exemplars()),
		// )
		ch <- m
	}
}

func addSumMetric(ch chan<- prometheus.Metric, m pmetric.Metric, ks, vs [2]string, name string, resourceKV keyVals) {
	sum := m.Sum()
	valueType := prometheus.CounterValue
	if !sum.IsMonotonic() {
		valueType = prometheus.GaugeValue
	}

	dataPoints := sum.DataPoints()
	for dpIdx := 0; dpIdx < dataPoints.Len(); dpIdx++ {
		dp := dataPoints.At(dpIdx)
		keys, values := getMapAttrs(dp.Attributes(), ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description(), keys, nil)

		var value float64
		switch dp.ValueType() {
		case pmetric.NumberDataPointValueTypeDouble:
			value = dp.DoubleValue()
		case pmetric.NumberDataPointValueTypeInt:
			value = float64(dp.IntValue())
		default:
			value = 0.0
		}

		m, err := prometheus.NewConstMetric(desc, valueType, value, values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		m = addExemplars(m, dp.Exemplars())
		// m = prometheus.NewMetricWithTimestamp(
		// 	dp.StartTimestamp().AsTime(),
		// 	addExemplars(m, dp.Exemplars()),
		// )
		ch <- m
	}
}

func addGaugeMetric(ch chan<- prometheus.Metric, m pmetric.Metric, ks, vs [2]string, name string, resourceKV keyVals) {
	gauge := m.Gauge()
	dataPoints := gauge.DataPoints()

	for dpIdx := 0; dpIdx < dataPoints.Len(); dpIdx++ {
		dp := dataPoints.At(dpIdx)
		keys, values := getMapAttrs(dp.Attributes(), ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description(), keys, nil)

		var value float64
		switch dp.ValueType() {
		case pmetric.NumberDataPointValueTypeDouble:
			value = dp.DoubleValue()
		case pmetric.NumberDataPointValueTypeInt:
			value = float64(dp.IntValue())
		default:
			value = 0.0
		}
		m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, value, values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		m = addExemplars(m, dp.Exemplars())
		// m = prometheus.NewMetricWithTimestamp(
		// 	dp.Timestamp().AsTime(),
		// 	addExemplars(m, dp.Exemplars()),
		// )
		ch <- m
	}
}

func addSummaryMetric(ch chan<- prometheus.Metric, m pmetric.Metric, ks, vs [2]string, name string, resourceKV keyVals) {
	summary := m.Summary()
	dataPoints := summary.DataPoints()

	for dpIdx := 0; dpIdx < dataPoints.Len(); dpIdx++ {
		dp := dataPoints.At(dpIdx)
		keys, values := getMapAttrs(dp.Attributes(), ks, vs, resourceKV)

		desc := prometheus.NewDesc(name, m.Description(), keys, nil)

		quartiles := map[float64]float64{}
		qv := dp.QuantileValues()
		for qvIdx := 0; qvIdx < qv.Len(); qvIdx++ {
			q := qv.At(qvIdx)
			quartiles[q.Quantile()] = q.Value()
		}

		m, err := prometheus.NewConstSummary(desc, dp.Count(), dp.Sum(), quartiles, values...)
		if err != nil {
			otel.Handle(err)
			continue
		}
		// m = prometheus.NewMetricWithTimestamp(
		// 	dp.StartTimestamp().AsTime(),
		// 	m,
		// )
		ch <- m
	}
}

// getMapAttrs parses the attribute.Set to two lists of matching Prometheus-style
// keys and values.
func getMapAttrs(attrs pcommon.Map, ks, vs [2]string, resourceKV keyVals) ([]string, []string) {
	keys := make([]string, 0, attrs.Len())
	values := make([]string, 0, attrs.Len())

	keysMap := make(map[string][]string)
	attrs.Range(func(k string, v pcommon.Value) bool {
		if model.NameValidationScheme == model.UTF8Validation {
			// Do not perform sanitization if prometheus supports UTF-8.
			keysMap[k] = []string{v.AsString()}
		} else {
			// It sanitizes invalid characters and handles duplicate keys
			// (due to sanitization) by sorting and concatenating the values following the spec.
			key := model.EscapeName(string(k), model.NameEscapingScheme)
			if _, ok := keysMap[key]; !ok {
				keysMap[key] = []string{v.AsString()}
			} else {
				// if the sanitized key is a duplicate, append to the list of keys
				keysMap[key] = append(keysMap[key], v.AsString())
			}
		}

		// visit all the keys and values
		return true
	})

	for key, vals := range keysMap {
		keys = append(keys, key)
		slices.Sort(vals)
		values = append(values, strings.Join(vals, ";"))
	}

	if ks[0] != "" {
		keys = append(keys, ks[:]...)
		values = append(values, vs[:]...)
	}

	for idx := range resourceKV.keys {
		keys = append(keys, resourceKV.keys[idx])
		values = append(values, resourceKV.vals[idx])
	}

	return keys, values
}

// getAttrs parses the attribute.Set to two lists of matching Prometheus-style
// keys and values.
func getAttrs(attrs attribute.Set, ks, vs [2]string, resourceKV keyVals) ([]string, []string) {
	keys := make([]string, 0, attrs.Len())
	values := make([]string, 0, attrs.Len())
	itr := attrs.Iter()

	if model.NameValidationScheme == model.UTF8Validation {
		// Do not perform sanitization if prometheus supports UTF-8.
		for itr.Next() {
			kv := itr.Attribute()
			keys = append(keys, string(kv.Key))
			values = append(values, kv.Value.Emit())
		}
	} else {
		// It sanitizes invalid characters and handles duplicate keys
		// (due to sanitization) by sorting and concatenating the values following the spec.
		keysMap := make(map[string][]string)
		for itr.Next() {
			kv := itr.Attribute()
			key := model.EscapeName(string(kv.Key), model.NameEscapingScheme)
			if _, ok := keysMap[key]; !ok {
				keysMap[key] = []string{kv.Value.Emit()}
			} else {
				// if the sanitized key is a duplicate, append to the list of keys
				keysMap[key] = append(keysMap[key], kv.Value.Emit())
			}
		}
		for key, vals := range keysMap {
			keys = append(keys, key)
			slices.Sort(vals)
			values = append(values, strings.Join(vals, ";"))
		}
	}

	if ks[0] != "" {
		keys = append(keys, ks[:]...)
		values = append(values, vs[:]...)
	}

	for idx := range resourceKV.keys {
		keys = append(keys, resourceKV.keys[idx])
		values = append(values, resourceKV.vals[idx])
	}

	return keys, values
}

func createInfoMetric(name, description string, res *resource.Resource) (prometheus.Metric, error) {
	keys, values := getAttrs(*res.Set(), [2]string{}, [2]string{}, keyVals{})
	desc := prometheus.NewDesc(name, description, keys, nil)
	return prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(1), values...)
}

func createScopeInfoMetric(scope pcommon.InstrumentationScope) (prometheus.Metric, error) {
	keys := scopeInfoKeys[:]
	desc := prometheus.NewDesc(scopeInfoMetricName, scopeInfoDescription, keys, nil)
	return prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(1), scope.Name(), scope.Version())
}

var unitSuffixes = map[string]string{
	// Time
	"d":   "_days",
	"h":   "_hours",
	"min": "_minutes",
	"s":   "_seconds",
	"ms":  "_milliseconds",
	"us":  "_microseconds",
	"ns":  "_nanoseconds",

	// Bytes
	"By":   "_bytes",
	"KiBy": "_kibibytes",
	"MiBy": "_mebibytes",
	"GiBy": "_gibibytes",
	"TiBy": "_tibibytes",
	"KBy":  "_kilobytes",
	"MBy":  "_megabytes",
	"GBy":  "_gigabytes",
	"TBy":  "_terabytes",

	// SI
	"m": "_meters",
	"V": "_volts",
	"A": "_amperes",
	"J": "_joules",
	"W": "_watts",
	"g": "_grams",

	// Misc
	"Cel": "_celsius",
	"Hz":  "_hertz",
	"1":   "_ratio",
	"%":   "_percent",
}

// getName returns the sanitized name, prefixed with the namespace and suffixed with unit.
func (c *collector) getName(m pmetric.Metric, typ *dto.MetricType) string {
	name := m.Name()
	if model.NameValidationScheme != model.UTF8Validation {
		// Only sanitize if prometheus does not support UTF-8.
		name = model.EscapeName(name, model.NameEscapingScheme)
	}
	addCounterSuffix := *typ == dto.MetricType_COUNTER
	if addCounterSuffix {
		// Remove the _total suffix here, as we will re-add the total suffix
		// later, and it needs to come after the unit suffix.
		name = strings.TrimSuffix(name, counterSuffix)
	}
	if suffix, ok := unitSuffixes[m.Unit()]; ok && !strings.HasSuffix(name, suffix) {
		name += suffix
	}
	if addCounterSuffix {
		name += counterSuffix
	}
	return name
}

func (c *collector) metricType(m pmetric.Metric) *dto.MetricType {
	switch m.Type() {
	case pmetric.MetricTypeHistogram:
		return dto.MetricType_HISTOGRAM.Enum()
	case pmetric.MetricTypeSum:
		if m.Sum().IsMonotonic() {
			return dto.MetricType_COUNTER.Enum()
		}
		return dto.MetricType_GAUGE.Enum()
	case pmetric.MetricTypeGauge:
		return dto.MetricType_GAUGE.Enum()
	}
	return nil
}

func (c *collector) createResourceAttributes(res *resource.Resource) {
	c.collectMutex.Lock()
	defer c.collectMutex.Unlock()

	resourceKeys, resourceValues := getAttrs(*res.Set(), [2]string{}, [2]string{}, keyVals{})
	c.resourceKeyVals = keyVals{keys: resourceKeys, vals: resourceValues}
}

func (c *collector) scopeInfo(scope pcommon.InstrumentationScope) (prometheus.Metric, error) {
	c.collectMutex.Lock()
	defer c.collectMutex.Unlock()

	scopeInfo, ok := c.scopeInfos[scope]
	if ok {
		return scopeInfo, nil
	}

	if _, ok := c.scopeInfosInvalid[scope]; ok {
		return nil, errScopeInvalid
	}

	scopeInfo, err := createScopeInfoMetric(scope)
	if err != nil {
		c.scopeInfosInvalid[scope] = struct{}{}
		return nil, fmt.Errorf("cannot create scope info metric: %w", err)
	}

	c.scopeInfos[scope] = scopeInfo

	return scopeInfo, nil
}

func (c *collector) validateMetrics(name, description string, metricType *dto.MetricType) (drop bool, help string) {
	c.collectMutex.Lock()
	defer c.collectMutex.Unlock()

	emf, exist := c.metricFamilies[name]

	if !exist {
		c.metricFamilies[name] = &dto.MetricFamily{
			Name: proto.String(name),
			Help: proto.String(description),
			Type: metricType,
		}
		return false, ""
	}

	if emf.GetType() != *metricType {
		// global.Error(
		// 	errors.New("instrument type conflict"),
		// 	"Using existing type definition.",
		// 	"instrument", name,
		// 	"existing", emf.GetType(),
		// 	"dropped", *metricType,
		// )
		return true, ""
	}
	if emf.GetHelp() != description {
		// global.Info(
		// 	"Instrument description conflict, using existing",
		// 	"instrument", name,
		// 	"existing", emf.GetHelp(),
		// 	"dropped", description,
		// )
		return false, emf.GetHelp()
	}

	return false, ""
}

func addExemplars(m prometheus.Metric, exemplars pmetric.ExemplarSlice) prometheus.Metric {
	exemplarsLen := exemplars.Len()
	if exemplarsLen == 0 {
		return m
	}
	promExemplars := make([]prometheus.Exemplar, exemplarsLen)

	for i := 0; i < exemplarsLen; i++ {
		e := exemplars.At(i)
		labels := attributesToLabels(e.FilteredAttributes())
		// Overwrite any existing trace ID or span ID attributes
		if traceID := e.TraceID(); !traceID.IsEmpty() {
			labels[traceIDExemplarKey] = hex.EncodeToString(traceID[:])
		}
		if spanID := e.SpanID(); !spanID.IsEmpty() {
			labels[spanIDExemplarKey] = hex.EncodeToString(spanID[:])
		}

		var value float64
		switch e.ValueType() {
		case pmetric.ExemplarValueTypeDouble:
			value = e.DoubleValue()
		case pmetric.ExemplarValueTypeInt:
			value = float64(e.IntValue())
		}

		promExemplars[i] = prometheus.Exemplar{
			Value:     value,
			Timestamp: e.Timestamp().AsTime(),
			Labels:    labels,
		}
	}
	metricWithExemplar, err := prometheus.NewMetricWithExemplars(m, promExemplars...)
	if err != nil {
		// If there are errors creating the metric with exemplars, just warn
		// and return the metric without exemplars.
		otel.Handle(err)
		return m
	}
	return metricWithExemplar
}

func attributesToLabels(attrs pcommon.Map) prometheus.Labels {
	labels := make(map[string]string)
	attrs.Range(func(k string, v pcommon.Value) bool {
		labels[k] = v.AsString()
		return true
	})
	return labels
}

func attributesToKey(attrs pcommon.Map) string {
	parts := []string{}
	attrs.Range(func(k string, v pcommon.Value) bool {
		parts = append(parts, fmt.Sprintf("%s:%s", k, v.AsString()))
		return true
	})
	slices.Sort(parts)
	return strings.Join(parts, ";")
}
