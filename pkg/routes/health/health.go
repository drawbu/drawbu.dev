package health

import (
	"io"
	"net/http"

	"app/pkg/app"

	"github.com/a-h/templ"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	io.WriteString(w, "OK")
	return nil, nil
}
