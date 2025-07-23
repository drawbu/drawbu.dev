package blog

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/url"
	"strings"
	"time"

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	mdParser "github.com/yuin/goldmark/parser"
)

type article struct {
	Title   string
	Date    time.Time
	Content []byte
	Uri     string
}

// TODO: Add metadata validation

func NewArticle(file fs.File) (*article, error) {

	var buf bytes.Buffer
	context, err := fileToMarkdown(file, &buf)
	if err != nil {
		return nil, fmt.Errorf("Could not convert to markdown: %s", err)
	}

	metadata := meta.Get(context)

	date, err := time.Parse("2006-01-02", metadata["date"].(string))
	if err != nil {
		return nil, fmt.Errorf("Could not parse as date : %s", err)
	}

	uri, err := makeArticleUri(date, metadata)
	if err != nil {
		return nil, fmt.Errorf("Could make article uri: %s", err)
	}

	return &article{
		Title:   metadata["title"].(string),
		Date:    date,
		Content: buf.Bytes(),
		Uri: uri,
	}, nil
}

func fileToMarkdown(file fs.File, buf *bytes.Buffer) (mdParser.Context, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Error getting file info: %s", err)
	}
	log.Info("Opening new article", "name", info.Name)

	content := make([]byte, info.Size())
	if _, err = file.Read(content); err != nil {
		return nil, fmt.Errorf("Could not read file: %s", err)
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
	context := mdParser.NewContext()
	if err = markdown.Convert(content, buf, mdParser.WithContext(context)); err != nil {
		return nil, fmt.Errorf("Could not converting file to markdown: %s", err)
	}
	return context, nil
}

func makeArticleUri(date time.Time, metadata map[string]any) (string, error) {

	metaUri := metadata["uri"].(string)

	if metaUri == "" {
		metaTitle := metadata["title"].(string)
		if metaTitle == "" {
			return "", fmt.Errorf("Either one of 'uri' or 'title' needs to be set")
		}
		metaUri = url.QueryEscape(strings.ReplaceAll(strings.ToLower(metaTitle), " ", "-"))
	}

	return fmt.Sprintf("%s/%s", date.Format("2006/01"), metaUri), nil
}
