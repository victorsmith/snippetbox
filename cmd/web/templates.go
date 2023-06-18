package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.victorsmith.dev/internal/models"
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
}

// filename: []ts
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize new map to act as the cache
	cache := map[string]*template.Template{}

	// Get slice of paths which match the provided pattern
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// register the template.FuncMap, and then parse the file as normal.
		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// parse files into template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Cache the ts. Use name as key
		cache[name] = ts
	}
	return cache, nil
}
