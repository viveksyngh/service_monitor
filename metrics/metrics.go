package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//MetricOptions metric options for external endpoints
type MetricOptions struct {
	ExternalURLStatus       *prometheus.GaugeVec
	ExternalURLResponseTime *prometheus.HistogramVec
}

//NewMetricOptions builds a new metric options
func NewMetricOptions() MetricOptions {
	urlStatus := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "sample",
			Subsystem: "external",
			Name:      "url_up",
			Help:      "URL status",
		},
		[]string{"url"},
	)

	urlResponseHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sample",
			Subsystem: "external",
			Name:      "url_response_time_ms",
			Help:      "URL response time in milli seconds",
		},
		[]string{"url"},
	)

	metricOptions := MetricOptions{
		ExternalURLStatus:       urlStatus,
		ExternalURLResponseTime: urlResponseHistogram,
	}

	return metricOptions
}

//PrometheusHandler prometheus metrics handler
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
