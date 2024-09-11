// chap6/complex-server/handlers/register.go
package handlers

import (
	"net/http"
	"github.com/justsushant/practical-go/chap6/complex-server/config"
)

func Register(mux *http.ServeMux, conf config.AppConfig) {
	mux.Handle("/healthz", &app{conf: conf, handler: healthCheckHandler})
	mux.Handle("/api", &app{conf: conf, handler: apiHandler})
	mux.Handle("/panic", &app{conf: conf, handler: panicHandler})
}