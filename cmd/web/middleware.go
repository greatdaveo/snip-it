package main

import (
	"fmt"
	"net/http"
)

// To add the two HTTP headers to every response:
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// To get the user IP address, URL and method requested
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

// This is for panic recovery
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// To create a deferred function that runs in the event of a panic as Go unwinds the stack
		defer func() {
			// To check if there has been a panic or not
			if err := recover(); err != nil {
				// Set a connection "close" header on the response
				w.Header().Set("Connection", "close")
				// To return a 500 error
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})

}
