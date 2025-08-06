package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Handled request %s %s in %v\n", r.Method, r.URL.Path, duration)
	})
}
