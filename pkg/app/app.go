package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"server/pkg/components"
	"strconv"
)

type Server struct {
	Port      int16
	AssetsDir string
	Blog      blogInfo
}

type blogInfo struct {
	GithubRepo string
	RepoPath   string
	LastLookup int64
	Articles   []Article
}

type Article struct {
	Title string
	Date  string
	Path  string
	URI   string
}

func (serv *Server) Run() {
	fmt.Printf("Listening on localhost:%d\n", serv.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(int(serv.Port)), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) error) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		err := handler(serv, w, r)
		if err == nil {
			return
		}
		fmt.Printf("error getting articles: %s\n", err)
		components.Template(components.Error(err.Error(), r.RequestURI)).Render(context.Background(), w)
	})
}

func (serv *Server) ServeRoute(route string, path string) {
	http.Handle(route, http.StripPrefix(path, http.FileServer(http.Dir(serv.AssetsDir))))
}
