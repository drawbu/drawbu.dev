package resume

import (
	"context"
	"net/http"

	"server/pkg/components"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	components.Template(resume()).Render(context.Background(), w)
}
