package circuitbreakermetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const MetricPrefix = "gobreaker"

func GetMetricName(metricName string) string {
	return fmt.Sprintf("%s_%s", MetricPrefix, metricName)
}

var (
	circuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: GetMetricName("state"),
			Help: "The states of the circuit breaker. 0=Not Active, 1=Active. state=['open','half-open','closed']",
		},
		[]string{"name", "state"},
	)

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: GetMetricName("requests_total"),
			Help: "Total number of requests executed through the circuit breaker",
		},
		[]string{"name"},
	)

	successCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: GetMetricName("success_total"),
			Help: "Total number of successful requests",
		},
		[]string{"name"},
	)

	failureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: GetMetricName("failure_total"),
			Help: "Total number of failed requests",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(circuitBreakerState)
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(successCounter)
	prometheus.MustRegister(failureCounter)
}

func RegisterGoBreakerMetrics(metrics GoBreakerMetricExporter, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		updateMetrics(metrics)
	}
}

func updateMetrics(metrics GoBreakerMetricExporter) {
	state := metrics.State()
	name := metrics.Name()
	count := metrics.Count()

	circuitBreakerState.WithLabelValues(name, StateOpen.String()).Set(state.Matches(StateOpen))
	circuitBreakerState.WithLabelValues(name, StateClosed.String()).Set(state.Matches(StateClosed))
	circuitBreakerState.WithLabelValues(name, StateHalfOpen.String()).Set(state.Matches(StateHalfOpen))

	requestCounter.WithLabelValues(name).Add(float64(count.Requests))
	successCounter.WithLabelValues(name).Add(float64(count.TotalSuccesses))
	failureCounter.WithLabelValues(name).Add(float64(count.TotalFailures))

	metrics.UpdateCustomMetrics()
}
