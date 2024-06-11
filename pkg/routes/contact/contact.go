package contact

import (
	"context"
	"net/http"

	"app/pkg/app"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	return serv.Template(contact()).Render(context.Background(), w)
}
