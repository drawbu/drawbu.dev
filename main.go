package main

import (
	"app/pkg/app"
	"app/pkg/routes/blog"
	"app/pkg/routes/contact"
	"app/pkg/routes/homepage"
	"app/pkg/routes/resume"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	staticDir   string = "./static/"
	articlesDir string = "./articles/"
	prod        string
)

func main() {
	serv := app.Server{Port: 8080, Prod: prod == "true"}

	home := homepage.Handler{StaticDir: staticDir}
	serv.AddRoute("GET /", home.Render)
	blog := blog.Handler(articlesDir)
	serv.AddRoute("GET /blog/{article...}", blog.Render)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
