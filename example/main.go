package main

import (
	"github.com/sony/gobreaker/v2"
	"gobreaker-metric/gobreakerexporter"
	"time"
)

func main() {
	circuitBreaker := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: ""})
	gobreakerexporter.Of(circuitBreaker).Start(time.Second * 1)
}
