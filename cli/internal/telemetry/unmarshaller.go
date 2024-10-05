package telemetry

import (
	"errors"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type unmarshaller struct {
	tracesUnmarshaler  ptrace.Unmarshaler
	metricsUnmarshaler pmetric.Unmarshaler
	logsUnmarshaler    plog.Unmarshaler
}

func newUnmarshaller() *unmarshaller {
	return &unmarshaller{
		tracesUnmarshaler:  &ptrace.JSONUnmarshaler{},
		metricsUnmarshaler: &pmetric.JSONUnmarshaler{},
		logsUnmarshaler:    &plog.JSONUnmarshaler{},
	}
}

func (m *unmarshaller) traces(buf []byte) (*ptrace.Traces, error) {
	if m.tracesUnmarshaler == nil {
		return nil, errors.New("traces are not supported by encoding")
	}
	td, err := m.tracesUnmarshaler.UnmarshalTraces(buf)
	if err != nil {
		return nil, err
	}
	return &td, nil
}

func (m *unmarshaller) metrics(buf []byte) (*pmetric.Metrics, error) {
	if m.metricsUnmarshaler == nil {
		return nil, errors.New("metrics are not supported by encoding")
	}
	md, err := m.metricsUnmarshaler.UnmarshalMetrics(buf)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

func (m *unmarshaller) logs(buf []byte) (*plog.Logs, error) {
	if m.logsUnmarshaler == nil {
		return nil, errors.New("logs are not supported by encoding")
	}
	ld, err := m.logsUnmarshaler.UnmarshalLogs(buf)
	if err != nil {
		return nil, err
	}
	return &ld, nil
}
