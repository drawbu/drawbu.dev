package root

import (
	"net/http"

	"app/pkg/app"
	"app/static"

	"github.com/a-h/templ"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	switch r.URL.Path {
	case "/":
		return homepage(), nil
	default:
        http.FileServerFS(static.Files).ServeHTTP(w, r)
		return nil, nil
	}
}
