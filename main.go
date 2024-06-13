package main

import (
	"app/pkg/app"
	"app/pkg/routes/blog"
	"app/pkg/routes/contact"
	"app/pkg/routes/resume"
	"app/pkg/routes/root"
)

func main() {
	serv := app.Server{Port: 8080}

	serv.AddRoute("GET /", root.Handler)
	serv.AddRoute("GET /blog/{article...}", blog.Handler)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
