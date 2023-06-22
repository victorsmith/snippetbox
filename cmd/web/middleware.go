package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// This mw will act on all routes
func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) appLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer func will always run last before program exit
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				// Go’s HTTP server will automatically close the current connection after a response has been sent
				w.Header().Set("Connection", "close")

				// Call the app.serverError helper method to return a 500 - internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		// If yes => go to next middleware call
		// If no 	=> stop the chain and goto login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			// Return stops mw proceeding
			return
		}

		// Set cache-control header to "no-store" s.t pages which require auth aren't stored
		// in users browser cache / or other intermediary caches
		// TODO: learn more about cache-control header
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}


func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authenticatedUserId from session
		// If not present => 0 returned + mw chain stops
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
		if id == 0 {
			// Still pass unchanged request and let the middleware handle it
			next.ServeHTTP(w, r)
			return
		}

		// Check if user exists
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		
		// If a matching user is found, we know that the request is 
		// coming from an authenticated user who exists in our database. We 
		// create a new copy of the request (with an isAuthenticatedContextKey 
		// value of true in the request context) and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		fmt.Println("gets here 1")

		// Call the next handler in the chain. next.ServeHTTP(w, r)
		next.ServeHTTP(w, r)
	})
}
