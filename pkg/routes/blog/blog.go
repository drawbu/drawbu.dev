package blog

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"app/articles"
	"app/pkg/app"

	"github.com/a-h/templ"
	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

var Articles = make(map[string]article)

func init() {
	for _, a := range getArticles(articles.Articles) {
		Articles[a.Title] = a
	}
}

type article struct {
	Title   string
	Date    time.Time
	Content []byte
}

func Handler(serv *app.Server, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil || article_name == "" {
		serv.Cache_route(w, r, 3600)
		return blog(getSortedArticles(Articles)), nil
	}

	a, ok := Articles[article_name]
	if !ok {
		return nil, errors.New("Article not found")
	}
	serv.Cache_route(w, r, 3600)
	return articleShow(a), nil
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
		a, err := parseMarkdownArticle(file)
		if err == nil || a != nil {
			result = append(result, *a)
		}
	}
	return result
}

func parseMarkdownArticle(file fs.File) (*article, error) {
	info, err := file.Stat()
	if err != nil {
		log.Warn("Error getting file info:", "reason", err)
		return nil, err
	}
	content := make([]byte, info.Size())
	if _, err = file.Read(content); err != nil {
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
	if err = markdown.Convert(content, &buf, parser.WithContext(context)); err != nil {
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
