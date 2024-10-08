package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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

// To ensure unauthenticated users can't create snippet
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect user to the login page and return the middleware chain so that subsequent handlers won't execute
		if app.authenticatedUser(r) == 0 {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		// Otherwise call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// This is a func that uses a customized CSRF cookie with the Secure, Path and HTTPOnly Flags set
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
