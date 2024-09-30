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

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	// To create an instance of a templateData struct holding the slice of snippets
	data := &templateData{Snippets: s}

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
	err = ts.Execute(w, data) // second parameter is the data to be pass in i.e templateData
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

	// To use the render helper function
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

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

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

// To add a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

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

	// To redirect the user to the relevant page of the snippet using semantic URL style
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}
