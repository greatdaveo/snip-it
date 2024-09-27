package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"snippet-box/pkg/models"
	"strconv"
)

// Changed the signature of the home handler so it is defined as a method against the application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// To check if the request URL path matches "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// To initialize a slice containing the paths to the two files
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// To read the template file & store it in a template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// To write the template content as the response body
	err = ts.Execute(w, nil) // second parameter is the data to be pass in
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// Changed the signature of the showSnippet handler so it is defined as a method against *application & // To show snippet
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// To get the id from the URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// To fetch the snippet data from the DB
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// To create an instance of a templateData struct holding the snippet data
	data := &templateData{Snippet: s}

	// To initialize a slice containing the paths to the template files
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// To load the templates from the file paths
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// To execute the template set
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// To add a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Some variables with a dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!"
	expires := "7"
	// To pass the data to the SnippetModel.Insert() method, by receiving the ID of the bew record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// To redirect the user to the relevant page of the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}
