package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_QueryURLs_200(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	}))

	defer ts.Close()

	expectedResult := []QueryResult{
		{
			URL:          ts.URL,
			Status:       UP,
			ResponseTime: time.Second * 1,
		},
	}

	results := QueryURLs(ts.Client(), []string{ts.URL}, 1)

	if !compareQueryResult(expectedResult, results) {
		t.Errorf("Failed, Want: %v, got: %v", expectedResult, results)
	}

}

func Test_QueryURLs_Non200(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))

	defer ts.Close()

	expectedResult := []QueryResult{
		{
			URL:          ts.URL,
			Status:       DOWN,
			ResponseTime: time.Second * 1,
		},
	}

	results := QueryURLs(ts.Client(), []string{ts.URL}, 1)

	if !compareQueryResult(expectedResult, results) {
		t.Errorf("Failed, Want: %v, got: %v", expectedResult, results)
	}

}

func compareQueryResult(expectedResults, actualResults []QueryResult) bool {
	if len(expectedResults) != len(actualResults) {
		return false
	}
	for i := 0; i < len(expectedResults); i++ {
		//Ignoring response time check as it is not deterministic
		if expectedResults[i].URL != actualResults[i].URL || expectedResults[i].Status != actualResults[i].Status {
			return false
		}
	}
	return true
}
