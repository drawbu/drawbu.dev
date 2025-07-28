package main

import (
	"app/pkg/app"
	"app/pkg/routes/blog"
	"app/pkg/routes/contact"
	"app/pkg/routes/health"
	"app/pkg/routes/resume"
	"app/pkg/routes/root"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/pschou/go-params"
)

var (
	rev         string = "foo"
	defaultPort int    = 8080
)

func main() {
	params.Usage = func() {
		fmt.Fprintf(os.Stderr, "drawbu.dev, build rev: %s\n\n"+
			"Usage: %s [options...]\n\n", rev, os.Args[0])
		params.PrintDefaults()
	}
	var port = params.Int("port", defaultPort, "Web server's listening port", "Number")
	params.Parse()

	log.Info("Starting server", "port", *port, "rev", rev)

	serv := app.Server{Port: int16(*port), Rev: rev}

	serv.AddRoute("GET /", root.Handler)
	serv.AddRoute("GET /health", health.Handler)
	serv.AddRoute("GET /blog/rss.xml", blog.RssHandler)
	serv.AddRoute("GET /blog/atom.xml", blog.AtomHandler)
	serv.AddRoute("GET /blog/{article...}", blog.Handler)
	serv.AddRoute("GET /contact", contact.Handler)
	serv.AddRoute("GET /resume", resume.Handler)

	serv.Run()
}
