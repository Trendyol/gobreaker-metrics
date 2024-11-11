package gobreakerexporter

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/sony/gobreaker/v2"
	"gobreaker-metric/circuitbreakermetrics"
	assert "gobreaker-metric/test"
	"testing"
	"time"
)

func TestGoBreakerMetricExporter_Count(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: "test_breaker", MaxRequests: 40,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.9
		}})

	exporter := Of(cb)

	simulateSuccess(cb, 20)
	simulateFailure(cb, 10)

	count := exporter.Count()

	assert.Equal(t, uint32(30), count.Requests)
	assert.Equal(t, uint32(20), count.TotalSuccesses)
	assert.Equal(t, uint32(10), count.TotalFailures)
}

func TestGoBreakerMetricExporter_State(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: "test_breaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.9
		},
	})

	exporter := Of(cb)

	assert.Equal(t, circuitbreakermetrics.StateClosed, exporter.State())

	simulateFailure(cb, 10)

	assert.Equal(t, circuitbreakermetrics.StateOpen, exporter.State())

}

func TestGoBreakerMetricExporter_UpdateCustomMetrics(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: "test_breaker", MaxRequests: 40,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.9
		}})

	exporter := Of(cb)

	simulateSuccess(cb, 12)

	exporter.UpdateCustomMetrics()

	successCount := testutil.ToFloat64(consecutiveSuccessesCounter.WithLabelValues("test_breaker"))
	assert.Equal(t, 12.0, successCount)
}

func TestGoBreakerMetricExporter_Start(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: "test_breaker"})
	exporter := Of(cb)
	exporter.Start(time.Second)
}

func simulateSuccess(cb *gobreaker.CircuitBreaker[any], count int) {
	for i := 0; i < count; i++ {
		_, _ = cb.Execute(func() (any, error) {
			return nil, nil
		})
	}
}

func simulateFailure(cb *gobreaker.CircuitBreaker[any], count int) {
	for i := 0; i < count; i++ {
		_, _ = cb.Execute(func() (any, error) {
			return nil, errors.New("test error")
		})
	}
}
