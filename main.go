package main

import (
	"app/pkg/app"
	"app/pkg/routes/blog"
	"app/pkg/routes/contact"
	"app/pkg/routes/health"
	"app/pkg/routes/resume"
	"app/pkg/routes/root"
)

var (
	hash string = "dev";
)

func main() {
	serv := app.Server{Port: 8080, Hash: hash}

	serv.AddRoute("GET /", root.Handler)
	serv.AddRoute("GET /health", health.Handler)
	serv.AddRoute("GET /blog/{article...}", blog.Handler)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
