package root

import (
	"context"
	"net/http"

	"app/pkg/app"
	"app/static"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	switch r.URL.Path {
	case "/":
		return serv.Template(homepage()).Render(context.Background(), w)
	default:
        http.FileServerFS(static.Files).ServeHTTP(w, r)
		return nil
	}
}
