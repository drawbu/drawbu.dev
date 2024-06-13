package main

import (
	"app/pkg/app"
	"app/pkg/routes/blog"
	"app/pkg/routes/contact"
	"app/pkg/routes/root"
	"app/pkg/routes/resume"
)

// These values may be set by the build script via the ldflags argument
var (
	staticDir   string = "./static/"
	articlesDir string = "./articles/"
)

func main() {
	serv := app.Server{Port: 8080}

	home := root.Handler(staticDir)
	serv.AddRoute("GET /", home.Render)
	blog := blog.Handler(articlesDir)
	serv.AddRoute("GET /blog/{article...}", blog.Render)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
