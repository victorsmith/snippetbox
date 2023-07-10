package ui

import (
	"embed"
)

// This is not a comment -> It's a special directive

//go:embed "html" "static"
var Files embed.FS

