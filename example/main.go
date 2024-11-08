package example

import (
	"github.com/sony/gobreaker/v2"
	"gobreaker-metric/gobreakerexporter"
	"time"
)

func main() {

	circuitBreaker := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: ""})

	goBreakerMetricExporter := gobreakerexporter.Of(circuitBreaker)
	goBreakerMetricExporter.Start(time.Second * 1)

}
