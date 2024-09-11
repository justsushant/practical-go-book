package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

// custom HTTP handler middleware to log how long it took to process request
// this middleware won't log 404 requests because requests never reaches this due to the lack of corresponsing handler function (and app object)
// func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	startTime := time.Now()
// 	a.handler(w, r, a.config)
// 	a.config.logger.Printf("path=%s method=%s duration=%f", r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds())
// }

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "OK")
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, World!")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
}

// logging middleware with http.HandlerFunc type
func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("path=%s method=%s duration=%f", r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds())
	})
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	
	config := appConfig {
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)
	// log.Fatal(http.ListenAndServe(listenAddr, mux))

	// below two lines, we wrap the ServeMux object with outer http.Handler type (i.e; loggingMiddleware)
	// after this, http.ServeMux object will be referred by the wrapped handler
	m := loggingMiddleware(mux)
	log.Fatal(http.ListenAndServe(listenAddr, m))

	// the value returned by the loggingMiddleware() function implements the http.Handler interface,
	// we specify that as the handler when calling the ListenAndServe() function
}