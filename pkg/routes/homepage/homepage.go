package homepage

import (
	"context"
	"net/http"
	"strings"

	"app/pkg/app"
	"app/pkg/components"
)

type Handler struct {
	StaticDir     string
	staticHandler http.Handler
}

func (h *Handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	switch strings.ToLower(r.URL.Path) {
	case "/":
		return components.Template(homepage()).Render(context.Background(), w)
	default:
		if h.staticHandler == nil {
			h.staticHandler = http.FileServer(http.Dir(h.StaticDir))
		}
		h.staticHandler.ServeHTTP(w, r)
		return nil
	}
}
