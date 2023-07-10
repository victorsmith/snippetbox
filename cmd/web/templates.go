package main

import (
	"html/template"
	"path/filepath"
	"time"
	"io/fs"

	"snippetbox.victorsmith.dev/internal/models"
	"snippetbox.victorsmith.dev/ui"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// { string: function } map => used to fetch functions in template
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Make a holding structure for incoming data
// Can expand if we wish to add additional data later on
type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// filename: []ts
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize new map to act as the cache
	cache := map[string]*template.Template{}


	pages, err := fs.Glob(ui.Files, "html/pages/*.html")	
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we 
		// want to parse.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files 
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Cache the ts. Use name as key
		cache[name] = ts
	}
	return cache, nil
}
