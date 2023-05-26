package main

import (
	"net/http"

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

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Strips '/static' leaving only /. That way the file server
	// doesn't process an unwanted resource (/static) if it happens to exist
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// /*filepath => this is a catch all
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// New middleware chain for stateful routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Auth routes
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// middlware chaining using Alice
	standard := alice.New(app.recoverPanic, app.appLogger, secureHeaders)
	return standard.Then(router)
}
