package main

import (
	"snippetbox.victorsmith.dev/internal/models"
)

// Make a holding structure for incoming data
// Can expand if we wish to add additional data later on
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
