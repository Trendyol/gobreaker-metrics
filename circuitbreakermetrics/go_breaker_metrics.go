package circuitbreakermetrics

import "time"

type GoBreakerMetricExporter interface {
	Name() string
	State() State
	Count() Count
	UpdateCustomMetrics()
	Start(interval time.Duration)
}
