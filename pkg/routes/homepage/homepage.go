package homepage

import (
	"context"
	"net/http"

	"app/pkg/app"
	"app/pkg/components"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	return components.Template(homepage()).Render(context.Background(), w)
}
