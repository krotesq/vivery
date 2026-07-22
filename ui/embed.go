//go:build !noui

package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

func FS() (fs.FS, error) {
	return fs.Sub(files, "dist")
}
