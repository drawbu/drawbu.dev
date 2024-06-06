package homepage

import (
	"context"
	"errors"
	"net/http"

	"app/pkg/app"
	"app/pkg/components"
)

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
    if r.URL.Path != "/" {
        return errors.New("Page not found")
    }
	return components.Template(homepage()).Render(context.Background(), w)
}
