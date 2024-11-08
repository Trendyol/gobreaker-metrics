package gobreakerexporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker/v2"
	"gobreaker-metric/circuitbreakermetrics"
	"time"
)

var (
	consecutiveSuccessesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: circuitbreakermetrics.GetMetricName("consecutive_success_total"),
			Help: "Total number of consecutive successful requests",
		},
		[]string{"name"},
	)

	consecutiveFailuresCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: circuitbreakermetrics.GetMetricName("consecutive_failure_total"),
			Help: "Total number of consecutive failed requests",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(consecutiveSuccessesCounter)
	prometheus.MustRegister(consecutiveFailuresCounter)
}

type GoBreakerMetricExporter[T any] struct {
	cb *gobreaker.CircuitBreaker[T]
}

func (m *GoBreakerMetricExporter[T]) Start(interval time.Duration) {
	go circuitbreakermetrics.RegisterGoBreakerMetrics(m, interval)
}

func Of[T any](cb *gobreaker.CircuitBreaker[T]) circuitbreakermetrics.GoBreakerMetricExporter {
	return &GoBreakerMetricExporter[T]{cb: cb}
}

func (m *GoBreakerMetricExporter[T]) Count() circuitbreakermetrics.Count {
	count := m.cb.Counts()
	return circuitbreakermetrics.Count{
		Requests:       count.Requests,
		TotalSuccesses: count.TotalSuccesses,
		TotalFailures:  count.TotalFailures,
	}
}

func (m *GoBreakerMetricExporter[T]) Name() string {
	return m.cb.Name()
}

func (m *GoBreakerMetricExporter[T]) State() circuitbreakermetrics.State {
	switch m.cb.State() {
	case gobreaker.StateOpen:
		return circuitbreakermetrics.StateOpen
	case gobreaker.StateClosed:
		return circuitbreakermetrics.StateClosed
	case gobreaker.StateHalfOpen:
		return circuitbreakermetrics.StateHalfOpen
	default:
		return circuitbreakermetrics.StateClosed
	}
}

func (m *GoBreakerMetricExporter[T]) UpdateCustomMetrics() {
	name := m.cb.Name()
	count := m.cb.Counts()
	consecutiveSuccessesCounter.WithLabelValues(name).Add(float64(count.ConsecutiveSuccesses))
	consecutiveFailuresCounter.WithLabelValues(name).Add(float64(count.ConsecutiveFailures))
}
