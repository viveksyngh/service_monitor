package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/viveksyngh/service_monitor/metrics"
)

func main() {
	router := mux.NewRouter()
	port := 8082
	readTimeout := 5 * time.Second
	writeTimeout := 5 * time.Second

	router.Handle("/metrics", metrics.PrometheusHandler())

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	log.Fatal(s.ListenAndServe())
}
