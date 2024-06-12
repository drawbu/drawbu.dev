package blog

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"app/pkg/app"

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Handler(articlesDir string) *handler {
	articles := make(map[string]article)
	for _, a := range getArticles(articlesDir) {
		articles[a.Title] = a
	}
	return &handler{
		articles: articles,
	}
}

type handler struct {
	articles map[string]article
}

type article struct {
	Title   string
	Date    time.Time
	Content []byte
}

func (h *handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil || article_name == "" {
		return serv.Template(blog(getSortedArticles(h.articles))).Render(context.Background(), w)
	}

	a, ok := h.articles[article_name]
	if !ok {
		return errors.New("Article not found")
	}
	return serv.Template(articleShow(a)).Render(context.Background(), w)
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

func getArticles(path string) []article {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []article{}
	}

	articles := []article{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			articles = append(articles, getArticles(path+"/"+entry.Name())...)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			a, err := parseMarkdownArticle(path + "/" + entry.Name())
			if err == nil || a != nil {
				articles = append(articles, *a)
			}
		}
	}
	return articles
}

func parseMarkdownArticle(path string) (*article, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("catppuccin-mocha"),
				highlighting.WithFormatOptions(chroma.WithLineNumbers(true)),
			),
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	metaData := meta.Get(context)
	date, err := time.Parse("2006-01-02", metaData["date"].(string))
	return &article{
		Title:   metaData["title"].(string),
		Date:    date,
		Content: buf.Bytes(),
	}, nil
}
