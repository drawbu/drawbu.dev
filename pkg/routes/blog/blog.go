package blog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"

	"app/pkg/app"
	"app/pkg/components"
)

type Handler struct{
	GithubRepo string
	RepoPath   string
	LastLookup time.Time
	Articles   []article
}

type article struct {
	Title   string
	Date    string
	Path    string
	URI     string
	Content []byte
}

func (h* Handler) Render (serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "/blog" || r.URL.Path == "/blog/" {
		return components.Template(blog(h.Articles)).Render(context.Background(), w)
	}

	a, err := findArticle(h.Articles, r.URL.Path)
	if err != nil {
		return err
	}
	return components.Template(articleShow(a)).Render(context.Background(), w)
}


func (h* Handler) FetchArticles() {
    ticker := time.NewTicker(1 * time.Hour)
    h.fetch()
    for range ticker.C {
        h.fetch()
    }
}

func (h* Handler) fetch() {
	if h.RepoPath == "" {
		h.GithubRepo = "https://github.com/drawbu/Notes"
		path, err := cloneRepo(h.GithubRepo)
		if err != nil {
            fmt.Println(err)
			return
		}
		h.RepoPath = path
		h.Articles = getArticles(h.RepoPath, h.RepoPath)
        return
	}

    fmt.Printf("Pulling from %s\n", h.RepoPath)
    err := exec.Command("git", "-C", h.RepoPath, "pull").Run()
    if err != nil {
        fmt.Println(err)
        return
    }
    h.LastLookup = time.Now()
    h.Articles = getArticles(h.RepoPath, h.RepoPath)
}

func findArticle(articles []article, path string) (article, error) {
	for _, a := range articles {
		if a.URI == path {
			return a, nil
		}
	}
	return article{}, errors.New("Article not found")
}

func cloneRepo(repo string) (string, error) {
	// Set path
	dirname := os.TempDir()
	if !strings.HasSuffix(dirname, "/") {
		dirname += "/"
	}
	dirname += "drawbu-blog"

	// Check if path exists
	_, err := os.Stat(dirname)
	if err == nil {
		err = os.RemoveAll(dirname)
		if err != nil {
			return "", err
		}
	}

	// Clone
	fmt.Printf("Cloning to %s\n", dirname)
	err = exec.Command("git", "clone", repo, dirname).Run()
	if err != nil {
		return "", errors.New("git: " + err.Error())
	}
	return dirname, nil
}

func getArticles(path string, basepath string) []article {
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
			articles = append(articles, getArticles(path+"/"+entry.Name(), basepath)...)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			uri := "/blog" + strings.TrimPrefix(path, basepath) + "/" + name
			filepath := path + "/" + entry.Name()
			articles = append(articles, article{Title: name, Path: filepath, URI: uri, Content: getContent(filepath)})
		}
	}
	return articles
}

func getContent(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	// create markdown parser with extensions
	renderer := newCustomizedRender()
	return markdown.ToHTML(file, nil, renderer)
}
