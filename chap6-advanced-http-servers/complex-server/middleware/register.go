// chap6/complex-server/middleware/register.go
package middleware

import (
	"net/http"
	"github.com/justsushant/practical-go/chap6/complex-server/config"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, c), c)
}