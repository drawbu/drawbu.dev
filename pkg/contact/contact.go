package contact

import (
	"context"
	"net/http"

	"server/pkg/components"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	components.Template(contact()).Render(context.Background(), w)
}
