package middleware

import (
	"net/http"
	"time"
)

func (mid *Middleware) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		mid.logger.Printf("New request | method %s | remote_addr %s | url %s | time %s\n", request.Method,
			request.RemoteAddr, request.URL.Path, time.Since(start))
	})
}
