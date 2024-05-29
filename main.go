package main

import (
	"server/pkg/app"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	assetsDir string = "./assets/"
)

func main() {
	server := app.App{Port: 8080, AssetsDir: assetsDir}
	server.Run()
}
