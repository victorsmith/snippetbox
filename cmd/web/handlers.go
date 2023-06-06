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
	Title                string     `form:"title"`
	Content              string     `form:"content"`
	Expires              int        `form:"expires"`
	validators.Validator `form:"-"` // tells the decoder to completely ignore a field during decoding.
}

type userSignupForm struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

type userLoginForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
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
	// Initalize empty form
	var form snippetCreateForm

	// If there is a problem, we return a 400 Bad Request response to the client.
	err := app.decodePostError(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet Succesfully Created!")

	// Redirect to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Auth Handlers
// Fetch user login page
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, data, http.StatusOK, "login.html")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Initialize empty form
	var form userLoginForm

	// Decode form errors => ** how does this work exactly?
	err := app.decodePostError(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form fields
	form.CheckField(validators.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validators.Matches(form.Email, validators.EmailRegexp), "email", "Field must be a valid email address")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")

	// Redirect to login page if the form contains any errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, data, http.StatusUnprocessableEntity, "login.html")
		return
	}

	// Check credential validity
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, data, http.StatusUnprocessableEntity, "login.html")
			return
		} else {
			// Catch other errors
			app.serverError(w, err)
			return
		}
	}

	// Use the RenewToken() method on the current session to change the session ID (generate a new id).
	// This should be done if: a) auth state changes or b) privelages state changes for the user
	// Do at login and logout
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add user id to session
	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)

	// Redirect user
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, data, http.StatusOK, "signup.html")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostError(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form content using helpers
	// Validate the form contents using our helper functions.

	form.CheckField(validators.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validators.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validators.Matches(form.Email, validators.EmailRegexp), "email", "This field must be a valid email address")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// Return to signup page and display errors in template
		app.render(w, data, http.StatusUnprocessableEntity, "signup.html")
		return
	}

	// Insert valid data
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		// Email duplicate error
		// => Render signup page again with errors in appropriate fields
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, data, http.StatusUnprocessableEntity, "signup.html")
		} else {
			// Throw a server error otherwise
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Change the session token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove authenticated userId
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	// Add a flash message to communicate the user has been logged out
	app.sessionManager.Put(r.Context(), "flash", "User has been logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
