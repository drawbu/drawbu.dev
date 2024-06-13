package homepage

import (
	"context"
	"net/http"
	"strings"

	"app/pkg/app"
)

type handler struct {
	staticDir http.Handler
}

func Handler(staticDir string) *handler {
	return &handler{
		staticDir: http.FileServer(http.Dir(staticDir)),
	}
}

func (h *handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	switch strings.ToLower(r.URL.Path) {
	case "/":
		return serv.Template(homepage()).Render(context.Background(), w)
	default:
		h.staticDir.ServeHTTP(w, r)
		return nil
	}
}
