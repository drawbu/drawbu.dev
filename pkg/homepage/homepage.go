package homepage

import (
	"context"
	"net/http"

	"server/pkg/components"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	components.Template(homepage()).Render(context.Background(), w)
}
