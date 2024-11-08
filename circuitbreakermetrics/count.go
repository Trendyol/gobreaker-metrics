package circuitbreakermetrics

type Count struct {
	Requests       uint32
	TotalSuccesses uint32
	TotalFailures  uint32
}
