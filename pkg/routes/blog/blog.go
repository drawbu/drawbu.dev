package blog

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"app/articles"
	"app/pkg/app"

	"github.com/a-h/templ"
	"github.com/charmbracelet/log"
)

var Articles = make(map[string]article)

func init() {
	for _, a := range getArticles(articles.Articles) {
		Articles[strings.ToLower(a.Uri)] = a
	}
}

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil {
		return nil, fmt.Errorf("Failed to get article: %s", err)
	}

	// root /blog
	if article_name == "" {
		serv.Cache_route(w, r, 3600)
		return serv.Template(blog(getSortedArticles(Articles))), nil
	}

	// query article
	a, ok := Articles[strings.ToLower(article_name)]
	if !ok {
		return nil, errors.New("Article not found")
	}
	serv.Cache_route(w, r, 3600)
	return serv.Template(articleShow(a)), nil
}

func getSortedArticles(articles map[string]article) []article {
	result := make([]article, 0, len(articles))
	for _, a := range articles {
		result = append(result, a)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})
	return result
}

func getArticles(filesystem fs.ReadDirFS) []article {
	entries, err := filesystem.ReadDir(".")
	if err != nil {
		log.Warn("Error reading directory:", "reason", err)
		return []article{}
	}

	result := []article{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") ||
			!entry.Type().IsRegular() ||
			!strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		file, err := filesystem.Open(entry.Name())
		if err != nil {
			log.Warn("Error opening file:", "reason", err)
			continue
		}
		a, err := NewArticle(file)
		if err != nil {
			log.Fatal("Error building article", "name", entry.Name(), "reason", err)
			continue
		}
		result = append(result, *a)
	}
	return result
}
