package homepage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"app/pkg/app"
	"app/pkg/components"
)

type Handler struct {
	Assets        string
	assetsHandler http.Handler
}

func (h *Handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	switch strings.ToLower(r.URL.Path) {
	case "/":
		return components.Template(homepage()).Render(context.Background(), w)
	case "/robots.txt":
		fmt.Fprintf(w, "User-agent: *\nallow: /")
		return nil
	default:
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			if h.assetsHandler == nil {
				h.assetsHandler = http.StripPrefix("/assets/", http.FileServer(http.Dir(h.Assets)))
			}
			h.assetsHandler.ServeHTTP(w, r)
			return nil
		}
		return errors.New("not found")
	}
}
