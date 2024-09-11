// chap6/complex-server/handlers/handler_test.go
package handlers

import (
	"strings"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justsushant/practical-go/chap6/complex-server/config"
)

func TestApiHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	w := httptest.NewRecorder()

	b := new(bytes.Buffer)
	c := config.InitConfig(b)

	apiHandler(w, r, c)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading the response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status: %v, Got: %v\n", http.StatusOK, resp.StatusCode)
	}

	expectedResponseBody := "Hello, world!"

	if string(body) != expectedResponseBody {
		t.Errorf("Expected response: %s, Got: %s\n", expectedResponseBody, string(body))
	}
}

func TestHealthCheckHandler(t *testing.T) {
	testCases := []struct{
		name string
		method string
		expectedOut string
		expectedStatusCode int
	}{
		{"get request", http.MethodGet, "ok", http.StatusOK},
		{"post request", http.MethodPost, "Method not allowed", http.StatusMethodNotAllowed},
	}

	b := new(bytes.Buffer)
	c := config.InitConfig(b)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/healthz", nil)
			w := httptest.NewRecorder()

			healthCheckHandler(w, r, c)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading the response body: %v", err)
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("Expected response status: %v, Got: %v\n", tc.expectedStatusCode, resp.StatusCode)
			}

			// http.Error adds an extra character to response for some reason (maybe a newline character)
			actualBody := strings.TrimSpace(string(body))

			if actualBody != tc.expectedOut {
				fmt.Println(len(string(body)))
				fmt.Println(len(tc.expectedOut))
				t.Errorf("Expected response: %s, Got: %s\n", tc.expectedOut, string(body))
			}			
		})
	}
}