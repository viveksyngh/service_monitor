package metrics

import "github.com/prometheus/client_golang/prometheus"

//Status status of URL
type Status int

const (
	UP   Status = 1
	DOWN Status = 0
)

//URL external URL for query
type URL struct {
	Endpoint     string
	Status       Status
	ResponseTime int
}

//Exporter a prometheus exporter
type Exporter struct {
	metricOptions MetricOptions
	urls          []URL
}

//NewExporter creates a new exporter
func NewExporter(options MetricOptions) *Exporter {
	return &Exporter{
		metricOptions: options,
	}
}

//Describe describe the metrics for prometheus
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	e.metricOptions.ExternalURLStatus.Describe(ch)
	e.metricOptions.ExternalURLResponseTime.Describe(ch)
}

//Collect collects data to be consumed by prometheus
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.metricOptions.ExternalURLStatus.Collect(ch)
	e.metricOptions.ExternalURLResponseTime.Collect(ch)
}
