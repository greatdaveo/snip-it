package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// The standard flow of control: panicRecovery -> logRequest -> secureHeaders -> servemux -> application handler
func (app *application) routes() http.Handler {

	// To create a middleware chain containing all standard middleware, which will be used for every request the application uses
	standardMiddleware := alice.New(app.recoverPanic, app.recoverPanic, secureHeaders)
	// To create a middleware chain containing the middleware specific to our dynamic application routes
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm)) // To display the form
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))    // To submit the form
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
