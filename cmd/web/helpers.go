package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog and sends 500 Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// To report is the file name and line number one step back in the stack trace,
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// This sends a 400 Error when theres is a Bad Request from the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// This sends a 404 Error to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// An helper method to easily render the templates from the cache
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// To initialize a new buffer to catch runtime errors (By writing the template to the buffer, instead to the http.ResponseWriter) & if there is err it calls the serveError helper
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data
	err := ts.Execute(buf, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// To write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}
