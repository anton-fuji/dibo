package templates

import (
	"embed"
)

/*
 FS holds the embedded .dockerignore templates.
*/

//go:embed *.dockerignore
var FS embed.FS
