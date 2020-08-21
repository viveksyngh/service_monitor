package client

import (
	"log"
	"net"
	"net/http"
	"time"
)

//QueryResult result of the URL query
type QueryResult struct {
	URL          string
	Status       int
	ResponseTime time.Duration
}

//QueryURLs query urls and return result
func QueryURLs(urls []string) []QueryResult {
	client := NewClient()
	queryResults := make(chan QueryResult)
	urlsChan := make(chan string)

	// var wg sync.WaitGroup
	var results []QueryResult

	for i := 0; i < 2; i++ {
		go queryURL(client, urlsChan, queryResults)
	}

	for _, url := range urls {
		urlsChan <- url
	}

	close(urlsChan)

	// wg.Wait()

	for i := 0; i < len(urls); i++ {
		results = append(results, <-queryResults)
	}

	return results
}

func queryURL(client *http.Client, input <-chan string, result chan<- QueryResult) {
	for url := range input {
		req, _ := http.NewRequest(http.MethodGet, url, nil)

		start := time.Now()

		res, err := client.Do(req)

		_ = res.Body.Close()

		duration := time.Since(start)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var status int
		if res.StatusCode == http.StatusOK {
			status = 1
		} else {
			status = 0
		}

		result <- QueryResult{
			URL:          url,
			Status:       status,
			ResponseTime: duration,
		}
	}
}

//NewClient returns a new http client
func NewClient() *http.Client {

	timeout := 10 * time.Second
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
