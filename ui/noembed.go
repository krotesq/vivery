//go:build noui

package ui

import "io/fs"

func FS() (fs.FS, error) { return nil, nil }
