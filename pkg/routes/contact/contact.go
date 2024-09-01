package contact

import (
	"net/http"

	"app/pkg/app"

	"github.com/a-h/templ"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return contact(), nil
}
