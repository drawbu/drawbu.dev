package contact

import (
	"context"
	"net/http"

	"server/pkg/components"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := components.Template(contact())
	component.Render(context.Background(), w)
}
