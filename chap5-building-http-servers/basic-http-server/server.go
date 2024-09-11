// chap5/basic-http-server/server.go
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

// The default value of ":8080" means that the server will listen on all network interfaces on the port 8080. 
// If you wanted the server only to be reachable on the computer where your application is running, 
// you would set the environment variable LISTEN_ADDR to "127.0.0.1:8080" and then start the application.