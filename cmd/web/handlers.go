package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.victorsmith.dev/internal/models"
	"snippetbox.victorsmith.dev/internal/validators"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	// FieldErrors map[string]string
	validators.Validator
}

// Make the home handler a method for the application struct to introduce dependency injection?
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, data, http.StatusOK, "home.html")
}

// Returns page containing detials of snippet with :id
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// The :id param is passed via the request context
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// We render the individual snippers under the view template
	app.render(w, data, http.StatusOK, "view.html")
}

// Fetch form page
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, data, http.StatusOK, "create.html")
}

// Creates new snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// Stores the request body in PostForm map as key value pairs
	// If the data is bad, or there is too much data => PostForm map remains blank
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Convert to int
	// Send 400 if conversion fails (invalid date)
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Create isntance of snippetCreateForm
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Validate errors
	// 1) Check that the title and content fields are not empty.
	// 2) Check that the title field is not more than 100 characters long.
	// 3) Check that the expires value exactly matches one of our permitted values ( 1 , 7 or 365 days).

	// Embedding of validators.Validator allows for a direct call to the Validator method(s)
	form.CheckField(validators.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validators.MaxChars(form.Title, 100), "title", "this field cannot be 100 chars long")
	form.CheckField(validators.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validators.PermittedInt(form.Expires, 1, 7, 365), "expires", "Value must be 1, 7 or 365")

	// use the HTTP status code 422 Unprocessable Entity to indicate bad data in fomr
	// pass the snippetCreateForm object to the template
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, data, http.StatusUnprocessableEntity, "create.html")
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
