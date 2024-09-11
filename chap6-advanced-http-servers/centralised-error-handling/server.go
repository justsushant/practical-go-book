package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type MethodNotAllowedError struct {
	mssg string
	code int
}

func (e MethodNotAllowedError) Error() string {
	return e.mssg
}

var ErrMethodNotAllowed = MethodNotAllowedError{mssg: "Method not allowed", code: http.StatusMethodNotAllowed}

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig) error
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := a.handler(w, r, a.config)
	if err != nil {
		if errType, ok := err.(MethodNotAllowedError); ok {
			// can log the error here or some other stuff
			a.config.logger.Printf("Error \"%s\" for %s with %s method\n", errType.mssg, r.URL, r.Method)
			http.Error(w, errType.Error(), errType.code)
			return
		} 
			
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)	
	}	
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) error {
	if r.Method != http.MethodGet {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ErrMethodNotAllowed
	}

	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "OK")
	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) error {
	config.logger.Println("Handling API request")
	fmt.Fprintf(w, "Hello, World!")
	return nil
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config: config, handler: healthCheckHandler})
	mux.Handle("/api", &app{config: config, handler: apiHandler})
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

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}