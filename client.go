package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

//URLs list of URLs to query
var URLs []string = []string{"https://httpstat.us/503", "https://httpstat.us/200"}

//QueryResult result of the URL query
type QueryResult struct {
	URL          string
	Status       int
	ResponseTime time.Duration
}

//QueryURLs query urls and return result
func QueryURLs(urls []string) []QueryResult {
	client := NewClient()
	// queryResults := make(chan QueryResult, len(urls))
	// urlsChan := make(chan string, len(urls))

	// var wg sync.WaitGroup
	var results []QueryResult

	for _, url := range urls {
		// wg.Add(1)
		// go func() {
		// defer wg.Done()
		// for url := range urlsChan {

		req, _ := http.NewRequest(http.MethodGet, url, nil)

		start := time.Now()

		res, err := client.Do(req)

		duration := time.Since(start)
		if err != nil {
			log.Println(err)
		}

		var status int
		if res.StatusCode == http.StatusOK {
			status = 1
		} else {
			status = 0
		}

		result := QueryResult{
			URL:          url,
			Status:       status,
			ResponseTime: duration,
		}

		results = append(results, result)
		// queryResults <- result

		// }
		// }()
	}

	// for _, url := range urls {
	// 	urlsChan <- url
	// }

	// close(urlsChan)

	// wg.Wait()

	// for queryResult := range queryResults {
	// 	result = append(result, queryResult)
	// }

	return results
}

//NewClient returns a new http client
func NewClient() *http.Client {

	timeout := 3 * time.Second
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 0,
			}).DialContext,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
	}

	return client
}
