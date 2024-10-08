package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

// To add default data to the templateData struct anytime a template is rendered
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()

	// To add add the flash message to the template data, if on exists
	td.Flash = app.session.PopString(r, "flash")

	return td
}

// An helper method to easily render the templates from the cache
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// To add default data (like current year) to the template data
	td = app.addDefaultData(td, r)

	// To get templates in cache
	// fmt.Println("Available templates in cache:")
	// for k := range app.templateCache {
	// 	fmt.Println(k)
	// }

	// To fetch the template from the cache
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// To initialize a new buffer to catch runtime errors (By writing the template to the buffer, instead to the http.ResponseWriter) & if there is err it calls the serveError helper
	buf := new(bytes.Buffer)

	// To execute the template set, by passing in the dynamic data (td)
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// To write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// To return the ID of the current user from the session, and zero if it is not from an authenticated user
func (app *application) authenticatedUser(r *http.Request) int {
	return app.session.GetInt(r, "userID")
}
