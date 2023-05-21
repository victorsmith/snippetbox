package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.victorsmith.dev/internal/models"
	"github.com/julienschmidt/httprouter"
)

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
	w.Write([]byte("snippetCreate..."))
}

// Creates new snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
