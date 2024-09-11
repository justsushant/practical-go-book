// chap5/http-serve-mux/server.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Key details to log are URL, Request type, Request body size, and Protocol
	// Each log line should be a JSON formatted string?? incomplete
	log.Printf("%s %s %d %s\n", r.URL, r.Method, r.ContentLength, r.Proto)
	fmt.Fprint(w, "Hello, World!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api", apiHandler)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)


	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

// The default value of ":8080" means that the server will listen on all network interfaces on the port 8080. 
// If you wanted the server only to be reachable on the computer where your application is running, 
// you would set the environment variable LISTEN_ADDR to "127.0.0.1:8080" and then start the application.


// When a request comes in and a handler function is available to handle it,
// the handler function is executed in a separate goroutine. 
// Once the processing completes, the goroutine is terminated 