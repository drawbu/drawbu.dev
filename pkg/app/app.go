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
	Rev string
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
		log_fmt := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)

		w.Header().Add("Cache-Control", "no-cache, must-revalidate")
		w.Header().Add("Vary", "HX-Request")
		is_htmx := r.Header.Get("HX-Request") == "true"

		// Cache control
		hash := r.Header.Get("If-None-Match")
		expected_hash := serv.Rev
		if is_htmx {
			expected_hash = "htmx-" + expected_hash
		}
		if hash == expected_hash {
			// Already cached
			log.Info(fmt.Sprintf("Cached %s", log_fmt))
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Add("ETag", expected_hash)

		comp, err := handler(serv, w, r)

		if err == nil {
			log.Info(log_fmt)
		} else {
			log.Warn(log_fmt, "reason", err)
			comp = Error(err.Error(), r.RequestURI)
		}

		// Already served
		if comp == nil {
			return
		}

		if !is_htmx {
			comp = serv.Template(comp)
		}
		comp.Render(context.Background(), w)
	})
}
