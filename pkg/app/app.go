package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Port      int16
	AssetsDir string
	Blog      blogInfo
}

type blogInfo struct {
	githubRepo string
	repoPath   string
	lastLookup int64
	articles   []Article
}

type Article struct {
	Title string
	Date  string
	Path  string
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

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		handler(serv, w, r)
	})
}

func (serv *Server) ServeRoute(route string, path string) {
	http.Handle(route, http.StripPrefix(path, http.FileServer(http.Dir(serv.AssetsDir))))
}
