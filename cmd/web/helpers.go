package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog and sends 500 Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// To report is the file name and line number one step back in the stack trace,
	app.errorLog.Output(2,trace)

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