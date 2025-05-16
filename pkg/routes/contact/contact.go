package contact

import (
	"net/http"

	"app/pkg/app"

	"github.com/a-h/templ"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	serv.Cache_route(w, r, 3600)
	return serv.Template(contact()), nil
}
