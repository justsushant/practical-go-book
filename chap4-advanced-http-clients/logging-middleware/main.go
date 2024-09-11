// chap4/logging-middleware/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io"
)

// can be used to emit metrics such as request latency or non-200 errors
// can be used to look up the request in a cache automatically, for example, to avoid making the call
// can be used to add one or more HTTP headers to every outgoing request like, sending an authentication header, propagating a request ID etc

type LoggingClient struct {
	log *log.Logger
}

func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)
	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)
	return resp, err
}

// can be used like below
// myTransport := LoggingClient{}
// client := http.Client{
//  Timeout: 10 * time.Second,
//  Transport: &myTransport,
// }

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

func fetchRemoteResource(client *http.Client, url string)  ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return io.ReadAll(r.Body)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify a HTTP URL to get data from")
		os.Exit(1)
	}

	myTransport := LoggingClient{}
	l := log.New(os.Stdout, "", log.LstdFlags)
	myTransport.log = l

	client := createHTTPClientWithTimeout(15*time.Second)
	client.Transport = &myTransport

	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Bytes in response: %d\n", len(body))
}

// There are two things you will have to keep in mind while implementing a
// custom RoundTripper:
// 1. The RoundTripper must be implemented with the assumption that there
// may be more than one instance of it running at any given point of time.
// Hence, if you are manipulating any data structure, the data structure must
// be concurrency safe.
// 2. The RoundTripper must not mutate the request or response or return an
// error.