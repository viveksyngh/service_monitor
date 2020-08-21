package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/viveksyngh/service_monitor/client"
	"github.com/viveksyngh/service_monitor/metrics"
)

//URLs list of URLs to query
var URLs []string = []string{"https://httpstat.us/503", "https://httpstat.us/200"}

//healthzHandler healthz hanlder
func healthzHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		break

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func makeQueryHandler(e *metrics.Exporter) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		results := client.QueryURLs(URLs)
		e.QueryResults = results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func main() {

	queryInterval := 15 * time.Second

	//Setup metric options and exporter
	metricOptions := metrics.NewMetricOptions()
	exporter := metrics.NewExporter(metricOptions)
	metrics.RegisterExporter(exporter)
	exporter.StartURLWatcher(URLs, queryInterval)

	//Create Router
	router := mux.NewRouter()
	port := 8080
	readTimeout := 5 * time.Second
	writeTimeout := 5 * time.Second

	//Register handlers
	router.Handle("/metrics", metrics.PrometheusHandler())
	router.HandleFunc("/", makeQueryHandler(exporter))
	router.HandleFunc("/healthz", healthzHandler)

	//Configure the HTTP server and start it
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	log.Fatal(s.ListenAndServe())
}
