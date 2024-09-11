package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) >= 1 {
		return errors.New("stopped after 1 redirect")
	}
	return nil
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d, CheckRedirect: redirectPolicyFunc}
	return &client
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify a HTTP URL to get data from")
		os.Exit(1)
	}

	client := createHTTPClientWithTimeout(15 * time.Second)
	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", body)
}

// customizing requests
func createHTTPGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req, err
}

// sending custom requests
// client := createHTTPClientWithTimeout(20 * time.Millisecond)
// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
// defer cancel()
// req, err := createHTTPGetRequest(ctx, ts.URL+"/api/packages", nil)
// resp, err := client.Do(req)