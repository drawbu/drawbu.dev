package blog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"

	"app/pkg/app"
)

type Handler struct {
	GithubRepo string
	RepoPath   string
	Articles   []article
}

type article struct {
	Title   string
	Date    string
	Path    string
	Content []byte
}

func (h *Handler) Render(serv *app.Server, w http.ResponseWriter, r *http.Request) error {
	article_name, err := url.PathUnescape(r.PathValue("article"))
	if err != nil || article_name == "" {
		return serv.Template(blog(h.Articles)).Render(context.Background(), w)
	}

	a, err := findArticle(h.Articles, "/"+article_name)
	if err != nil {
		return err
	}
	return serv.Template(articleShow(a)).Render(context.Background(), w)
}

func (h *Handler) FetchArticles() {
	ticker := time.NewTicker(1 * time.Hour)
	h.fetch()
	for range ticker.C {
		h.fetch()
	}
}

func (h *Handler) fetch() {
	if h.RepoPath == "" {
		h.GithubRepo = "https://github.com/drawbu/Notes"
		path, err := cloneRepo(h.GithubRepo)
		if err != nil {
			fmt.Println(err)
			return
		}
		h.RepoPath = path
	} else {
		fmt.Printf("Pulling from %s\n", h.RepoPath)
		err := exec.Command("git", "-C", h.RepoPath, "pull").Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	h.Articles = getArticles("", h.RepoPath)
}

func findArticle(articles []article, path string) (article, error) {
	for _, a := range articles {
		if a.Path == path {
			return a, nil
		}
	}
	return article{}, errors.New("Article not found")
}

func cloneRepo(repo string) (string, error) {
	// Set path
	dirname := os.TempDir()
	err := os.Mkdir(dirname, 0777)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	if !strings.HasSuffix(dirname, "/") {
		dirname += "/"
	}
	dirname += "drawbu-blog"

	// Check if path exists
	_, err = os.Stat(dirname)
	if err == nil {
		err = os.RemoveAll(dirname)
		if err != nil {
			return "", err
		}
	}

	// Clone
	fmt.Printf("Cloning to %s\n", dirname)
	out, err := exec.Command("git", "clone", repo, dirname).CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return "", errors.New("git: " + err.Error())
	}
	return dirname, nil
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
	// create markdown parser with extensions
	renderer := newCustomizedRender()
	return markdown.ToHTML(file, nil, renderer)
}
