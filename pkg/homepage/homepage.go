package homepage

import (
	"context"
	"net/http"

	"server/pkg/components"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := components.Template(homepage())
	component.Render(context.Background(), w)
}
