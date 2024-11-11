package circuitbreakermetrics

import (
	"github.com/prometheus/client_golang/prometheus/testutil"
	assert "gobreaker-metric/test"
	"testing"
	"time"
)

type MockGoBreakerMetricExporter struct {
	state       State
	name        string
	requests    uint32
	successes   uint32
	failures    uint32
	customCalls int
}

func (m *MockGoBreakerMetricExporter) Start(time.Duration) {
}

func (m *MockGoBreakerMetricExporter) State() State {
	return m.state
}

func (m *MockGoBreakerMetricExporter) Name() string {
	return m.name
}

func (m *MockGoBreakerMetricExporter) Count() Count {
	return Count{
		Requests:       m.requests,
		TotalSuccesses: m.successes,
		TotalFailures:  m.failures,
	}
}

func (m *MockGoBreakerMetricExporter) UpdateCustomMetrics() {
	m.customCalls++
}

func TestGetMetricName(t *testing.T) {
	metricName := GetMetricName("requests_total")
	expectedName := "gobreaker_requests_total"
	assert.Equal(t, expectedName, metricName)
}

func TestUpdateMetrics(t *testing.T) {
	mockMetrics := &MockGoBreakerMetricExporter{
		state:     StateOpen,
		name:      "test_circuit",
		requests:  5,
		successes: 3,
		failures:  2,
	}

	updateMetrics(mockMetrics)

	openState := testutil.ToFloat64(circuitBreakerState.WithLabelValues("test_circuit", StateOpen.String()))
	closedState := testutil.ToFloat64(circuitBreakerState.WithLabelValues("test_circuit", StateClosed.String()))
	halfOpenState := testutil.ToFloat64(circuitBreakerState.WithLabelValues("test_circuit", StateHalfOpen.String()))
	assert.Equal(t, 1.0, openState)
	assert.Equal(t, 0.0, closedState)
	assert.Equal(t, 0.0, halfOpenState)

	requests := testutil.ToFloat64(requestCounter.WithLabelValues("test_circuit"))
	assert.Equal(t, 5.0, requests)

	successes := testutil.ToFloat64(successCounter.WithLabelValues("test_circuit"))
	assert.Equal(t, 3.0, successes)

	failures := testutil.ToFloat64(failureCounter.WithLabelValues("test_circuit"))
	assert.Equal(t, 2.0, failures)

	assert.Equal(t, 1, mockMetrics.customCalls)
}

func TestRegisterGoBreakerMetrics(t *testing.T) {
	mockMetrics := &MockGoBreakerMetricExporter{
		state:     StateClosed,
		name:      "test_circuit",
		requests:  10,
		successes: 8,
		failures:  2,
	}

	go func() {
		RegisterGoBreakerMetrics(mockMetrics, 100*time.Millisecond)
	}()

	time.Sleep(300 * time.Millisecond)

	requests := testutil.ToFloat64(requestCounter.WithLabelValues("test_circuit"))
	assert.GreaterOrEqual(t, 10.0, requests)

	successes := testutil.ToFloat64(successCounter.WithLabelValues("test_circuit"))
	assert.GreaterOrEqual(t, 8.0, successes)

	failures := testutil.ToFloat64(failureCounter.WithLabelValues("test_circuit"))
	assert.GreaterOrEqual(t, 2.0, failures)

	closedState := testutil.ToFloat64(circuitBreakerState.WithLabelValues("test_circuit", StateClosed.String()))
	assert.Equal(t, closedState, 1.0)
}
