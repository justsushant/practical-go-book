// chap5/http-serve-mux/server_test.go
package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	tests := []struct{
		name string
		path string
		expected string
	}{
		{name: "index", path: "/api", expected: "Hello, World!"},
		{name: "healthCheck", path: "/healthz", expected: "OK"},
	}

	mux := http.NewServeMux()
	setupHandlers(mux)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.path)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			if string(respBody) != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, string(respBody))
			}
		})
	}
}