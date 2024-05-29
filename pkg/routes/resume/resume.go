package resume

import (
	"context"
	"net/http"

	"server/pkg/app"
	"server/pkg/components"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	return components.Template(resume()).Render(context.Background(), w)
}
