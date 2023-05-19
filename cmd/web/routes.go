package main

import (
	"net/http"
	
	"github.com/justinas/alice" // for better middleware chaining
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Strips '/static' leaving only /. That way the file server
	// doesn't process an unwanted resource (/static) if it happens to exist
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Pass mux as a http.Handler into custom mw
	// recoverPanic => appLogger => secureHeaders => servemux => handler
	// return app.recoverPanic(app.appLogger(secureHeaders(mux)))

	// middlware chaining using Alice
	standard := alice.New(app.recoverPanic, app.appLogger, secureHeaders)
	return standard.Then(mux)
}
