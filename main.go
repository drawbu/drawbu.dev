package main

import (
	"server/pkg/app"
	"server/pkg/blog"
	"server/pkg/contact"
	"server/pkg/homepage"
	"server/pkg/resume"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	assetsDir string = "./assets/"
)

func main() {
	serv := app.Server{Port: 8080, AssetsDir: assetsDir}

	serv.AddRoute("/", homepage.Handler)
	serv.AddRoute("/blog/{article...}", blog.Handler)
	serv.AddRoute("/contact", contact.Handler)
	serv.AddRoute("/resume", resume.Handler)
	serv.ServeRoute("/assets/", "/assets/")

	serv.Run()
}
