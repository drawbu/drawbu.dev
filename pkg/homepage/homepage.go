package homepage

import (
	"context"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := homepage()
	component.Render(context.Background(), w)
}
