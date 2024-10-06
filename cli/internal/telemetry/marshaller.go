// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Statement of changes: modified to only support JSON marshalling without compression

package telemetry

import (
	"errors"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type marshaller struct {
	tracesMarshaler  ptrace.Marshaler
	metricsMarshaler pmetric.Marshaler
	logsMarshaler    plog.Marshaler
}

func newMarshaller() *marshaller {
	return &marshaller{
		tracesMarshaler:  &ptrace.JSONMarshaler{},
		metricsMarshaler: &pmetric.JSONMarshaler{},
		logsMarshaler:    &plog.JSONMarshaler{},
	}
}

func (m *marshaller) traces(td ptrace.Traces) ([]byte, error) {
	if m.tracesMarshaler == nil {
		return nil, errors.New("traces are not supported by encoding")
	}
	buf, err := m.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (m *marshaller) metrics(md pmetric.Metrics) ([]byte, error) {
	if m.metricsMarshaler == nil {
		return nil, errors.New("metrics are not supported by encoding")
	}
	buf, err := m.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (m *marshaller) logs(ld plog.Logs) ([]byte, error) {
	if m.logsMarshaler == nil {
		return nil, errors.New("logs are not supported by encoding")
	}
	buf, err := m.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
