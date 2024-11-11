# gobreaker-metrics

This library enables you to monitor the state of your [sony/gobreaker](https://github.com/sony/gobreaker) circuit breaker in real-time using Prometheus. It provides essential metrics like circuit state, failure rate, success rate, and others to track the health and performance of your circuit breakers.

The library is inspired by [resilience4j-micrometer](https://resilience4j.readme.io/docs/micrometer) and acts as a similar tool but specifically designed for sony/gobreaker. Currently, no other metric library exists for gobreaker, making this tool a unique addition to the monitoring ecosystem.

### Features
- Monitor Circuit State: Track if the circuit breaker is in a closed, open, or half-open state.
- Track Success and Failure Rates: Understand how your circuit breaker performs under different load conditions.
- Prometheus Integration: Seamlessly export metrics to Prometheus for powerful visualization and alerting.
- Lightweight and Easy-to-Integrate: Designed to work with your existing Go-Sony Breaker setup without much overhead.
Getting Started

### Prerequisites
- gobreaker: Ensure that your project is already using sony/gobreaker.
- Prometheus: This library is built to export metrics to Prometheus. You need a Prometheus server running to scrape the metrics.


### Installation
You can add the library to your Go module by running:

    go get github.com/Trendyol/gobreaker-metrics

### Basic Setup
Below is a sample code snippet to demonstrate how to start monitoring your Go-Sony Breaker circuit breaker using this library.

```go
package main

import (
	"github.com/sony/gobreaker/v2"
	"gobreaker-metric/gobreakerexporter" "time")

func main() {

	circuitBreaker := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{Name: ""})   # Step1

	goBreakerMetricExporter := gobreakerexporter.Of(circuitBreaker)  # Step2
	goBreakerMetricExporter.Start(time.Second * 1)  # Step3

}
```

### Explanation of Code
- Step 1: Configure and initialize a circuit breaker using sony/goreaker settings.
- Step 2: Create a new MetricsExporter instance from your circuit breaker. This will automatically track metrics for the breaker.
- Step 3: Start metric exporter. The metrics on an HTTP endpoint.

### Prometheus Configuration
To scrape the metrics from the exporter, make sure prometheus '**/metrics**' endpoint active. Add the following Prometheus configuration to your code if you use the Gin web framework:

```go
router.GET("/metrics", prometheusHandler())

func prometheusHandler() gin.HandlerFunc {  
    h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true})  
  
    return func(c *gin.Context) {  
       h.ServeHTTP(c.Writer, c.Request)  
    }  
}
```

### Available Metrics
The following metrics are exposed by the library:

- gobreaker_state: Indicates the current state of the circuit breaker.
- gobreaker_requests_total: Total number of requests made through the circuit breaker.
- gobreaker_failure_total: Total number of failed requests handled by the circuit breaker.
- gobreaker_success_total: Total number of successful requests.

These metrics provide a comprehensive view of the circuit breaker's health and allow you to set up alerting rules in Prometheus to be notified of potential issues.

### Grafana

[Here](https://gist.github.com/osemrt/92b92329067e2dd2a633f9607a3d4460) is a pre-built Grafana dashboard that can be easily imported to get you started quickly.

![img.png](https://github.com/user-attachments/assets/6fd79cf0-1a8e-436f-aee5-ae2686816dac)

![img.png](https://github.com/user-attachments/assets/30cfb9e9-3540-4030-81fc-f86cedb0bf21)


### Troubleshooting
If you encounter issues, here are a few common checks:

- Verify Prometheus Target: Ensure Prometheus can access the /metrics endpoint.
- Check Circuit Breaker Configurations: Ensure your Go-Sony Breaker is set up correctly.
- Debug Logs: Enable debug logs to get more insights into the metrics exporter.

### Contributing
We welcome contributions! Please open issues or pull requests for any bugs, feature requests, or improvements.

### License
This project is licensed under the MIT License.

