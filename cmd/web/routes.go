package main

import (
	"net/http"

	"snippetbox.victorsmith.dev/ui"

	"github.com/julienschmidt/httprouter" // for better routing
	"github.com/justinas/alice"           // for better middleware chaining
)

func (app *application) routes() http.Handler {
	// router init
	router := httprouter.New()

	// 404 handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	
	// Old approach
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	
	// Take ui.Files and covert to http.FS type to satisfy the http.Filesystem interface
	fileServer := http.FileServer(http.FS(ui.Files))

	// now longer need to strip the prefix from the request URL 
	// any requests that start with /static/ can just be passed 
	// directly to the file server and the corresponding static 
	// file will be served (so long as it exists).
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// New middleware chain for stateful routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	
	// Protected (authenticated-only) application routes, using a new "protected" 
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	// Protected Routes
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	
	// middlware chaining using Alice
	standard := alice.New(app.recoverPanic, app.appLogger, secureHeaders)
	return standard.Then(router)
}
