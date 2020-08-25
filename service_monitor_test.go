package main

import "testing"

func Test_parsedURLS(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		Output []string
	}{
		{
			Name:   "empty",
			Input:  "",
			Output: []string{},
		},

		{
			Name:   "single value",
			Input:  "https://www.google.com",
			Output: []string{"https://www.google.com"},
		},
		{
			Name:   "multiple values",
			Input:  "https://www.google.com,https://www.example.com",
			Output: []string{"https://www.google.com", "https://www.example.com"},
		},
		{
			Name:   "invalid URL",
			Input:  "https://www.google.com,https://www.example.com,example2.com, ",
			Output: []string{"https://www.google.com", "https://www.example.com"},
		},
	}

	for _, testCase := range testCases {
		actualOutput := parseURLs(testCase.Input)
		if !compareSlices(actualOutput, testCase.Output) {
			t.Errorf("Testcase %s failed. Want: %v, Got: %v", testCase.Name, testCase.Output, actualOutput)
		}
	}
}

func compareSlices(t1 []string, t2 []string) bool {
	if len(t1) != len(t2) {
		return false
	}

	for i := 0; i < len(t1); i++ {
		if t1[i] != t2[i] {
			return false
		}
	}
	return true
}
