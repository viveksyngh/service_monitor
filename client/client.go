package client

import (
	"log"
	"net"
	"net/http"
	"time"
)

//Status status of URL
type Status int

const (
	UP   Status = 1
	DOWN Status = 0
)

//QueryResult result of the URL query
type QueryResult struct {
	URL          string        `json:"url"`
	Status       Status        `json:"status"`
	ResponseTime time.Duration `json:"response_time"`
}

//QueryURLs query urls and return result
func QueryURLs(client *http.Client, urls []string, workerCount int) []QueryResult {

	queryResults := make(chan QueryResult)
	urlsChan := make(chan string)

	var results []QueryResult

	for i := 0; i < workerCount; i++ {
		go queryURL(client, urlsChan, queryResults)
	}

	for _, url := range urls {
		urlsChan <- url
	}

	close(urlsChan)

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

		duration := time.Since(start)
		if err != nil {
			log.Println(err.Error())
			result <- QueryResult{
				URL:          url,
				Status:       DOWN,
				ResponseTime: duration,
			}
			continue
		}

		_ = res.Body.Close()

		var status Status
		if res.StatusCode == http.StatusOK {
			status = UP
		} else {
			status = DOWN
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

	timeout := 5 * time.Second
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
