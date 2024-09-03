package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/charmbracelet/log"
)

type Server struct {
	Port int16
	Hash string
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

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) (templ.Component, error)) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		comp, err := handler(serv, w, r)
		req_fmt := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)

		if err == nil {
			log.Info(req_fmt)
		} else {
			log.Warn(req_fmt, "reason", err)
			comp = Error(err.Error(), r.RequestURI)
		}

		// Already served
		if comp == nil {
			return
		}

		if r.Header.Get("HX-Request") != "true" {
			comp = serv.Template(comp)
		}
		comp.Render(context.Background(), w)
	})
}
