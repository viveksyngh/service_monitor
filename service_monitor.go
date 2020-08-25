package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/viveksyngh/service_monitor/client"
	"github.com/viveksyngh/service_monitor/metrics"
)

//URLs list of URLs to query

var URLs []string
var defaultURLs []string = []string{"https://httpstat.us/503", "https://httpstat.us/200"}
var defaultWorkerCount = 2

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
		results := client.QueryURLs(URLs, defaultWorkerCount)
		e.QueryResults = results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func parseURLs(urls string) []string {
	var urlList []string

	splitURLs := strings.Split(urls, ",") // Split URLs my comma

	//Parse and check for valid URLs
	for _, u := range splitURLs {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			continue
		}
		urlList = append(urlList, u.String())
	}
	return urlList
}

func main() {

	queryInterval := 15 * time.Second

	//Read URLs to monitor from environment variable
	var urls = os.Getenv("urls")

	if len(urls) == 0 {
		URLs = defaultURLs // Use default if nothing configured
	} else {
		URLs = parseURLs(urls)
	}

	fmt.Printf("URLs to monitor: %v\n", URLs)

	//Setup metric options and exporter
	metricOptions := metrics.NewMetricOptions()
	exporter := metrics.NewExporter(metricOptions)
	metrics.RegisterExporter(exporter)
	exporter.StartURLWatcher(URLs, queryInterval, defaultWorkerCount)

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
