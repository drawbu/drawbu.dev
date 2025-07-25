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
	"gopkg.in/yaml.v3"
)

type article struct {
	Title       string
	Date        time.Time
	Content     []byte
	Uri         string
	Author      articleAuthor
	Description string
}

type articleAuthor struct {
	Name  string
	Email string
}

type articleMetadata struct {
	Title  string      `yaml:"title"`
	Date   ArticleDate `yaml:"date"`
	Uri    string      `yaml:"uri"`
	Author struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email,omitempty"`
	} `yaml:"author"`
	Description string `yaml:"description"`
}

func NewArticle(file fs.File) (*article, error) {

	var buf bytes.Buffer
	context, err := fileToMarkdown(file, &buf)
	if err != nil {
		return nil, fmt.Errorf("Could not convert to markdown: %s", err)
	}

	metadata, err := getMetadata(context)
	if err != nil {
		return nil, fmt.Errorf("Could not get article metadata: %s", err)
	}

	uri, err := makeArticleUri(*metadata)
	if err != nil {
		return nil, fmt.Errorf("Could make article uri: %s", err)
	}

	return &article{
		Title:   metadata.Title,
		Date:    metadata.Date.Time,
		Content: buf.Bytes(),
		Uri:     uri,
		Author: articleAuthor{
			Email: metadata.Author.Email,
			Name:  metadata.Author.Name,
		},
		Description: metadata.Description,
	}, nil
}

func getMetadata(context mdParser.Context) (*articleMetadata, error) {
	metadata, err := meta.TryGet(context)
	if err != nil {
		return nil, fmt.Errorf("Could not get article metadata: %s", err)
	}

	out, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("Could not Marshal metadata: %s", err)
	}
	var final articleMetadata
	if err = yaml.Unmarshal(out, &final); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal metadata: %s", err)
	}

	return &final, nil
}

func fileToMarkdown(file fs.File, buf *bytes.Buffer) (mdParser.Context, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Error getting file info: %s", err)
	}
	log.Info("Opening new article", "name", info.Name())

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

func makeArticleUri(metadata articleMetadata) (string, error) {
	uri := metadata.Uri
	if uri == "" {
		if metadata.Title == "" {
			return "", fmt.Errorf("Either one of 'uri' or 'title' needs to be set")
		}
		uri = url.QueryEscape(strings.ReplaceAll(strings.ToLower(metadata.Title), " ", "-"))
	}

	return fmt.Sprintf("%s/%s", metadata.Date.Time.Format("2006/01"), uri), nil
}

type ArticleDate struct {
	Time time.Time
}

func (t *ArticleDate) UnmarshalYAML(unmarshal func(any) error) error {

	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return fmt.Errorf("Could not unmarshal date: %s", err)
	}

	tt, err := time.Parse("2006-01-02", strings.TrimSpace(buf))
	if err != nil {
		return fmt.Errorf("Could not parse date: %s", err)
	}
	t.Time = tt
	return nil
}
