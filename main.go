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
	assetsDir string = "./assets/"
)

func main() {
	serv := app.Server{Port: 8080}

	home := homepage.Handler{Assets: assetsDir}
	serv.AddRoute("GET /", home.Render)
	blog := blog.Handler{}
	go blog.FetchArticles()
	serv.AddRoute("GET /blog/{article...}", blog.Render)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
