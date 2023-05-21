package main

import (
	"net/http"
	
	"github.com/justinas/alice" // for better middleware chaining
	"github.com/julienschmidt/httprouter" // for better routing
)

func (app *application) routes() http.Handler {
	// router init 
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Strips '/static' leaving only /. That way the file server
	// doesn't process an unwanted resource (/static) if it happens to exist
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// /*filepath => this is a catch all  
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate) // fetches form
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// middlware chaining using Alice
	standard := alice.New(app.recoverPanic, app.appLogger, secureHeaders)
	return standard.Then(router)
}
