package middlewares

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details about each incoming HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Capture the start time
        start := time.Now()

        // Log the incoming request details
        log.Printf("Started %s %s \n", r.Method, r.URL.Path)

        // Pass the request to the next handler in the chain
        next.ServeHTTP(w, r)

        // Calculate the duration and log it
        duration := time.Since(start)
        log.Printf("Completed %s %s in %v \n", r.Method, r.URL.Path, duration)
    })
}