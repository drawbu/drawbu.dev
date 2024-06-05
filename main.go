package main

import (
	"server/pkg/app"
	"server/pkg/routes/blog"
	"server/pkg/routes/contact"
	"server/pkg/routes/homepage"
	"server/pkg/routes/resume"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	assetsDir string = "./assets/"
)

func main() {
	serv := app.Server{Port: 8080, AssetsDir: assetsDir}

	serv.AddRoute("/", homepage.Handler)
	blog := blog.Handler{}
	go blog.FetchArticles()
	serv.AddRoute("/blog/{article...}", blog.Render)
	serv.AddRoute("/contact", contact.Handler)
	serv.AddRoute("/resume", resume.Handler)
	serv.ServeRoute("/assets/", "/assets/")

	serv.Run()
}
