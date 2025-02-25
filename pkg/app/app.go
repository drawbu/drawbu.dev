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

func is_htmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func compute_etag(serv *Server, r *http.Request) string {
	etag := serv.Rev

	if etag == "dev" {
		return ""
	}
	if is_htmx(r) {
		return "htmx-" + etag
	}
	return etag
}

func (serv *Server) AddRoute(route string, handler func(app *Server, w http.ResponseWriter, r *http.Request) (templ.Component, error)) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		log_fmt := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)

		expected_etag := compute_etag(serv, r)
		if expected_etag != "" {
			// Cache duration is one hour
			w.Header().Add("Cache-Control", "max-age=3600")
			w.Header().Add("Vary", "HX-Request")
			if expected_etag == r.Header.Get("If-None-Match") {
				// Already cached
				log.Info("Cached " + log_fmt)
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Add("ETag", expected_etag)
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
			log.Info(log_fmt)
			return
		}

		if !is_htmx(r) {
			comp = serv.Template(comp)
		}
		comp.Render(context.Background(), w)
	})
}
