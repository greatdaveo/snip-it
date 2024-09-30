package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The flow of control: panicRecovery -> logRequest -> secureHeaders -> servemux -> application handler
func (app *application) routes() http.Handler {

	// To create a middleware chain containing all standard middleware, which will be used for every request the application uses
	standardMiddleware := alice.New(app.recoverPanic, app.recoverPanic, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
