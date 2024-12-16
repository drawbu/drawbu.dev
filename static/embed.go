package static

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var embeded_files embed.FS

var Files fs.FS

func init() {
	var err error
	Files, err = fs.Sub(embeded_files, "dist")

	if err != nil {
		panic(err)
	}
}
