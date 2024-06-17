package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
)

type Server struct {
	Port int16
}

func (serv *Server) Run() {
	log.Infof("Listening on localhost:%d", serv.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(int(serv.Port)), nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("App closed")
	} else if err != nil {
		log.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) error) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		err := handler(serv, w, r)
        req_fmt := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)
		if err == nil {
			log.Info(req_fmt)
		} else {
			log.Warn(req_fmt, "reason", err)
			serv.Template(Error(err.Error(), r.RequestURI)).Render(context.Background(), w)
		}
	})
}
