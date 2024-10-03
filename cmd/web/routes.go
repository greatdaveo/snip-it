package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// The flow of control: panicRecovery -> logRequest -> secureHeaders -> servemux -> application handler
func (app *application) routes() http.Handler {

	// To create a middleware chain containing all standard middleware, which will be used for every request the application uses
	standardMiddleware := alice.New(app.recoverPanic, app.recoverPanic, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm)) // To display the form
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))    // To submit the form
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
