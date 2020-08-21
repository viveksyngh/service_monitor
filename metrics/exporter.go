package metrics

import (
	"log"
	"sync"
	"time"

	"github.com/viveksyngh/service_monitor/client"

	"github.com/prometheus/client_golang/prometheus"
)

//Exporter a prometheus exporter
type Exporter struct {
	metricOptions MetricOptions
	QueryResults  []client.QueryResult
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

	for _, queryResult := range e.QueryResults {
		e.metricOptions.ExternalURLStatus.
			WithLabelValues(queryResult.URL).
			Set(float64(queryResult.Status))

		e.metricOptions.ExternalURLResponseTime.
			WithLabelValues(queryResult.URL).
			Observe(queryResult.ResponseTime.Seconds())
	}
	e.metricOptions.ExternalURLStatus.Collect(ch)
	e.metricOptions.ExternalURLResponseTime.Collect(ch)
}

//For syncronisation to make sure MustRegister only called once
var once = sync.Once{}

//RegisterExporter registers with prometheus for tracking
func RegisterExporter(e *Exporter) {
	once.Do(func() {
		prometheus.MustRegister(e)
	})
}

//StartURLWatcher start worker to watch URLS
func (e *Exporter) StartURLWatcher(urls []string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Querying URLs: %v\n", urls)
				queryResults := client.QueryURLs(urls)
				e.QueryResults = queryResults
				break
			case <-quit:
				return
			}
		}
	}()
}
