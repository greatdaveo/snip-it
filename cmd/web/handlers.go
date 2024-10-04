package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"snippet-box/pkg/forms"
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
	//  To extract the ID from the URL path
	idStr := r.URL.Path[len("/snippet/"):]
	// To convert to int
	id, err := strconv.Atoi(idStr)
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

// To render the snippet form page
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// To pass a new empty forms.forms object to the templte
		Form: forms.New(nil),
	})
}

// To add a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// To parse the form data in POST, PUT & PATCH request bodies to the r.PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// To create a new forms.Form struct containing the POSTed data from the form, and using the validation method to check the content
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "367", "7", "1")

	// If the form isn't valid, redisplay the template passing in the form.Form object as the data
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// To retrieve the validated fields values from the form
	title := form.Get("title")
	content := form.Get("content")
	expiresStr := form.Get("expires")
	//To convert expires to int
	expires, err := strconv.Atoi(expiresStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// To insert the snippet validated data in the DB
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// To add a confirmation message to the user's session data
	app.session.Put(r, "flash", "Snippet successfully created!")

	// To redirect the user to the relevant page of the snippet using semantic URL style
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
