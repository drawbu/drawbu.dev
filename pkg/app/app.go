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
	Rev  string
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

func (serv *Server) is_htmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (serv *Server) compute_etag(r *http.Request) string {
	etag := serv.Rev

	if etag == "dev" {
		return ""
	}
	if serv.is_htmx(r) {
		return "htmx-" + etag
	}
	return etag
}

func (serv *Server) Cache_route(w http.ResponseWriter, r *http.Request, max_age int32) {
	expected_etag := serv.compute_etag(r)
	if expected_etag == "" {
		return
	}

	w.Header().Add("Cache-Control", "max-age="+fmt.Sprint(max_age))
	w.Header().Add("Vary", "HX-Request")
	w.Header().Add("ETag", expected_etag)
}

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) (templ.Component, error)) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		log_fmt := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)

		etag := serv.compute_etag(r)
		if etag != "" && etag == r.Header.Get("If-None-Match") {
			log.Info("Cached " + log_fmt)
			w.WriteHeader(http.StatusNotModified)
			return
		}

		comp, err := handler(serv, w, r)

		if err == nil {
			log.Info(log_fmt)
		} else {
			log.Warn(log_fmt, "reason", err)
			comp = Error(err.Error(), r.RequestURI)
		}

		// Already served by handler
		if comp == nil {
			return
		}
		if !serv.is_htmx(r) {
			comp = serv.Template(comp)
		}
		log.Info(log_fmt)
		comp.Render(context.Background(), w)
	})
}
