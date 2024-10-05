package telemetry

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var _ = Context("otlpAccumulator", func() {
	It("should consolidate, aggregate, and collect", func() {
		otlpMetricsStrings := []string{
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159626585935000","timeUnixNano":"1728159626589295000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159626589202000","timeUnixNano":"1728159626589297000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159626589247000","timeUnixNano":"1728159626589298000","count":"1","sum":3,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159647250277000","timeUnixNano":"1728159647254903000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159647254849000","timeUnixNano":"1728159647254905000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159647254881000","timeUnixNano":"1728159647254905000","count":"1","sum":4,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159664127157000","timeUnixNano":"1728159664132456000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159664132321000","timeUnixNano":"1728159664132459000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159664132380000","timeUnixNano":"1728159664132459000","count":"1","sum":5,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159664753828000","timeUnixNano":"1728159664756862000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159664756783000","timeUnixNano":"1728159664756864000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159664756822000","timeUnixNano":"1728159664756864000","count":"1","sum":3,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159665432425000","timeUnixNano":"1728159665436168000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159665436095000","timeUnixNano":"1728159665436169000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159665436135000","timeUnixNano":"1728159665436174000","count":"1","sum":3,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
			`{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159678084770000","timeUnixNano":"1728159678087502000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159678087418000","timeUnixNano":"1728159678087504000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159678087454000","timeUnixNano":"1728159678087504000","count":"1","sum":2,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`,
		}

		unmarshaller := newUnmarshaller()
		otlpMetrics := []pmetric.Metrics{}
		for _, otlpMetric := range otlpMetricsStrings {
			metrics, err := unmarshaller.metricsUnmarshaler.UnmarshalMetrics([]byte(otlpMetric))
			if err != nil {
				panic(err)
			}
			otlpMetrics = append(otlpMetrics, metrics)
		}

		accumulated := accumulate(otlpMetrics)
		Expect(len(accumulated)).To(Equal(1))

		expectedMetricsString := `{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"opctl"}}]},"scopeMetrics":[{"scope":{"name":"opctl"},"metrics":[{"name":"opctl.command.usage","sum":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159678084770000","timeUnixNano":"1728159678087502000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.errors","sum":{"dataPoints":[{"attributes":[{"key":"error_type","value":{"stringValue":"*fmt.wrapErrors"}}],"startTimeUnixNano":"1728159678087418000","timeUnixNano":"1728159678087504000","asInt":"1"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"opctl.command.execution_time","histogram":{"dataPoints":[{"attributes":[{"key":"command","value":{"stringValue":"run"}}],"startTimeUnixNano":"1728159678087454000","timeUnixNano":"1728159678087504000","count":"1","sum":2,"bucketCounts":["0","1","0","0","0","0","0","0","0","0","0","0","0","0","0","0"],"explicitBounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000]}],"aggregationTemporality":2}}]}]}]}`
		expectedMetricsOTLP, err := unmarshaller.metricsUnmarshaler.UnmarshalMetrics([]byte(expectedMetricsString))

		Expect(err).To(BeNil())
		Expect(accumulated[0].ResourceMetrics().Len()).To(Equal(expectedMetricsOTLP.ResourceMetrics().Len()))

		resourceMetric := accumulated[0].ResourceMetrics().At(0)
		expectedMetric := expectedMetricsOTLP.ResourceMetrics().At(0)

		Expect(resourceMetric.Resource().Attributes()).To(Equal(expectedMetric.Resource().Attributes()))

		Expect(resourceMetric.ScopeMetrics().Len()).To(Equal(expectedMetric.ScopeMetrics().Len()))

		scopeMetric := resourceMetric.ScopeMetrics().At(0)
		expectedScopeMetric := expectedMetric.ScopeMetrics().At(0)

		Expect(scopeMetric.Scope().Name()).To(Equal(expectedScopeMetric.Scope().Name()))
		Expect(scopeMetric.Scope().Version()).To(Equal(expectedScopeMetric.Scope().Version()))

		metrics := scopeMetric.Metrics()
		expectedMetrics := expectedScopeMetric.Metrics()
		Expect(metrics.Len()).To(Equal(expectedMetrics.Len()))

		for e := 0; e < expectedMetrics.Len(); e++ {
			expectedMetric := expectedMetrics.At(e)

			var metric pmetric.Metric
			for m := 0; m < metrics.Len(); m++ {
				if metrics.At(m).Name() == expectedMetric.Name() && metrics.At(m).Type() == expectedMetric.Type() {
					metric = metrics.At(m)
					break
				}
			}

			Expect(metric.Name()).To(Equal(expectedMetric.Name()))
			Expect(metric.Type()).To(Equal(expectedMetric.Type()))

			switch metric.Type() {
			case pmetric.MetricTypeSum:
				Expect(metric.Sum().AggregationTemporality()).To(Equal(expectedMetric.Sum().AggregationTemporality()))
				Expect(metric.Sum().IsMonotonic()).To(Equal(expectedMetric.Sum().IsMonotonic()))
				Expect(metric.Sum().DataPoints().Len()).To(Equal(expectedMetric.Sum().DataPoints().Len()))
				for d := 0; d < metric.Sum().DataPoints().Len(); d++ {
					dataPoint := metric.Sum().DataPoints().At(d)
					expectedDataPoint := expectedMetric.Sum().DataPoints().At(d)
					Expect(dataPoint.Attributes()).To(Equal(expectedDataPoint.Attributes()))
					Expect(dataPoint.StartTimestamp()).To(Equal(expectedDataPoint.StartTimestamp()))
					Expect(dataPoint.Timestamp()).To(Equal(expectedDataPoint.Timestamp()))
					Expect(dataPoint.IntValue()).To(Equal(expectedDataPoint.IntValue()))
				}
			case pmetric.MetricTypeGauge:
				Expect(metric.Gauge().DataPoints().Len()).To(Equal(expectedMetric.Gauge().DataPoints().Len()))
				for d := 0; d < metric.Gauge().DataPoints().Len(); d++ {
					dataPoint := metric.Gauge().DataPoints().At(d)
					expectedDataPoint := expectedMetric.Gauge().DataPoints().At(d)
					Expect(dataPoint.Attributes()).To(Equal(expectedDataPoint.Attributes()))
					Expect(dataPoint.StartTimestamp()).To(Equal(expectedDataPoint.StartTimestamp()))
					Expect(dataPoint.Timestamp()).To(Equal(expectedDataPoint.Timestamp()))
					Expect(dataPoint.DoubleValue()).To(Equal(expectedDataPoint.DoubleValue()))
				}
			case pmetric.MetricTypeHistogram:
				Expect(metric.Histogram().DataPoints().Len()).To(Equal(expectedMetric.Histogram().DataPoints().Len()))
				for d := 0; d < metric.Histogram().DataPoints().Len(); d++ {
					dataPoint := metric.Histogram().DataPoints().At(d)
					expectedDataPoint := expectedMetric.Histogram().DataPoints().At(d)
					Expect(dataPoint.Attributes()).To(Equal(expectedDataPoint.Attributes()))
					Expect(dataPoint.StartTimestamp()).To(Equal(expectedDataPoint.StartTimestamp()))
					Expect(dataPoint.Timestamp()).To(Equal(expectedDataPoint.Timestamp()))
					Expect(dataPoint.Count()).To(Equal(expectedDataPoint.Count()))
					Expect(dataPoint.Sum()).To(Equal(expectedDataPoint.Sum()))
					Expect(dataPoint.BucketCounts()).To(Equal(expectedDataPoint.BucketCounts()))
					Expect(dataPoint.ExplicitBounds()).To(Equal(expectedDataPoint.ExplicitBounds()))
				}
			}
		}
	})
})
