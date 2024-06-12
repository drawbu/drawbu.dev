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
)

func Handler(articlesDir string) *handler {
	return &handler{
		articles: getArticles("", articlesDir),
	}
}

type handler struct {
	articles []article
}

type article struct {
	Title   string
	Date    string
	Path    string
	Content []byte
}

func (h *handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil || article_name == "" {
		return serv.Template(blog(h.articles)).Render(context.Background(), w)
	}

	a, err := findArticle(h.articles, "/"+article_name)
	if err != nil {
		return err
	}
	return serv.Template(articleShow(a)).Render(context.Background(), w)
}

func findArticle(articles []article, path string) (article, error) {
	for _, a := range articles {
		if a.Path == path {
			return a, nil
		}
	}
	return article{}, errors.New("Article not found")
}

func getArticles(path string, repo_path string) []article {
	fullpath := repo_path + "/" + path
	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return []article{}
	}

	articles := []article{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if entry.IsDir() {
			articles = append(articles, getArticles(path+"/"+entry.Name(), repo_path)...)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			articles = append(articles, article{
				Title:   name,
				Path:    path + "/" + name,
				Content: parseMarkdownArticle(fullpath + "/" + entry.Name()),
			})
		}
	}
	return articles
}

func parseMarkdownArticle(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("catppuccin-mocha"),
				highlighting.WithFormatOptions(
					chroma.WithLineNumbers(true),
				),
			),
		),
	)
	var buf bytes.Buffer
	if err := markdown.Convert(file, &buf); err != nil {
		return []byte{}
	}
	return buf.Bytes()
}
