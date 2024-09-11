package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type requestIDKey struct {}
type requestIDVal struct {
	requestID string
}

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

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

func panicHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("I panicked")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
	mux.Handle("/panic", &app{config: config, handler: panicHandler})
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		var reqID string
		if v, ok := r.Context().Value(requestIDKey{}).(requestIDVal); ok {
			reqID = v.requestID
		}
		log.Printf("path=%s method=%s duration=%f requestID=%v", r.URL.Path, r.Method, time.Since(startTime).Seconds(), reqID)
	})
}

func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Printf("inside the defer of panic middleware")
			if v, ok := r.Context().Value(requestIDKey{}).(requestIDVal); ok {
				log.Printf("Processing request inside panic middleware: %s", v.requestID)
			}
			if rValue := recover(); rValue != nil {
				log.Println("panic detected", rValue)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected server error")
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := "requestID123"
		c := requestIDVal{requestID: requestID}

		newCtx := context.WithValue(r.Context(), requestIDKey{}, c)	
		if v, ok := newCtx.Value(requestIDKey{}).(requestIDVal); ok {
			log.Printf("Processing request: %s", v.requestID)
		}
		h.ServeHTTP(w, r.WithContext(newCtx))
		log.Println("After ServeHTTP of requestID middleware")
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

	m := requestIDMiddleware(loggingMiddleware(panicMiddleware(mux)))
	log.Fatal(http.ListenAndServe(listenAddr, m))
}