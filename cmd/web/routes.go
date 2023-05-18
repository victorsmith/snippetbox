package main

import "net/http"

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
	// appLogger => secureHeaders => servemux => handler
	return app.appLogger(secureHeaders(mux))
}
