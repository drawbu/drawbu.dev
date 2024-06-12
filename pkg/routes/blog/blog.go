package blog

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"app/pkg/app"

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Handler(articlesDir string) *handler {
	return &handler{
		articles: getArticles(articlesDir),
	}
}

type handler struct {
	articles []article
}

type article struct {
	Title   string
	Date    string
	Content []byte
}

func (h *handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil || article_name == "" {
		return serv.Template(blog(h.articles)).Render(context.Background(), w)
	}

	a, err := findArticle(h.articles, article_name)
	if err != nil {
		return err
	}
	return serv.Template(articleShow(*a)).Render(context.Background(), w)
}

func findArticle(articles []article, title string) (*article, error) {
	for _, a := range articles {
		if a.Title == title {
			return &a, nil
		}
	}
	return nil, errors.New("Article not found")
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
	return &article{
		Title:   metaData["title"].(string),
		Date:    metaData["date"].(string),
		Content: buf.Bytes(),
	}, nil
}
