package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Port int16
}

func (serv *Server) Run() {
	fmt.Printf("Listening on localhost:%d\n", serv.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(int(serv.Port)), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("app closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) error) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		err := handler(serv, w, r)
		if err == nil {
			fmt.Printf("[%s] %s\n", r.Method, r.RequestURI)
		} else {
			fmt.Printf("[%s] %s: %s\n", r.Method, r.RequestURI, err)
			serv.Template(Error(err.Error(), r.RequestURI)).Render(context.Background(), w)
		}
	})
}
