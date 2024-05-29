package blog

import (
	"context"
	"net/http"
	"strings"

	"server/pkg/app"
	"server/pkg/components"
)

func getArticles(serv *app.Server) ([]app.Article, error) {
	articles := []app.Article{}
	return articles, nil
}

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) {
	_, article_path, found := strings.Cut(r.RequestURI, "/blog/")

	if !found || article_path == "" {
		components.Template(blog()).Render(context.Background(), w)
		return
	}
	components.Template(article()).Render(context.Background(), w)
}
